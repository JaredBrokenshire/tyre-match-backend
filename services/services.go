package services

import (
	"tyre-match-backend/db/repositories"
	"tyre-match-backend/pkg/dependencies"
)

type Services struct {
	TyreModel       *TyreModelService
	TyreImpression  *TyreImpressionService
	ImageProcessing *ImageProcessingService
	File            *FileService
}

func NewServices(repos *repositories.Repos, dependencies *dependencies.DependencyService) *Services {
	file := NewFileService(dependencies.GetFileStore())
	imageProcessingService := NewImageProcessingService(repos.TyreImpression, repos.TyreModel, file, repos.File, dependencies.GetFileStore())

	return &Services{
		TyreImpression:  NewTyreImpressionService(repos.TyreImpression, repos.File, dependencies.GetFileStore(), NewUploadedFileValidator(), imageProcessingService),
		TyreModel:       NewTyreModelService(repos.TyreModel, repos.File, dependencies.GetFileStore(), NewUploadedFileValidator(), imageProcessingService),
		ImageProcessing: imageProcessingService,
		File:            file,
	}
}
