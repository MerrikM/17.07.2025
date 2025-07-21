package service

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"workmate_test_project/internal/model"
	"workmate_test_project/internal/util"
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

// GetTaskStatusById возвращает задачу по её ID.
// Если задача с таким ID не найдена, возвращается ошибка.
func (service *TaskService) GetTaskStatusById(ctx context.Context, taskId int) (*model.Task, error) {
	task, exist := service.tasks[taskId]
	if exist == false {
		return nil, fmt.Errorf("задача с id = %d не найдена", taskId)
	}

	return task, nil
}

// CreateTask создает новую задачу с архивом ZIP в указанном пути и имени.
// Метод использует контекст для отмены операции и ограничивает
// количество одновременно создаваемых задач через канал tasksSlot.
// Возвращает созданную задачу или ошибку, если архив не удалось создать
// или сервер занят.
func (service *TaskService) CreateTask(ctx context.Context, zipArchivePath string, zipArchiveName string) (*model.Task, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case service.tasksSlot <- struct{}{}:
		service.mutex.Lock()
		defer service.mutex.Unlock()

		service.id++

		archiveFile, zipWriter, err := util.CreateZIPArchive(zipArchivePath, zipArchiveName)
		if err != nil {
			<-service.tasksSlot
			return nil, fmt.Errorf("ошибка создания архива: %w", err)
		}

		task := &model.Task{
			ID:               service.id,
			Files:            []string{},
			FileCountChannel: make(chan struct{}, 3),
			DoneChannel:      make(chan struct{}),
			ArchiveFile:      archiveFile,
			ArchiveWriter:    zipWriter,
			Status:           "создана",
			ArchiveLink:      zipArchivePath + "/" + zipArchiveName + ".zip",
		}
		service.tasks[task.ID] = task

		return task, nil

	default:
		return nil, fmt.Errorf("сервер в данный момент занят")
	}
}

// AddFileToTask добавляет один файл к задаче с заданным taskId.
// Метод проверяет расширение файла и контролирует максимальное количество файлов,
// обновляет статус задачи и формирует zip архив с добавленными файлами.
// Ограничение на обработку файлов реализовано через канал FileCountChannel,
// для соблюдения услвоия, что для 1 задачи не более 3 файлов.
//
// Метод является синхронным, но безопасно может вызываться из отдельной горутины,
// чтобы, например, вызывать его в отдельном методе для загрузки сразу нескольких файлов.
// Содержимое case task.FileCountChannel <- struct{}{}: можно обернуть в отдельную горутину,
// чтобы сделать выполнение полностью асинхронным. Однако в этом случае управление и обработка ошибок
// станут менее контролируемыми.
func (service *TaskService) AddFileToTask(ctx context.Context, taskId int, fileURL string, fileName string) error {
	extension := filepath.Ext(fileURL)
	if _, exist := fileExtension[extension]; exist == false {
		return fmt.Errorf("не поддерживаемое расширение файла")
	}

	var task *model.Task

	service.mutex.Lock()
	task, err := service.GetTaskStatusById(ctx, taskId)

	if err != nil {
		service.mutex.Unlock()
		return fmt.Errorf("не удалось найти задачу: %w", err)
	}

	if task.FilesAdded >= 3 {
		service.mutex.Unlock()
		return fmt.Errorf("достигнут максимальный лимит файлов в задаче")
	}

	task.Status = "выполняется"
	service.mutex.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case task.FileCountChannel <- struct{}{}:
		defer func() {
			<-task.FileCountChannel
		}()

		if err := util.DownloadAndAddToZip(task.ArchiveWriter, fileURL, fileName); err != nil {
			return fmt.Errorf("ошибка обработки файла: %v", err)
		}

		service.mutex.Lock()
		task.FilesAdded++
		task.Files = append(task.Files, fileName)

		if task.FilesAdded == 3 {
			task.Status = "завершена"
			if err := task.ArchiveWriter.Close(); err != nil {
				service.mutex.Unlock()
				return fmt.Errorf("ошибка закрытия архива: %v", err)
			}
			if err := task.ArchiveFile.Close(); err != nil {
				service.mutex.Unlock()
				return fmt.Errorf("ошибка закрытия файла: %v", err)
			}
			<-service.tasksSlot
		}

		service.mutex.Unlock()

		return nil
	default:
		return fmt.Errorf("одновременно может обрабатываться только 3 файла")
	}
}
