package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:all
)

var (
	ErrFromFileEqualToFile   = errors.New("from file path equal to file path")
	ErrFromFileIsDir         = errors.New("error to file is dir")
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
	toFileInfo, err := os.Stat(toPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	if os.SameFile(fileInfo, toFileInfo) {
		return ErrFromFileEqualToFile
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
		return err
	}
	defer fromFile.Close()
	// Проверяем что есть такое смещение
	if _, err = fromFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}
	// Создаем файл куда копировать
	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error %s file create : %w", toPath, err)
	}
	// Закрываем всегда при выходе
	defer toFile.Close()
	// Копируем
	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	barReader := bar.NewProxyReader(fromFile)
	_, err = io.CopyN(toFile, barReader, limit)
	if err != nil {
		os.Remove(toPath)
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
		return nil, err
	}
	// Проверяем что это обычный файл
	if !fromFileInfo.Mode().IsRegular() {
		return nil, ErrUnsupportedFile
	}
	if fromFileInfo.IsDir() {
		return nil, ErrFromFileIsDir
	}

	if fromFileInfo.Size() < offset {
		return nil, ErrOffsetExceedsFileSize
	}
	return fromFileInfo, nil
}
