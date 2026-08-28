package helpers

import cv "gocv.io/x/gocv"

func Empty() *cv.Mat {
	res := cv.NewMat()
	return &res
}

func CloneGray(src cv.Mat) *cv.Mat {
	dst := src.Clone()
	return &dst
}

func SolidGray(width, height int) *cv.Mat {
	img := cv.NewMatWithSize(height, width, cv.MatTypeCV8UC1)
	img.SetTo(cv.NewScalar(120, 120, 120, 0))

	return &img
}

func GradientGray(width, height int) *cv.Mat {
	img := cv.NewMatWithSize(height, width, cv.MatTypeCV8UC1)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			value := uint8(0)

			if width > 1 {
				value = uint8(x * 255 / (width - 1))
			}

			img.SetUCharAt(y, x, value)
		}
	}

	return &img
}

func CheckerGray(width, height int) *cv.Mat {
	img := cv.NewMatWithSize(height, width, cv.MatTypeCV8UC1)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			value := uint8(0)

			if (x+y)%2 == 1 {
				value = 255
			}

			img.SetUCharAt(y, x, value)
		}
	}

	return &img
}

func ImpulseGray(width, height int) *cv.Mat {
	img := SolidGray(width, height)

	img.SetUCharAt(height/2, width/2, 255)

	return img
}

func StripesGray(width, height int) *cv.Mat {
	img := cv.NewMatWithSize(height, width, cv.MatTypeCV8UC1)

	for y := 0; y < height; y++ {
		value := uint8(0)

		if y%2 == 1 {
			value = 255
		}

		for x := 0; x < width; x++ {
			img.SetUCharAt(y, x, value)
		}
	}

	return &img
}

func OffsetGray(width, height int) *cv.Mat {
	img := cv.NewMatWithSize(height, width, cv.MatTypeCV8UC1)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			value := uint8(0)

			if width > 1 {
				value = uint8(x * 255 / (width - 1))
			}

			if y%2 == 1 {
				value /= 2
			}

			img.SetUCharAt(y, x, value)
		}
	}

	return &img
}

func Noise(variance uint8) *cv.Mat {
	result := SolidGray(100, 100)

	for y := 0; y < result.Rows(); y++ {
		for x := 0; x < result.Cols(); x++ {
			noise := uint8((x*37 + y*53) % int(variance*2+1))
			value := int(result.GetUCharAt(y, x)) + int(noise) - int(variance)

			if value < 0 {
				value = 0
			}
			if value > 255 {
				value = 255
			}

			result.SetUCharAt(y, x, uint8(value))
		}
	}

	return result
}

func LowContrast(width, height int) *cv.Mat {
	mat := cv.NewMatWithSize(height, width, cv.MatTypeCV8UC1)

	regionWidth := width / 4
	regionHeight := height / 4

	values := []uint8{
		90, 94, 98, 102,
		94, 98, 102, 106,
		98, 102, 106, 110,
		102, 106, 110, 114,
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			regionX := x / regionWidth
			regionY := y / regionHeight

			if regionX >= 4 {
				regionX = 3
			}
			if regionY >= 4 {
				regionY = 3
			}

			mat.SetUCharAt(
				y,
				x,
				values[regionY*4+regionX],
			)
		}
	}

	return &mat
}
