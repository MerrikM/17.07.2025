package main

import (
	"log"
	"workmate_test_project/internal/model"
	"workmate_test_project/internal/service"
)

func main() {
	taskService := service.NewTaskService()

	tasks := make([]*model.Task, 0)

	for i := 0; i < 4; i++ {
		task, err := taskService.CreateTask()
		if err != nil {
			log.Printf("ошибка создания задачи: %v", err)
		}
		tasks = append(tasks, task)
	}

	for i := 0; i < 3; i++ {
		err := taskService.AddFileToTask(tasks[i].ID, "file1.png")
		if err != nil {
			log.Printf("ошибка добавления файла к задаче: %v", err)
		}
	}

	log.Println("задача с 3-я файлами успешно создана")
}
