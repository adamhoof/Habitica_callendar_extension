package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TaskList struct {
	Success bool   `json:"success"`
	Data    []Task `json:"data"`
}

type Task struct {
	Text   string          `json:"text"`
	Repeat map[string]bool `json:"repeat"`
}

func parseTasksListFromResponse(response *http.Response) (TaskList, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing body")
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	var taskList TaskList
	err = json.Unmarshal(body, &taskList)

	return taskList, err
}
