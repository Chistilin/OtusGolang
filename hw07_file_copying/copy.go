package main

import (
	"errors"
	"io"
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
		panic(err)
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
