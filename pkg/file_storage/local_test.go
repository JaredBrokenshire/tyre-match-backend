package file_storage

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

type MockFileContext struct {
}

func NewMockFileContext() *MockFileContext {
	return &MockFileContext{}
}

func (c *MockFileContext) File(file string) error {
	_, err := os.Stat(file)
	return err
}

type SaveTestCase struct {
	TestName        string
	FilePath        string
	FileName        string
	DestinationPath string
}

func TestGet(t *testing.T) {
	localStorage := NewLocalStorage("../../assets/")
	fileContext := NewMockFileContext()

	err := localStorage.Get(fileContext, "example.jpg")
	assert.NoError(t, err, "Unexpected error when finding file")

	fmt.Println("Success")
}

func TestGetFileNotExist(t *testing.T) {
	localStorage := NewLocalStorage("../../assets/")
	fileContext := NewMockFileContext()

	err := localStorage.Get(fileContext, "/does-not-exist.jpg")
	assert.Error(t, err, "Expected error to be thrown for file that does not exist")

	fmt.Println("Success")
}

func TestGetFileFolderNotExist(t *testing.T) {
	localStorage := NewLocalStorage("../../assets/")
	fileContext := NewMockFileContext()

	err := localStorage.Get(fileContext, "/folder-dont-exist/does-not-exist.jpg")
	assert.Error(t, err, "Expected error to be thrown for file that does not exist")

	fmt.Println("Success")
}

func TestSave(t *testing.T) {
	localStorage := NewLocalStorage("../../assets/files/")

	cases := []SaveTestCase{
		{
			TestName:        "Can save jpg",
			FilePath:        "../../assets/",
			FileName:        "example.jpg",
			DestinationPath: "test-save/",
		},
		{
			TestName:        "Can save larger jpg",
			FilePath:        "../../assets/",
			FileName:        "example-large.jpg",
			DestinationPath: "test-save/",
		},
		{
			TestName:        "Can save txt",
			FilePath:        "../../assets/",
			FileName:        "example.txt",
			DestinationPath: "test-save/",
		},
		{
			TestName:        "Can save pdf",
			FilePath:        "../../assets/",
			FileName:        "example.pdf",
			DestinationPath: "test-save/",
		},
	}

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			// Load file
			src, err := os.OpenFile(test.FilePath+test.FileName, os.O_RDWR, 0755)
			if err != nil {
				t.Fatalf("Error opening test file: %v", err)
			}

			// Save file
			err = localStorage.Save(test.DestinationPath, src, test.FileName)
			assert.NoError(t, err, "Unexpected error")

			// Check uploaded file exists
			targetFilePath := localStorage.storage + test.DestinationPath + test.FileName
			_, err = os.Stat(targetFilePath)
			assert.NoError(t, err, "Expected file to exist")

			// Check files are the same
			assertFilesEqual(t, test.FilePath+test.FileName, targetFilePath)
		})
	}
}

func TestSaveDirectoryNotExist(t *testing.T) {
	localStorage := NewLocalStorage("../../assets/files/")

	// Get directory name
	dir := fmt.Sprintf("../../assets/files/test-save-dir-%v", time.Now().Unix())

	// Confirm directory does not already exist
	_, err := os.Stat(dir)
	assert.Error(t, err, "Expected an error")
	assert.Contains(t, err.Error(), "The system cannot find the file specified")

	// Save file to directory
	src, err := os.OpenFile("../../assets/example.jpg", os.O_RDWR, 0755)
	if err != nil {
		t.Fatalf("Error opening test file: %v", err)
	}

	err = localStorage.Save(dir, src, "example.jpg")
	assert.NoError(t, err, "Unexpected error")

	// Check directory now exists
	_, err = os.Stat(dir)
	assert.NoError(t, err, "Unexpected error - directory should exist")

	// Delete directory
	err = os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("Error removing created directory: %v", err)
	}
}

func TestSaveOverwrite(t *testing.T) {
	localStorage := NewLocalStorage("../../assets/files/")

	// Create file to overwrite
	createFile(t, localStorage.storage, "overwrite.txt", "version 1")

	// Load file to overwrite with
	src, err := os.OpenFile("../../assets/example.txt", os.O_RDWR, 0755)
	if err != nil {
		t.Fatalf("Error opening test file: %v", err)
	}

	// Save the file
	err = localStorage.Save(localStorage.storage, src, "overwrite.txt")
	assert.NoError(t, err, "Unexpected error")

	// Check files are the same
	assertFilesEqual(t, "../../assets/example.txt", "../../assets/files/overwrite.txt")
}

// Save - overwrite existing - overwrite an existing text file with a new string value and check they are not equal

func TestDelete(t *testing.T) {
	localStorage := NewLocalStorage("../../assets/files/")

	// Create file
	createFile(t, localStorage.storage, "delete-example.txt", "testing")

	// Delete file
	err := localStorage.Delete("delete-example.txt")
	assert.NoError(t, err, "Unexpected error when deleting the file")
	// Check file no longer exists
	_, err = os.Stat(localStorage.storage + "delete-example.txt")
	assert.Error(t, err, "Expected an error")
	assert.Contains(t, err.Error(), "The system cannot find the file specified")
}

