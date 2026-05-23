package structure

import (
	"encoding/xml"
	"strings"

	"github.com/niarbSoftHub/squash/pkg/interfaces"
	"github.com/niarbSoftHub/squash/pkg/properties"
)

// BaseElement provides common element fields
type BaseElement struct {
	ParentEl interfaces.Parent `xml:"-"`
	StyleRaw string            `xml:"styles,attr,omitempty"` // comma separated styles
	Font     string            `xml:"font,attr,omitempty"`   // element-specific font family
}

func (b *BaseElement) Parent() interfaces.Parent {
	return b.ParentEl
}

func (b *BaseElement) Styles() []string {
	if b.StyleRaw == "" {
		return nil
	}
	// Split by comma or space
	styles := strings.FieldsFunc(b.StyleRaw, func(c rune) bool {
		return c == ',' || c == ' '
	})
	return styles
}

// Text Element
type Text struct {
	XMLName xml.Name `xml:"text"`
	BaseElement
	Value string `xml:",chardata"`
}

func (t *Text) Type() properties.ElementType { return properties.ElementText }

// Number Element — renders an integer value (no decimal point), not monetary.
type Number struct {
	XMLName xml.Name `xml:"number"`
	BaseElement
	Value string `xml:",chardata"`
}

func (n *Number) Type() properties.ElementType { return properties.ElementNumber }

// Money Element — renders a monetary value with configurable decimals and comma formatting.
type Money struct {
	XMLName xml.Name `xml:"money"`
	BaseElement
	Value string `xml:",chardata"`
}

func (m *Money) Type() properties.ElementType { return properties.ElementMoney }

// Date Element — renders a date value with configurable output format.
type Date struct {
	XMLName xml.Name `xml:"date"`
	BaseElement
	Value string `xml:",chardata"`
}

func (d *Date) Type() properties.ElementType { return properties.ElementDate }

// Time Element — renders a time value with configurable output format.
type Time struct {
	XMLName xml.Name `xml:"time"`
	BaseElement
	Value string `xml:",chardata"`
}

func (ti *Time) Type() properties.ElementType { return properties.ElementTime }

// Float Element — renders a decimal number with configurable decimal places.
type Float struct {
	XMLName xml.Name `xml:"float"`
	BaseElement
	Value string `xml:",chardata"`
}

func (f *Float) Type() properties.ElementType { return properties.ElementFloat }

// Spacer Element
type Spacer struct {
	XMLName xml.Name `xml:"space"`
	BaseElement
	Height float64 `xml:"height,attr,omitempty"`
}

func (s *Spacer) Type() properties.ElementType { return properties.ElementSpacer }

// CText Element
type CText struct {
	XMLName xml.Name `xml:"c-text"`
	BaseElement
	Caption   string `xml:"caption,attr,omitempty"`
	Separator string `xml:"separator,attr,omitempty"`
	Value     string `xml:",chardata"`
}

func (ct *CText) Type() properties.ElementType { return properties.ElementCaptionText }

// Image Element
type Image struct {
	XMLName xml.Name `xml:"image"`
	BaseElement
	Width  float64 `xml:"width,attr,omitempty"`
	Height float64 `xml:"height,attr,omitempty"`
	Path   string  `xml:",chardata"`
}

func (i *Image) Type() properties.ElementType { return properties.ElementImage }

// Cell Element
type Cell struct {
	XMLName xml.Name `xml:"cell"`
	BaseElement
	Ratio   int            `xml:"ratio,attr,omitempty"`
	Element ElementWrapper `xml:",any"`
}

func (c *Cell) Type() properties.ElementType { return properties.ElementCell }
func (c *Cell) AddElement(e interfaces.Element) {
	c.Element = ElementWrapper{Element: e}
}

// VerticalLayout
type VerticalLayout struct {
	XMLName xml.Name `xml:"layout-v"`
	BaseElement
	Elements []ElementWrapper `xml:",any"`
}

func (v *VerticalLayout) Type() properties.ElementType { return properties.ElementVerticalLayout }
func (v *VerticalLayout) AddElement(e interfaces.Element) {
	v.Elements = append(v.Elements, ElementWrapper{Element: e})
}

// HorizontalLayout
type HorizontalLayout struct {
	XMLName xml.Name `xml:"layout-h"`
	BaseElement
	Height   float64          `xml:"height,attr,omitempty"` // fixed row height in mm; 0 = auto-grow
	Elements []ElementWrapper `xml:",any"`
}

func (h *HorizontalLayout) Type() properties.ElementType { return properties.ElementHorizontalLayout }
func (h *HorizontalLayout) AddElement(e interfaces.Element) {
	h.Elements = append(h.Elements, ElementWrapper{Element: e})
}

// Group
type Group struct {
	XMLName xml.Name `xml:"group"`
	BaseElement
	PageBreak bool            `xml:"pageBreak,attr,omitempty"`
	Header    *VerticalLayout `xml:"header>layout-v,omitempty"`
	Details   *DetailWrapper  `xml:"details"`
	Footer    *VerticalLayout `xml:"footer>layout-v,omitempty"`
}

func (g *Group) Type() properties.ElementType { return properties.ElementGroup }

// Report
type Report struct {
	XMLName xml.Name `xml:"report"`
	BaseElement
	Orientation  properties.Orientation `xml:"orientation,attr,omitempty"`
	PageSize     properties.PageSize    `xml:"pageSize,attr,omitempty"`
	MarginLeft   float32                `xml:"marginLeft,attr,omitempty"`
	MarginRight  float32                `xml:"marginRight,attr,omitempty"`
	MarginTop    float32                `xml:"marginTop,attr,omitempty"`
	MarginBottom float32                `xml:"marginBottom,attr,omitempty"`

	ReportHeader *VerticalLayout `xml:"reportHeader>layout-v,omitempty"`
	PageHeader   *VerticalLayout `xml:"pageHeader>layout-v,omitempty"`
	Details      *DetailWrapper  `xml:"details"`
	ReportFooter *VerticalLayout `xml:"reportFooter>layout-v,omitempty"`
	PageFooter   *VerticalLayout `xml:"pageFooter>layout-v,omitempty"`
}

func (r *Report) Type() properties.ElementType { return properties.ElementReport }

// Signatures Element
type Signatures struct {
	XMLName xml.Name `xml:"signatures"`
	BaseElement
	Items []string `xml:"signature"`
}

func (s *Signatures) Type() properties.ElementType { return properties.ElementSignatures }
