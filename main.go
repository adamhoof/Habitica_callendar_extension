package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"net/http"
)

func requestTasks(id *string, apiKey *string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://habitica.com/api/v3/tasks/user?type=dailys", nil)

	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("x-api-user", *id)
	request.Header.Add("x-api-key", *apiKey)

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

func fetchTasksButtonHandler(day string, id *string, apiKey *string) (tomorrowsTasksListString string) {
	response, err := requestTasks(id, apiKey)

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
	a := app.New()
	credW := a.NewWindow("Enter credentials")
	credW.Show()

	idInput := widget.NewEntry()
	idInput.SetPlaceHolder("username")
	apiKeyInput := widget.NewEntry()
	apiKeyInput.SetPlaceHolder("password")
	apiKeyInput.Password = true
	content := container.NewVBox(idInput, apiKeyInput, widget.NewButton("Login", func() {
		w := a.NewWindow("Habitica Callendar Extension")

		var tasksListString string
		dayOptions := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		tg := widget.NewTextGrid()

		w.SetContent(container.NewVBox(
			widget.NewSelect(dayOptions, func(day string) {
				tasksListString = fetchTasksButtonHandler(day, &idInput.Text, &apiKeyInput.Text)
				tg.SetText(tasksListString)
			}),
			tg,
		))
		w.Show()
		w.SetMaster()
		credW.Close()
	}))
	credW.SetContent(content)
	credW.ShowAndRun()
}