func TestDeleteFileNotExist(t *testing.T) {
	localStorage := NewLocalStorage("../../assets/files/")

	err := localStorage.Delete("does-not-exist.txt")
	assert.Error(t, err, "Expected an error")
	assert.Contains(t, err.Error(), "The system cannot find the file specified")
}

func TestDeleteFileFolderNotExist(t *testing.T) {
	localStorage := NewLocalStorage("../../assets/files/")

	err := localStorage.Delete("no-folder/does-not-exist.txt")
	assert.Error(t, err, "Expected an error")
	assert.Contains(t, err.Error(), "The system cannot find the path specified")
}

func TestDeleteDirectory(t *testing.T) {
	localStorage := NewLocalStorage("../../assets/files/")

	// Create file
	createFile(t, localStorage.storage+"test-directory/", "example.txt", "testing")

	// Attempt to delete whole directory
	err := localStorage.Delete("test-directory/")
	assert.Error(t, err, "Expected an error")
	assert.Contains(t, err.Error(), "directory is not empty")

	// Check file still exists
	_, err = os.Stat(localStorage.storage + "test-directory/example.txt")
	assert.NoError(t, err, "Expected file to still exist")
}

func TestReadFile(t *testing.T) {
	localStorage := NewLocalStorage("../../assets/")

	testCases := []struct {
		Name                 string
		Path                 string
		ExpectedOutput       []byte
		ExpectedErrorMessage string
	}{
		{
			Name:                 "Can read file",
			Path:                 "example.txt",
			ExpectedOutput:       []byte("A file used for testing"),
			ExpectedErrorMessage: "",
		},
		{
			Name:                 "Cannot read file that does not exist",
			Path:                 "does-not-exist.txt",
			ExpectedOutput:       nil,
			ExpectedErrorMessage: "The system cannot find the file specified",
		},
		{
			Name:                 "Cannot read directory",
			Path:                 "files",
			ExpectedOutput:       nil,
			ExpectedErrorMessage: "Incorrect function",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			res, err := localStorage.ReadFile(test.Path)
			if test.ExpectedErrorMessage == "" {
				assert.NoError(t, err, "Expected no error")
			} else {
				assert.ErrorContains(t, err, test.ExpectedErrorMessage, "Expected error to contain message")
			}
			assert.Equal(t, test.ExpectedOutput, res, "Expected output to match")
		})
	}
}

func TestLocalStorage_DeleteDir(t *testing.T) {
	localStorage := NewLocalStorage("../../assets/files/")

	testCases := []struct {
		Name          string
		Path          string
		ExpectedError string
		Setup         func(t *testing.T)
		Callback      func(t *testing.T)
	}{
		{
			Name:          "Cannot delete directory that does not exist",
			Path:          "does-not-exist/",
			ExpectedError: "",
			Setup:         nil,
			Callback:      nil,
		},
		{
			Name:          "Can delete a file",
			Path:          "test-dir-file-delete/example.txt",
			ExpectedError: "",
			Setup: func(t *testing.T) {
				t.Helper()
				// Create file
				createFile(t, localStorage.storage+"test-dir-file-delete/", "example.txt", "testing")
			},
			Callback: func(t *testing.T) {
				// Check directory no longer exists
				_, err := os.Stat(localStorage.storage + "test-dir-file-delete/example.txt")
				if !os.IsNotExist(err) {
					assert.Fail(t, "Expected directory to not exist")
				}
			},
		},
		{
			Name:          "Can delete directory",
			Path:          "test-dir-delete/",
			ExpectedError: "",
			Setup: func(t *testing.T) {
				t.Helper()
				// Create file
				createFile(t, localStorage.storage+"test-dir-delete/", "example.txt", "testing")
			},
			Callback: func(t *testing.T) {
				// Check directory no longer exists
				_, err := os.Stat(localStorage.storage + "test-dir-delete/")
				if !os.IsNotExist(err) {
					assert.Fail(t, "Expected directory to not exist")
				}
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			if test.Setup != nil {
				test.Setup(t)
			}
			err := localStorage.DeleteDir(test.Path)
			if test.ExpectedError == "" {
				assert.NoError(t, err, "Expected no error")
			} else {
				assert.ErrorContains(t, err, test.ExpectedError, "Expected error to contain message")
			}
			if test.Callback != nil {
				test.Callback(t)
			}
		})
	}
}

func createFile(t *testing.T, path string, filename string, contents string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		t.Fatalf("Error creating directories: %v", err)
	}

	data := []byte(contents)
	err = os.WriteFile(path+filename, data, 0755)
	if err != nil {
		t.Fatalf("Error creating file: %v", err)
	}
}

func assertFilesEqual(t *testing.T, srcPath string, targetPath string) {
	f1, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatalf("Unable to read source file to compare: %v", err)
	}
	f2, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("Unable to read target file to compare: %v", err)
	}

	assert.True(t, bytes.Equal(f1, f2), "Expected files to be equal")
}
