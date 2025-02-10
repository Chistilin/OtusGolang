package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	//Валидируем перед открытием файла
	fileInfo, err := validate(fromPath, offset, limit)
	if err != nil {
		return err
	}

	if (fileInfo.Size() - offset) < limit {
		limit = fileInfo.Size() - offset
	}
	//Открываем файл
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return ErrUnsupportedFile
	}
	//Проверяем что есть такое смещение
	if _, err = fromFile.Seek(offset, 0); err != nil {
		return ErrOffsetExceedsFileSize
	}
	//Создаем файл куда копировать
	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error %s file create : %v", toPath, err)
	}
	toFileInfo, err := toFile.Stat()
	if err != nil {
		return fmt.Errorf("failed about info To file: %s", toPath)
	}
	if toFileInfo.IsDir() {
		return ErrUnsupportedFile
	}

	//Закрываем всегда при выходе
	defer func(fromFile *os.File, toFile *os.File) {
		errFrom := fromFile.Close()
		if errFrom != nil {
			log.Panicf("error %s file close : %v", fromPath, err)
		}
		errTo := toFile.Close()
		if errTo != nil {
			log.Panicf("error %s file close : %v", toPath, err)
		}
	}(fromFile, toFile)

	//Копируем
	buf := bufio.NewReaderSize(fromFile, int(fileInfo.Size()))
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(buf)
	_, err = io.CopyN(toFile, barReader, limit)
	bar.Finish()
	return nil
}

func validate(fromPath string, offset, limit int64) (os.FileInfo, error) {
	if offset < 0 {
		return nil, ErrOffsetExceedsFileSize
	}
	if limit < 0 {
		return nil, fmt.Errorf("Limit < 0")
	}
	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return nil, fmt.Errorf("failed about info From file: %s", fromPath)
	}
	if fromFileInfo.IsDir() {
		return nil, fmt.Errorf("File is dir: %s", fromPath)
	}
	if fromFileInfo.Size() < offset {
		return nil, fmt.Errorf("File from %s size %s < offset", fromFileInfo.Size(), fromPath)
	}
	return fromFileInfo, nil
}
