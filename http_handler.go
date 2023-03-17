package TaskService

import (
	"TaskService/processing"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
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

func (app *AppHandler) handleCreateTask(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	type createTaskRequest struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
}
