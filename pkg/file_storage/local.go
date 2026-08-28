package file_storage

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileContext interface {
	File(file string) error
}

type LocalStorage struct {
	storage string
}

var folderLocations = []string{"../", "../../"}

func NewLocalStorage(localStoragePath string) *LocalStorage {

	// Let's find the main folder
	locationParts := strings.Split(localStoragePath, "/")

	if _, err := os.Stat(locationParts[0]); os.IsNotExist(err) {
		// Loop through possible folder locations checking for the folder
		for _, folderLocation := range folderLocations {
			if _, err := os.Stat(folderLocation + locationParts[0]); !os.IsNotExist(err) {
				// If this folder exists then update the localStoragePath
				localStoragePath = folderLocation + localStoragePath
				log.Printf("Found location: updating location storage path to %v\n", localStoragePath)
				break
			}
		}
	}

	// make sure the storage directory exists
	if _, err := os.Stat(localStoragePath); os.IsNotExist(err) {
		err := os.Mkdir(localStoragePath, 0755)
		if err != nil {
			log.Fatalf("FATAL: Error starting local storage service. Unable to create folder for local handler: %v", err)
		}
	}
	return &LocalStorage{
		storage: localStoragePath,
	}
}

func (s *LocalStorage) GetStorageLocation() string {
	return s.storage
}

func (s *LocalStorage) Get(c FileContext, path string) error {
	return c.File(s.storage + path)
}

func (s *LocalStorage) ReadFile(path string) ([]byte, error) {
	fullPath := filepath.Join(s.storage, path)
	res, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *LocalStorage) Save(path string, file io.Reader, filename string) error {
	uploadPath := filepath.Join(s.storage, path)
	err := os.MkdirAll(uploadPath, 0755)
	if err != nil {
		return err
	}

	// Destination
	dst, err := os.Create(filepath.Join(uploadPath, filename))
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, file); err != nil {
		return err
	}

	return nil
}

func (s *LocalStorage) Delete(path string) error {
	err := os.Remove(filepath.Join(s.storage, path))
	return err
}

func (s *LocalStorage) DeleteDir(path string) error {
	err := os.RemoveAll(filepath.Join(s.storage, path))
	return err
}
