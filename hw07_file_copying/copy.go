package main

import (
	"errors"
	"io"
	"log"
	"os"
)

type fileParams struct {
	offset string
	limit  int64
}

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	//Валидируем оба пути
	validate(fromPath, toPath, offset, limit)
	//Открываем файл
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Panicf("failed %s to read: %v", fromFile, err)
	}
	//Проверяем что есть такое смещение
	if _, err = fromFile.Seek(offset, 0); err != nil {
		log.Panic(ErrOffsetExceedsFileSize)
	}
	//Создаем файл куда копировать
	toFile, err := os.Create(toPath)
	if err != nil {
		log.Panicf("error %s file create : %v", toPath, err)
	}
	toFileInfo, err := os.Stat(toPath)
	if err != nil {
		log.Panicf("failed about info To file: %s", toPath)
	}
	if toFileInfo.IsDir() {
		log.Panic(ErrUnsupportedFile)
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
	_, err = io.CopyN(toFile, fromFile, offset)
	if err != nil {
		log.Panicf("error copy file : %v, %s", fromPath, err)
	}

	return nil
}

func validate(fromPath, toPath string, offset, limit int64) {
	if offset < 0 {
		log.Panic("Offset < 0")
	}
	if limit < 0 {
		log.Panic("Limit < 0")
	}
	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		log.Panicf("failed about info From file: %s", fromPath)
	}
	if fromFileInfo.IsDir() {
		log.Panicf("File is dir: %s", fromPath)
	}
	if fromFileInfo.Size() < offset {
		log.Panicf("File from %s size %s < offset", fromFileInfo.Size(), fromPath)
	}
}
