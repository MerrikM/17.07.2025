package model

// Task - структура задачи
// ID - идентификатор задачи
// Files - массив файлов
// FileCountChannel - буферизированный канал, ограничивающий максимальное количество файлов в одной задаче
// DoneChannel - канал-сигнал завершения (используется для сигнала о том, что архив с файлами готов)
// ArchiveLink - ссылка на созданный архив с файлами
type Task struct {
	ID               int
	Files            []string
	FileCountChannel chan string
	DoneChannel      chan struct{}
	ArchiveLink      string
}
