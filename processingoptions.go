package imgproxyurl

import (
	"errors"
	"fmt"
	"strconv"
)

//ResizingType defines how imgproxy will resize the source image.
type ResizingType string

const (
	//ResizingTypeFit resizes the image while keeping aspect ratio to fit given size
	ResizingTypeFit ResizingType = "fit"
	//ResizingTypeFill resizes the image while keeping aspect ratio to fill given size and cropping projecting parts
	ResizingTypeFill ResizingType = "fill"
	//if both source and resulting dimensions have the same orientation (portrait or landscape), imgproxy will use ResizingTypeFill. Otherwise, it will use ResizingTypeFit
	ResizingTypeAuto ResizingType = "auto"
)

// SetResizingType defines how imgproxy will resize the source image
func (url *Url) SetResizingType(resizingType ResizingType) *Url {
	url.setOption("rt", string(resizingType))
	return url
}

// ResizingAlgorithm defines the algorithm that imgproxy will use for resizing
type ResizingAlgorithm string

const (
	ResizingAlgorithmNearest  ResizingAlgorithm = "nearest"
	ResizingAlgorithmLinear   ResizingAlgorithm = "linear"
	ResizingAlgorithmCubic    ResizingAlgorithm = "cubic"
	ResizingAlgorithmLanczos2 ResizingAlgorithm = "lanczos2"
	ResizingAlgorithmLanczos3 ResizingAlgorithm = "lanczos3"
)

// SetResizingAlgorithm defines the algorithm that imgproxy will use for resizing
func (url *Url) SetResizingAlgorithm(resizingAlgorithm ResizingAlgorithm) *Url {
	url.setOption("ra", string(resizingAlgorithm))
	return url
}

// SetWidth defines the width of the resulting image.
// When set to 0, imgproxy will calculate the resulting width using the defined height and source aspect ratio.
func (url *Url) SetWidth(width int) *Url {
	url.setOption("w", strconv.Itoa(width))
	return url
}

// SetHeight defines the height of the resulting image.
// When set to 0, imgproxy will calculate resulting height using the defined width and source aspect ratio.
func (url *Url) SetHeight(height int) *Url {
	url.setOption("h", strconv.Itoa(height))
	return url
}

// When set, imgproxy will multiply the image dimensions according to this factor for HiDPI (Retina) devices.
// The value must be greater than 0.
func (url *Url) SetDpr(dpr int) *Url {
	url.setOption("dpr", strconv.Itoa(dpr))
	return url
}

// When set, imgproxy will enlarge the image if it is smaller than the given size.
func (url *Url) SetEnlarge(enlarge bool) *Url {
	if enlarge {
		url.setOption("el", "1")
	} else {
		url.unsetOption("el")
	}

	return url
}

type GravityType string

const (
	//GravityTypeDefault lets imgproxy use the default gravity. This is equal to not specifying the gravity on the url
	GravityTypeDefault    GravityType = ""
	GravityTypeNorth      GravityType = "no"
	GravityTypeSouth      GravityType = "so"
	GravityTypeEast       GravityType = "ea"
	GravityTypeWest       GravityType = "we"
	GravityTypeNorthEast  GravityType = "noea"
	GravityTypeNorthWest  GravityType = "nowe"
	GravityTypeSouthEast  GravityType = "soea"
	GravityTypeSouthWest  GravityType = "sowe"
	GravityTypeCenter     GravityType = "ce"
	GravityTypeSmart      GravityType = "sm"
	GravityTypeFocusPoint GravityType = "fp"
)

type GravityOffsets interface {
	IsGravityOffset() bool
}

type GravityIntegerOffsets struct {
	X int
	Y int
}

func (g GravityIntegerOffsets) IsGravityOffset() bool {
	return true
}

type GravityFloatOffsets struct {
	X float64
	Y float64
}

func (g GravityFloatOffsets) IsGravityOffset() bool {
	return true
}

func (url *Url) getGravityArguments(gravityType GravityType, offsets GravityOffsets) []string {
	if gravityType == GravityTypeDefault {
		url.setError(errors.New("specific gravity type is required"))
	}

	var gravityOffsetsTypeInteger bool
	if offsets != nil {
		switch offsets.(type) {
		case GravityIntegerOffsets:
			gravityOffsetsTypeInteger = true
		case GravityFloatOffsets:
			gravityOffsetsTypeInteger = false
			if offsets.(GravityFloatOffsets).X < 0 ||
				offsets.(GravityFloatOffsets).X > 1 ||
				offsets.(GravityFloatOffsets).Y < 0 ||
				offsets.(GravityFloatOffsets).Y > 1 {
				url.setError(errors.New("float offsets must within (0,1) range"))
			}
		}
	}

	arguments := []string{string(gravityType)}
	if gravityType == GravityTypeSmart {
		if offsets != nil {
			url.setError(errors.New("offsets are not applicable for smart gravity"))
		}
	} else if gravityType == GravityTypeFocusPoint {
		if offsets == nil {
			url.setError(errors.New("offsets are required for focus point gravity"))
		} else if gravityOffsetsTypeInteger {
			url.setError(errors.New("focus point gravity requires floating-point offsets"))
		} else {
			arguments = append(arguments, fmt.Sprintf("%.3f", offsets.(GravityFloatOffsets).X), fmt.Sprintf("%.3f", offsets.(GravityFloatOffsets).Y))
		}
	} else {
		if offsets != nil && !gravityOffsetsTypeInteger {
			url.setError(errors.New("integer offsets are required"))
		} else if offsets != nil {
			arguments = append(arguments, strconv.Itoa(offsets.(GravityIntegerOffsets).X), strconv.Itoa(offsets.(GravityIntegerOffsets).Y))
		}
	}

	return arguments
}

