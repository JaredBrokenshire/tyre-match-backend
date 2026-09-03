package service_tests

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	cv "gocv.io/x/gocv"
	"testing"
	m "tyre-match-backend/db/models"
	"tyre-match-backend/services"
	"tyre-match-backend/tests/helpers"
	"tyre-match-backend/tests/mocks"
)

func TestImageProcessingService_ProcessTyreImpression(t *testing.T) {}

func TestImageProcessingService_ProcessTyreModel(t *testing.T) {}

func TestImageProcessingService_SaveStage(t *testing.T) {
	ts.ClearTable("files")

	fileStoreMock := mocks.NewFileStoreMock()
	ts.SetFileStore(fileStoreMock)

	service := services.NewImageProcessingService(
		ts.S.Repos.TyreImpression,
		ts.S.Repos.TyreModel,
		ts.S.Services.File,
		ts.S.Repos.File,
		ts.S.Dependencies.GetFileStore(),
	)

	setup := func(t *testing.T) {
		fileStoreMock.Reset()
	}

	emptyImage := cv.NewMat()
	grayImage := helpers.SolidGray(10, 10)

	type Expected struct {
		Error         error
		Callback      func()
		DatabaseCheck *helpers.DatabaseCheck
	}

	cases := []struct {
		Name     string
		Setup    func(*testing.T)
		ID       uint
		Model    string
		FileType string
		Image    *cv.Mat
		Expected Expected
	}{
		{
			Name:     "Can't save stage with nil image",
			ID:       1,
			Model:    m.FileModelTyreModel,
			FileType: m.FileTypeOriginal,
			Image:    nil,
			Expected: Expected{
				Error: fmt.Errorf("%w: can not save empty image", services.ProcessingError),
			},
		},
		{
			Name:     "Can't save stage with empty image",
			ID:       1,
			Model:    m.FileModelTyreModel,
			FileType: m.FileTypeOriginal,
			Image:    &emptyImage,
			Expected: Expected{
				Error: fmt.Errorf("%w: can not save empty image", services.ProcessingError),
			},
		},
		{
			Name:     "Can't save stage with unsupported model",
			ID:       1,
			Model:    "test file model",
			FileType: m.FileTypeOriginal,
			Image:    grayImage,
			Expected: Expected{
				Error: fmt.Errorf("%w: unsupported model test file model", services.ProcessingError),
			},
		},
		{
			Name: "Can't save stage if there is an error from the file store",
			Setup: func(t *testing.T) {
				setup(t)
				fileStoreMock.SaveError = errors.New("save error")
			},
			ID:       1,
			Model:    m.FileModelTyreModel,
			FileType: m.FileTypeOriginal,
			Image:    grayImage,
			Expected: Expected{
				Error: errors.New("error saving image processing stage: error saving file to file store: save error"),
			},
		},
		{
			Name:     "Can save tyre impression stage",
			Setup:    setup,
			ID:       1,
			Model:    m.FileModelTyreImpression,
			FileType: m.FileTypeOriginal,
			Image:    grayImage,
			Expected: Expected{
				Error: nil,
				Callback: func() {
					// Ensure file store was called correctly
					assert.Equal(t, 1, len(fileStoreMock.SaveCalls))
					assert.Equal(t, "tyre-impressions/1/original", fileStoreMock.SaveCalls[0].Path)
					assert.Contains(t, fileStoreMock.SaveCalls[0].FileName, ".png")
				},
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "file record was created",
					Model: m.File{
						Model:    m.FileModelTyreImpression,
						ModelId:  1,
						Location: "tyre-impressions/1/original",
					},
					CountExpected: 1,
				},
			},
		},
		{
			Name:     "Can save tyre model stage",
			Setup:    setup,
			ID:       1,
			Model:    m.FileModelTyreModel,
			FileType: m.FileTypeOriginal,
			Image:    grayImage,
			Expected: Expected{
				Error: nil,
				Callback: func() {
					// Ensure file store was called correctly
					assert.Equal(t, 1, len(fileStoreMock.SaveCalls))
					assert.Equal(t, "tyre-models/1/original", fileStoreMock.SaveCalls[0].Path)
					assert.Contains(t, fileStoreMock.SaveCalls[0].FileName, ".png")
				},
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "file record was created",
					Model: m.File{
						Model:    m.FileModelTyreModel,
						ModelId:  1,
						Location: "tyre-models/1/original",
					},
					CountExpected: 1,
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			if test.Setup != nil {
				test.Setup(t)
			}

			err := service.SaveStage(test.ID, test.Model, test.FileType, test.Image)
			if test.Expected.Error != nil {
				assert.Equal(t, test.Expected.Error.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if test.Expected.Callback != nil {
				test.Expected.Callback()
			}

			if test.Expected.DatabaseCheck != nil {
				CheckDatabase(test.Expected.DatabaseCheck)
			}

		})
	}
}
