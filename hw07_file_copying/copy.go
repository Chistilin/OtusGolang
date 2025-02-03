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
	ErrCopyFile              = errors.New("error copy file")
	ErrFileClose             = errors.New("error file closed")
	ErrFileOpen              = errors.New("error file open")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			log.Panicf("File not found: %s", fromFile)
		}

		log.Panicf("failed %s to read: %v", fromFile, err)
	}
	//buf := make([]byte, offset) // подготавливаем буфер нужного размера
	for offset < 0 {
		//read, err := fromFile.Read(buf[offset:])
		//offset += read
		if err == io.EOF {
			// что если не дочитали ?
			break
		}
		/*		if err != nil {
				log.Panicf("failed to read: %v", err)
			}*/
	}
	toFile, err := os.Create(toPath)
	if err != nil {
		panic(ErrFileOpen)
	}

	defer func(fromFile *os.File, toFile *os.File) {
		errFrom := fromFile.Close()
		if errFrom != nil {
			panic(ErrFileClose)
		}
		errTo := toFile.Close()
		if errTo != nil {
			panic(ErrFileClose)
		}
	}(fromFile, toFile)

	_, err = io.CopyN(fromFile, toFile, limit)
	if err != nil {
		panic(ErrCopyFile)
	}

	return nil
}

/*func validateFile(file string) {
}*/
