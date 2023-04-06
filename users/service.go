package users

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"net/mail"
	"strings"
	"time"
)

var (
	ErrorDuplicated       = errors.New("already exist")
	ErrorNotFound         = errors.New("not found")
	ErrorPass             = errors.New("wrong password")
	ErrorInvalidEmail     = errors.New("invalid email")
	ErrorInvalidOperation = errors.New("something go wrong")
	ErrorResetCode        = errors.New("wrong reset code")
)

const (
	conflictErrMessage = "pq: duplicate key value violates unique constraint \"users_username_key\""
)

type (
	Service struct {
		db        *sql.DB
		key       string
		tableName string
	}

	User struct {
		ID   uuid.UUID
		Name string
	}
)

func NewService(db *sql.DB, tableName, key string) *Service {
	return &Service{
		db:        db,
		key:       key,
		tableName: tableName,
	}
}

func (s *Service) Create(username, pass string) (*User, error) {
	if !validateEmail(username) {
		return nil, ErrorInvalidEmail
	}
	id := uuid.New()
	hashedPass := hashPass(s.key, pass)
	_, err := s.db.Exec(fmt.Sprintf("INSERT INTO %s (id,username,password) VALUES ($1, $2, $3);", s.tableName), id.String(), username, hashedPass)
	if err != nil {
		if err.Error() == conflictErrMessage {
			return nil, ErrorDuplicated
		}
		return nil, fmt.Errorf("can't creat user in db %v", err)
	}

	return &User{
		ID:   id,
		Name: username,
	}, nil
}

func (s *Service) GetUser(username, pass string) (*User, error) {
	if !validateEmail(username) {
		return nil, ErrorInvalidEmail
	}
	hashedPass := hashPass(s.key, pass)
	res := s.db.QueryRow(fmt.Sprintf("SELECT id, username, password FROM %s WHERE username = $1;", s.tableName), username)
	type RawUser struct {
		ID       uuid.UUID
		Username string
		Password string
	}
	u := RawUser{}
	err := res.Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("can't get user %v", err)
	}
	if hashedPass != u.Password {
		return nil, ErrorPass
	}

	return &User{
		ID:   u.ID,
		Name: u.Username,
	}, nil
}

func (s *Service) GetAllUsers() ([]User, error) {
	all, err := s.db.Query(fmt.Sprintf("SELECT id, username FROM %s;", s.tableName))
	if err != nil {
		return nil, fmt.Errorf("can't get users %v", err)
	}

	u := User{}
	res := []User{}
	for all.Next() {
		err = all.Scan(&u.ID, &u.Name)

		if err != nil {
			return nil, fmt.Errorf("can't read storage response %v", err)
		}

		res = append(res, u)
	}

	return res, nil
}

func (s *Service) ResetCode(username string) (string, error) {
	if !validateEmail(username) {
		return "", ErrorInvalidEmail
	}
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // Например "ExcbsVQs"
	fmt.Println(str)
	hashedPass := hashPass(str, username)
	_, err := s.db.Exec(fmt.Sprintf("UPDATE  %s SET resetpasswordtoken = $1 WHERE username = $2;", s.tableName), hashedPass, username)
	if err != nil {
		return "", ErrorInvalidOperation
	}
	return str, nil
}

func (s *Service) ChangePassword(username, code, newpass string) (string, error) {
	if !validateEmail(username) {
		return "", ErrorInvalidEmail
	}
	if len(code) != 8 {
		log.Println(code)
		return "", ErrorResetCode
	}
	var hashCode string
	chekCode := hashPass(code, username)
	rows, err := s.db.Query(fmt.Sprintf("SELECT resetpasswordtoken FROM %s WHERE username = $1;", s.tableName), username)
	if err != nil {
		return "", fmt.Errorf("can't get code %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&hashCode); err != nil {
			log.Fatal(err)
		}
	}
	if chekCode == hashCode {
		hashedPass := hashPass(s.key, newpass)
		_, err := s.db.Exec(fmt.Sprintf("UPDATE %s SET password = $1,resetpasswordtoken = $2 WHERE username = $3;", s.tableName), hashedPass, "", username)
		if err != nil {
			return "", ErrorInvalidOperation
		}
	}
	return "Пароль успешно изменён", nil
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func hashPass(key, pass string) string {
	hash := md5.Sum([]byte(key + pass))
	return hex.EncodeToString(hash[:])
}
