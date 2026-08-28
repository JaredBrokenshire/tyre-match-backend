package mocks

import (
	"io"
	"tyre-match-backend/pkg/file_storage"
)

type FileStoreMock struct {
	GetCalls []string
	GetError error

	ReadFileCalls    []string
	ReadFileResponse []byte
	ReadFileError    error

	SaveCalls []struct{ Path, FileName string }
	SaveError error

	DeleteCalls []string
	DeleteError error

	DeleteDirCalls []string
	DeleteDirError error
}

func NewFileStoreMock() *FileStoreMock {
	m := FileStoreMock{}
	m.Reset()
	return &m
}

func (m *FileStoreMock) GetStorageLocation() string {
	return "test-storage-location"
}

func (m *FileStoreMock) Get(_ file_storage.FileContext, file string) error {
	m.GetCalls = append(m.GetCalls, file)
	return m.GetError
}

func (m *FileStoreMock) ReadFile(path string) ([]byte, error) {
	m.ReadFileCalls = append(m.ReadFileCalls, path)
	return m.ReadFileResponse, m.ReadFileError
}

func (m *FileStoreMock) Save(path string, _ io.Reader, fileName string) error {
	m.SaveCalls = append(m.SaveCalls, struct{ Path, FileName string }{path, fileName})
	return m.SaveError
}

func (m *FileStoreMock) Delete(path string) error {
	m.DeleteCalls = append(m.DeleteCalls, path)
	return m.DeleteError
}

func (m *FileStoreMock) DeleteDir(path string) error {
	m.DeleteDirCalls = append(m.DeleteDirCalls, path)
	return m.DeleteDirError
}

func (m *FileStoreMock) Reset() {
	m.GetCalls = []string{}
	m.GetError = nil

	m.ReadFileCalls = []string{}
	m.ReadFileResponse = []byte("A file used for testing")
	m.ReadFileError = nil

	m.SaveCalls = []struct{ Path, FileName string }{}
	m.SaveError = nil

	m.DeleteCalls = []string{}
	m.DeleteError = nil

	m.DeleteDirCalls = []string{}
	m.DeleteDirError = nil
}
