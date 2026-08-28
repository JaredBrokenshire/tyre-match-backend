package services

import "errors"

var (
	NotFoundError      = errors.New("resource not found")
	AlreadyExistsError = errors.New("resource already exists")
	InvalidUploadError = errors.New("invalid upload")
	FileStoreError     = errors.New("file store error")
	ProcessingError    = errors.New("image processing error")
	DatabaseError      = errors.New("database error")
)
