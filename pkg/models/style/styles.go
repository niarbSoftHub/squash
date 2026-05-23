package style

import (
	"encoding/json"
	"strings"

	"github.com/niarbSoftHub/squash/pkg/properties"
)

// BaseStyle provides the style type for polymorphic deserialization.
type BaseStyle struct {
	StyleType properties.StyleType `json:"styleType"`
}

type StyleSheet struct {
	// Text and layout styles
	TextStyles             map[string]TextStyle             `json:"textStyles,omitempty"`
	CTextStyles            map[string]CTextStyle            `json:"cTextStyles,omitempty"`
	CellStyles             map[string]CellStyle             `json:"cellStyles,omitempty"`
	HorizontalLayoutStyles map[string]HorizontalLayoutStyle `json:"horizontalLayoutStyle,omitempty"`
	VerticalLayoutStyles   map[string]VerticalLayoutStyle   `json:"verticalLayoutStyles,omitempty"`

	// Typed value element styles
	NumberStyles map[string]NumberStyle `json:"numberStyles,omitempty"`
	MoneyStyles  map[string]MoneyStyle  `json:"moneyStyles,omitempty"`
	DateStyles   map[string]DateStyle   `json:"dateStyles,omitempty"`
	TimeStyles   map[string]TimeStyle   `json:"timeStyles,omitempty"`
	FloatStyles  map[string]FloatStyle  `json:"floatStyles,omitempty"`

	// Default styles — applied when an element has no styles attribute declared
	DefaultTextStyle   *TextStyle   `json:"defaultTextStyle,omitempty"`
	DefaultNumberStyle *NumberStyle `json:"defaultNumberStyle,omitempty"`
	DefaultMoneyStyle  *MoneyStyle  `json:"defaultMoneyStyle,omitempty"`
	DefaultDateStyle   *DateStyle   `json:"defaultDateStyle,omitempty"`
	DefaultTimeStyle   *TimeStyle   `json:"defaultTimeStyle,omitempty"`
	DefaultFloatStyle  *FloatStyle  `json:"defaultFloatStyle,omitempty"`
}

type TextStyle struct {
	BaseStyle     `json:"-"`
	Font          properties.Font      `json:"font,omitempty"`
	FontSize      float64              `json:"fontSize,omitempty"`
	Color         string               `json:"color,omitempty"`
	Alignment     properties.Alignment `json:"alignment,omitempty"`
	Bold          bool                 `json:"bold,omitempty"`
	Italic        bool                 `json:"italic,omitempty"`
	Underline     bool                 `json:"underline,omitempty"`
	StrikeThrough bool                 `json:"strikeThrough,omitempty"`
	Superscript   bool                 `json:"superscript,omitempty"`
	Subscript     bool                 `json:"subscript,omitempty"`
	Rotation      float64              `json:"rotation,omitempty"`
}

// NumberStyle — for <number> elements (integers, no decimal point, not monetary).
type NumberStyle struct {
	BaseStyle `json:"-"`
	Font      properties.Font      `json:"font,omitempty"`
	FontSize  float64              `json:"fontSize,omitempty"`
	Color     string               `json:"color,omitempty"`
	Alignment properties.Alignment `json:"alignment,omitempty"`
	Bold      bool                 `json:"bold,omitempty"`
	Italic    bool                 `json:"italic,omitempty"`
	UseComma  bool                 `json:"useComma,omitempty"` // e.g. 1,234
}

// MoneyStyle — for <money> elements with optional decimals, commas, and currency symbol.
type MoneyStyle struct {
	BaseStyle      `json:"-"`
	Font           properties.Font      `json:"font,omitempty"`
	FontSize       float64              `json:"fontSize,omitempty"`
	Color          string               `json:"color,omitempty"`
	Alignment      properties.Alignment `json:"alignment,omitempty"`
	Bold           bool                 `json:"bold,omitempty"`
	Italic         bool                 `json:"italic,omitempty"`
	Decimals       int                  `json:"decimals"`                 // 0 = no decimal, 2 = two decimal places
	UseComma       bool                 `json:"useComma,omitempty"`       // e.g. 1,234.00
	Symbol         string               `json:"symbol,omitempty"`         // e.g. "$", "৳", "€"
	SymbolPosition string               `json:"symbolPosition,omitempty"` // "prefix" or "suffix"
}

