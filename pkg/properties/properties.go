package properties

// Alignment represents text alignment.
type Alignment string

const (
	AlignmentLeft   Alignment = "left"
	AlignmentRight  Alignment = "right"
	AlignmentCenter Alignment = "center"
)

// VAlignment represents vertical alignment (for cells).
type VAlignment string

const (
	VAlignmentTop    VAlignment = "top"
	VAlignmentCenter VAlignment = "center"
	VAlignmentBottom VAlignment = "bottom"
)

// BorderStyle represents the style of a cell border.
type BorderStyle string

const (
	BorderDotted BorderStyle = "dotted"
	BorderSolid  BorderStyle = "solid"
	BorderDouble BorderStyle = "double"
	BorderDashed BorderStyle = "dashed"
	BorderNone   BorderStyle = "none"
)

// DataDetailsType represents the type of a data detail block.
type DataDetailsType string

const (
	DataGroup DataDetailsType = "group"
	DataTable DataDetailsType = "table"
)

// DetailsType represents the type of a structural detail block.
type DetailsType string

const (
	DetailsGroup DetailsType = "group"
	DetailsRow   DetailsType = "layout-h"
)

// Font represents the available fonts.
type Font string

const (
	FontArial         Font = "arial"
	FontTimesNewRoman Font = "times-new-roman"
)

// Orientation represents the page orientation.
type Orientation int

const (
	OrientationPortrait  Orientation = 0
	OrientationLandscape Orientation = 1
)

// PageSize represents the size of the page.
type PageSize int

const (
	// ISO A Series (most common)
	PageSizeA0  PageSize = 0  // 841 x 1189 mm
	PageSizeA1  PageSize = 1  // 594 x 841 mm
	PageSizeA2  PageSize = 2  // 420 x 594 mm
	PageSizeA3  PageSize = 3  // 297 x 420 mm
	PageSizeA4  PageSize = 4  // 210 x 297 mm (default)
	PageSizeA5  PageSize = 5  // 148 x 210 mm
	PageSizeA6  PageSize = 6  // 105 x 148 mm
	PageSizeA7  PageSize = 7  // 74 x 105 mm
	PageSizeA8  PageSize = 8  // 52 x 74 mm
	PageSizeA9  PageSize = 9  // 37 x 52 mm
	PageSizeA10 PageSize = 10 // 26 x 37 mm

	// North American Sizes
	PageSizeLetter  PageSize = 11 // 8.5 x 11 inches (215.9 x 279.4 mm)
	PageSizeLegal   PageSize = 12 // 8.5 x 14 inches (215.9 x 355.6 mm)
	PageSizeTabloId PageSize = 13 // 11 x 17 inches (279.4 x 431.8 mm)

	// ISO C Series (Envelopes)
	PageSizeC0  PageSize = 14 // 917 x 1297 mm
	PageSizeC1  PageSize = 15 // 648 x 917 mm
	PageSizeC2  PageSize = 16 // 458 x 648 mm
	PageSizeC3  PageSize = 17 // 324 x 458 mm
	PageSizeC4  PageSize = 18 // 229 x 324 mm
	PageSizeC5  PageSize = 19 // 162 x 229 mm
	PageSizeC6  PageSize = 20 // 114 x 162 mm
	PageSizeC7  PageSize = 21 // 81 x 114 mm
	PageSizeC8  PageSize = 22 // 57 x 81 mm
	PageSizeC9  PageSize = 23 // 40 x 57 mm
	PageSizeC10 PageSize = 24 // 28 x 40 mm

	// ISO B Series (Between A and C)
	PageSizeB0  PageSize = 25 // 1000 x 1414 mm
	PageSizeB1  PageSize = 26 // 707 x 1000 mm
	PageSizeB2  PageSize = 27 // 500 x 707 mm
	PageSizeB3  PageSize = 28 // 353 x 500 mm
	PageSizeB4  PageSize = 29 // 250 x 353 mm
	PageSizeB5  PageSize = 30 // 176 x 250 mm
	PageSizeB6  PageSize = 31 // 125 x 176 mm
	PageSizeB7  PageSize = 32 // 88 x 125 mm
	PageSizeB8  PageSize = 33 // 62 x 88 mm
	PageSizeB9  PageSize = 34 // 44 x 62 mm
	PageSizeB10 PageSize = 35 // 31 x 44 mm

	// DL Envelope (common for documents)
	PageSizeDL PageSize = 36 // 110 x 220 mm
)

// SeparatorBehavior controls how the separator acts in a CText component.
type SeparatorBehavior string

const (
	SeparatorDefault        SeparatorBehavior = "default"
	SeparatorRepel          SeparatorBehavior = "repel"
	SeparatorRepelLeftHalf  SeparatorBehavior = "repel-left-half"
	SeparatorRepelRightHalf SeparatorBehavior = "repel-right-half"
	SeparatorRepelLeftFull  SeparatorBehavior = "repel-left-full"
	SeparatorRepelRightFull SeparatorBehavior = "repel-right-full"
)

// StyleType represents the type of style component.
type StyleType string

const (
	StyleCText            StyleType = "ctext-style"
	StyleText             StyleType = "text-style"
	StyleNumber           StyleType = "number-style"
	StyleMoney            StyleType = "money-style"
	StyleDate             StyleType = "date-style"
	StyleTime             StyleType = "time-style"
	StyleFloat            StyleType = "float-style"
	StyleHorizontalLayout StyleType = "horizontal-layout-style"
	StyleVerticalLayout   StyleType = "vertical-layout-style"
	StyleCell             StyleType = "cell-style"
)

// ElementType represents the fundamental structural element type.
type ElementType string

const (
	ElementText             ElementType = "text"
	ElementCaptionText      ElementType = "c-text"
	ElementNumber           ElementType = "number"
	ElementMoney            ElementType = "money"
	ElementDate             ElementType = "date"
	ElementTime             ElementType = "time"
	ElementFloat            ElementType = "float"
	ElementCell             ElementType = "cell"
	ElementGroup            ElementType = "group"
	ElementDetail           ElementType = "detail"
	ElementHorizontalLayout ElementType = "layout-h"
	ElementVerticalLayout   ElementType = "layout-v"
	ElementReport           ElementType = "root"
	ElementSpacer           ElementType = "space"
	ElementImage            ElementType = "image"
	ElementSignatures       ElementType = "signatures"
)
