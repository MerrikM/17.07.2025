package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CreateZIPArchive(archivePath string, archiveName string) (*os.File, *zip.Writer, error) {
	archive, err := os.Create(filepath.Join(archivePath, archiveName))
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка создания .zip архива: %w", err)
	}
	zipWriter := zip.NewWriter(archive)

	return archive, zipWriter, nil
}

func AddFileToZIPArchive(zipWriter *zip.Writer, filePath string) error {
	fileToZip, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return fmt.Errorf("не удалось получить информацию о файле: %w", err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("не удалось создать заголовк для файла: %w", err)
	}

	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("не удалось создать файл внутри архива: %w", err)
	}

	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		return fmt.Errorf("ошибка при записи файла в архив: %w", err)
	}

	return nil
}
