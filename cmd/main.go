package main

import (
	"TaskService"
	"TaskService/comments"
	_ "TaskService/comments"
	"TaskService/processing"
	"TaskService/storage"
	"TaskService/tasks"
	_ "TaskService/tasks"
	"TaskService/users"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type (
	Config struct {
		DB                storage.ConnectionConfig `json:"db"`
		Key               string                   `json:"key"`
		UsersTableName    string                   `json:"users_table_name"`
		TasksTableName    string                   `json:"tasks_table_name"`
		CommentsTableName string                   `json:"comments_table_name"`
		HttpPort          int                      `json:"http_port"`
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

	userSvc := users.NewService(db, config.UsersTableName, config.Key)
	taskSvc := tasks.NewService(db, config.TasksTableName)
	commentsSvc := comments.NewService(db, config.CommentsTableName)

	proc := processing.NewService(userSvc, taskSvc, commentsSvc)
	httpHandlerSvc := TaskService.NewAppHandler(proc)

	mux := http.NewServeMux()
	httpHandlerSvc.SetHandlersToMux(mux)

	http.ListenAndServe(fmt.Sprintf(":%d", config.HttpPort), mux)
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
