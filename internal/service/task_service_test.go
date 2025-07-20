package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTask_Success(t *testing.T) {
	service := NewTaskService()

	task, err := service.CreateTask()
	assert.NoError(t, err, "ошибка не должна возникать при создании задачи")
	assert.NotNil(t, task, "задача не должна быть nil")
	assert.Equal(t, 1, task.ID, "первая задача должна иметь ID = 1")
	assert.Equal(t, 0, len(task.Files), "у новой задачи не должно быть файлов")
	assert.NotNil(t, task.FileCountChannel, "канал FileCountChannel должен быть создан")
	assert.NotNil(t, task.DoneChannel, "канал DoneChannel должен быть создан")

	service.mutex.Lock()
	_, ok := service.tasks[task.ID]
	service.mutex.Unlock()
	assert.True(t, ok, "задача должна быть сохранена в сервисе")
}

func TestCreateTask_ExceedsLimit(t *testing.T) {
	taskService := NewTaskService()

	for i := 0; i < 3; i++ {
		_, err := taskService.CreateTask()
		assert.NoError(t, err)
	}

	task, err := taskService.CreateTask()
	assert.Nil(t, task, "если превышен лимит задач, задача должна быть nil")
	assert.Error(t, err, "ожидается ошибка при создании 4-й задачи, по требованию максимум 3")
	assert.Equal(t, "сервер в данный момент занят", err.Error())
}
