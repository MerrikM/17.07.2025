package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"workmate_test_project/internal/service"
)

type TaskHandler struct {
	*service.TaskService
}

// TaskStatusResponse представляет собой ответ сервера со статусом задачи.
//
// Используется в ответе на GET-запрос получения статуса задачи.
// ArchiveLink будет непустым, только если задача завершена.
type TaskStatusResponse struct {
	TaskID      int    `json:"taskID" example:"1"`
	Status      string `json:"status" example:"завершена"`
	ArchiveLink string `json:"archiveLink" example:"G:/GithubRepo/17.07.2025/internal/util/task_1.zip"`
}

// CreateTaskRequest содержит путь и имя архива, который будет создан для задачи.
type CreateTaskRequest struct {
	ZipArchivePath string `json:"zipArchivePath" example:"G:/GithubRepo/17.07.2025/internal/util"`
	ZipArchiveName string `json:"zipArchiveName" example:"test1"`
}

// CreateTaskResponse возвращает ID созданной задачи.
type CreateTaskResponse struct {
	Message string `json:"message" example:"id вашей задачи: "`
	TaskID  int    `json:"taskID" example:"1"`
}

// AddFileToTaskRequest содержит параметры запроса для добавления файла к задаче.
type AddFileToTaskRequest struct {
	TaskID   int    `json:"taskID" example:"1"`
	FileURL  string `json:"fileURL" example:"https://example.com/file.pdf"`
	FileName string `json:"fileName" example:"test3"`
}

// AddFileToTaskResponse содержит ответ после успешного добавления файла.
type AddFileToTaskResponse struct {
	Message string `json:"message" example:"файлы успешно добавлен к вашей задаче"`
	TaskID  int    `json:"taskID" example:"1"`
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService}
}

// GetTaskStatusById возвращает статус задачи по её ID.
//
// @Summary Получить статус задачи
// @Description Возвращает статус задачи и ссылку на архив (если все файлы добавлены).
// @Tags tasks
// @Accept json
// @Produce json
// @Param task-id query int true "ID задачи"
// @Success 200 {object} TaskStatusResponse
// @Failure 400 {string} string "некорректный ID задачи или задача не найдена"
// @Router /get [get]
func (handler *TaskHandler) GetTaskStatusById(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 4*time.Second)
	defer cancel()

	taskIdStr := request.URL.Query().Get("task-id")

	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		http.Error(writer, "некорректный ID задачи", http.StatusBadRequest)
		return
	}

	task, err := handler.TaskService.GetTaskStatusById(ctx, taskId)
	if err != nil {
		log.Printf("задача не найден: %v", err)
		http.Error(writer, "задача не была найдена", http.StatusBadRequest)
		return
	}

	response := &TaskStatusResponse{
		TaskID: task.ID,
		Status: task.Status,
	}

	if task.FilesAdded == 3 {
		response = &TaskStatusResponse{
			TaskID:      task.ID,
			Status:      task.Status,
			ArchiveLink: task.ArchiveLink,
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(&response)
}

// CreateTask создаёт новую задачу с указанным путем и именем архива.
//
// @Summary      Создание новой задачи
// @Description  Создаёт задачу, для которой можно добавлять файлы в ZIP архив.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body CreateTaskRequest true "Путь и имя архива"
// @Success      200 {object} CreateTaskResponse "Успешный ответ с ID созданной задачи"
// @Failure      400 {string} string "Неверный формат JSON"
// @Failure      503 {string} string "Ошибка создания задачи"
// @Router       /create-task [post]
func (handler *TaskHandler) CreateTask(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 4*time.Second)
	defer cancel()

	var createTaskRequest CreateTaskRequest
	if err := json.NewDecoder(request.Body).Decode(&createTaskRequest); err != nil {
		http.Error(writer, "неверный формат json", http.StatusBadRequest)
		return
	}

	task, err := handler.TaskService.CreateTask(ctx, createTaskRequest.ZipArchivePath, createTaskRequest.ZipArchiveName)
	if err != nil {
		log.Printf("ошибка создания задачи: %v", err)
		http.Error(writer, "сервер в данный момент занят", http.StatusServiceUnavailable)
		return
	}

	response := &CreateTaskResponse{
		Message: "id вашей задачи: ",
		TaskID:  task.ID,
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(&response)
}

// AddFileToTask добавляет файл к задаче по её ID.
//
// @Summary      Добавить файл к задаче
// @Description  Добавляет файл в архив задачи, ограничение — максимум 3 файла на задачу.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body AddFileToTaskRequest true "Данные для добавления файла"
// @Success      200 {object} AddFileToTaskResponse "Файл успешно добавлен к задаче"
// @Failure      400 {string} string "Неверный формат JSON или превышен лимит файлов"
// @Router       /add-file-to-task [post]
func (handler *TaskHandler) AddFileToTask(writer http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), 4*time.Second)
	defer cancel()

	var addFileToTaskRequest AddFileToTaskRequest
	if err := json.NewDecoder(request.Body).Decode(&addFileToTaskRequest); err != nil {
		http.Error(writer, "неверный формат json", http.StatusBadRequest)
		return
	}

	err := handler.TaskService.AddFileToTask(
		ctx, addFileToTaskRequest.TaskID, addFileToTaskRequest.FileURL, addFileToTaskRequest.FileName,
	)
	if err != nil {
		log.Printf("ошибка добавления файла к задаче: %v", err)
		http.Error(writer, "в задаче может быть максимум 3 файла", http.StatusBadRequest)
		return
	}

	response := AddFileToTaskResponse{
		Message: "файлы успешно добавлен к вашей задаче",
		TaskID:  addFileToTaskRequest.TaskID,
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(&response)
}
