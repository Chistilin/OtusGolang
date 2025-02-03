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
			break
		}
		/*		if err != nil {
				log.Panicf("failed to read: %v", err)
			}*/
	}
	toFile, err := os.Create(toPath)
	if err != nil {
		log.Panicf("error %s file create : %v", toPath, err)
	}

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

	_, err = io.CopyN(fromFile, toFile, limit)
	if err != nil {
		log.Panicf("error copy file : %v", fromPath)
	}

	return nil
}

/*func validateFile(file string) {
}*/
