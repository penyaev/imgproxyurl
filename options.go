package imgproxyurl

import (
	"fmt"
	"strings"
)

type Option interface{}

type ProcessingOption interface {
	Key() string
	String() string
}

func format(key string, arguments ...interface{}) string {
	var ss []string
	for _, argument := range arguments {
		ss = append(ss, fmt.Sprint(argument))
	}
	return strings.Join(ss, ":")
}

// Width defines the width of the resulting image.
// When set to 0, imgproxy will calculate the resulting width using the defined height and source aspect ratio.
type Width struct {
	W int
}

func (Width) Key() string {
	return "w"
}
func (o Width) String() string {
	return format(o.Key(), o.W)
}

// Height defines the height of the resulting image.
// When set to 0, imgproxy will calculate resulting height using the defined width and source aspect ratio.
type Height struct {
	H int
}

func (Height) Key() string {
	return "h"
}
func (o Height) String() string {
	return format(o.Key(), o.H)
}

//ResizingTypeName defines how imgproxy will resize the source image.
type ResizingTypeName string

const (
	//ResizingTypeFit resizes the image while keeping aspect ratio to fit given size
	ResizingTypeFit ResizingTypeName = "fit"
	//ResizingTypeFill resizes the image while keeping aspect ratio to fill given size and cropping projecting parts
	ResizingTypeFill ResizingTypeName = "fill"
	//if both source and resulting dimensions have the same orientation (portrait or landscape), imgproxy will use ResizingTypeFill. Otherwise, it will use ResizingTypeFit
	ResizingTypeAuto ResizingTypeName = "auto"
)

// ResizingType defines how imgproxy will resize the source image
type ResizingType struct {
	ResizingType ResizingTypeName
}

func (ResizingType) Key() string {
	return "rt"
}
func (o ResizingType) String() string {
	return format(o.Key(), o.ResizingType)
}

type ResizingAlgorithmName string

const (
	ResizingAlgorithmNearest  ResizingAlgorithmName = "nearest"
	ResizingAlgorithmLinear   ResizingAlgorithmName = "linear"
	ResizingAlgorithmCubic    ResizingAlgorithmName = "cubic"
	ResizingAlgorithmLanczos2 ResizingAlgorithmName = "lanczos2"
	ResizingAlgorithmLanczos3 ResizingAlgorithmName = "lanczos3"
)

// ResizingAlgorithm defines the algorithm that imgproxy will use for resizing
type ResizingAlgorithm struct {
	ResizingAlgorithm ResizingAlgorithmName
}

func (ResizingAlgorithm) Key() string {
	return "ra"
}
func (o ResizingAlgorithm) String() string {
	return format(o.Key(), o.ResizingAlgorithm)
}

// When set, imgproxy will multiply the image dimensions according to this factor for HiDPI (Retina) devices.
// The value must be greater than 0.
type Dpr struct {
	Dpr int
}

func (Dpr) Key() string {
	return "dpr"
}
func (o Dpr) String() string {
	return format(o.Key(), o.Dpr)
}

// When set, imgproxy will enlarge the image if it is smaller than the given size.
type Enlarge struct {
	Enlarge bool
}

func (Enlarge) Key() string {
	return "el"
}
func (o Enlarge) String() string {
	return format(o.Key(), o.Enlarge)
}

//When extend is set to true, imgproxy will extend the image if it is smaller than the given size.
type Extend struct {
	Extend  bool
	Gravity *Gravity
}

func (Extend) Key() string {
	return "ex"
}
func (o Extend) String() string {
	var arguments = []interface{}{o.Extend}
	if o.Gravity != nil {
		arguments = append(arguments, o.Gravity)
	}
	return format(o.Key(), arguments...)
}

//Defines an area of the image to be processed (crop before resize).
//
//Width and height define the size of the area:
//
//When width or height is greater than or equal to 1, imgproxy treats it as an absolute value.
//
//When width or height is less than 1, imgproxy treats it as a relative value.
//
//When width or height is set to 0, imgproxy will use the full width/height of the source image.
type Crop struct {
	Width   float64
	Height  float64
	Gravity *Gravity
}

