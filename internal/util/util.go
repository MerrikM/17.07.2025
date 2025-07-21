package util

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

// CreateZIPArchive создаёт новый ZIP-архив по указанному пути с заданным именем.
//
// Шаги функции:
// 1. Создаёт файл на диске по пути archivePath/archiveName.
// 2. Создает zip.Writer для записи данных в архив.
//
// Возвращает: *os.File; zip.Writer для добавления файлов; ошибку.
//
// Учтите, что вызов этой функции не закрывает ни файл, ни zip.Writer.
// После завершения работы нужно самостоятельно их закрыть.
func CreateZIPArchive(archivePath string, archiveName string) (*os.File, *zip.Writer, error) {
	archive, err := os.Create(filepath.Join(archivePath, archiveName+".zip"))
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка создания .zip архива: %w", err)
	}
	zipWriter := zip.NewWriter(archive)

	return archive, zipWriter, nil
}

// DownloadAndAddToZip скачивает файл по переданному URL и сохраняет его внутрь zip-архива.
// Имя файла внутри архива формируется как prefix + оригинальное расширение файла.
// Пример: если URL заканчивается на `image.jpeg` и prefix = "file1", то в архиве будет "file1.jpeg".
//
// Шаги функции:
// 1. Извлекается имя и расширение файла из URL (например, ".jpeg").
// 2. Скачивается содержимое файла по HTTP GET.
// 3. Внутри zip-архива создаётся новый файл с нужным именем.
// 4. Содержимое скачанного файла копируется прямо в архив (файл на диск не сохраняется).
//
// Возвращет: ошибку
func DownloadAndAddToZip(zipWriter *zip.Writer, fileURL string, filename string) error {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("ошибка парсинга URL: %w", err)
	}

	originalName := path.Base(parsedURL.Path)
	extension := filepath.Ext(originalName)
	if extension == "" {
		return fmt.Errorf("у файла отсутствует расширение: %s", fileURL)
	}

	filenameInZip := filename + extension

	response, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("ошибка скачивания файла: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("сервер вернул ошибку: %s", response.Status)
	}

	zipFile, err := zipWriter.Create(filenameInZip)
	if err != nil {
		return fmt.Errorf("ошибка создания файла в zip архиве: %w", err)
	}

	_, err = io.Copy(zipFile, response.Body)
	if err != nil {
		return fmt.Errorf("ошибка копирования данных в zip архив: %w", err)
	}

	return nil
}
