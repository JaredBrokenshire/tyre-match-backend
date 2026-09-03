package tests

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"tyre-match-backend/api/requests"
	m "tyre-match-backend/db/models"
	"tyre-match-backend/services"
	"tyre-match-backend/tests/factories"
	"tyre-match-backend/tests/helpers"
	"tyre-match-backend/tests/mocks"
)

func TestTyreImpression_List(t *testing.T) {
	ts.ClearTable("tyre_impressions")

	tyreImpression := &m.TyreImpression{}
	factories.NewTyreImpression(ts.S.Db, tyreImpression)
	tyreImpression2 := &m.TyreImpression{}
	factories.NewTyreImpression(ts.S.Db, tyreImpression2)
	tyreImpression3 := &m.TyreImpression{}
	factories.NewTyreImpression(ts.S.Db, tyreImpression3)

	getRequest := func(query string) helpers.Request {
		return helpers.Request{
			Method: http.MethodGet,
			Url:    fmt.Sprintf("/tyre-impressions%v", query),
		}
	}

	cases := []helpers.TestCase{
		{
			Name:    "Can get tyre impressions",
			Request: getRequest(""),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"id":%v`, tyreImpression.ID),
					fmt.Sprintf(`"id":%v`, tyreImpression2.ID),
					fmt.Sprintf(`"id":%v`, tyreImpression2.ID),
					`"total_count":3`,
				},
			},
		},
		{
			Name:    "Can get page 0 tyre impressions",
			Request: getRequest("?page=0&page_size=1"),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"id":%v`, tyreImpression.ID),
					`"total_count":3`,
				},
				BodyPartsMissing: []string{
					fmt.Sprintf(`"id":%v`, tyreImpression2.ID),
					fmt.Sprintf(`"id":%v`, tyreImpression3.ID),
				},
			},
		},
		{
			Name:    "Can get page 1 tyre impressions",
			Request: getRequest("?page=1&page_size=1"),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"id":%v`, tyreImpression2.ID),
					`"total_count":3`,
				},
				BodyPartsMissing: []string{
					fmt.Sprintf(`"id":%v`, tyreImpression.ID),
					fmt.Sprintf(`"id":%v`, tyreImpression3.ID),
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			RunTestCase(t, test)
		})
	}
}

func TestTyreImpression_Get(t *testing.T) {
	ts.ClearTable("tyre_impressions")
	ts.ClearTable("files")

	tyreImpression := &m.TyreImpression{ID: 1}
	factories.NewTyreImpression(ts.S.Db, tyreImpression)
	tyreImpression2 := &m.TyreImpression{ID: 100}
	factories.NewTyreImpression(ts.S.Db, tyreImpression2)

	file1 := &m.File{Model: m.FileModelTyreImpression, ModelId: tyreImpression.ID}
	factories.NewFile(ts.S.Db, file1)
	file2 := &m.File{Model: m.FileModelTyreImpression, ModelId: tyreImpression.ID, FileType: m.FileTypeEnhanced}
	factories.NewFile(ts.S.Db, file2)
	otherImpressionFile := &m.File{Model: m.FileModelTyreImpression, ModelId: tyreImpression2.ID}
	factories.NewFile(ts.S.Db, otherImpressionFile)

	getRequest := func(id interface{}) helpers.Request {
		return helpers.Request{
			Method: http.MethodGet,
			Url:    fmt.Sprintf("/tyre-impressions/%v", id),
		}
	}

	cases := []helpers.TestCase{
		{
			Name:    "Can't get tyre impression that doesn't exist",
			Request: getRequest(1000),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusNotFound,
				BodyPart:   "TyreImpression not found",
			},
		},
		{
			Name:    "Can't get tyre impression with invalid id",
			Request: getRequest("invalid-id"),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "Invalid ID",
			},
		},
		{
			Name:    "Can get tyre impression",
			Request: getRequest(tyreImpression.ID),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"id":%v`, tyreImpression.ID),
					fmt.Sprintf(`"status":"%v"`, tyreImpression.Status),
					fmt.Sprintf(`"name":"%v"`, file1.Name),
					fmt.Sprintf(`"name":"%v"`, file2.Name),
				},
				BodyPartsMissing: []string{
					fmt.Sprintf(`"id":%v`, tyreImpression2.ID),
					fmt.Sprintf(`"name":"%v"`, otherImpressionFile.Name),
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			RunTestCase(t, test)
		})
	}
}

