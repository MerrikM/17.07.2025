package model

import (
	"archive/zip"
	"os"
)

// Task - структура задачи
// ID - идентификатор задачи
// Files - массив файлов
// FileCountChannel - буферизированный канал, ограничивающий максимальное количество файлов в одной задаче
// DoneChannel - канал-сигнал завершения (используется для сигнала о том, что архив с файлами готов)
// ArchiveWriter - для записи файлов в архив
// ArchiveFile - для закрытия
// ArchiveLink - ссылка на созданный архив с файлами
type Task struct {
	ID               int
	Files            []string
	FileCountChannel chan struct{}
	DoneChannel      chan struct{}
	ArchiveWriter    *zip.Writer
	ArchiveFile      *os.File
	ArchiveLink      string
	Status           string
	FilesAdded       int
}
