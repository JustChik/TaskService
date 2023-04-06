package TaskService

import (
	"TaskService/comments"
	"TaskService/processing"
	"TaskService/tasks"
	"TaskService/users"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type AppHandler struct {
	proc *processing.Service
}

func NewAppHandler(proc *processing.Service) *AppHandler {
	return &AppHandler{
		proc: proc,
	}
}

func (app *AppHandler) SetHandlersToMux(mux *http.ServeMux) {
	mux.HandleFunc("/users/create", app.handleCreateUser)
	mux.HandleFunc("/users/get", app.handleGetUser)
	mux.HandleFunc("/users/get/all", app.handleGetAllUsers)
	mux.HandleFunc("/tasks/create", app.handleCreateTask)
	mux.HandleFunc("/tasks/change/task", app.handleChangeTask)
	mux.HandleFunc("/tasks/change/status", app.handleChangeStatus)
	mux.HandleFunc("/tasks/get/task", app.handleGetTask)
	mux.HandleFunc("/tasks/create/comment", app.handleCreateComment)
	mux.HandleFunc("/tasks/get/comment", app.handleGetComments)
	mux.HandleFunc("/users/reset/gen", app.handleResetCodeGen)
	mux.HandleFunc("/users/reset/pass", app.handleResetPass)
}

func (app *AppHandler) handleCreateUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	type createUserRequest struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	createReq := createUserRequest{}
	if err = json.Unmarshal(rawBody, &createReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := app.proc.CreateUser(processing.UserInfo{
		UserName: createReq.UserName,
		UserPass: createReq.Password,
	})

	if err != nil {
		fmt.Printf("Can not create user: %v", err)
		if errors.Is(users.ErrorInvalidEmail, err) {
			http.Error(w, "invalid email", http.StatusBadRequest)
			log.Printf("invalid email:%s", createReq.UserName)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	type createUserResponse struct {
		ID uuid.UUID `json:"id"`
	}

	r := createUserResponse{
		ID: result.ID,
	}

	response, err := json.Marshal(&r)
	if err != nil {
		fmt.Printf("Can not marshal response whith creating user: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (app *AppHandler) handleGetUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	type getUserRequest struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	getReq := getUserRequest{}
	if err = json.Unmarshal(rawBody, &getReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := app.proc.GetUser(processing.UserInfo{
		UserName: getReq.UserName,
		UserPass: getReq.Password,
	})

	if err != nil {
		fmt.Printf("Can not get user: %v", err)
		if errors.Is(users.ErrorInvalidEmail, err) {
			http.Error(w, "invalid email", http.StatusBadRequest)
			log.Printf("invalid email:%s", getReq.UserName)
			return
		}
		if errors.Is(users.ErrorNotFound, err) {
			http.Error(w, "user doesn't exist", http.StatusNotFound)
			log.Printf("user doesn't exist:%s", getReq.UserName)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	type getUserResponse struct {
		ID       uuid.UUID `json:"id"`
		UserName string
	}

	r := getUserResponse{
		ID:       result.ID,
		UserName: result.Name,
	}
	response, err := json.Marshal(&r)
	if err != nil {
		fmt.Printf("Can not marshal response whith getting user: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (app *AppHandler) handleGetAllUsers(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	result, err := app.proc.GetAllUsers()

	if err != nil {
		fmt.Printf("Can not get users: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(&result)
	if err != nil {
		fmt.Printf("Can not marshal response whith getting users: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (app *AppHandler) handleCreateTask(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	auth := req.Header.Get("Authorization")
	auth = strings.ReplaceAll(auth, "Basic ", "")
	decodedAuth, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	authSlice := strings.Split(string(decodedAuth), ":")
	if len(authSlice) < 2 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	username := authSlice[0]
	pass := authSlice[1]

	user, err := app.proc.GetUser(processing.UserInfo{
		UserName: username,
		UserPass: pass,
	})

	if err != nil {
		fmt.Printf("Can not get user: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	type createTaskRequest struct {
		Tittle      string `json:"tittle"`
		Description string `json:"description"`
	}

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	createTaskReq := createTaskRequest{}
	if err := json.Unmarshal(rawBody, &createTaskReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := app.proc.CreateTask(processing.UserInfo{
		UserName: username,
		UserPass: pass,
	}, processing.TaskInfo{
		UserID:      user.ID,
		Tittle:      createTaskReq.Tittle,
		Description: createTaskReq.Description,
	})
	if err != nil {
		fmt.Printf("Can not create task: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	type createTaskResponse struct {
		Id          uuid.UUID
		Tittle      string
		Description string
		Status      tasks.Status
	}
	r := createTaskResponse{
		Id:          result.Id,
		Tittle:      result.Tittle,
		Description: result.Description,
		Status:      result.Status,
	}
	response, err := json.Marshal(&r)
	if err != nil {
		fmt.Printf("Can not marshal response whith creating task: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (app *AppHandler) handleChangeTask(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	auth := req.Header.Get("Authorization")
	auth = strings.ReplaceAll(auth, "Basic ", "")
	decodedAuth, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	authSlice := strings.Split(string(decodedAuth), ":")
	if len(authSlice) < 2 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	username := authSlice[0]
	pass := authSlice[1]

	type changeTaskRequest struct {
		UserID      uuid.UUID `json:"user_id"`
		TaskID      uuid.UUID `json:"task_id"`
		Tittle      string    `json:"tittle"`
		Description string    `json:"description"`
	}

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	changeTaskReq := changeTaskRequest{}
	if err := json.Unmarshal(rawBody, &changeTaskReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := app.proc.ChangeTask(processing.UserInfo{
		UserName: username,
		UserPass: pass,
	}, processing.ChangeTaskRequest{
		IDUser:      changeTaskReq.UserID,
		IDTask:      changeTaskReq.TaskID,
		Tittle:      changeTaskReq.Tittle,
		Description: changeTaskReq.Description,
	})
	if err != nil {
		fmt.Printf("Can not change task: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	type changeTaskResponse struct {
		Id          uuid.UUID
		Tittle      string
		Description string
	}
	r := changeTaskResponse{
		Id:          result.Id,
		Tittle:      result.Tittle,
		Description: result.Description,
	}
	response, err := json.Marshal(&r)
	if err != nil {
		fmt.Printf("Can not marshal response whith creating task: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (app *AppHandler) handleChangeStatus(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	auth := req.Header.Get("Authorization")
	auth = strings.ReplaceAll(auth, "Basic ", "")
	decodedAuth, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	authSlice := strings.Split(string(decodedAuth), ":")
	if len(authSlice) < 2 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	username := authSlice[0]
	pass := authSlice[1]

	type changeStatusRequest struct {
		TaskID uuid.UUID    `json:"task_id"`
		Status tasks.Status `json:"status"`
	}

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	changeStatusReq := changeStatusRequest{}
	if err := json.Unmarshal(rawBody, &changeStatusReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := app.proc.ChangeTaskStatus(processing.UserInfo{
		UserName: username,
		UserPass: pass,
	}, processing.ChangeStatusRequest{
		TaskID: changeStatusReq.TaskID,
		Status: changeStatusReq.Status,
	})
	if err != nil {
		fmt.Printf("Can not change status: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	type changeStatusResponse struct {
		Status tasks.Status
	}
	r := changeStatusResponse{
		Status: result,
	}
	response, err := json.Marshal(&r)
	if err != nil {
		fmt.Printf("Can not marshal response whith creating task: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (app *AppHandler) handleGetTask(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	auth := req.Header.Get("Authorization")
	auth = strings.ReplaceAll(auth, "Basic ", "")
	decodedAuth, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	authSlice := strings.Split(string(decodedAuth), ":")
	if len(authSlice) < 2 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	username := authSlice[0]
	pass := authSlice[1]

	type getTaskRequest struct {
		Tittle      string       `json:"tittle"`
		Description string       `json:"description"`
		Status      tasks.Status `json:"status"`
	}

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	getReq := getTaskRequest{}

	if err = json.Unmarshal(rawBody, &getReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := app.proc.GetTasks(processing.UserInfo{
		UserName: username,
		UserPass: pass,
	}, processing.GetTasksRequest{
		Tittle:      getReq.Tittle,
		Description: getReq.Description,
		Status:      getReq.Status,
	})

	if err != nil {
		fmt.Printf("Can not get task: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	all := []tasks.Task{}
	for _, el := range result {
		all = append(all, el)
	}
	response, err := json.Marshal(&all)
	if err != nil {
		fmt.Printf("Can not marshal response whith getting task: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (app *AppHandler) handleCreateComment(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	auth := req.Header.Get("Authorization")
	auth = strings.ReplaceAll(auth, "Basic ", "")
	decodedAuth, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	authSlice := strings.Split(string(decodedAuth), ":")
	if len(authSlice) < 2 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	username := authSlice[0]
	pass := authSlice[1]

	type createCommentRequest struct {
		TaskID uuid.UUID `json:"task_id"`
		Text   string    `json:"text"`
	}

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	createReq := createCommentRequest{}
	if err = json.Unmarshal(rawBody, &createReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := app.proc.CreateComment(processing.UserInfo{
		UserName: username,
		UserPass: pass,
	}, processing.CreateCommentRequest{
		TaskID: createReq.TaskID,
		Text:   createReq.Text,
	})

	if err != nil {
		fmt.Printf("Can not create comment: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	type createCommentResponse struct {
		Id     uuid.UUID
		Text   string
		Data   time.Time
		TaskID uuid.UUID
	}

	r := createCommentResponse{
		Id:     result.Id,
		Text:   result.Text,
		Data:   result.Data,
		TaskID: result.TaskID,
	}

	response, err := json.Marshal(&r)
	if err != nil {
		fmt.Printf("Can not marshal response whith creating comment: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (app *AppHandler) handleGetComments(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	auth := req.Header.Get("Authorization")
	auth = strings.ReplaceAll(auth, "Basic ", "")
	decodedAuth, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	authSlice := strings.Split(string(decodedAuth), ":")
	if len(authSlice) < 2 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	username := authSlice[0]
	pass := authSlice[1]

	type getCommentRequest struct {
		TaskID uuid.UUID `json:"task_id"`
	}

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	getReq := getCommentRequest{}
	if err = json.Unmarshal(rawBody, &getReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := app.proc.GetComments(processing.UserInfo{
		UserName: username,
		UserPass: pass,
	}, processing.GetCommentRequst{
		TaskID: getReq.TaskID})

	if err != nil {
		fmt.Printf("Can not get comment: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	all := []comments.GetCommentRequest{}

	for _, el := range result {
		all = append(all, el)
	}

	response, err := json.Marshal(&all)
	if err != nil {
		fmt.Printf("Can not marshal response whith getting comment: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (app *AppHandler) handleResetCodeGen(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	type resetCodeRequest struct {
		UserName string `json:"user_name"`
	}
	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	resetReq := resetCodeRequest{}
	if err = json.Unmarshal(rawBody, &resetReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := app.proc.ResetCodeUser(processing.UserInfo{
		UserName: resetReq.UserName,
	})
	if err != nil {
		fmt.Printf("Can not generate reset code: %v", err)
		if errors.Is(users.ErrorInvalidEmail, err) {
			http.Error(w, "invalid email", http.StatusBadRequest)
			log.Printf("invalid email:%s", resetReq.UserName)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	type resetCodeRespose struct {
		Code string `json:"code"`
	}
	r := resetCodeRespose{
		Code: result,
	}
	response, err := json.Marshal(&r)
	if err != nil {
		fmt.Printf("Can not marshal response whith generate reset code: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (app *AppHandler) handleResetPass(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	type resetPassRequest struct {
		UserName string `json:"user_name"`
		Code     string `json:"code"`
		NewPass  string `json:"new_pass"`
	}

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	resetReq := resetPassRequest{}
	if err = json.Unmarshal(rawBody, &resetReq); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	result, err := app.proc.ChangePasswordUser(processing.ChangePass{
		UserName: resetReq.UserName,
		Code:     resetReq.Code,
		NewPass:  resetReq.NewPass,
	})
	if err != nil {
		fmt.Printf("Can not reset password: %v", err)
		if errors.Is(users.ErrorInvalidEmail, err) {
			http.Error(w, "invalid email", http.StatusBadRequest)
			log.Printf("invalid email:%s", resetReq.UserName)
			return
		}
		if errors.Is(users.ErrorResetCode, err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("invalid code:%s", resetReq.Code)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	type resetCodeRespose struct {
		Successfully string `json:"successfully"`
	}
	r := resetCodeRespose{
		Successfully: result,
	}
	response, err := json.Marshal(&r)
	if err != nil {
		fmt.Printf("Can not marshal response whith reset password: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