func TestTyreImpression_Create(t *testing.T) {
	ts.ClearTable("tyre_impressions")

	request := helpers.Request{
		Method: http.MethodPost,
		Url:    "/tyre-impressions",
	}

	cases := []helpers.TestCase{
		{
			Name:        "Can't create tyre impression without required fields",
			Request:     request,
			RequestBody: requests.CreateTyreImpressionRequest{},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyParts: []string{
					"Required fields are empty or not valid",
					"PixelsPerInch is a required field",
				},
			},
		},
		{
			Name:    "Can't create tyre impression if fields are below minimum value",
			Request: request,
			RequestBody: requests.CreateTyreImpressionRequest{
				PixelsPerInch: -1,
			},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyParts: []string{
					"Required fields are empty or not valid",
					"PixelsPerInch must be 0 or greater",
				},
			},
		},
		{
			Name:    "Can create tyre impression",
			Request: request,
			RequestBody: requests.CreateTyreImpressionRequest{
				PixelsPerInch: 60,
				ROITop:        100,
				ROILeft:       50,
				ROIRight:      800,
				ROIBottom:     600,
			},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusCreated,
				BodyParts: []string{
					`"pixels_per_inch":60`,
					fmt.Sprintf(`"status":"%v"`, m.ProcessingStatusUploaded),
				},
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "TyreImpression was created",
					Model: m.TyreImpression{
						PixelsPerInch: 60,
						Status:        m.ProcessingStatusUploaded,
					},
					CountExpected: 1,
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			RunTestCase(t, test)
		})
	}
}

