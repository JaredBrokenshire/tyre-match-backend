package feature_extraction

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"

	cv "gocv.io/x/gocv"
)

const (
	canonicalWidth        = 64
	canonicalHeight       = 32
	profileBins           = 8
	gridRows              = 2
	gridCols              = 4
	minimumComponentCells = 3
	extractionMethod      = "named_multiscale_tread_fingerprint_v3"
)

type Feature struct {
	Name        string  `json:"name"`
	Group       string  `json:"group"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	Value       float32 `json:"value"`
	Weight      float32 `json:"weight"`
}

type QualityMetric struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Unit        string  `json:"unit"`
	Value       float32 `json:"value"`
}

type Fingerprint struct {
	ExtractionMethod     string          `json:"extraction_method"`
	InputWidth           int             `json:"input_width"`
	InputHeight          int             `json:"input_height"`
	ForegroundConvention string          `json:"foreground_convention"`
	PixelsPerInch        float32         `json:"pixels_per_inch"`
	Features             []Feature       `json:"features"`
	Quality              []QualityMetric `json:"quality"`
}

type SimilarityBreakdown struct {
	Final          float32
	Pattern        float32
	Longitudinal   float32
	StructureCount float32
	Scale          float32
}

type Extractor struct{}

func NewExtractor() *Extractor { return &Extractor{} }

func (e *Extractor) Extract(binary *cv.Mat, ppi float32) (*Fingerprint, error) {
	if binary == nil || binary.Empty() {
		return nil, errors.New("feature extraction received an empty image")
	}
	if binary.Type() != cv.MatTypeCV8UC1 {
		return nil, fmt.Errorf("feature extraction requires an 8-bit single-channel binary image, got type %v", binary.Type())
	}
	if ppi <= 0 || math.IsNaN(float64(ppi)) || math.IsInf(float64(ppi), 0) {
		return nil, errors.New("feature extraction requires a positive pixels-per-inch value")
	}

	imageWidth, imageHeight := binary.Cols(), binary.Rows()
	imageSize := imageWidth * imageHeight

	source := binary.ToBytes()

	input := make([]float32, imageSize)
	for i, value := range source[:imageSize] {
		input[i] = float32(value)
	}

	// The tread mask is expected to use white foreground and black background.
	// We reduce it to a fixed occupancy grid before calculating structural
	// descriptors. This makes the features tolerant of substrate noise while
	// retaining the repeating geometry of the tread pattern.
	coarse := occupancyGrid(input, imageWidth, imageHeight, canonicalWidth, canonicalHeight)
	if countForeground(coarse) == 0 {
		return nil, errors.New("feature extraction input contains no foreground pixels")
	}
	clean := removeSmallComponents(coarse, canonicalWidth, canonicalHeight, minimumComponentCells)
	if countForeground(clean) == 0 {
		return nil, errors.New("feature extraction input contains no usable foreground after cleanup")
	}

	bbox := foregroundBBox(clean, canonicalWidth, canonicalHeight)
	if bbox.empty() {
		return nil, errors.New("feature extraction could not determine the foreground bounding box")
	}
	canonical, cw, ch := cropMask(clean, canonicalWidth, bbox)
	if cw < 2 || ch < 2 {
		return nil, errors.New("feature extraction foreground is too small for structural analysis")
	}

	features := make([]Feature, 0, 45)
	add := func(group, name, description, unit string, value, weight float32) {
		if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
			return
		}
		features = append(features, Feature{Name: name, Group: group, Description: description, Unit: unit, Value: value, Weight: weight})
	}

	foreground := float64(countForeground(canonical)) / float64(cw*ch)
	centroidX, centroidY := foregroundCentroid(canonical, cw, ch)
	physicalWidthMm := float64(cw) / float64(canonicalWidth) * float64(imageWidth) / float64(ppi) * 25.4
	physicalHeightMm := float64(ch) / float64(canonicalHeight) * float64(imageHeight) / float64(ppi) * 25.4

	add("global_geometry", "foreground_area_ratio", "Proportion of the cleaned canonical tread region occupied by white tread/contact pixels.", "ratio", float32(foreground), 1.0)
	add("global_geometry", "foreground_centroid_x", "Horizontal centroid of the cleaned tread/contact region expressed as a normalised coordinate.", "normalised coordinate", float32(centroidX), 0.6)
	add("global_geometry", "foreground_centroid_y", "Vertical centroid of the cleaned tread/contact region expressed as a normalised coordinate.", "normalised coordinate", float32(centroidY), 0.6)
	add("physical_geometry", "foreground_width_mm", "Estimated physical width of the cleaned foreground tread region using the submitted image scale.", "mm", float32(physicalWidthMm), 0.8)
	add("physical_geometry", "foreground_height_mm", "Estimated physical height of the cleaned foreground tread region using the submitted image scale.", "mm", float32(physicalHeightMm), 0.8)
	add("physical_geometry", "foreground_aspect_ratio", "Width-to-height ratio of the cleaned tread region.", "ratio", safeFloat32(float64(cw)/float64(ch)), 0.5)

	longitudinal := profileCols(canonical, cw, ch, profileBins)
	transverse := profileRows(canonical, cw, ch, profileBins)
	for i, value := range longitudinal {
		add("longitudinal_structure", fmt.Sprintf("longitudinal_density_bin_%02d", i+1), fmt.Sprintf("Mean tread/contact occupancy in longitudinal band %02d of %d across the horizontal tread direction.", i+1, profileBins), "ratio", value, 1.0)
	}
	for i, value := range transverse {
		add("transverse_structure", fmt.Sprintf("transverse_density_bin_%02d", i+1), fmt.Sprintf("Mean tread/contact occupancy in transverse band %02d of %d across the vertical tread direction.", i+1, profileBins), "ratio", value, 0.75)
	}

	grid := occupancyGrid(canonical, cw, ch, gridCols, gridRows)
	for y := 0; y < gridRows; y++ {
		for x := 0; x < gridCols; x++ {
			add("spatial_pattern", fmt.Sprintf("spatial_density_row_%02d_col_%02d", y+1, x+1), fmt.Sprintf("Foreground occupancy of spatial tread cell row %02d, column %02d in a %dx%d grid.", y+1, x+1, gridRows, gridCols), "ratio", grid[y*gridCols+x], 1.0)
		}
	}

	add("profile_shape", "longitudinal_profile_contrast", "Standard deviation of longitudinal foreground occupancy, describing repeated density changes along the tread direction.", "ratio", profileStdDev(longitudinal), 0.8)
	add("profile_shape", "transverse_profile_contrast", "Standard deviation of transverse foreground occupancy, describing density changes across tread width.", "ratio", profileStdDev(transverse), 0.6)

	longRepeat, longStrength := periodicity(longitudinal)
	transRepeat, transStrength := periodicity(transverse)
	longRepeatMm := longRepeat / float64(profileBins) * physicalWidthMm
	transRepeatMm := transRepeat / float64(profileBins) * physicalHeightMm
	add("periodicity", "longitudinal_repeat_distance_mm", "Dominant repeated structural distance along the horizontal tread direction, converted to millimetres using PPI.", "mm", safeFloat32(longRepeatMm), 1.4)
	add("periodicity", "longitudinal_periodicity_strength", "Strength of the dominant repeated structure along the horizontal tread direction.", "ratio", float32(longStrength), 1.2)
	add("periodicity", "transverse_repeat_distance_mm", "Dominant repeated structural distance across the vertical tread direction, converted to millimetres using PPI.", "mm", safeFloat32(transRepeatMm), 0.8)
	add("periodicity", "transverse_periodicity_strength", "Strength of the dominant repeated structure across the vertical tread direction.", "ratio", float32(transStrength), 0.7)

	components := componentStats(canonical, cw, ch)
	add("tread_topology", "substantial_component_count", "Number of substantial connected foreground regions after coarse-scale noise suppression.", "count", float32(components.Count), 0.8)
	add("tread_topology", "largest_component_area_ratio", "Area of the largest substantial connected foreground region relative to the cleaned canonical region.", "ratio", float32(components.LargestAreaRatio), 0.6)
	add("tread_topology", "component_area_cv", "Coefficient of variation of substantial connected-component areas.", "ratio", float32(components.AreaCV), 0.5)
	add("tread_topology", "component_fragmentation", "Number of substantial foreground components relative to their total foreground area, describing fragmentation of the tread representation.", "ratio", float32(components.Fragmentation), 0.5)

	dominantAngle, orientationConcentration := orientationSummary(canonical, cw, ch)
	add("orientation", "dominant_boundary_orientation_deg", "Dominant local boundary orientation estimated from transitions between tread/contact and background regions.", "degrees", float32(dominantAngle), 0.7)
	add("orientation", "boundary_orientation_concentration", "Proportion of boundary-transition energy concentrated in the dominant orientation sector.", "ratio", float32(orientationConcentration), 0.5)

	before := countForeground(coarse)
	after := countForeground(clean)
	quality := []QualityMetric{
		{Name: "input_foreground_ratio", Description: "Foreground proportion of the coarse tread mask before deterministic small-component cleanup.", Unit: "ratio", Value: safeFloat32(float64(before) / float64(len(coarse)))},
		{Name: "cleanup_removed_ratio", Description: "Proportion of coarse foreground cells removed as isolated small components.", Unit: "ratio", Value: safeFloat32(float64(before-after) / math.Max(float64(before), 1))},
		{Name: "canonical_width_cells", Description: "Width of the cleaned foreground region used for structural comparison.", Unit: "cells", Value: float32(cw)},
		{Name: "canonical_height_cells", Description: "Height of the cleaned foreground region used for structural comparison.", Unit: "cells", Value: float32(ch)},
		{Name: "scale_pixels_per_inch", Description: "User-supplied image scale used for physical feature conversion.", Unit: "pixels/inch", Value: ppi},
	}

	return &Fingerprint{
		ExtractionMethod: extractionMethod,
		InputWidth:       binary.Cols(), InputHeight: binary.Rows(),
		ForegroundConvention: "white_foreground_black_background",
		PixelsPerInch:        ppi,
		Features:             features, Quality: quality,
	}, nil
}

func (e *Extractor) Similarity(a, b *Fingerprint) float32 {
	return e.Compare(a, b).Final
}

func (e *Extractor) Compare(a, b *Fingerprint) SimilarityBreakdown {
	if a == nil || b == nil || a.ExtractionMethod != b.ExtractionMethod || a.ForegroundConvention != b.ForegroundConvention {
		return SimilarityBreakdown{}
	}
	am, bm := featureMap(a), featureMap(b)
	pattern := groupedSimilarity(am, bm, "spatial_pattern", "global_geometry")
	longitudinal := groupedSimilarity(am, bm, "longitudinal_structure", "periodicity")
	structure := groupedSimilarity(am, bm, "tread_topology", "orientation", "transverse_structure")
	scale := scaleSimilarity(am, bm)
	final := 0.60*pattern + 0.25*longitudinal + 0.10*structure + 0.05*scale
	return SimilarityBreakdown{Final: float32(clamp01(final)), Pattern: float32(pattern), Longitudinal: float32(longitudinal), StructureCount: float32(structure), Scale: float32(scale)}
}

func Encode(f *Fingerprint) (string, error) {
	if f == nil {
		return "", errors.New("cannot encode nil fingerprint")
	}
	if _, err := validateFingerprint(f); err != nil {
		return "", err
	}
	data, err := json.Marshal(f)
	if err != nil {
		return "", fmt.Errorf("marshal fingerprint: %w", err)
	}
	return string(data), nil
}

func Decode(encoded string) (*Fingerprint, error) {
	if encoded == "" {
		return nil, errors.New("empty fingerprint")
	}
	var f Fingerprint
	if err := json.Unmarshal([]byte(encoded), &f); err != nil {
		return nil, fmt.Errorf("unmarshal fingerprint: %w", err)
	}
	if _, err := validateFingerprint(&f); err != nil {
		return nil, err
	}
	return &f, nil
}

func validateFingerprint(f *Fingerprint) (map[string]struct{}, error) {
	if f.ExtractionMethod == "" {
		return nil, errors.New("fingerprint has no extraction method")
	}
	if f.ForegroundConvention == "" {
		return nil, errors.New("fingerprint has no foreground convention")
	}
	if f.PixelsPerInch <= 0 || math.IsNaN(float64(f.PixelsPerInch)) || math.IsInf(float64(f.PixelsPerInch), 0) {
		return nil, errors.New("invalid pixels-per-inch")
	}
	if len(f.Features) == 0 {
		return nil, errors.New("fingerprint has no named features")
	}
	names := make(map[string]struct{}, len(f.Features))
	for _, feature := range f.Features {
		if feature.Name == "" || feature.Group == "" || feature.Description == "" || feature.Unit == "" {
			return nil, errors.New("feature has missing descriptive metadata")
		}
		if feature.Weight <= 0 || math.IsNaN(float64(feature.Weight)) || math.IsInf(float64(feature.Weight), 0) {
			return nil, fmt.Errorf("invalid weight for feature %s", feature.Name)
		}
		if _, exists := names[feature.Name]; exists {
			return nil, fmt.Errorf("duplicate feature %s", feature.Name)
		}
		names[feature.Name] = struct{}{}
	}
	return names, nil
}

func occupancyGrid(mask []float32, width, height, gridWidth, gridHeight int) []float32 {
	out := make([]float32, gridWidth*gridHeight)
	for gy := 0; gy < gridHeight; gy++ {
		y0, y1 := gy*height/gridHeight, (gy+1)*height/gridHeight
		if y1 <= y0 {
			y1 = y0 + 1
		}
		for gx := 0; gx < gridWidth; gx++ {
			x0, x1 := gx*width/gridWidth, (gx+1)*width/gridWidth
			if x1 <= x0 {
				x1 = x0 + 1
			}
			on, total := 0, 0
			for y := y0; y < y1 && y < height; y++ {
				for x := x0; x < x1 && x < width; x++ {
					total++
					if mask[y*width+x] != 0 {
						on++
					}
				}
			}
			if total > 0 {
				out[gy*gridWidth+gx] = float32(on) / float32(total)
			}
		}
	}
	return out
}

func countForeground(mask []float32) int {
	n := 0
	for _, v := range mask {
		if v >= 0.5 {
			n++
		}
	}
	return n
}

func removeSmallComponents(mask []float32, width, height, minimum int) []float32 {
	out := append([]float32(nil), mask...)
	visited := make([]bool, len(mask))
	for i, v := range mask {
		if v < 0.5 || visited[i] {
			continue
		}
		queue := []int{i}
		visited[i] = true
		component := make([]int, 0)
		for len(queue) > 0 {
			p := queue[0]
			queue = queue[1:]
			component = append(component, p)
			x, y := p%width, p/width
			for _, d := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
				nx, ny := x+d[0], y+d[1]
				if nx < 0 || nx >= width || ny < 0 || ny >= height {
					continue
				}
				n := ny*width + nx
				if mask[n] >= 0.5 && !visited[n] {
					visited[n] = true
					queue = append(queue, n)
				}
			}
		}
		if len(component) < minimum {
			for _, p := range component {
				out[p] = 0
			}
		}
	}
	return out
}

type boundingBox struct{ left, top, right, bottom int }

func (b boundingBox) empty() bool { return b.right <= b.left || b.bottom <= b.top }

func foregroundBBox(mask []float32, width, height int) boundingBox {
	b := boundingBox{left: width, top: height, right: -1, bottom: -1}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if mask[y*width+x] < 0.5 {
				continue
			}
			if x < b.left {
				b.left = x
			}
			if y < b.top {
				b.top = y
			}
			if x >= b.right {
				b.right = x + 1
			}
			if y >= b.bottom {
				b.bottom = y + 1
			}
		}
	}
	return b
}

func cropMask(mask []float32, width int, b boundingBox) ([]float32, int, int) {
	if b.empty() {
		return nil, 0, 0
	}
	nw, nh := b.right-b.left, b.bottom-b.top
	out := make([]float32, nw*nh)
	for y := 0; y < nh; y++ {
		copy(out[y*nw:(y+1)*nw], mask[(b.top+y)*width+b.left:(b.top+y)*width+b.right])
	}
	return out, nw, nh
}

func profileCols(mask []float32, width, height, bins int) []float32 {
	out := make([]float32, bins)
	for i := 0; i < bins; i++ {
		x0, x1 := i*width/bins, (i+1)*width/bins
		if x1 <= x0 {
			x1 = x0 + 1
		}
		sum := 0.0
		count := 0
		for y := 0; y < height; y++ {
			for x := x0; x < x1 && x < width; x++ {
				sum += float64(mask[y*width+x])
				count++
			}
		}
		if count > 0 {
			out[i] = float32(sum / float64(count))
		}
	}
	return out
}

func profileRows(mask []float32, width, height, bins int) []float32 {
	out := make([]float32, bins)
	for i := 0; i < bins; i++ {
		y0, y1 := i*height/bins, (i+1)*height/bins
		if y1 <= y0 {
			y1 = y0 + 1
		}
		sum := 0.0
		count := 0
		for y := y0; y < y1 && y < height; y++ {
			for x := 0; x < width; x++ {
				sum += float64(mask[y*width+x])
				count++
			}
		}
		if count > 0 {
			out[i] = float32(sum / float64(count))
		}
	}
	return out
}

func foregroundCentroid(mask []float32, width, height int) (float64, float64) {
	var sx, sy, total float64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			v := float64(mask[y*width+x])
			if v < 0.5 {
				continue
			}
			sx += float64(x) * v
			sy += float64(y) * v
			total += v
		}
	}
	if total == 0 {
		return 0.5, 0.5
	}
	return sx / total / math.Max(float64(width-1), 1), sy / total / math.Max(float64(height-1), 1)
}

func profileStdDev(profile []float32) float32 {
	if len(profile) == 0 {
		return 0
	}
	mean := 0.0
	for _, v := range profile {
		mean += float64(v)
	}
	mean /= float64(len(profile))
	variance := 0.0
	for _, v := range profile {
		d := float64(v) - mean
		variance += d * d
	}
	variance /= float64(len(profile))
	return float32(math.Sqrt(variance))
}

func periodicity(profile []float32) (float64, float64) {
	if len(profile) < 4 {
		return 0, 0
	}
	mean := 0.0
	for _, v := range profile {
		mean += float64(v)
	}
	mean /= float64(len(profile))
	var energy float64
	for _, v := range profile {
		d := float64(v) - mean
		energy += d * d
	}
	if energy <= 1e-9 {
		return 0, 0
	}
	bestLag, best := 0, 0.0
	for lag := 1; lag <= len(profile)/2; lag++ {
		corr := 0.0
		for i := 0; i < len(profile)-lag; i++ {
			corr += (float64(profile[i]) - mean) * (float64(profile[i+lag]) - mean)
		}
		corr /= energy
		if corr > best {
			best = corr
			bestLag = lag
		}
	}
	return float64(bestLag), clamp01(best)
}

type componentSummary struct {
	Count            int
	LargestAreaRatio float64
	AreaCV           float64
	Fragmentation    float64
}

func componentStats(mask []float32, width, height int) componentSummary {
	visited := make([]bool, len(mask))
	areas := []int{}
	for i, v := range mask {
		if v < 0.5 || visited[i] {
			continue
		}
		q := []int{i}
		visited[i] = true
		area := 0
		for len(q) > 0 {
			p := q[0]
			q = q[1:]
			area++
			x, y := p%width, p/width
			for _, d := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
				nx, ny := x+d[0], y+d[1]
				if nx < 0 || nx >= width || ny < 0 || ny >= height {
					continue
				}
				n := ny*width + nx
				if mask[n] >= 0.5 && !visited[n] {
					visited[n] = true
					q = append(q, n)
				}
			}
		}
		if area >= minimumComponentCells {
			areas = append(areas, area)
		}
	}
	if len(areas) == 0 {
		return componentSummary{}
	}
	total, maxArea := 0, 0
	for _, a := range areas {
		total += a
		if a > maxArea {
			maxArea = a
		}
	}
	mean := float64(total) / float64(len(areas))
	variance := 0.0
	for _, a := range areas {
		d := float64(a) - mean
		variance += d * d
	}
	variance /= float64(len(areas))
	return componentSummary{Count: len(areas), LargestAreaRatio: float64(maxArea) / float64(width*height), AreaCV: math.Sqrt(variance) / math.Max(mean, 1), Fragmentation: float64(len(areas)) / math.Max(float64(total), 1)}
}

func orientationSummary(mask []float32, width, height int) (float64, float64) {
	hist := make([]float64, 8)
	total := 0.0
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			left, right := mask[y*width+x-1], mask[y*width+x+1]
			up, down := mask[(y-1)*width+x], mask[(y+1)*width+x]
			dx := float64(right - left)
			dy := float64(down - up)
			mag := math.Hypot(dx, dy)
			if mag < 0.1 {
				continue
			}
			angle := math.Atan2(dy, dx)
			if angle < 0 {
				angle += math.Pi
			}
			bin := int(angle / (math.Pi / 8))
			if bin >= 8 {
				bin = 7
			}
			hist[bin] += mag
			total += mag
		}
	}
	if total == 0 {
		return 0, 0
	}
	bestBin, best := 0, 0.0
	for i, v := range hist {
		if v > best {
			best = v
			bestBin = i
		}
	}
	return float64(bestBin) * 180.0 / 8.0, best / total
}

func featureMap(f *Fingerprint) map[string]Feature {
	m := make(map[string]Feature, len(f.Features))
	for _, v := range f.Features {
		m[v.Name] = v
	}
	return m
}
func groupedSimilarity(a, b map[string]Feature, groups ...string) float64 {
	wanted := map[string]bool{}
	for _, g := range groups {
		wanted[g] = true
	}
	num, den := 0.0, 0.0
	for name, af := range a {
		if !wanted[af.Group] {
			continue
		}
		bf, ok := b[name]
		if !ok {
			continue
		}
		diff := normalisedDifference(af.Value, bf.Value)
		w := math.Max(float64(af.Weight), 0)
		num += w * diff
		den += w
		if name == "dominant_boundary_orientation_deg" {
			_ = name
		}
	}
	if den == 0 {
		return 0
	}
	return clamp01(1 - num/den)
}
func scaleSimilarity(a, b map[string]Feature) float64 {
	pairs := [][2]string{{"foreground_width_mm", "foreground_width_mm"}, {"foreground_height_mm", "foreground_height_mm"}, {"longitudinal_repeat_distance_mm", "longitudinal_repeat_distance_mm"}, {"transverse_repeat_distance_mm", "transverse_repeat_distance_mm"}}
	num, den := 0.0, 0.0
	for _, p := range pairs {
		af, aok := a[p[0]]
		bf, bok := b[p[1]]
		if !aok || !bok {
			continue
		}
		w := math.Max(float64(af.Weight), 0.1)
		num += w * normalisedDifference(af.Value, bf.Value)
		den += w
	}
	if den == 0 {
		return 0
	}
	return clamp01(1 - num/den)
}
func normalisedDifference(a, b float32) float64 {
	af, bf := float64(a), float64(b)
	scale := math.Abs(af) + math.Abs(bf)
	if scale < 1e-6 {
		return 0
	}
	return math.Abs(af-bf) / scale
}
func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
func safeFloat32(v float64) float32 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return float32(v)
}

// Keep deterministic ordering for callers that display the fingerprint.
func SortFeatures(f *Fingerprint) {
	if f == nil {
		return
	}
	sort.Slice(f.Features, func(i, j int) bool { return f.Features[i].Name < f.Features[j].Name })
}
