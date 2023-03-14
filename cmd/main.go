package main

import (
	"TaskService/storage"
	"TaskService/tasks"
	_ "TaskService/tasks"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type (
	Config struct {
		DB             storage.ConnectionConfig `json:"db"`
		Key            string                   `json:"key"`
		UsersTableName string                   `json:"users_table_name"`
	}
)

func main() {
	config, err := GetConfig()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := storage.Connect(config.DB)
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}
	taskService := tasks.NewTask(db, config.UsersTableName)
	uuidID, err := uuid.Parse("2fd134a9-85ec-44cd-9606-020354f03e9f")
	if err != nil {
		fmt.Println("Ошибка")
	}
	r := tasks.CreateTaskRequest{
		User_id:     uuidID,
		Tittle:      "Первая заметка",
		Description: "Сделать такс менеджер",
	}
	_, err = taskService.CreateNewTask(r)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Заметка создана")

	//userService := users.NewService(db, config.UsersTableName, config.Key)
	//taskService := tasks.NewTask(db, "tasks")
	//
	//mux := http.NewServeMux()
	//mux.HandleFunc("/users/create", func(writer http.ResponseWriter, request *http.Request) {
	//	if request.Method != http.MethodPost {
	//		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	//		return
	//	}
	//
	//	rawBody, err := io.ReadAll(request.Body)
	//	defer request.Body.Close()
	//	if err != nil {
	//		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	//		return
	//	}
	//
	//	type CreateUserRequest struct {
	//		Name string `json:"username"`
	//		Pass string `json:"password"`
	//	}
	//	r := CreateUserRequest{}
	//
	//	if err = json.Unmarshal(rawBody, &r); err != nil {
	//		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	//		return
	//	}
	//
	//	result, err := userService.Create(r.Name, r.Pass)
	//	if err != nil {
	//		fmt.Println(err)
	//		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//		return
	//	}
	//
	//	fmt.Printf("New user created! %s\n", result.ID.String())
	//})
	//
	//mux.HandleFunc("/users/get_all", func(writer http.ResponseWriter, request *http.Request) {
	//	result, err := userService.GetAllUsers()
	//	if err != nil {
	//		fmt.Println(err)
	//		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//		return
	//	}
	//
	//	res, err := json.Marshal(result)
	//	if err != nil {
	//		fmt.Println(err)
	//		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//		return
	//	}
	//	writer.Write(res)
	//})
	//
	//mux.HandleFunc("/tasks/create", func(writer http.ResponseWriter, request *http.Request) {
	//	if request.Method != http.MethodPost {
	//		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	//		return
	//	}
	//	rawBody, err := io.ReadAll(request.Body)
	//	defer request.Body.Close()
	//	if err != nil {
	//		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	//		return
	//	}
	//
	//	r := tasks.CreateTaskRequest{}
	//	if err = json.Unmarshal(rawBody, &r); err != nil {
	//		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	//		return
	//	}
	//
	//	result, err := taskService.CreateNewTask(r)
	//	if err != nil {
	//		fmt.Println(err)
	//		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//		return
	//	}
	//
	//	fmt.Printf("New task created! %s\n", result.Id.String())
	//
	//})

	//mux.HandleFunc("/tasks/change_tasks", func(writer http.ResponseWriter, request *http.Request) {
	//	result, err := userService.GetAllUsers()
	//	if err != nil {
	//		fmt.Println(err)
	//		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//		return
	//	}
	//
	//	res, err := json.Marshal(result)
	//	if err != nil {
	//		fmt.Println(err)
	//		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//		return
	//	}
	//	writer.Write(res)
	//})

	//http.ListenAndServe(":8000", mux)

}

func GetConfig() (Config, error) {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}

	c := Config{}
	if err = json.Unmarshal(file, &c); err != nil {
		return Config{}, err
	}

	return c, nil
}