// DateStyle — for <date> elements with a configurable output format.
// Input is expected as ISO 8601: "2006-01-02". Output format uses Go-style tokens:
//
//	yyyy=2006, MM=01, dd=02, yy=06, M=1, d=2
type DateStyle struct {
	BaseStyle `json:"-"`
	Font      properties.Font      `json:"font,omitempty"`
	FontSize  float64              `json:"fontSize,omitempty"`
	Color     string               `json:"color,omitempty"`
	Alignment properties.Alignment `json:"alignment,omitempty"`
	Bold      bool                 `json:"bold,omitempty"`
	Italic    bool                 `json:"italic,omitempty"`
	Format    string               `json:"format,omitempty"` // e.g. "dd/MM/yyyy", "MM-dd-yy"
}

// TimeStyle — for <time> elements with a configurable output format.
// Input is expected as "15:04" or "15:04:05". Output format uses Go-style tokens:
//
//	HH=15, mm=04, ss=05
type TimeStyle struct {
	BaseStyle `json:"-"`
	Font      properties.Font      `json:"font,omitempty"`
	FontSize  float64              `json:"fontSize,omitempty"`
	Color     string               `json:"color,omitempty"`
	Alignment properties.Alignment `json:"alignment,omitempty"`
	Bold      bool                 `json:"bold,omitempty"`
	Italic    bool                 `json:"italic,omitempty"`
	Format    string               `json:"format,omitempty"` // e.g. "HH:mm", "HH:mm:ss"
}

// FloatStyle — for <float> elements with configurable decimal precision and comma separation.
type FloatStyle struct {
	BaseStyle `json:"-"`
	Font      properties.Font      `json:"font,omitempty"`
	FontSize  float64              `json:"fontSize,omitempty"`
	Color     string               `json:"color,omitempty"`
	Alignment properties.Alignment `json:"alignment,omitempty"`
	Bold      bool                 `json:"bold,omitempty"`
	Italic    bool                 `json:"italic,omitempty"`
	Decimals  int                  `json:"decimals"`           // number of decimal places
	UseComma  bool                 `json:"useComma,omitempty"` // e.g. 1,234.56
}

type CTextStyle struct {
	BaseStyle            `json:"-"`
	CaptionFont          properties.Font              `json:"captionFont,omitempty"`
	CaptionFontSize      float64                      `json:"captionFontSize,omitempty"`
	CaptionColor         string                       `json:"captionColor,omitempty"`
	CaptionBold          bool                         `json:"captionBold,omitempty"`
	CaptionItalic        bool                         `json:"captionItalic,omitempty"`
	CaptionUnderline     bool                         `json:"captionUnderline,omitempty"`
	CaptionStrikeThrough bool                         `json:"captionStrikeThrough,omitempty"`
	CaptionSuperscript   bool                         `json:"captionSuperscript,omitempty"`
	CaptionSubscript     bool                         `json:"captionSubscript,omitempty"`
	Separator            string                       `json:"separator,omitempty"`
	SeparatorBehavior    properties.SeparatorBehavior `json:"separatorBehavior,omitempty"`
	ReverseOrder         bool                         `json:"reverseOrder,omitempty"`
	ValueFont            properties.Font              `json:"valueFont,omitempty"`
	ValueFontSize        float64                      `json:"valueFontSize,omitempty"`
	ValueColor           string                       `json:"valueColor,omitempty"`
	ValueBold            bool                         `json:"valueBold,omitempty"`
	ValueItalic          bool                         `json:"valueItalic,omitempty"`
	ValueUnderline       bool                         `json:"valueUnderline,omitempty"`
	ValueStrikeThrough   bool                         `json:"valueStrikeThrough,omitempty"`
	ValueSuperscript     bool                         `json:"valueSuperscript,omitempty"`
	ValueSubscript       bool                         `json:"valueSubscript,omitempty"`
	Rotation             float64                      `json:"rotation,omitempty"`
	Alignment            properties.Alignment         `json:"alignment,omitempty"`
}