func TestTyreImpression_Upload(t *testing.T) {
	ts.ClearTable("tyre_impressions")
	ts.ClearTable("files")

	fileStoreMock := mocks.NewFileStoreMock()
	ts.SetFileStore(fileStoreMock)
	imageProcessingServiceMock := mocks.NewImageProcessingServiceMock()
	validator := services.NewUploadedFileValidator()

	setup := func(test *helpers.TestCase) {
		ts.ClearTable("files")
		fileStoreMock.Reset()
		imageProcessingServiceMock.Reset()
		ts.S.Services.TyreImpression = services.NewTyreImpressionService(
			ts.S.Repos.TyreImpression,
			ts.S.Repos.File,
			fileStoreMock,
			validator,
			imageProcessingServiceMock,
		)
	}

	tyreImpression := &m.TyreImpression{}
	factories.NewTyreImpression(ts.S.Db, tyreImpression)
	tyreImpressionWithExistingImage := &m.TyreImpression{}
	factories.NewTyreImpression(ts.S.Db, tyreImpressionWithExistingImage)

	existingImage := &m.File{Model: m.FileModelTyreImpression, FileType: m.FileTypeOriginal, ModelId: tyreImpressionWithExistingImage.ID}
	factories.NewFile(ts.S.Db, existingImage)

	// PNG files
	pngBody, pngMw := createMultipartFile(t, "file", "../assets/example.png")
	pngBody2, pngMw2 := createMultipartFile(t, "file", "../assets/example.png")
	pngBody3, pngMw3 := createMultipartFile(t, "file", "../assets/example.png")
	pngBody4, pngMw4 := createMultipartFile(t, "file", "../assets/example.png")
	pngBody5, pngMw5 := createMultipartFile(t, "file", "../assets/example.png")
	pngBody6, pngMw6 := createMultipartFile(t, "file", "../assets/example.png")
	// JPG file
	jpgBody, jpgMw := createMultipartFile(t, "file", "../assets/example.jpg")
	// JPEG file
	jpegBody, jpegMw := createMultipartFile(t, "file", "../assets/example.jpeg")
	// WEBP file
	webpBody, webpMw := createMultipartFile(t, "file", "../assets/example.webp")
	// PDF file
	pdfBody, pdfMw := createMultipartFile(t, "file", "../assets/example.pdf")

	getRequest := func(id interface{}) helpers.Request {
		return helpers.Request{
			Method: http.MethodPost,
			Url:    fmt.Sprintf("/tyre-impressions/%v", id),
		}
	}

	cases := []helpers.TestCase{
		{
			Name:    "Can't upload image for tyre impression if no file is provided",
			Request: getRequest(tyreImpression.ID),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "Error parsing file from request",
			},
		},
		{
			Name:               "Can't upload image for tyre impression with invalid file type",
			Request:            getRequest(tyreImpression.ID),
			RequestReader:      pdfBody,
			RequestContentType: pdfMw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "Invalid image upload for TyreImpression",
			},
		},
		{
			Name:               "Can't upload image for tyre impression that doesn't exist",
			Request:            getRequest(1000),
			RequestReader:      pngBody3,
			RequestContentType: pngMw3.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusNotFound,
				BodyPart:   "TyreImpression not found",
			},
		},
		{
			Name:               "Can't upload image for tyre impression with invalid id",
			Request:            getRequest("invalid-id"),
			RequestReader:      pngBody4,
			RequestContentType: pngMw4.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "Invalid ID",
			},
		},
		{
			Name:               "Can't upload image for tyre impression that already has images",
			Request:            getRequest(tyreImpressionWithExistingImage.ID),
			RequestReader:      pngBody5,
			RequestContentType: pngMw5.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "TyreImpression image already exists",
			},
		},
		{
			Name: "Can't upload image for tyre impression if there is an error from the file store",
			Setup: func(test *helpers.TestCase) {
				setup(test)

				fileStoreMock.SaveError = errors.New("test error")
			},
			Request:            getRequest(tyreImpression.ID),
			RequestReader:      pngBody6,
			RequestContentType: pngMw6.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusInternalServerError,
				BodyPart:   "Error saving image to file store",
			},
		},
		{
			Name: "Can't upload image when image processing fails",
			Setup: func(test *helpers.TestCase) {
				setup(test)
				imageProcessingServiceMock.ProcessTyreImpressionError = services.ProcessingError
			},
			Request:            getRequest(tyreImpression.ID),
			RequestReader:      pngBody2,
			RequestContentType: pngMw2.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode:    http.StatusInternalServerError,
				BodyPart:      "Error processing TyreImpression image",
				DatabaseCheck: &helpers.DatabaseCheck{Name: "Tyre impression is marked failed", Model: m.TyreImpression{ID: tyreImpression.ID, Status: m.ProcessingStatusFailed}, CountExpected: 1},
				ExpectedCallBack: func(res *httptest.ResponseRecorder) {
					assert.Equal(t, 1, len(imageProcessingServiceMock.ProcessTyreImpressionCalls))
					assert.Equal(t, tyreImpression.ID, imageProcessingServiceMock.ProcessTyreImpressionCalls[0].ID)
					assert.Equal(t, m.ProcessingStatusFailed, imageProcessingServiceMock.ProcessTyreImpressionCalls[0].Status)
				},
			},
		},
		{
			Name:               "Can upload png",
			Setup:              setup,
			Request:            getRequest(tyreImpression.ID),
			RequestReader:      pngBody,
			RequestContentType: pngMw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"id":%v`, tyreImpression.ID),
					fmt.Sprintf(`"model_id":%v`, tyreImpression.ID),
					fmt.Sprintf(`"file_type":"%v"`, m.FileTypeOriginal),
				},
				ExpectedCallBack: func(res *httptest.ResponseRecorder) {
					// Ensure file store was called correctly
					assert.Equal(t, 1, len(fileStoreMock.SaveCalls))
					assert.Contains(t, fileStoreMock.SaveCalls[0].Path, fmt.Sprintf("tyre-impressions/%v/%v", tyreImpression.ID, m.FileTypeOriginal))
					assert.Contains(t, fileStoreMock.SaveCalls[0].FileName, ".png")

					// Ensure image processing service was called correctly
					assert.Equal(t, 1, len(imageProcessingServiceMock.ProcessTyreImpressionCalls))
					assert.Equal(t, tyreImpression.ID, imageProcessingServiceMock.ProcessTyreImpressionCalls[0].ID)
				},
			},
		},
		{
			Name:               "Can upload jpg",
			Setup:              setup,
			Request:            getRequest(tyreImpression.ID),
			RequestReader:      jpgBody,
			RequestContentType: jpgMw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"id":%v`, tyreImpression.ID),
					fmt.Sprintf(`"model_id":%v`, tyreImpression.ID),
					fmt.Sprintf(`"file_type":"%v"`, m.FileTypeOriginal),
				},
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "File was uploaded",
					Model: m.File{
						Model:    m.FileModelTyreImpression,
						ModelId:  tyreImpression.ID,
						FileType: m.FileTypeOriginal,
					},
					CountExpected: 1,
				},
				ExpectedCallBack: func(res *httptest.ResponseRecorder) {
					// Ensure file store was called correctly
					assert.Equal(t, 1, len(fileStoreMock.SaveCalls))
					assert.Contains(t, fileStoreMock.SaveCalls[0].Path, fmt.Sprintf("tyre-impressions/%v/%v", tyreImpression.ID, m.FileTypeOriginal))
					assert.Contains(t, fileStoreMock.SaveCalls[0].FileName, ".jpg")

					// Ensure image processing service was called correctly
					assert.Equal(t, 1, len(imageProcessingServiceMock.ProcessTyreImpressionCalls))
					assert.Equal(t, tyreImpression.ID, imageProcessingServiceMock.ProcessTyreImpressionCalls[0].ID)
				},
			},
		},
		{
			Name:               "Can upload jpeg",
			Setup:              setup,
			Request:            getRequest(tyreImpression.ID),
			RequestReader:      jpegBody,
			RequestContentType: jpegMw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"id":%v`, tyreImpression.ID),
					fmt.Sprintf(`"model_id":%v`, tyreImpression.ID),
					fmt.Sprintf(`"file_type":"%v"`, m.FileTypeOriginal),
				},
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "File was uploaded",
					Model: m.File{
						Model:    m.FileModelTyreImpression,
						ModelId:  tyreImpression.ID,
						FileType: m.FileTypeOriginal,
					},
					CountExpected: 1,
				},
				ExpectedCallBack: func(res *httptest.ResponseRecorder) {
					// Ensure file store was called correctly
					assert.Equal(t, 1, len(fileStoreMock.SaveCalls))
					assert.Contains(t, fileStoreMock.SaveCalls[0].Path, fmt.Sprintf("tyre-impressions/%v/%v", tyreImpression.ID, m.FileTypeOriginal))
					assert.Contains(t, fileStoreMock.SaveCalls[0].FileName, ".jpeg")

					// Ensure image processing service was called correctly
					assert.Equal(t, 1, len(imageProcessingServiceMock.ProcessTyreImpressionCalls))
					assert.Equal(t, tyreImpression.ID, imageProcessingServiceMock.ProcessTyreImpressionCalls[0].ID)
				},
			},
		},
		{
			Name:               "Can upload webp",
			Setup:              setup,
			Request:            getRequest(tyreImpression.ID),
			RequestReader:      webpBody,
			RequestContentType: webpMw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"id":%v`, tyreImpression.ID),
					fmt.Sprintf(`"model_id":%v`, tyreImpression.ID),
					fmt.Sprintf(`"file_type":"%v"`, m.FileTypeOriginal),
				},
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "File was uploaded",
					Model: m.File{
						Model:    m.FileModelTyreImpression,
						ModelId:  tyreImpression.ID,
						FileType: m.FileTypeOriginal,
					},
					CountExpected: 1,
				},
				ExpectedCallBack: func(res *httptest.ResponseRecorder) {
					// Ensure file store was called correctly
					assert.Equal(t, 1, len(fileStoreMock.SaveCalls))
					assert.Contains(t, fileStoreMock.SaveCalls[0].Path, fmt.Sprintf("tyre-impressions/%v/%v", tyreImpression.ID, m.FileTypeOriginal))
					assert.Contains(t, fileStoreMock.SaveCalls[0].FileName, ".webp")

					// Ensure image processing service was called correctly
					assert.Equal(t, 1, len(imageProcessingServiceMock.ProcessTyreImpressionCalls))
					assert.Equal(t, tyreImpression.ID, imageProcessingServiceMock.ProcessTyreImpressionCalls[0].ID)
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			RunTestCase(t, test)
		})
	}
}

