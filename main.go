package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type TaskList struct {
	Success bool   `json:"success"`
	Data    []Task `json:"data"`
}

type Task struct {
	Text   string          `json:"text"`
	Repeat map[string]bool `json:"repeat"`
}

func requestTasks() (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://habitica.com/api/v3/tasks/user?type=dailys", nil)

	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("x-api-user", os.Getenv("id"))
	request.Header.Add("x-api-key", os.Getenv("key"))

	return client.Do(request)
}

func parseTasksList(response *http.Response) (TaskList, error) {
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

func getTomorrowsDayNumber() uint8 {
	currentTime := time.Now()

	tomorrow := currentTime.Add(72 * time.Hour)

	dayNumber := uint8(tomorrow.Weekday())

	if dayNumber == 0 {
		dayNumber = 7
	}
	return dayNumber

}

func printTomorrowsTasks(list *TaskList) {
	for _, task := range list.Data {
		fmt.Printf("Task Name: %s\n", task.Text)
		fmt.Printf("Repeat: %v\n", task.Repeat)
	}
}

func main() {
	response, err := requestTasks()

	if err != nil {
		log.Fatalf("Request failed: %s\n", err.Error())
	}

	taskList, err := parseTasksList(response)

	if err != nil {
		log.Printf("Error parsing task list: %s\n", err.Error())
	}

	fmt.Println(getTomorrowsDayNumber())
	printTomorrowsTasks(&taskList)
}