// Border properties shared across Cell and Layouts
type BorderProps struct {
	BorderLeft        properties.BorderStyle `json:"borderLeft,omitempty"`
	BorderLeftWidth   float64                `json:"borderLeftWidth,omitempty"`
	BorderLeftColor   string                 `json:"borderLeftColor,omitempty"`
	BorderRight       properties.BorderStyle `json:"borderRight,omitempty"`
	BorderRightWidth  float64                `json:"borderRightWidth,omitempty"`
	BorderRightColor  string                 `json:"borderRightColor,omitempty"`
	BorderTop         properties.BorderStyle `json:"borderTop,omitempty"`
	BorderTopWidth    float64                `json:"borderTopWidth,omitempty"`
	BorderTopColor    string                 `json:"borderTopColor,omitempty"`
	BorderBottom      properties.BorderStyle `json:"borderBottom,omitempty"`
	BorderBottomWidth float64                `json:"borderBottomWidth,omitempty"`
	BorderBottomColor string                 `json:"borderBottomColor,omitempty"`
	BackgroundColor   string                 `json:"backgroundColor,omitempty"`
	Rotation          float64                `json:"rotation,omitempty"`
}

type CellStyle struct {
	BaseStyle     `json:"-"`
	BorderProps   `json:",inline"`
	PaddingLeft   float64               `json:"paddingLeft,omitempty"`
	PaddingRight  float64               `json:"paddingRight,omitempty"`
	PaddingTop    float64               `json:"paddingTop,omitempty"`
	PaddingBottom float64               `json:"paddingBottom,omitempty"`
	VAlignment    properties.VAlignment `json:"vAlignment,omitempty"`
}

type HorizontalLayoutStyle struct {
	BaseStyle   `json:"-"`
	BorderProps `json:",inline"`
}

type VerticalLayoutStyle struct {
	BaseStyle   `json:"-"`
	BorderProps `json:",inline"`
}

// --- MarshalJSON implementations ---

func (s TextStyle) MarshalJSON() ([]byte, error) {
	type Alias TextStyle
	return json.Marshal(&struct {
		StyleType properties.StyleType `json:"styleType"`
		Alias
	}{
		StyleType: properties.StyleText,
		Alias:     (Alias)(s),
	})
}

func (s NumberStyle) MarshalJSON() ([]byte, error) {
	type Alias NumberStyle
	return json.Marshal(&struct {
		StyleType properties.StyleType `json:"styleType"`
		Alias
	}{
		StyleType: properties.StyleNumber,
		Alias:     (Alias)(s),
	})
}

func (s MoneyStyle) MarshalJSON() ([]byte, error) {
	type Alias MoneyStyle
	return json.Marshal(&struct {
		StyleType properties.StyleType `json:"styleType"`
		Alias
	}{
		StyleType: properties.StyleMoney,
		Alias:     (Alias)(s),
	})
}

func (s DateStyle) MarshalJSON() ([]byte, error) {
	type Alias DateStyle
	return json.Marshal(&struct {
		StyleType properties.StyleType `json:"styleType"`
		Alias
	}{
		StyleType: properties.StyleDate,
		Alias:     (Alias)(s),
	})
}

func (s TimeStyle) MarshalJSON() ([]byte, error) {
	type Alias TimeStyle
	return json.Marshal(&struct {
		StyleType properties.StyleType `json:"styleType"`
		Alias
	}{
		StyleType: properties.StyleTime,
		Alias:     (Alias)(s),
	})
}

func (s FloatStyle) MarshalJSON() ([]byte, error) {
	type Alias FloatStyle
	return json.Marshal(&struct {
		StyleType properties.StyleType `json:"styleType"`
		Alias
	}{
		StyleType: properties.StyleFloat,
		Alias:     (Alias)(s),
	})
}

func (s CTextStyle) MarshalJSON() ([]byte, error) {
	type Alias CTextStyle
	return json.Marshal(&struct {
		StyleType properties.StyleType `json:"styleType"`
		Alias
	}{
		StyleType: properties.StyleCText,
		Alias:     (Alias)(s),
	})
}

func (s CellStyle) MarshalJSON() ([]byte, error) {
	type Alias CellStyle
	return json.Marshal(&struct {
		StyleType properties.StyleType `json:"styleType"`
		Alias
	}{
		StyleType: properties.StyleCell,
		Alias:     (Alias)(s),
	})
}

func (s HorizontalLayoutStyle) MarshalJSON() ([]byte, error) {
	type Alias HorizontalLayoutStyle
	return json.Marshal(&struct {
		StyleType properties.StyleType `json:"styleType"`
		Alias
	}{
		StyleType: properties.StyleHorizontalLayout,
		Alias:     (Alias)(s),
	})
}

