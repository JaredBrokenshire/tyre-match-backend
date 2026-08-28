package requests

type CreateTyreImpressionRequest struct {
	PixelsPerInch float32 `json:"pixels_per_inch" validate:"required,min=0"`
}
