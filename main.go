package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	MON = 1
	TUE = 2
	WED = 3
	THU = 4
	FRI = 5
	SAT = 6
	SUN = 0
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

func getTomorrowsDayNumber() uint8 {
	currentTime := time.Now()

	tomorrow := currentTime.Add(24 * time.Hour)

	return uint8(tomorrow.Weekday())
}

func printTomorrowsTasks(tomorrowDayNumber uint8, list *TaskList) {
	fmt.Println(tomorrowDayNumber)

	for _, task := range list.Data {
		fmt.Printf("Task Name: %s\n", task.Text)
		fmt.Printf("Repeat: %v\n", task.Repeat)
	}
}

func main() {
	err := godotenv.Load("env.list")
	response, err := requestTasks()

	if err != nil {
		log.Fatalf("Request failed: %s\n", err.Error())
	}

	if response.StatusCode != 200 {
		log.Fatalf("Request failed: %s\n", response.Status)
	}

	taskList, err := parseTasksListFromResponse(response)

	if err != nil {
		log.Printf("Error parsing task list: %s\n", err.Error())
	}

	printTomorrowsTasks(getTomorrowsDayNumber(), &taskList)
}