func TestTyreImpression_Delete(t *testing.T) {
	ts.ClearTable("tyre_impressions")

	tyreImpression := &m.TyreImpression{}
	factories.NewTyreImpression(ts.S.Db, tyreImpression)
	tyreImpression2 := &m.TyreImpression{}
	factories.NewTyreImpression(ts.S.Db, tyreImpression2)

	getRequest := func(id interface{}) helpers.Request {
		return helpers.Request{
			Method: http.MethodDelete,
			Url:    fmt.Sprintf("/tyre-impressions/%v", id),
		}
	}

	cases := []helpers.TestCase{
		{
			Name:    "Can't delete tyre impression that does not exist",
			Request: getRequest(1000),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusNotFound,
				BodyPart:   "TyreImpression not found",
			},
		},
		{
			Name:    "Can't delete tyre impression with invalid id",
			Request: getRequest("invalid-id"),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "Invalid ID",
			},
		},
		{
			Name:    "Can delete tyre impression",
			Request: getRequest(tyreImpression.ID),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyPart:   "TyreImpression deleted successfully",
				DatabaseChecks: []*helpers.DatabaseCheck{
					{
						Name: "Tyre model was deleted",
						Model: m.TyreImpression{
							ID: tyreImpression.ID,
						},
						CountExpected: 0,
					},
					{
						Name: "Other tyre impression was not deleted",
						Model: m.TyreImpression{
							ID: tyreImpression2.ID,
						},
						CountExpected: 1,
					},
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			RunTestCase(t, test)
		})
	}
}
