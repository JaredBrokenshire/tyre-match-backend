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

func TestTyreModel_List(t *testing.T) {
	ts.ClearTable("tyre_models")

	tyreModel := &m.TyreModel{ModelName: "000"}
	factories.NewTyreModel(ts.S.Db, tyreModel)
	tyreModel2 := &m.TyreModel{ModelName: "001"}
	factories.NewTyreModel(ts.S.Db, tyreModel2)
	manufacturerTyreModel := &m.TyreModel{Manufacturer: "Test Manufacturer"}
	factories.NewTyreModel(ts.S.Db, manufacturerTyreModel)
	modelNameTyreModel := &m.TyreModel{ModelName: "Test Model Name"}
	factories.NewTyreModel(ts.S.Db, modelNameTyreModel)

	getRequest := func(query string) helpers.Request {
		return helpers.Request{
			Method: http.MethodGet,
			Url:    fmt.Sprintf("/tyre-models%v", query),
		}
	}

	cases := []helpers.TestCase{
		{
			Name:    "Can get tyre models",
			Request: getRequest(""),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"manufacturer":"%v"`, tyreModel.Manufacturer),
					fmt.Sprintf(`"manufacturer":"%v"`, tyreModel2.Manufacturer),
					fmt.Sprintf(`"manufacturer":"%v"`, manufacturerTyreModel.Manufacturer),
					fmt.Sprintf(`"manufacturer":"%v"`, modelNameTyreModel.Manufacturer),
					`"total_count":4`,
				},
			},
		},
		{
			Name:    "Can get page 0 tyre models",
			Request: getRequest("?page=0&page_size=1"),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"manufacturer":"%v"`, tyreModel.Manufacturer),
					`"total_count":4`,
				},
				BodyPartsMissing: []string{
					fmt.Sprintf(`"manufacturer":"%v"`, tyreModel2.Manufacturer),
					fmt.Sprintf(`"manufacturer":"%v"`, manufacturerTyreModel.Manufacturer),
					fmt.Sprintf(`"manufacturer":"%v"`, modelNameTyreModel.Manufacturer),
				},
			},
		},
		{
			Name:    "Can get page 1 tyre models",
			Request: getRequest("?page=1&page_size=1"),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"manufacturer":"%v"`, tyreModel2.Manufacturer),
					`"total_count":4`,
				},
				BodyPartsMissing: []string{
					fmt.Sprintf(`"manufacturer":"%v"`, tyreModel.Manufacturer),
					fmt.Sprintf(`"manufacturer":"%v"`, manufacturerTyreModel.Manufacturer),
					fmt.Sprintf(`"manufacturer":"%v"`, modelNameTyreModel.Manufacturer),
				},
			},
		},
		{
			Name:    "Can filter by manufacturer",
			Request: getRequest("?search=Manufacturer"),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"manufacturer":"%v"`, manufacturerTyreModel.Manufacturer),
					`"total_count":1`,
				},
				BodyPartsMissing: []string{
					fmt.Sprintf(`"manufacturer":"%v"`, tyreModel.Manufacturer),
					fmt.Sprintf(`"manufacturer":"%v"`, tyreModel2.Manufacturer),
					fmt.Sprintf(`"manufacturer":"%v"`, modelNameTyreModel.Manufacturer),
				},
			},
		},
		{
			Name:    "Can filter by model name",
			Request: getRequest("?search=Model Name"),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"manufacturer":"%v"`, modelNameTyreModel.Manufacturer),
					`"total_count":1`,
				},
				BodyPartsMissing: []string{
					fmt.Sprintf(`"manufacturer":"%v"`, tyreModel.Manufacturer),
					fmt.Sprintf(`"manufacturer":"%v"`, tyreModel2.Manufacturer),
					fmt.Sprintf(`"manufacturer":"%v"`, manufacturerTyreModel.Manufacturer),
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

func TestTyreModel_Get(t *testing.T) {
	ts.ClearTable("tyre_models")
	ts.ClearTable("files")

	tyreModel := &m.TyreModel{}
	factories.NewTyreModel(ts.S.Db, tyreModel)
	tyreModel2 := &m.TyreModel{}
	factories.NewTyreModel(ts.S.Db, tyreModel2)

	file1 := &m.File{Model: m.FileModelTyreModel, ModelId: tyreModel.ID}
	factories.NewFile(ts.S.Db, file1)
	otherModelFile := &m.File{Model: m.FileModelTyreModel, ModelId: tyreModel2.ID}
	factories.NewFile(ts.S.Db, otherModelFile)

	getRequest := func(id interface{}) helpers.Request {
		return helpers.Request{
			Method: http.MethodGet,
			Url:    fmt.Sprintf("/tyre-models/%v", id),
		}
	}

	cases := []helpers.TestCase{
		{
			Name:    "Can't get tyre model that doesn't exist",
			Request: getRequest(1000),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusNotFound,
				BodyPart:   "TyreModel not found",
			},
		},
		{
			Name:    "Can't get tyre model with invalid id",
			Request: getRequest("invalid-id"),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "Invalid ID",
			},
		},
		{
			Name:    "Can get tyre model",
			Request: getRequest(tyreModel.ID),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"manufacturer":"%v"`, tyreModel.Manufacturer),
					fmt.Sprintf(`"model_name":"%v"`, tyreModel.ModelName),
					fmt.Sprintf(`"width_mm":%v`, tyreModel.WidthMm),
					fmt.Sprintf(`"aspect_ratio":%v`, tyreModel.AspectRatio),
					fmt.Sprintf(`"rim_diameter_inches":%v`, tyreModel.RimDiameterInches),
					fmt.Sprintf(`"groove_count":%v`, tyreModel.GrooveCount),
					fmt.Sprintf(`"name":"%v"`, file1.Name),
				},
				BodyPartsMissing: []string{
					fmt.Sprintf(`"manufacturer":"%v"`, tyreModel2.Manufacturer),
					fmt.Sprintf(`"model_name":"%v"`, tyreModel2.ModelName),
					fmt.Sprintf(`"name":"%v"`, otherModelFile.Name),
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

func TestTyreModel_Create(t *testing.T) {
	ts.ClearTable("tyre_models")

	request := helpers.Request{
		Method: http.MethodPost,
		Url:    "/tyre-models",
	}

	cases := []helpers.TestCase{
		{
			Name:        "Can't create tyre model without required fields",
			Request:     request,
			RequestBody: requests.CreateTyreModelRequest{},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyParts: []string{
					"Required fields are empty or not valid",
					"Manufacturer is a required field",
					"ModelName is a required field",
				},
			},
		},
		{
			Name:    "Can't create tyre model if fields exceed max length",
			Request: request,
			RequestBody: requests.CreateTyreModelRequest{
				Manufacturer: string(make([]byte, 201)),
				ModelName:    string(make([]byte, 201)),
			},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyParts: []string{
					"Required fields are empty or not valid",
					"ModelName must be a maximum of 200 characters",
				},
			},
		},
		{
			Name:    "Can't create tyre model if fields are below minimum value",
			Request: request,
			RequestBody: requests.CreateTyreModelRequest{
				Manufacturer:      "Test Manufacturer",
				ModelName:         "Test Model Name",
				WidthMm:           -1,
				RimDiameterInches: -1,
				GrooveCount:       -1,
				PixelsPerInch:     -1,
			},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyParts: []string{
					"Required fields are empty or not valid",
					"WidthMm must be 0 or greater",
					"RimDiameterInches must be 0 or greater",
					"GrooveCount must be 0 or greater",
					"PixelsPerInch must be 0 or greater",
				},
			},
		},
		{
			Name:    "Can create tyre model",
			Request: request,
			RequestBody: requests.CreateTyreModelRequest{
				Manufacturer:      "Test Manufacturer",
				ModelName:         "Test ModelName",
				WidthMm:           205,
				RimDiameterInches: 17,
				GrooveCount:       4,
				PixelsPerInch:     200.1,
				ROITop:            200,
				ROIRight:          200,
				ROIBottom:         200,
				ROILeft:           200,
			},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusCreated,
				BodyParts: []string{
					`"manufacturer":"Test Manufacturer"`,
					`"model_name":"Test ModelName"`,
					`"width_mm":205`,
					`"rim_diameter_inches":17`,
					`"groove_count":4`,
					`"pixels_per_inch":200.1`,
					`"roi_top":200`,
					`"roi_right":200`,
					`"roi_bottom":200`,
					`"roi_left":200`,
				},
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "TyreModel was created",
					Model: m.TyreModel{
						Manufacturer:      "Test Manufacturer",
						ModelName:         "Test ModelName",
						WidthMm:           205,
						RimDiameterInches: 17,
						GrooveCount:       4,
						PixelsPerInch:     200.1,
						ROITop:            200,
						ROIRight:          200,
						ROIBottom:         200,
						ROILeft:           200,
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

func TestTyreModel_Update(t *testing.T) {
	ts.ClearTable("tyre_models")

	tyreModel := &m.TyreModel{Manufacturer: "Test", ModelName: "Test"}
	factories.NewTyreModel(ts.S.Db, tyreModel)

	getRequest := func(id interface{}) helpers.Request {
		return helpers.Request{
			Method: http.MethodPut,
			Url:    fmt.Sprintf("/tyre-models/%v", id),
		}
	}

	cases := []helpers.TestCase{
		{
			Name:        "Can't update tyre model without required fields",
			Request:     getRequest(tyreModel.ID),
			RequestBody: requests.CreateTyreModelRequest{},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyParts: []string{
					"Required fields are empty or not valid",
					"Manufacturer is a required field",
					"ModelName is a required field",
				},
			},
		},
		{
			Name:    "Can't update tyre model if fields exceed max length",
			Request: getRequest(tyreModel.ID),
			RequestBody: requests.CreateTyreModelRequest{
				Manufacturer: string(make([]byte, 201)),
				ModelName:    string(make([]byte, 201)),
			},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyParts: []string{
					"Required fields are empty or not valid",
					"ModelName must be a maximum of 200 characters",
				},
			},
		},
		{
			Name:    "Can't update tyre model if fields are below minimum value",
			Request: getRequest(tyreModel.ID),
			RequestBody: requests.CreateTyreModelRequest{
				Manufacturer:      "Test Manufacturer",
				ModelName:         "Test Model Name",
				WidthMm:           -1,
				RimDiameterInches: -1,
				GrooveCount:       -1,
			},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyParts: []string{
					"Required fields are empty or not valid",
					"WidthMm must be 0 or greater",
					"RimDiameterInches must be 0 or greater",
					"GrooveCount must be 0 or greater",
				},
			},
		},
		{
			Name:    "Can't update tyre model that does not exist",
			Request: getRequest(1000),
			RequestBody: requests.CreateTyreModelRequest{
				Manufacturer: "New Manufacturer",
				ModelName:    "New Model Name",
			},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusNotFound,
				BodyPart:   "TyreModel not found",
			},
		},
		{
			Name:    "Can't update tyre model with invalid id",
			Request: getRequest("invalid-id"),
			RequestBody: requests.CreateTyreModelRequest{
				Manufacturer: "New Manufacturer",
				ModelName:    "New Model Name",
			},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "Invalid ID",
			},
		},
		{
			Name:    "Can update tyre model",
			Request: getRequest(tyreModel.ID),
			RequestBody: requests.CreateTyreModelRequest{
				Manufacturer:      "New Manufacturer",
				ModelName:         "New ModelName",
				WidthMm:           205,
				RimDiameterInches: 17,
				GrooveCount:       4,
			},
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					`"manufacturer":"New Manufacturer"`,
					`"model_name":"New ModelName"`,
					`"width_mm":205`,
					`"rim_diameter_inches":17`,
					`"groove_count":4`,
				},
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "TyreModel was updated",
					Model: m.TyreModel{
						Manufacturer:      "New Manufacturer",
						ModelName:         "New ModelName",
						WidthMm:           205,
						RimDiameterInches: 17,
						GrooveCount:       4,
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

func TestTyreModel_Upload(t *testing.T) {
	ts.ClearTable("files")
	ts.ClearTable("tyre_models")

	fileStoreMock := mocks.NewFileStoreMock()
	ts.SetFileStore(fileStoreMock)
	imageProcessingServiceMock := mocks.NewImageProcessingServiceMock()
	validator := services.NewUploadedFileValidator()

	setup := func(test *helpers.TestCase) {
		ts.ClearTable("files")
		fileStoreMock.Reset()
		imageProcessingServiceMock.Reset()
		ts.S.Services.TyreModel = services.NewTyreModelService(
			ts.S.Repos.TyreModel,
			ts.S.Repos.File,
			fileStoreMock,
			validator,
			imageProcessingServiceMock,
		)
	}

	tyreModel := &m.TyreModel{}
	factories.NewTyreModel(ts.S.Db, tyreModel)
	tyreModelWithExistingImage := &m.TyreModel{}
	factories.NewTyreModel(ts.S.Db, tyreModelWithExistingImage)

	// PNG files
	pngBody, pngMw := createMultipartFile(t, "file", "../assets/example.png")
	png2Body, png2Mw := createMultipartFile(t, "file", "../assets/example.png")
	png3Body, png3Mw := createMultipartFile(t, "file", "../assets/example.png")
	png4Body, png4Mw := createMultipartFile(t, "file", "../assets/example.png")
	png5Body, png5Mw := createMultipartFile(t, "file", "../assets/example.png")
	png6Body, png6Mw := createMultipartFile(t, "file", "../assets/example.png")
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
			Url:    fmt.Sprintf("/tyre-models/%v", id),
		}
	}

	cases := []helpers.TestCase{
		{
			Name:    "Can't upload image for tyre model if no file is provided",
			Request: getRequest(tyreModel.ID),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "Error parsing file from request",
			},
		},
		{
			Name:               "Can't upload image for tyre model that doesn't exist",
			Request:            getRequest(1000),
			RequestReader:      png3Body,
			RequestContentType: png3Mw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusNotFound,
				BodyPart:   "TyreModel not found",
			},
		},
		{
			Name:               "Can't upload image for tyre model with invalid id",
			Request:            getRequest("invalid-id"),
			RequestReader:      png4Body,
			RequestContentType: png4Mw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "Invalid ID",
			},
		},
		{
			Name:               "Can't upload image for tyre model with invalid file type",
			Request:            getRequest(tyreModel.ID),
			RequestReader:      pdfBody,
			RequestContentType: pdfMw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "Invalid image upload for TyreModel",
			},
		},
		{
			Name: "Can't upload image for tyre model with existing image",
			Setup: func(test *helpers.TestCase) {
				setup(test)

				existingImage := &m.File{ModelId: tyreModelWithExistingImage.ID, Model: m.FileModelTyreModel, Name: "existing-logo.png"}
				factories.NewFile(ts.S.Db, existingImage)
			},
			Request:            getRequest(tyreModelWithExistingImage.ID),
			RequestReader:      png2Body,
			RequestContentType: png2Mw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "TyreModel image already exists",
			},
		},
		{
			Name: "Can't upload image if there is an error from the file store",
			Setup: func(test *helpers.TestCase) {
				setup(test)

				fileStoreMock.SaveError = errors.New("test-error")
			},
			Request:            getRequest(tyreModelWithExistingImage.ID),
			RequestReader:      png5Body,
			RequestContentType: png5Mw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusInternalServerError,
				BodyPart:   "Error saving image to file store",
			},
		},
		{
			Name: "Can't upload image when image processing fails",
			Setup: func(test *helpers.TestCase) {
				setup(test)
				imageProcessingServiceMock.ProcessTyreModelError = services.ProcessingError
			},
			Request:            getRequest(tyreModel.ID),
			RequestReader:      png6Body,
			RequestContentType: png6Mw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusInternalServerError,
				BodyPart:   "Error processing TyreModel image",
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "Tyre model is marked failed",
					Model: m.TyreModel{
						ID:     tyreModel.ID,
						Status: m.ProcessingStatusFailed,
					},
					CountExpected: 1,
				},
				ExpectedCallBack: func(res *httptest.ResponseRecorder) {
					// Ensure image processing service was called correctly
					assert.Equal(t, 1, len(imageProcessingServiceMock.ProcessTyreModelCalls))
					assert.Equal(t, tyreModel.ID, imageProcessingServiceMock.ProcessTyreModelCalls[0].ID)
					assert.Equal(t, m.ProcessingStatusFailed, imageProcessingServiceMock.ProcessTyreModelCalls[0].Status)
				},
			},
		},
		{
			Name:               "Can upload png",
			Setup:              setup,
			Request:            getRequest(tyreModel.ID),
			RequestReader:      pngBody,
			RequestContentType: pngMw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"id":%v`, tyreModel.ID),
					fmt.Sprintf(`"model_id":%v`, tyreModel.ID),
					fmt.Sprintf(`"file_type":"%v"`, m.FileTypeOriginal),
				},
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "File was uploaded",
					Model: m.File{
						Model:    m.FileModelTyreModel,
						ModelId:  tyreModel.ID,
						FileType: m.FileTypeOriginal,
					},
					CountExpected: 1,
				},
				ExpectedCallBack: func(res *httptest.ResponseRecorder) {
					// Ensure file store was called correctly
					assert.Equal(t, 1, len(fileStoreMock.SaveCalls))
					assert.Contains(t, fileStoreMock.SaveCalls[0].Path, fmt.Sprintf("tyre-models/%v/%v", tyreModel.ID, m.FileTypeOriginal))
					assert.Contains(t, fileStoreMock.SaveCalls[0].FileName, ".png")

					// Ensure image processing service was called correctly
					assert.Equal(t, 1, len(imageProcessingServiceMock.ProcessTyreModelCalls))
					assert.Equal(t, tyreModel.ID, imageProcessingServiceMock.ProcessTyreModelCalls[0].ID)
				},
			},
		},
		{
			Name:               "Can upload jpg",
			Setup:              setup,
			Request:            getRequest(tyreModel.ID),
			RequestReader:      jpgBody,
			RequestContentType: jpgMw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"id":%v`, tyreModel.ID),
					fmt.Sprintf(`"model_id":%v`, tyreModel.ID),
					fmt.Sprintf(`"file_type":"%v"`, m.FileTypeOriginal),
				},
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "File was uploaded",
					Model: m.File{
						Model:    m.FileModelTyreModel,
						ModelId:  tyreModel.ID,
						FileType: m.FileTypeOriginal,
					},
					CountExpected: 1,
				},
				ExpectedCallBack: func(res *httptest.ResponseRecorder) {
					// Ensure file store was called correctly
					assert.Equal(t, 1, len(fileStoreMock.SaveCalls))
					assert.Contains(t, fileStoreMock.SaveCalls[0].Path, fmt.Sprintf("tyre-models/%v/%v", tyreModel.ID, m.FileTypeOriginal))
					assert.Contains(t, fileStoreMock.SaveCalls[0].FileName, ".jpg")

					// Ensure image processing service was called correctly
					assert.Equal(t, 1, len(imageProcessingServiceMock.ProcessTyreModelCalls))
					assert.Equal(t, tyreModel.ID, imageProcessingServiceMock.ProcessTyreModelCalls[0].ID)
				},
			},
		},
		{
			Name:               "Can upload jpeg",
			Setup:              setup,
			Request:            getRequest(tyreModel.ID),
			RequestReader:      jpegBody,
			RequestContentType: jpegMw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"id":%v`, tyreModel.ID),
					fmt.Sprintf(`"model_id":%v`, tyreModel.ID),
					fmt.Sprintf(`"file_type":"%v"`, m.FileTypeOriginal),
				},
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "File was uploaded",
					Model: m.File{
						Model:    m.FileModelTyreModel,
						ModelId:  tyreModel.ID,
						FileType: m.FileTypeOriginal,
					},
					CountExpected: 1,
				},
				ExpectedCallBack: func(res *httptest.ResponseRecorder) {
					// Ensure file store was called correctly
					assert.Equal(t, 1, len(fileStoreMock.SaveCalls))
					assert.Contains(t, fileStoreMock.SaveCalls[0].Path, fmt.Sprintf("tyre-models/%v/%v", tyreModel.ID, m.FileTypeOriginal))
					assert.Contains(t, fileStoreMock.SaveCalls[0].FileName, ".jpeg")

					// Ensure image processing service was called correctly
					assert.Equal(t, 1, len(imageProcessingServiceMock.ProcessTyreModelCalls))
					assert.Equal(t, tyreModel.ID, imageProcessingServiceMock.ProcessTyreModelCalls[0].ID)
				},
			},
		},
		{
			Name:               "Can upload webp",
			Setup:              setup,
			Request:            getRequest(tyreModel.ID),
			RequestReader:      webpBody,
			RequestContentType: webpMw.FormDataContentType(),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyParts: []string{
					fmt.Sprintf(`"id":%v`, tyreModel.ID),
					fmt.Sprintf(`"model_id":%v`, tyreModel.ID),
					fmt.Sprintf(`"file_type":"%v"`, m.FileTypeOriginal),
				},
				DatabaseCheck: &helpers.DatabaseCheck{
					Name: "File was uploaded",
					Model: m.File{
						Model:    m.FileModelTyreModel,
						ModelId:  tyreModel.ID,
						FileType: m.FileTypeOriginal,
					},
					CountExpected: 1,
				},
				ExpectedCallBack: func(res *httptest.ResponseRecorder) {
					// Ensure file store was called correctly
					assert.Equal(t, 1, len(fileStoreMock.SaveCalls))
					assert.Contains(t, fileStoreMock.SaveCalls[0].Path, fmt.Sprintf("tyre-models/%v/%v", tyreModel.ID, m.FileTypeOriginal))
					assert.Contains(t, fileStoreMock.SaveCalls[0].FileName, ".webp")

					// Ensure image processing service was called correctly
					assert.Equal(t, 1, len(imageProcessingServiceMock.ProcessTyreModelCalls))
					assert.Equal(t, tyreModel.ID, imageProcessingServiceMock.ProcessTyreModelCalls[0].ID)
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

func TestTyreModel_Delete(t *testing.T) {
	ts.ClearTable("tyre_models")

	tyreModel := &m.TyreModel{}
	factories.NewTyreModel(ts.S.Db, tyreModel)
	tyreModel2 := &m.TyreModel{}
	factories.NewTyreModel(ts.S.Db, tyreModel2)

	getRequest := func(id interface{}) helpers.Request {
		return helpers.Request{
			Method: http.MethodDelete,
			Url:    fmt.Sprintf("/tyre-models/%v", id),
		}
	}

	cases := []helpers.TestCase{
		{
			Name:    "Can't delete tyre model that does not exist",
			Request: getRequest(1000),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusNotFound,
				BodyPart:   "TyreModel not found",
			},
		},
		{
			Name:    "Can't delete tyre model with invalid id",
			Request: getRequest("invalid-id"),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusBadRequest,
				BodyPart:   "Invalid ID",
			},
		},
		{
			Name:    "Can delete tyre model",
			Request: getRequest(tyreModel.ID),
			Expected: helpers.ExpectedResponse{
				StatusCode: http.StatusOK,
				BodyPart:   "TyreModel deleted successfully",
				DatabaseChecks: []*helpers.DatabaseCheck{
					{
						Name: "Tyre model was deleted",
						Model: m.TyreModel{
							ID: tyreModel.ID,
						},
						CountExpected: 0,
					},
					{
						Name: "Other tyre model was not deleted",
						Model: m.TyreModel{
							ID: tyreModel2.ID,
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
