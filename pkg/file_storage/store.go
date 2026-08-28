package file_storage

import (
	"io"
)

type Store interface {
	GetStorageLocation() string
	Get(c FileContext, file string) error
	ReadFile(path string) ([]byte, error)
	Save(path string, file io.Reader, filename string) error
	Delete(path string) error
	DeleteDir(path string) error
}