func (s VerticalLayoutStyle) MarshalJSON() ([]byte, error) {
	type Alias VerticalLayoutStyle
	return json.Marshal(&struct {
		StyleType properties.StyleType `json:"styleType"`
		Alias
	}{
		StyleType: properties.StyleVerticalLayout,
		Alias:     (Alias)(s),
	})
}

// UnmarshalJSON handles smart routing based on styleType and naming patterns.
// Automatically routes styles to appropriate maps and defaults.
// Supports two patterns:
//  1. "default-XXX-style" → sets Default*Style field
//  2. Any other root-level style with styleType → adds to appropriate typed map
func (ss *StyleSheet) UnmarshalJSON(data []byte) error {
	type Alias StyleSheet
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(ss),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Capture root-level styles that aren't in typed maps
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err == nil {
		// Process root-level styles with styleType
		for key, val := range rawData {
			// Skip known top-level sections
			if key == "textStyles" || key == "cTextStyles" || key == "cellStyles" ||
				key == "horizontalLayoutStyle" || key == "verticalLayoutStyles" ||
				key == "numberStyles" || key == "moneyStyles" || key == "dateStyles" ||
				key == "timeStyles" || key == "floatStyles" {
				continue
			}

			// This is a potential root-level style, check if it has styleType
			styleMap, ok := val.(map[string]interface{})
			if !ok {
				continue
			}

			styleTypeVal, hasType := styleMap["styleType"]
			if !hasType {
				continue
			}

			styleTypeStr, ok := styleTypeVal.(string)
			if !ok {
				continue
			}

			styleJSON, _ := json.Marshal(styleMap)

			// Check if this is a "default-XXX-style" pattern
			isDefaultStyle := strings.HasPrefix(key, "default-") && strings.HasSuffix(key, "-style")

			switch properties.StyleType(styleTypeStr) {
			case properties.StyleNumber:
				var ns NumberStyle
				if err := json.Unmarshal(styleJSON, &ns); err == nil {
					if isDefaultStyle {
						ss.DefaultNumberStyle = &ns
					} else {
						if ss.NumberStyles == nil {
							ss.NumberStyles = make(map[string]NumberStyle)
						}
						ss.NumberStyles[key] = ns
					}
				}
			case properties.StyleMoney:
				var ms MoneyStyle
				if err := json.Unmarshal(styleJSON, &ms); err == nil {
					if isDefaultStyle {
						ss.DefaultMoneyStyle = &ms
					} else {
						if ss.MoneyStyles == nil {
							ss.MoneyStyles = make(map[string]MoneyStyle)
						}
						ss.MoneyStyles[key] = ms
					}
				}
			case properties.StyleDate:
				var ds DateStyle
				if err := json.Unmarshal(styleJSON, &ds); err == nil {
					if isDefaultStyle {
						ss.DefaultDateStyle = &ds
					} else {
						if ss.DateStyles == nil {
							ss.DateStyles = make(map[string]DateStyle)
						}
						ss.DateStyles[key] = ds
					}
				}
			case properties.StyleTime:
				var ts TimeStyle
				if err := json.Unmarshal(styleJSON, &ts); err == nil {
					if isDefaultStyle {
						ss.DefaultTimeStyle = &ts
					} else {
						if ss.TimeStyles == nil {
							ss.TimeStyles = make(map[string]TimeStyle)
						}
						ss.TimeStyles[key] = ts
					}
				}
			case properties.StyleFloat:
				var fs FloatStyle
				if err := json.Unmarshal(styleJSON, &fs); err == nil {
					if isDefaultStyle {
						ss.DefaultFloatStyle = &fs
					} else {
						if ss.FloatStyles == nil {
							ss.FloatStyles = make(map[string]FloatStyle)
						}
						ss.FloatStyles[key] = fs
					}
				}
			case properties.StyleText:
				var ts TextStyle
				if err := json.Unmarshal(styleJSON, &ts); err == nil {
					if isDefaultStyle {
						ss.DefaultTextStyle = &ts
					} else {
						if ss.TextStyles == nil {
							ss.TextStyles = make(map[string]TextStyle)
						}
						ss.TextStyles[key] = ts
					}
				}
			}
		}
	}

	return nil
}