func (Crop) Key() string {
	return "c"
}
func (o Crop) String() string {
	var arguments = []interface{}{o.Width, o.Height}
	if o.Gravity != nil {
		arguments = append(arguments, o.Gravity)
	}
	return format(o.Key(), arguments...)
}

//Defines padding size in css manner. All arguments are optional but at least one dimension must be set. Padded space is filled according to background option.
type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

func (Padding) Key() string {
	return "pd"
}
func (o Padding) String() string {
	return format(o.Key(), o.Top, o.Right, o.Bottom, o.Left)
}

type GravityType string

const (
	//north (top edge)
	GravityTypeNorth GravityType = "no"

	//south (bottom edge)
	GravityTypeSouth GravityType = "so"

	//east (right edge)
	GravityTypeEast GravityType = "ea"

	//west (left edge)
	GravityTypeWest GravityType = "we"

	//north-east (top-right corner)
	GravityTypeNorthEast GravityType = "noea"

	//north-west (top-left corner)
	GravityTypeNorthWest GravityType = "nowe"

	//south-east (bottom-right corner)
	GravityTypeSouthEast GravityType = "soea"

	//south-west (bottom-left corner)
	GravityTypeSouthWest GravityType = "sowe"

	//center
	GravityTypeCenter GravityType = "ce"

	//Smart gravity. libvips detects the most “interesting” section of the image and considers it as the center of the resulting image. Offsets are not applicable here;
	GravityTypeSmart GravityType = "sm"

	//Focus point gravity. Offsets are floating point numbers between 0 and 1 that define the coordinates of the center of the resulting image. Treat 0 and 1 as right/left for x and top/bottom for y.
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

func (g GravityIntegerOffsets) String() string {
	return fmt.Sprintf("%d:%d", g.X, g.Y)
}

type GravityFloatOffsets struct {
	X float64
	Y float64
}

func (g GravityFloatOffsets) IsGravityOffset() bool {
	return true
}

func (g GravityFloatOffsets) String() string {
	return fmt.Sprintf("%v:%v", g.X, g.Y)
}

//When imgproxy needs to cut some parts of the image, it is guided by the gravity.
type Gravity struct {
	Type    GravityType
	Offsets GravityOffsets
}

func (Gravity) Key() string {
	return "g"
}
func (o Gravity) String() string {
	var arguments = []interface{}{o.Type}

	if o.Offsets != nil {
		arguments = append(arguments, o.Offsets)
	}

	return format(o.Key(), arguments...)
}

//When set, imgproxy will apply the sharpen filter to the resulting image
//
//As an approximate guideline, use 0.5 sigma for 4 pixels/mm (display resolution), 1.0 for 12 pixels/mm and 1.5 for 16 pixels/mm (300 dpi == 12 pixels/mm).
type Sharpen struct {
	//Sigma is the size of a mask imgproxy will use.
	Sigma float64
}

func (Sharpen) Key() string {
	return "sh"
}
func (o Sharpen) String() string {
	return format(o.Key(), o.Sigma)
}

//Redefines quality of the resulting image, percentage. When 0, quality is assumed based on IMGPROXY_QUALITY and IMGPROXY_FORMAT_QUALITY.
type Quality struct {
	Quality int
}

func (Quality) Key() string {
	return "q"
}
func (o Quality) String() string {
	return format(o.Key(), o.Quality)
}

//When set, imgproxy automatically degrades the quality of the image until the image is under the specified amount of bytes.
//
//Note: Applicable only to jpg, webp, heic, and tiff.
type MaxBytes struct {
	MaxBytes int
}

func (MaxBytes) Key() string {
	return "mb"
}
func (o MaxBytes) String() string {
	return format(o.Key(), o.MaxBytes)
}

//When set, imgproxy will fill the resulting image background with the specified color. HexColor is a hex-coded value of the color. Useful when you convert an image with alpha-channel to JPEG.
type BackgroundHex struct {
	HexColor string
}

func (BackgroundHex) Key() string {
	return "bg"
}
func (o BackgroundHex) String() string {
	return format(o.Key(), o.HexColor)
}

//When set, imgproxy will fill the resulting image background with the specified color. R, G, and B are red, green and blue channel values of the background color (0-255). Useful when you convert an image with alpha-channel to JPEG.
type BackgroundRGB struct {
	R byte
	G byte
	B byte
}

func (BackgroundRGB) Key() string {
	return "bg"
}
func (o BackgroundRGB) String() string {
	return format(o.Key(), o.R, o.G, o.B)
}

//Adds alpha channel to background. alpha is a positive floating point number between 0 and 1.
type BackgroundAlpha struct {
	Alpha float64
}

func (BackgroundAlpha) Key() string {
	return "bga"
}
func (o BackgroundAlpha) String() string {
	return format(o.Key(), o.Alpha)
}

//Defines a list of presets to be used by imgproxy. Feel free to use as many presets in a single URL as you need.
type Presets struct {
	Presets []string
}

func (Presets) Key() string {
	return "pr"
}
func (o Presets) String() string {
	return format(o.Key(), strings.Join(o.Presets, ":"))
}

//Removes surrounding background.
type Trim struct {
	// Color similarity tolerance.
	Threshold int
	// Hex-coded value of the color that needs to be cut off
	Color string
	// When set, imgproxy will cut only equal parts from left and right sides. That means that if 10px of background can be cut off from left and 5px from right then 5px will be cut off from both sides. For example, it can be useful if objects on your images are centered but have non-symmetrical shadow.
	EqualHor bool
	// Acts like EqualHor but for top/bottom sides.
	EqualVer bool
}

func (Trim) Key() string {
	return "t"
}
func (o Trim) String() string {
	return format(o.Key(), o.Threshold, o.Color, o.EqualHor, o.EqualVer)
}

//Rotates the image on the specified angle. The orientation from the image metadata is applied before the rotation unless autorotation is disabled.
type Rotate struct {
	//Only 0/90/180/270/etc degrees angles are supported.
	Angle int
}

func (Rotate) Key() string {
	return "rot"
}
func (o Rotate) String() string {
	return format(o.Key(), o.Angle)
}

//When set, imgproxy will apply the gaussian blur filter to the resulting image
type Blur struct {
	//Sigma defines the size of a mask imgproxy will use.
	Sigma int
}

func (Blur) Key() string {
	return "bl"
}
func (o Blur) String() string {
	return format(o.Key(), o.Sigma)
}

//When set, imgproxy will automatically rotate images based onon the EXIF Orientation parameter (if available in the image meta data). The orientation tag will be removed from the image anyway. Normally this is controlled by the IMGPROXY_AUTO_ROTATE configuration but this procesing option allows the configuration to be set for each request.
type AutoRotate struct {
	AutoRotate bool
}

func (AutoRotate) Key() string {
	return "ar"
}
func (o AutoRotate) String() string {
	return format(o.Key(), o.AutoRotate)
}

//Defines a filename for Content-Disposition header. When not specified, imgproxy will get filename from the source url.
type Filename struct {
	Filename string
}

func (Filename) Key() string {
	return "fn"
}
func (o Filename) String() string {
	return format(o.Key(), o.Filename)
}

type Raw struct {
	OptionKey  string
	Parameters []interface{}
}

func (o Raw) Key() string {
	return o.OptionKey
}
func (o Raw) String() string {
	return format(o.Key(), o.Parameters...)
}

type Format struct {
	Format string
}

type SourceUrl struct {
	Url string
}

type PlainSourceUrl struct {
	Plain bool
}

type Key struct {
	Key string
}

type Salt struct {
	Salt string
}

type KeyRaw struct {
	KeyRaw []byte
}

type SaltRaw struct {
	SaltRaw []byte
}

type Endpoint struct {
	Endpoint string
}

type SignatureSize struct {
	SignatureSize int
}
