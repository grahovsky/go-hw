package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileIn, err := os.OpenFile(fromPath, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	fileIn.Seek(offset, 0)
	defer fileIn.Close()

	fileOut, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fileOut.Close()

	bs, err := BarSize(fileIn)
	if err != nil {
		return err
	}
	bar := pb.StartNew(bs)

	bufSize := 1 * 1024
	buf := make([]byte, bufSize)

	readed := 0
	needBreak := false

	for {
		bufSizeActual := bufSize
		if limit != 0 && ((readed + bufSize) > int(limit)) {
			bufSizeActual = int(limit) - readed
			needBreak = true
		}

		r, err := fileIn.Read(buf[:bufSizeActual])
		if err != nil && err != io.EOF {
			return err
		}
		readed += r
		if r < bufSizeActual {
			bufSizeActual = r
		}

		if err == io.EOF {
			break
		}

		fileOut.Write(buf[:bufSizeActual])

		bar.Add(bufSizeActual)

		if needBreak {
			break
		}
	}

	bar.Finish()

	return nil
}

func BarSize(file *os.File) (int, error) {
	var barSize int
	if fi, err := file.Stat(); err == nil {
		fileSize := fi.Size()

		if offset > fileSize {
			return 0, ErrOffsetExceedsFileSize
		}
		if fi.Size() == 0 || (limit != 0 && fi.Size()-offset > limit) {
			barSize = int(limit)
		} else {
			barSize = int(fi.Size() - offset)
		}
	} else {
		return 0, err
	}

	return barSize, nil
}
