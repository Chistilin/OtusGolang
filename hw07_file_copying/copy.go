package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:all
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeOffsetSize    = errors.New("offset negative")
	ErrNegativeLimit         = errors.New("limit negative")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Валидируем перед открытием файла
	fileInfo, err := validate(fromPath, offset, limit)
	if err != nil {
		return err
	}
	if (fileInfo.Size() - offset) < limit {
		limit = fileInfo.Size() - offset
	}
	if limit == 0 && offset == 0 {
		limit = fileInfo.Size()
	}
	// Открываем файл
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer fromFile.Close()
	// Проверяем что есть такое смещение
	if _, err = fromFile.Seek(offset, io.SeekStart); err != nil {
		return ErrOffsetExceedsFileSize
	}
	// Создаем файл куда копировать
	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error %s file create : %w", toPath, err)
	}
	toFileInfo, err := toFile.Stat()
	if err != nil {
		return fmt.Errorf("failed about info To file: %s", toPath)
	}
	// Закрываем всегда при выходе
	defer toFile.Close()
	if toFileInfo.IsDir() {
		return ErrUnsupportedFile
	}
	// Копируем
	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	barReader := bar.NewProxyReader(fromFile)
	_, err = io.CopyN(toFile, barReader, limit)
	if err != nil {
		return err
	}
	return nil
}

func validate(fromPath string, offset, limit int64) (os.FileInfo, error) {
	if offset < 0 {
		return nil, ErrNegativeOffsetSize
	}
	if limit < 0 {
		return nil, ErrNegativeLimit
	}
	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return nil, ErrUnsupportedFile
	}
	if fromFileInfo.IsDir() || fromFileInfo.Size() < offset {
		return nil, ErrUnsupportedFile
	}
	if fromFileInfo.Size() == 0 {
		return nil, ErrUnsupportedFile
	}
	return fromFileInfo, nil
}
