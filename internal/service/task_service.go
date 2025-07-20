package service

import (
	"fmt"
	"path/filepath"
	"sync"
	"workmate_test_project/internal/model"
)

// TaskService - сервис для работы с задачами, он состоит из:
// id - счётчик для генерации уникальных идентификаторов задач
// tasksSlot - буферизированный канал, ограничивающий максимальное количество активных задач (3 по ТЗ)
// tasks - мапа активных задач, где ключ — id задачи, а значение — указатель на структуру model.Task
// mutex - мьютекс для защиты от гонки данных
type TaskService struct {
	id        int
	tasksSlot chan struct{}
	tasks     map[int]*model.Task
	mutex     sync.Mutex
}

var fileExtension = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".webp": {},
	".pdf":  {},
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasksSlot: make(chan struct{}, 3),
		tasks:     make(map[int]*model.Task),
		mutex:     sync.Mutex{},
	}
}

func (service *TaskService) CreateTask() (*model.Task, error) {
	select {
	case service.tasksSlot <- struct{}{}:
		service.mutex.Lock()

		service.id++
		task := &model.Task{
			ID:               service.id,
			Files:            []string{},
			FileCountChannel: make(chan string, 3),
			DoneChannel:      make(chan struct{}),
		}
		service.tasks[task.ID] = task
		service.mutex.Unlock()

		return task, nil

	default:
		return nil, fmt.Errorf("сервер в данный момент занят")
	}
}

func (service *TaskService) AddFileToTask(taskId int, file string) error {
	extension := filepath.Ext(file)
	if _, exist := fileExtension[extension]; exist == false {
		return fmt.Errorf("не поддерживаемое расширение файла")
	}

	var task *model.Task

	service.mutex.Lock()
	task, exist := service.tasks[taskId]
	service.mutex.Unlock()

	if exist == false {
		return fmt.Errorf("задача с id = %d не найдена", taskId)
	}

	select {
	case task.FileCountChannel <- file:
		return nil
	default:
		return fmt.Errorf("одновременно может обрабатываться только 3 файла")
	}
}