//When imgproxy needs to cut some parts of the image, it is guided by the gravity.
func (url *Url) SetGravity(gravityType GravityType, offsets GravityOffsets) *Url {
	url.setOption("g", url.getGravityArguments(gravityType, offsets)...)
	return url
}

func (url *Url) SetExtendWithGravityOffsets(extend bool, gravityType GravityType, gravityOffsets GravityOffsets) *Url {
	if !extend {
		url.unsetOption("ex")
	} else {
		arguments := []string{"1"}

		if gravityType == GravityTypeDefault {
			if gravityOffsets != nil {
				url.setError(errors.New("offsets are not applicable for default gravity"))
			}
		} else if gravityType == GravityTypeSmart {
			url.setError(errors.New("smart gravity type is not applicable here"))
		} else {
			arguments = append(arguments, url.getGravityArguments(gravityType, gravityOffsets)...)
		}

		url.setOption("ex", arguments...)
	}

	return url
}
func (url *Url) SetExtendWithGravity(extend bool, gravityType GravityType) *Url {
	return url.SetExtendWithGravityOffsets(extend, gravityType, nil)
}
func (url *Url) SetExtend(extend bool) *Url {
	return url.SetExtendWithGravityOffsets(extend, GravityTypeDefault, nil)
}

func (url *Url) SetCropWithGravityOffsets(width int, height int, gravityType GravityType, gravityOffsets GravityOffsets) *Url {
	arguments := []string{strconv.Itoa(width), strconv.Itoa(height)}

	if gravityType == GravityTypeDefault {
		if gravityOffsets != nil {
			url.setError(errors.New("offsets are not applicable for default gravity"))
		}
	} else {
		arguments = append(arguments, url.getGravityArguments(gravityType, gravityOffsets)...)
	}

	url.setOption("c", arguments...)

	return url
}
func (url *Url) SetCropWithGravity(width int, height int, gravityType GravityType) *Url {
	return url.SetCropWithGravityOffsets(width, height, gravityType, nil)
}
func (url *Url) SetCrop(width int, height int) *Url {
	return url.SetCropWithGravityOffsets(width, height, GravityTypeDefault, nil)
}

func (url *Url) SetPadding(top int, right int, bottom int, left int) *Url {
	url.setOption("pd", strconv.Itoa(top), strconv.Itoa(right), strconv.Itoa(bottom), strconv.Itoa(left))
	return url
}

func (url *Url) SetPaddingAll(padding int) *Url {
	url.setOption("pd", strconv.Itoa(padding))
	return url
}

type TrimOption interface {
	IsTrimOptions() bool
}
type TrimOptionColor struct {
	Color string
}

func (t TrimOptionColor) IsTrimOptions() bool {
	return true
}

type TrimOptionEqualHor struct{}

func (t TrimOptionEqualHor) IsTrimOptions() bool {
	return true
}

type TrimOptionEqualVer struct{}

func (t TrimOptionEqualVer) IsTrimOptions() bool {
	return true
}

//SetTrim Removes surrounding background.
func (url *Url) SetTrim(threshold int, options ...TrimOption) *Url {
	arguments := []string{strconv.Itoa(threshold), "", "", ""}

	for _, option := range options {
		switch option.(type) {
		case TrimOptionColor:
			arguments[1] = option.(TrimOptionColor).Color
		case TrimOptionEqualHor:
			arguments[2] = "1"
		case TrimOptionEqualVer:
			arguments[3] = "1"
		}
	}

	url.setOption("t", arguments...)
	return url
}

//SetQuality Redefines quality of the resulting image, percentage.
func (url *Url) SetQuality(quality int) *Url {
	url.setOption("q", strconv.Itoa(quality))
	return url
}

//When set, imgproxy automatically degrades the quality of the image until the image is under the specified amount of bytes.
func (url *Url) SetMaxBytes(maxBytes int) *Url {
	url.setOption("mb", strconv.Itoa(maxBytes))
	return url
}

// When set, imgproxy will fill the resulting image background with the specified color.
// red, green and blue are channel values of the background color (0-255).
// Useful when you convert an image with alpha-channel to JPEG.
//
// When not set, disables any background manipulations.
func (url *Url) SetBackgroundRGB(red int, green int, blue int) *Url {
	url.setOption("bg", strconv.Itoa(red), strconv.Itoa(green), strconv.Itoa(blue))
	return url
}

// When set, imgproxy will fill the resulting image background with the specified color.
// color is a hex-coded value of the color..
// Useful when you convert an image with alpha-channel to JPEG.
//
// When not set, disables any background manipulations.
func (url *Url) SetBackgroundHex(color string) *Url {
	url.setOption("bg", color)
	return url
}

//SetBackgroundAlpha adds alpha channel to background. alpha is a positive floating point number between 0 and 1.
func (url *Url) SetBackgroundAlpha(alpha float64) *Url {
	url.setOption("bga", fmt.Sprintf("%.3f", alpha))
	return url
}
