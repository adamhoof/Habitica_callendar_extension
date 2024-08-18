package main

import (
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

func getTasksForDay(day string, list *TaskList) string {
	dayShortName := shortenDayName(day)
	var tasksString string

	for _, task := range list.Data {
		if task.Repeat[dayShortName] == false {
			continue
		}

		tasksString += task.Text
		tasksString += "\n"
	}
	return tasksString
}

func fetchTasksButtonHandler(day string) (tomorrowsTasksListString string) {
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

	return getTasksForDay(day, taskList)
}

func main() {
	err := godotenv.Load("env.list")

	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err.Error())
	}

	a := app.New()
	w := a.NewWindow("Habitica Callendar Extension")
	var tasksListString string
	dayOptions := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	tg := widget.NewTextGrid()

	w.SetContent(container.NewVBox(
		widget.NewSelect(dayOptions, func(day string) {
			tasksListString = fetchTasksButtonHandler(day)
			tg.SetText(tasksListString)
		}),
		tg,
	))
	w.ShowAndRun()
}
