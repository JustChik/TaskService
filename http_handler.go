package TaskService

import (
	"TaskService/processing"
	"TaskService/tasks"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
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
