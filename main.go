package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

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

func getTomorrowsTasks(tomorrowDayNumber uint8, list *TaskList) string {
	dayShortName := convertDayNumberToShortString(tomorrowDayNumber)
	var tasksString string
	fmt.Printf("Tomorrows (%s) tasks...\n", dayShortName)

	for _, task := range list.Data {
		if task.Repeat[dayShortName] == false {
			continue
		}

		tasksString += task.Text
		tasksString += "\n"
	}
	return tasksString
}

func fetchTasksButtonHandler() (tomorrowsTasksListString string) {
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

	return getTomorrowsTasks(getTomorrowsDayNumber(), &taskList)
}

func main() {
	err := godotenv.Load("env.list")

	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err.Error())
	}

	a := app.New()
	w := a.NewWindow("Habitica Callendar Extension")
	var tomorrowsTasksListString string
	tg := widget.NewTextGrid()

	w.SetContent(container.NewVBox(
		widget.NewButton("Fetch tasks", func() {
			tomorrowsTasksListString = fetchTasksButtonHandler()
			tg.SetText(tomorrowsTasksListString)
		}),
		tg,
	))
	w.ShowAndRun()
}
