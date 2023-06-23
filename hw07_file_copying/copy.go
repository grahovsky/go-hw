package main

import (
	"errors"
	"io"
	"os"
	"syscall"

	"github.com/cheggaaa/pb/v3"
)

var (
	// ErrUnsupportedFile          = errors.New("unsupported file")
	ErrSameFile                 = errors.New("same file")
	ErrOffsetExceedsFileSize    = errors.New("offset exceeds file size")
	ErrFileIsDir                = errors.New("file is directory")
	ErrNoLimitedDeviceOperation = errors.New("device operation no limited")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileIn, err := os.OpenFile(fromPath, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer fileIn.Close()

	fileIn.Seek(offset, 0)

	fileOut, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fileOut.Close()

	if err := checkRestrictions([]*os.File{fileIn, fileOut}, limit); err != nil {
		return err
	}

	bs, err := BarSize(fileIn, limit, offset)
	if err != nil {
		return err
	}

	bar := pb.StartNew(bs)
	// bar.Set(pb.Bytes, true)
	// bar.SetRefreshRate(time.Microsecond)
	// bar.Set(pb.SIBytesPrefix, true)
	// barWriter := bar.NewProxyWriter(fileOut)

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

func BarSize(file *os.File, limit, offset int64) (int, error) {
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

func checkRestrictions(files []*os.File, limit int64) error {
	isDevice := false
	cachedName := ""

	for _, file := range files {

		if file.Name() == cachedName {
			return ErrSameFile
		}
		cachedName = file.Name()

		stat, err := file.Stat()
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return ErrFileIsDir
		}

		if stat.Sys().(*syscall.Stat_t).Mode&syscall.S_IFBLK != 0 {
			isDevice = true
		}

	}

	if isDevice && limit == 0 {
		return ErrNoLimitedDeviceOperation
	}

	return nil
}
