package cv_tests_test

import (
	"encoding/json"
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cv "gocv.io/x/gocv"
	fe "tyre-match-backend/cv/feature_extraction"
	"tyre-match-backend/tests/helpers"
)

func TestFeatureExtractorExtractWithContextProducesNamedStructuralFeatures(t *testing.T) {
	extractor := fe.NewExtractor()
	mask := syntheticTreadMask(240, 160, false)
	defer mask.Close()

	fingerprint, err := extractor.Extract(&mask, 200)

	require.NoError(t, err)
	require.NotNil(t, fingerprint)
	assert.Equal(t, "named_multiscale_tread_fingerprint_v3", fingerprint.ExtractionMethod)
	assert.Equal(t, "white_foreground_black_background", fingerprint.ForegroundConvention)
	assert.Equal(t, float32(200), fingerprint.PixelsPerInch)
	assert.GreaterOrEqual(t, len(fingerprint.Features), 35)
	assert.LessOrEqual(t, len(fingerprint.Features), 45)
	assert.NotEmpty(t, fingerprint.Quality)

	names := make(map[string]bool)
	for _, feature := range fingerprint.Features {
		assert.NotEmpty(t, feature.Name)
		assert.NotEmpty(t, feature.Group)
		assert.NotEmpty(t, feature.Description)
		assert.NotEmpty(t, feature.Unit)
		assert.Greater(t, feature.Weight, float32(0))
		assert.False(t, names[feature.Name], "duplicate feature %s", feature.Name)
		names[feature.Name] = true
	}

	assert.Contains(t, names, "foreground_area_ratio")
	assert.Contains(t, names, "foreground_width_mm")
	assert.Contains(t, names, "longitudinal_density_bin_01")
	assert.Contains(t, names, "transverse_density_bin_01")
	assert.Contains(t, names, "spatial_density_row_01_col_01")
	assert.Contains(t, names, "longitudinal_repeat_distance_mm")
	assert.Contains(t, names, "longitudinal_periodicity_strength")
	assert.Contains(t, names, "substantial_component_count")
	assert.Contains(t, names, "dominant_boundary_orientation_deg")
}

func TestFeatureExtractorUsesOnlyProvidedImage(t *testing.T) {
	extractor := fe.NewExtractor()

	full := syntheticTreadMask(220, 160, false)
	defer full.Close()

	// The feature extractor receives an already-cropped ROI. Pixels outside
	// this region must therefore have no influence on the extracted features.
	roi := full.Region(imageRect(30, 20, 200, 140))
	roiClone := roi.Clone()
	roi.Close()
	defer roiClone.Close()

	before, err := extractor.Extract(&roiClone, 150)
	require.NoError(t, err)

	// Modify the original image outside the ROI. The already-cropped image
	// supplied to the extractor must remain unchanged.
	outsideTop := full.Region(imageRect(0, 0, 220, 20))
	outsideTop.SetTo(cv.NewScalar(0, 0, 0, 0))
	outsideTop.Close()

	after, err := extractor.Extract(&roiClone, 150)
	require.NoError(t, err)

	assert.Equal(t, featureValues(before), featureValues(after))
}

func TestFeatureExtractorSeparatesStructurallyDifferentPatterns(t *testing.T) {
	extractor := fe.NewExtractor()
	first := syntheticTreadMask(240, 160, false)
	second := syntheticTreadMask(240, 160, true)
	defer first.Close()
	defer second.Close()

	fa, err := extractor.Extract(&first, 200)
	require.NoError(t, err)
	fb, err := extractor.Extract(&second, 200)
	require.NoError(t, err)

	assert.InDelta(t, 1.0, extractor.Similarity(fa, fa), 0.0001)
	assert.Less(t, extractor.Similarity(fa, fb), float32(0.95))
}

func TestFeatureExtractorIgnoresSmallIsolatedNoise(t *testing.T) {
	extractor := fe.NewExtractor()
	clean := syntheticTreadMask(240, 160, false)
	noisy := syntheticTreadMask(240, 160, false)
	defer clean.Close()
	defer noisy.Close()

	for y := 35; y < 125; y += 7 {
		for x := 25; x < 215; x += 11 {
			noisy.SetUCharAt(y, x, 255)
		}
	}

	fc, err := extractor.Extract(&clean, 200)
	require.NoError(t, err)
	fn, err := extractor.Extract(&noisy, 200)
	require.NoError(t, err)

	assert.Greater(t, extractor.Similarity(fc, fn), float32(0.98))
}

