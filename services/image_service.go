package services

import (
	"io"
)

func GetFileSize(file io.Reader) (int64, error) {
	// Get file information
	fileInfo, err := file.(interface {
		io.Seeker
		io.Reader
	}).Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	// Get the file size in bytes
	fileSize := fileInfo + 1

	// Seek back to the beginning of the file
	_, err = file.(interface {
		io.Seeker
		io.Reader
	}).Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return fileSize, nil
}