func TestFeatureFingerprintEncodingRoundTripsDescriptiveFeatures(t *testing.T) {
	extractor := fe.NewExtractor()
	mask := syntheticTreadMask(200, 120, false)
	defer mask.Close()

	fingerprint, err := extractor.Extract(&mask, 240)
	require.NoError(t, err)

	encoded, err := fe.Encode(fingerprint)
	require.NoError(t, err)
	assert.NotContains(t, encoded, `"values"`)
	assert.Contains(t, encoded, `"description"`)
	assert.Contains(t, encoded, `"longitudinal_repeat_distance_mm"`)

	var document map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(encoded), &document))
	assert.Contains(t, document, "features")
	assert.Contains(t, document, "quality")

	decoded, err := fe.Decode(encoded)
	require.NoError(t, err)
	assert.Equal(t, fingerprint.ExtractionMethod, decoded.ExtractionMethod)
	assert.Equal(t, fingerprint.PixelsPerInch, decoded.PixelsPerInch)
	assert.Equal(t, featureValues(fingerprint), featureValues(decoded))
}

func TestFeatureFingerprintDecodeRejectsMalformedFingerprints(t *testing.T) {
	validFeature := `{"name":"x","group":"global_geometry","description":"d","unit":"ratio","value":0.1,"weight":1}`
	cases := []struct {
		name string
		json string
		err  string
	}{
		{"empty", "", "empty fingerprint"},
		{"invalid JSON", "{", "unmarshal fingerprint"},
		{"no named features", `{"extraction_method":"x","foreground_convention":"white_foreground_black_background","pixels_per_inch":100,"features":[]}`, "no named features"},
		{"missing extraction method", `{"foreground_convention":"white_foreground_black_background","pixels_per_inch":100,"features":[` + validFeature + `]}`, "no extraction method"},
		{"invalid pixels per inch", `{"extraction_method":"x","foreground_convention":"white_foreground_black_background","pixels_per_inch":0,"features":[` + validFeature + `]}`, "invalid pixels-per-inch"},
		{"missing descriptive metadata", `{"extraction_method":"x","foreground_convention":"white_foreground_black_background","pixels_per_inch":100,"features":[{"name":"x","group":"","description":"d","unit":"ratio","value":0.1,"weight":1}]}`, "missing descriptive metadata"},
		{"invalid weight", `{"extraction_method":"x","foreground_convention":"white_foreground_black_background","pixels_per_inch":100,"features":[{"name":"x","group":"g","description":"d","unit":"ratio","value":0.1,"weight":0}]}`, "invalid weight"},
		{"duplicate name", `{"extraction_method":"x","foreground_convention":"white_foreground_black_background","pixels_per_inch":100,"features":[` + validFeature + `,` + validFeature + `]}`, "duplicate feature"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := fe.Decode(test.json)
			require.Error(t, err)
			assert.Contains(t, err.Error(), test.err)
		})
	}
}

func TestFeatureExtractorRejectsInvalidInput(t *testing.T) {
	extractor := fe.NewExtractor()
	empty := cv.NewMat()
	defer empty.Close()
	colour := helpers.Colour()
	defer colour.Close()
	gray := helpers.StripesGray(20, 20)
	defer gray.Close()
	black := cv.NewMatWithSize(20, 20, cv.MatTypeCV8UC1)
	defer black.Close()

	cases := []struct {
		name  string
		image *cv.Mat
		ppi   float32
		err   string
	}{
		{"nil image", nil, 0, "feature extraction received an empty image"},
		{"empty image", &empty, 0, "feature extraction received an empty image"},
		{"colour image", colour, 100, "requires an 8-bit single-channel binary image"},
		{"zero pixels per inch", gray, 0, "requires a positive pixels-per-inch"},
		{"empty foreground", &black, 100, "input contains no foreground pixels"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := extractor.Extract(test.image, test.ppi)
			require.Error(t, err)
			assert.Contains(t, err.Error(), test.err)
		})
	}
}

func syntheticTreadMask(width, height int, alternate bool) cv.Mat {
	mask := cv.NewMatWithSize(height, width, cv.MatTypeCV8UC1)
	for y := 20; y < height-20; y++ {
		for x := 20; x < width-20; x++ {
			on := y%20 < 8
			if alternate {
				on = on || ((x+y)/18)%3 == 0
			} else {
				on = on || (x/30)%2 == 0
			}
			if on {
				mask.SetUCharAt(y, x, 255)
			}
		}
	}
	return mask
}

func imageRect(left, top, right, bottom int) image.Rectangle {
	return image.Rect(left, top, right, bottom)
}

func featureValues(f *fe.Fingerprint) map[string]float32 {
	values := make(map[string]float32, len(f.Features))
	for _, feature := range f.Features {
		values[feature.Name] = feature.Value
	}
	return values
}
