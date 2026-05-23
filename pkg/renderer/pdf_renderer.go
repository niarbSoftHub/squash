package renderer

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	pngenc "image/png"
	"math"
	"os"
	"strconv"
	"strings"
	gotime "time"

	"github.com/phpdave11/gofpdf"
	"golang.org/x/image/webp"

	"github.com/niarbSoftHub/squash/pkg/interfaces"
	"github.com/niarbSoftHub/squash/pkg/models/structure"
	"github.com/niarbSoftHub/squash/pkg/models/style"
	"github.com/niarbSoftHub/squash/pkg/properties"
)

// GeneratedDocument is the custom wrapper returned to keep public API identical.
type GeneratedDocument struct {
	pdf *gofpdf.Fpdf
}

var pageSizeMap = map[properties.PageSize]string{
	properties.PageSizeA0:      "A0",
	properties.PageSizeA1:      "A1",
	properties.PageSizeA2:      "A2",
	properties.PageSizeA3:      "A3",
	properties.PageSizeA4:      "A4",
	properties.PageSizeA5:      "A5",
	properties.PageSizeA6:      "A6",
	properties.PageSizeA7:      "A7",
	properties.PageSizeA8:      "A8",
	properties.PageSizeA9:      "A9",
	properties.PageSizeA10:     "A10",
	properties.PageSizeLetter:  "Letter",
	properties.PageSizeLegal:   "Legal",
	properties.PageSizeTabloId: "Tabloid",
	properties.PageSizeC0:      "C0",
	properties.PageSizeC1:      "C1",
	properties.PageSizeC2:      "C2",
	properties.PageSizeC3:      "C3",
	properties.PageSizeC4:      "C4",
	properties.PageSizeC5:      "C5",
	properties.PageSizeC6:      "C6",
	properties.PageSizeC7:      "C7",
	properties.PageSizeC8:      "C8",
	properties.PageSizeC9:      "C9",
	properties.PageSizeC10:     "C10",
	properties.PageSizeB0:      "B0",
	properties.PageSizeB1:      "B1",
	properties.PageSizeB2:      "B2",
	properties.PageSizeB3:      "B3",
	properties.PageSizeB4:      "B4",
	properties.PageSizeB5:      "B5",
	properties.PageSizeB6:      "B6",
	properties.PageSizeB7:      "B7",
	properties.PageSizeB8:      "B8",
	properties.PageSizeB9:      "B9",
	properties.PageSizeB10:     "B10",
	properties.PageSizeDL:      "DL",
}

func pageToGofpdfSize(pageSize properties.PageSize) string {
	if s, ok := pageSizeMap[pageSize]; ok {
		return s
	}
	return "A4" // Default to A4
}

// Generate implements the document production step, returning the saveable result.
func (gd *GeneratedDocument) Generate() (*DocumentResult, error) {
	return &DocumentResult{pdf: gd.pdf}, nil
}

// DocumentResult wraps the completed PDF document and provides the Save method.
type DocumentResult struct {
	pdf *gofpdf.Fpdf
}

// Save writes the rendered PDF file to the specified disk path.
func (dr *DocumentResult) Save(path string) error {
	return dr.pdf.OutputFileAndClose(path)
}

// renderContext encapsulates the layout state and options for rendering.
type renderContext struct {
	pdf                  *gofpdf.Fpdf
	sheet                *style.StyleSheet
	pageHeight           float64
	pageWidth            float64
	bottomMargin         float64
	originalBottomMargin float64
	leftMargin           float64
	topMargin            float64
	rightMargin          float64
	contentWidth         float64
	activeGroupHeader    *structure.VerticalLayout
	report               *structure.Report
	registeredFonts      map[string]bool
}

// addPage adds a new page, prints headers, and repeats active group header if present.
func (ctx *renderContext) addPage() {
	if ctx.pdf.PageNo() > 0 {
		ctx.drawFooter()
	}
	ctx.pdf.AddPage()
	if ctx.report.PageHeader != nil {
		ctx.pdf.SetY(ctx.topMargin)
		ctx.renderElement(ctx.report.PageHeader, ctx.leftMargin, ctx.contentWidth, false)
	}
	if ctx.activeGroupHeader != nil {
		ctx.renderElement(ctx.activeGroupHeader, ctx.leftMargin, ctx.contentWidth, false)
	}
}

// drawFooter prints the page footer and the page number at the bottom.
func (ctx *renderContext) drawFooter() {
	if ctx.report.PageFooter != nil {
		curY := ctx.pdf.GetY()
		footerHeight := ctx.measureElement(ctx.report.PageFooter, ctx.contentWidth)
		footerY := ctx.pageHeight - ctx.originalBottomMargin - footerHeight
		ctx.pdf.SetY(footerY)
		ctx.renderElement(ctx.report.PageFooter, ctx.leftMargin, ctx.contentWidth, false)
		ctx.pdf.SetY(curY)
	}

	// Draw page number in footer (bottom right)
	ctx.pdf.SetFont("Arial", "I", 8.0)
	ctx.pdf.SetTextColor(93, 109, 126) // Gray #5D6D7E
	pageStr := fmt.Sprintf("%d / {nb}", ctx.pdf.PageNo())
	ctx.pdf.SetXY(ctx.pageWidth-ctx.rightMargin-30, ctx.pageHeight-10.0)
	ctx.pdf.CellFormat(30, 4, pageStr, "", 0, "R", false, 0, "")
}

// Generate builds the complete document layout using a custom canvas-based engine.
func Generate(report *structure.Report, styleSheet *style.StyleSheet) (*GeneratedDocument, error) {
	orientation := "P"
	if report.Orientation == properties.OrientationLandscape {
		orientation = "L"
	}

	// Get page size from report, default to A4 if not specified
	pageSize := "A4"
	if report.PageSize != 0 {
		// A custom page size was specified (including A0)
		pageSize = pageToGofpdfSize(report.PageSize)
	}

	// Initialize gofpdf Canvas
	pdf := gofpdf.New(orientation, "mm", pageSize, "")
	pdf.AliasNbPages("{nb}") // Register alias for total pages count placeholder

	// Configure page margins
	pdf.SetMargins(float64(report.MarginLeft), float64(report.MarginTop), float64(report.MarginRight))

	// Turn off internal auto page break; we calculate and trigger breaks perfectly ourselves
	pdf.SetAutoPageBreak(false, float64(report.MarginBottom))

	pageWidth, pageHeight := pdf.GetPageSize()
	leftMargin, topMargin, rightMargin, _ := pdf.GetMargins()
	bottomMargin := float64(report.MarginBottom)
	if bottomMargin <= 0 {
		bottomMargin = 10.0
	}
	contentWidth := pageWidth - leftMargin - rightMargin

	ctx := &renderContext{
		pdf:                  pdf,
		sheet:                styleSheet,
		pageHeight:           pageHeight,
		pageWidth:            pageWidth,
		bottomMargin:         bottomMargin,
		originalBottomMargin: bottomMargin,
		leftMargin:           leftMargin,
		topMargin:            topMargin,
		rightMargin:          rightMargin,
		contentWidth:         contentWidth,
		report:               report,
		registeredFonts:      make(map[string]bool),
	}

	// Adjust bottom margin if page footer is present to prevent content/footer overlapping
	if report.PageFooter != nil {
		footerHeight := ctx.measureElement(report.PageFooter, contentWidth)
		ctx.bottomMargin = bottomMargin + footerHeight
	}

	// Trigger the first page addition
	ctx.addPage()

	// 1. Render ReportHeader
	if report.ReportHeader != nil {
		h := ctx.measureElement(report.ReportHeader, contentWidth)
		if pdf.GetY()+h > pageHeight-bottomMargin {
			ctx.addPage()
		}
		ctx.renderElement(report.ReportHeader, leftMargin, contentWidth, true)
	}

	// 2. Render Details & Groups
	if report.Details != nil && report.Details.Detail != nil {
		if vl, ok := report.Details.Detail.(*structure.VerticalLayout); ok {
			for i, child := range vl.Elements {
				if g, ok := child.Element.(*structure.Group); ok {
					if g.PageBreak && i > 0 {
						ctx.addPage()
					}
					ctx.renderElement(g, leftMargin, contentWidth, true)
				} else {
					isContainer := child.Element.Type() == properties.ElementGroup || child.Element.Type() == properties.ElementVerticalLayout
					ch := ctx.measureElement(child.Element, contentWidth)
					if !isContainer {
						if pdf.GetY()+ch > pageHeight-bottomMargin {
							ctx.addPage()
						}
					}
					ctx.renderElement(child.Element, leftMargin, contentWidth, true)
				}
			}
		} else {
			ctx.renderElement(report.Details.Detail, leftMargin, contentWidth, true)
		}
	}

	// 3. Render ReportFooter
	if report.ReportFooter != nil {
		h := ctx.measureElement(report.ReportFooter, contentWidth)
		if pdf.GetY()+h > pageHeight-bottomMargin {
			ctx.addPage()
		}
		ctx.renderElement(report.ReportFooter, leftMargin, contentWidth, true)
	}

	// Final page footer print
	ctx.drawFooter()

	return &GeneratedDocument{pdf: pdf}, nil
}

// measureElement recursively calculates the exact height needed by an element under a given width.
func (ctx *renderContext) measureHorizontalLayout(v *structure.HorizontalLayout, width float64) float64 {
	if v.Height > 0 {
		return v.Height
	}
	var totalRatio int
	cells := make([]*structure.Cell, 0)
	for _, ew := range v.Elements {
		if c, ok := ew.Element.(*structure.Cell); ok {
			cells = append(cells, c)
			ratio := c.Ratio
			if ratio == 0 {
				ratio = 1
			}
			totalRatio += ratio
		}
	}
	if len(cells) == 0 {
		return 0.0
	}
	maxH := 0.0
	for _, c := range cells {
		ratio := c.Ratio
		if ratio == 0 {
			ratio = 1
		}
		cellW := (float64(ratio) / float64(totalRatio)) * width
		cellH := ctx.measureElement(c, cellW)
		if cellH > maxH {
			maxH = cellH
		}
	}
	return maxH
}

func (ctx *renderContext) measureCell(v *structure.Cell, width float64) float64 {
	padL, padR := 1.0, 1.0
	padT, padB := 1.0, 1.0
	for _, sName := range v.Styles() {
		if cs, ok := ctx.sheet.CellStyles[sName]; ok {
			if cs.PaddingLeft > 0 || cs.PaddingRight > 0 {
				padL = cs.PaddingLeft
				padR = cs.PaddingRight
			}
			if cs.PaddingTop > 0 || cs.PaddingBottom > 0 {
				padT = cs.PaddingTop
				padB = cs.PaddingBottom
			}
		}
	}
	innerWidth := width - padL - padR
	if innerWidth <= 0 {
		innerWidth = width
	}
	innerHeight := ctx.measureElement(v.Element.Element, innerWidth)
	return innerHeight + padT + padB
}

func (ctx *renderContext) measureElement(el interfaces.Element, width float64) float64 {
	if el == nil {
		return 0.0
	}
	switch v := el.(type) {
	case *structure.Spacer:
		h := v.Height
		if h <= 0 {
			h = 10.0
		}
		return h
	case *structure.VerticalLayout:
		totalHeight := 0.0
		for _, child := range v.Elements {
			totalHeight += ctx.measureElement(child.Element, width)
		}
		return totalHeight
	case *structure.HorizontalLayout:
		return ctx.measureHorizontalLayout(v, width)
	case *structure.Cell:
		return ctx.measureCell(v, width)
	case *structure.Group:
		h := 0.0
		if v.Header != nil {
			h += ctx.measureElement(v.Header, width)
		}
		if v.Details != nil && v.Details.Detail != nil {
			h += ctx.measureElement(v.Details.Detail, width)
		}
		if v.Footer != nil {
			h += ctx.measureElement(v.Footer, width)
		}
		return h
	case *structure.Text:
		family, styleFlags, fontSize, _ := ctx.resolveTextFontProps(v)
		ctx.pdf.SetFont(family, styleFlags, fontSize)
		lines := strings.Split(v.Value, "\n")
		totalLinesCount := 0
		for _, line := range lines {
			splitLines := ctx.pdf.SplitText(line, width)
			totalLinesCount += len(splitLines)
		}
		lineHeight := fontSize*0.3527 + 1.2
		return float64(totalLinesCount) * lineHeight
	case *structure.Number:
		s := ctx.resolveNumberStyle(v)
		ctx.pdf.SetFont(ctx.resolveFont2(string(s.Font), s.Bold, s.Italic))
		lineHeight := s.FontSize*0.3527 + 1.2
		return lineHeight
	case *structure.Money:
		s := ctx.resolveMoneyStyle(v)
		ctx.pdf.SetFont(ctx.resolveFont2(string(s.Font), s.Bold, s.Italic))
		lineHeight := s.FontSize*0.3527 + 1.2
		return lineHeight
	case *structure.Date:
		s := ctx.resolveDateStyle(v)
		ctx.pdf.SetFont(ctx.resolveFont2(string(s.Font), s.Bold, s.Italic))
		lineHeight := s.FontSize*0.3527 + 1.2
		return lineHeight
	case *structure.Time:
		s := ctx.resolveTimeStyle(v)
		ctx.pdf.SetFont(ctx.resolveFont2(string(s.Font), s.Bold, s.Italic))
		lineHeight := s.FontSize*0.3527 + 1.2
		return lineHeight
	case *structure.Float:
		s := ctx.resolveFloatStyle(v)
		ctx.pdf.SetFont(ctx.resolveFont2(string(s.Font), s.Bold, s.Italic))
		lineHeight := s.FontSize*0.3527 + 1.2
		return lineHeight
	case *structure.CText:
		family, styleFlags, fontSize, _ := ctx.resolveCTextFontProps(v)
		ctx.pdf.SetFont(family, styleFlags, fontSize)
		fullText := v.Caption + v.Separator + v.Value
		lines := ctx.pdf.SplitText(fullText, width)
		lineHeight := fontSize*0.3527 + 1.2
		return float64(len(lines)) * lineHeight
	case *structure.Image:
		if v.Height > 0 {
			return v.Height
		}
		if v.Path == "" {
			return 15.0
		}
		wPix, hPix, err := getImageDimensions(v.Path)
		if err != nil || wPix == 0 {
			return 15.0
		}
		aspectRatio := float64(hPix) / float64(wPix)
		return width * aspectRatio
	case *structure.Signatures:
		return 15.0
	}
	return 0.0
}

// renderElement recursively renders each XML element at absolute coordinates.
func (ctx *renderContext) renderSimpleElement(el interfaces.Element, x, width float64, allowPageBreak bool, render func(float64, float64, float64, float64)) {
	h := ctx.measureElement(el, width)
	if allowPageBreak && ctx.pdf.GetY()+h > ctx.pageHeight-ctx.bottomMargin {
		ctx.addPage()
	}
	startY := ctx.pdf.GetY()
	render(x, startY, width, h)
	ctx.pdf.SetXY(x, startY+h)
}

func (ctx *renderContext) renderHorizontalLayout(v *structure.HorizontalLayout, x, width float64, allowPageBreak bool) {
	var h float64
	if v.Height > 0 {
		h = v.Height
	} else {
		h = ctx.measureElement(v, width)
	}
	if allowPageBreak && ctx.pdf.GetY()+h > ctx.pageHeight-ctx.bottomMargin {
		ctx.addPage()
	}
	var totalRatio int
	cells := make([]*structure.Cell, 0)
	for _, ew := range v.Elements {
		if c, ok := ew.Element.(*structure.Cell); ok {
			cells = append(cells, c)
			ratio := c.Ratio
			if ratio == 0 {
				ratio = 1
			}
			totalRatio += ratio
		}
	}
	if len(cells) == 0 {
		return
	}
	startY := ctx.pdf.GetY()
	currentX := x
	for _, c := range cells {
		ratio := c.Ratio
		if ratio == 0 {
			ratio = 1
		}
		cellW := (float64(ratio) / float64(totalRatio)) * width
		ctx.pdf.SetXY(currentX, startY)
		ctx.renderCell(c, currentX, startY, cellW, h, false)
		currentX += cellW
	}
	ctx.pdf.SetXY(x, startY+h)
}

func (ctx *renderContext) renderElement(el interfaces.Element, x, width float64, allowPageBreak bool) {
	if el == nil {
		return
	}
	switch v := el.(type) {
	case *structure.Spacer:
		h := v.Height
		if h <= 0 {
			h = 10.0
		}
		if allowPageBreak && ctx.pdf.GetY()+h > ctx.pageHeight-ctx.bottomMargin {
			ctx.addPage()
		}
		ctx.pdf.SetY(ctx.pdf.GetY() + h)
	case *structure.VerticalLayout:
		for _, child := range v.Elements {
			isContainer := child.Element.Type() == properties.ElementGroup || child.Element.Type() == properties.ElementVerticalLayout
			ch := ctx.measureElement(child.Element, width)
			if !isContainer {
				if allowPageBreak && ctx.pdf.GetY()+ch > ctx.pageHeight-ctx.bottomMargin {
					ctx.addPage()
				}
			}
			ctx.pdf.SetX(x)
			ctx.renderElement(child.Element, x, width, allowPageBreak)
		}
	case *structure.HorizontalLayout:
		ctx.renderHorizontalLayout(v, x, width, allowPageBreak)
	case *structure.Cell:
		h := ctx.measureElement(v, width)
		if allowPageBreak && ctx.pdf.GetY()+h > ctx.pageHeight-ctx.bottomMargin {
			ctx.addPage()
		}
		startY := ctx.pdf.GetY()
		ctx.renderCell(v, x, startY, width, h, allowPageBreak)
		ctx.pdf.SetXY(x, startY+h)
	case *structure.Group:
		if v.Header != nil {
			ch := ctx.measureElement(v.Header, width)
			if allowPageBreak && ctx.pdf.GetY()+ch > ctx.pageHeight-ctx.bottomMargin {
				ctx.addPage()
			}
			ctx.renderElement(v.Header, x, width, allowPageBreak)
		}

		oldHeader := ctx.activeGroupHeader
		ctx.activeGroupHeader = v.Header

		if v.Details != nil && v.Details.Detail != nil {
			ctx.renderElement(v.Details.Detail, x, width, allowPageBreak)
		}

		ctx.activeGroupHeader = oldHeader

		if v.Footer != nil {
			ch := ctx.measureElement(v.Footer, width)
			if allowPageBreak && ctx.pdf.GetY()+ch > ctx.pageHeight-ctx.bottomMargin {
				ctx.addPage()
			}
			ctx.renderElement(v.Footer, x, width, allowPageBreak)
		}
	case *structure.Text:
		ctx.renderSimpleElement(v, x, width, allowPageBreak, func(px, py, pw, ph float64) {
			ctx.renderText(v, px, py, pw, ph)
		})
	case *structure.CText:
		ctx.renderSimpleElement(v, x, width, allowPageBreak, func(px, py, pw, ph float64) {
			ctx.renderCText(v, px, py, pw, ph)
		})
	case *structure.Number:
		ctx.renderSimpleElement(v, x, width, allowPageBreak, func(px, py, pw, ph float64) {
			ctx.renderNumber(v, px, py, pw, ph)
		})
	case *structure.Money:
		ctx.renderSimpleElement(v, x, width, allowPageBreak, func(px, py, pw, ph float64) {
			ctx.renderMoney(v, px, py, pw, ph)
		})
	case *structure.Date:
		ctx.renderSimpleElement(v, x, width, allowPageBreak, func(px, py, pw, ph float64) {
			ctx.renderDate(v, px, py, pw, ph)
		})
	case *structure.Time:
		ctx.renderSimpleElement(v, x, width, allowPageBreak, func(px, py, pw, ph float64) {
			ctx.renderTime(v, px, py, pw, ph)
		})
	case *structure.Float:
		ctx.renderSimpleElement(v, x, width, allowPageBreak, func(px, py, pw, ph float64) {
			ctx.renderFloat(v, px, py, pw, ph)
		})
	case *structure.Image:
		ctx.renderSimpleElement(v, x, width, allowPageBreak, func(px, py, pw, ph float64) {
			ctx.renderImage(v, px, py, pw, ph)
		})
	case *structure.Signatures:
		ctx.renderSimpleElement(v, x, width, allowPageBreak, func(px, py, pw, ph float64) {
			ctx.renderSignatures(v, px, py, pw)
		})
	}
}

// renderCell renders background, vertical content alignment, and boundaries of a cell.
func (ctx *renderContext) renderCell(cell *structure.Cell, x, y, w, h float64, allowPageBreak bool) {
	var borderTop, borderBottom, borderLeft, borderRight properties.BorderStyle
	var borderTopColor, borderBottomColor, borderLeftColor, borderRightColor string
	var bgColor string
	var vAlign properties.VAlignment = properties.VAlignmentCenter // default: center
	padL, padR := 0.5, 0.5
	padT := 0.5
	padB := 0.5

	for _, sName := range cell.Styles() {
		if cs, ok := ctx.sheet.CellStyles[sName]; ok {
			borderTop = cs.BorderTop
			borderBottom = cs.BorderBottom
			borderLeft = cs.BorderLeft
			borderRight = cs.BorderRight
			borderTopColor = cs.BorderTopColor
			borderBottomColor = cs.BorderBottomColor
			borderLeftColor = cs.BorderLeftColor
			borderRightColor = cs.BorderRightColor
			bgColor = cs.BackgroundColor
			if cs.VAlignment != "" {
				vAlign = cs.VAlignment
			}
			if cs.PaddingLeft > 0 || cs.PaddingRight > 0 {
				padL = cs.PaddingLeft
				padR = cs.PaddingRight
			}
			if cs.PaddingTop > 0 {
				padT = cs.PaddingTop
			}
			if cs.PaddingBottom > 0 {
				padB = cs.PaddingBottom
			}
		}
	}

	// 1. Background Fill
	if bgColor != "" {
		r, g, b := parseHexColorRGB(bgColor)
		ctx.pdf.SetFillColor(r, g, b)
		ctx.pdf.Rect(x, y, w, h, "F")
	}

	// 2. Child Element Rendering
	if cell.Element.Element != nil {
		availableHeight := h - padT - padB
		contentHeight := ctx.measureElement(cell.Element.Element, w-padL-padR)

		// Calculate Y position based on vertical alignment
		var offsetY float64 = padT
		if vAlign == properties.VAlignmentCenter {
			if contentHeight < availableHeight {
				offsetY = padT + (availableHeight-contentHeight)/2
			}
		} else if vAlign == properties.VAlignmentBottom {
			if contentHeight < availableHeight {
				offsetY = padT + (availableHeight - contentHeight)
			}
		}
		// For VAlignmentTop, use padT (default)

		// If the child is a VerticalLayout, temporarily adjust the last HorizontalLayout's height
		// to distribute space properly, but restore it after rendering to avoid modifying the
		// original structure (which may be reused on subsequent pages, e.g., group headers)
		var lastLayout *structure.HorizontalLayout
		var originalHeight float64
		if vl, ok := cell.Element.Element.(*structure.VerticalLayout); ok && len(vl.Elements) > 0 {
			// Identify all child HorizontalLayout elements
			var hLayouts []*structure.HorizontalLayout
			for _, child := range vl.Elements {
				if hl, ok := child.Element.(*structure.HorizontalLayout); ok {
					hLayouts = append(hLayouts, hl)
				}
			}

			// Calculate the height for the last HorizontalLayout to absorb remaining space
			if len(hLayouts) > 1 {
				otherHeights := 0.0
				for i := 0; i < len(hLayouts)-1; i++ {
					otherHeights += ctx.measureElement(hLayouts[i], w-padL-padR)
				}
				calculatedLastHeight := availableHeight - otherHeights
				if calculatedLastHeight > 0 {
					lastLayout = hLayouts[len(hLayouts)-1]
					originalHeight = lastLayout.Height
					lastLayout.Height = calculatedLastHeight
				}
			}
		}

		ctx.pdf.SetXY(x+padL, y+offsetY)
		ctx.renderElement(cell.Element.Element, x+padL, w-padL-padR, allowPageBreak)

		// Restore the original height if we modified it
		if lastLayout != nil {
			lastLayout.Height = originalHeight
		}
	}

	// 3. Borders Drawing
	ctx.pdf.SetLineWidth(0.2) // Standard thin border line

	if borderTop == properties.BorderSolid || borderTop != "" {
		r, g, b := parseHexColorRGB(borderTopColor)
		ctx.pdf.SetDrawColor(r, g, b)
		ctx.pdf.Line(x, y, x+w, y)
	}
	if borderBottom == properties.BorderSolid || borderBottom != "" {
		r, g, b := parseHexColorRGB(borderBottomColor)
		ctx.pdf.SetDrawColor(r, g, b)
		ctx.pdf.Line(x, y+h, x+w, y+h)
	}
	if borderLeft == properties.BorderSolid || borderLeft != "" {
		r, g, b := parseHexColorRGB(borderLeftColor)
		ctx.pdf.SetDrawColor(r, g, b)
		ctx.pdf.Line(x, y, x, y+h)
	}
	if borderRight == properties.BorderSolid || borderRight != "" {
		r, g, b := parseHexColorRGB(borderRightColor)
		ctx.pdf.SetDrawColor(r, g, b)
		ctx.pdf.Line(x+w, y, x+w, y+h)
	}
}

// renderText draws wrapped block text formatted to style properties.
func (ctx *renderContext) renderText(textObj *structure.Text, x, y, w, h float64) {
	family, styleFlags, fontSize, fontColor := ctx.resolveTextFontProps(textObj)

	ctx.pdf.SetFont(family, styleFlags, fontSize)
	r, g, b := parseHexColorRGB(fontColor)
	ctx.pdf.SetTextColor(r, g, b)

	align := properties.AlignmentLeft
	for _, sName := range textObj.Styles() {
		if ts, ok := ctx.sheet.TextStyles[sName]; ok {
			if ts.Alignment != "" {
				align = ts.Alignment
			}
		}
	}

	isFloat := isFloatValue(textObj.Value)
	if isFloat && align == properties.AlignmentLeft {
		align = properties.AlignmentRight
	}

	alignStr := "L"
	if align == properties.AlignmentCenter {
		alignStr = "C"
	} else if align == properties.AlignmentRight {
		alignStr = "R"
	}

	lineHeight := fontSize*0.3527 + 1.2

	lines := strings.Split(textObj.Value, "\n")
	currentY := y
	for _, line := range lines {
		fmtLine := formatIfFloat(line)
		splitLines := ctx.pdf.SplitText(fmtLine, w)
		for _, sl := range splitLines {
			ctx.pdf.SetXY(x, currentY)
			ctx.pdf.CellFormat(w, lineHeight, sl, "", 0, alignStr, false, 0, "")
			currentY += lineHeight
		}
	}
}

// renderCText draws captioned text elements.
func (ctx *renderContext) renderCText(ctextObj *structure.CText, x, y, w, h float64) {
	family, styleFlags, fontSize, fontColor := ctx.resolveCTextFontProps(ctextObj)

	ctx.pdf.SetFont(family, styleFlags, fontSize)
	r, g, b := parseHexColorRGB(fontColor)
	ctx.pdf.SetTextColor(r, g, b)

	align := properties.AlignmentLeft
	for _, sName := range ctextObj.Styles() {
		if cts, ok := ctx.sheet.CTextStyles[sName]; ok {
			if cts.Alignment != "" {
				align = cts.Alignment
			}
		} else if ts, ok := ctx.sheet.TextStyles[sName]; ok {
			if ts.Alignment != "" {
				align = ts.Alignment
			}
		}
	}

	formattedVal := formatIfFloat(ctextObj.Value)
	fullText := ctextObj.Caption + ctextObj.Separator + formattedVal

	alignStr := "L"
	if align == properties.AlignmentCenter {
		alignStr = "C"
	} else if align == properties.AlignmentRight {
		alignStr = "R"
	}

	lineHeight := fontSize*0.3527 + 1.2

	splitLines := ctx.pdf.SplitText(fullText, w)
	currentY := y
	for _, sl := range splitLines {
		ctx.pdf.SetXY(x, currentY)
		ctx.pdf.CellFormat(w, lineHeight, sl, "", 0, alignStr, false, 0, "")
		currentY += lineHeight
	}
}

// parseHexColorRGB parses a hex string into R, G, B channels.
func parseHexColorRGB(hexStr string) (int, int, int) {
	if hexStr == "" {
		return 0, 0, 0
	}
	if hexStr[0] == '#' {
		hexStr = hexStr[1:]
	}
	if len(hexStr) != 6 {
		return 0, 0, 0
	}
	var r, g, b int
	_, err := fmt.Sscanf(hexStr, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return 0, 0, 0
	}
	return r, g, b
}

// formatIfFloat formats decimals into rounded integers.
func formatIfFloat(val string) string {
	if isFloatValue(val) {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return fmt.Sprintf("%.0f", math.Round(f))
		}
	}
	return val
}

// isFloatValue checks if a string is a standard float value.
func isFloatValue(val string) bool {
	if !strings.Contains(val, ".") {
		return false
	}
	_, err := strconv.ParseFloat(val, 64)
	return err == nil
}

func getImageDimensions(path string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	if strings.HasSuffix(strings.ToLower(path), ".webp") {
		cfg, err := webp.DecodeConfig(file)
		if err != nil {
			return 0, 0, err
		}
		return cfg.Width, cfg.Height, nil
	}

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}

func (ctx *renderContext) renderImage(imageObj *structure.Image, x, y, w, h float64) {
	if imageObj.Path == "" {
		return
	}

	renderW := w
	renderH := h

	if imageObj.Width > 0 || imageObj.Height > 0 {
		if imageObj.Width > 0 && imageObj.Height > 0 {
			renderW = imageObj.Width
			renderH = imageObj.Height
		} else {
			wPix, hPix, err := getImageDimensions(imageObj.Path)
			if err == nil && wPix > 0 {
				aspectRatio := float64(hPix) / float64(wPix) // h/w
				imageRatio := float64(wPix) / float64(hPix)  // w/h
				if imageObj.Height > 0 {
					renderH = imageObj.Height
					renderW = renderH * imageRatio
				} else if imageObj.Width > 0 {
					renderW = imageObj.Width
					renderH = renderW * aspectRatio
				}
			}
		}
	}

	offsetX := (w - renderW) / 2.0
	offsetY := (h - renderH) / 2.0
	if offsetX < 0 {
		offsetX = 0
	}
	if offsetY < 0 {
		offsetY = 0
	}
	renderX := x + offsetX
	renderY := y + offsetY

	isWebP := strings.HasSuffix(strings.ToLower(imageObj.Path), ".webp")
	if !isWebP {
		// Native formats (JPEG/PNG/GIF): let gofpdf load and render it natively
		ctx.pdf.Image(imageObj.Path, renderX, renderY, renderW, renderH, false, "", 0, "")
		return
	}

	// WebP transcode and render
	img, err := decodeImage(imageObj.Path)
	if err != nil {
		ctx.drawFallbackError(x, y, w, h)
		return
	}

	var buf bytes.Buffer
	if err := pngenc.Encode(&buf, img); err != nil {
		ctx.drawFallbackError(x, y, w, h)
		return
	}

	ctx.pdf.RegisterImageReader(imageObj.Path, "png", &buf)
	ctx.pdf.Image(imageObj.Path, renderX, renderY, renderW, renderH, false, "png", 0, "")
}

func (ctx *renderContext) drawFallbackError(x, y, w, h float64) {
	ctx.pdf.SetDrawColor(200, 200, 200)
	ctx.pdf.SetLineWidth(0.2)
	ctx.pdf.Rect(x, y, w, h, "D")
	ctx.pdf.SetFont("Arial", "", 6)
	ctx.pdf.SetTextColor(150, 150, 150)
	ctx.pdf.SetXY(x+1, y+h/2-2)
	ctx.pdf.Cell(w-2, 4, "Image Error")
}

func decodeImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if strings.HasSuffix(strings.ToLower(path), ".webp") {
		return webp.Decode(file)
	}

	img, _, err := image.Decode(file)
	return img, err
}

func (ctx *renderContext) renderSignatures(sigs *structure.Signatures, x, y, width float64) {
	n := len(sigs.Items)
	if n == 0 {
		return
	}
	colWidth := width / float64(n)
	for i := 0; i < n; i++ {
		colX := x + float64(i)*colWidth
		centerX := colX + colWidth/2.0

		// Premium signature line length
		lineWidth := 35.0
		if lineWidth > colWidth*0.8 {
			lineWidth = colWidth * 0.8
		}

		lineX1 := centerX - lineWidth/2.0
		lineX2 := centerX + lineWidth/2.0

		// Draw signature line
		ctx.pdf.SetDrawColor(100, 100, 100) // sleek dark gray line
		ctx.pdf.SetLineWidth(0.25)
		ctx.pdf.Line(lineX1, y+10.0, lineX2, y+10.0)

		// Draw centered label text below the line
		ctx.pdf.SetFont("Arial", "B", 7.5) // Sleek bold Arial font
		ctx.pdf.SetTextColor(50, 50, 50)   // Premium dark charcoal color
		ctx.pdf.SetXY(colX, y+11.2)
		ctx.pdf.CellFormat(colWidth, 4.0, sigs.Items[i], "", 0, "C", false, 0, "")
	}
}

func (ctx *renderContext) resolveTextFontProps(textObj *structure.Text) (string, string, float64, string) {
	fontSize := 8.0
	fontFamily := "Arial"
	fontStyleStr := ""
	fontColor := "#000000"

	// 1. Resolve from stylesheet if available
	for _, sName := range textObj.Styles() {
		if ts, ok := ctx.sheet.TextStyles[sName]; ok {
			if ts.FontSize > 0 {
				fontSize = ts.FontSize
			}
			if ts.Bold {
				fontStyleStr += "B"
			}
			if ts.Italic {
				fontStyleStr += "I"
			}
			if ts.Underline {
				fontStyleStr += "U"
			}
			if ts.Font != "" {
				fontFamily = string(ts.Font)
			}
			if ts.Color != "" {
				fontColor = ts.Color
			}
		}
	}

	// 2. Element-level direct font attribute takes final precedence
	if textObj.Font != "" {
		fontFamily = textObj.Font
	}

	// 3. Resolve using dynamic font loader/fallback
	family, styleFlags := ctx.resolveFont(fontFamily, fontStyleStr)
	return family, styleFlags, fontSize, fontColor
}

func (ctx *renderContext) resolveCTextFontProps(ctextObj *structure.CText) (string, string, float64, string) {
	fontSize := 8.0
	fontFamily := "Arial"
	fontStyleStr := ""
	fontColor := "#000000"

	// 1. Resolve from stylesheet if available
	for _, sName := range ctextObj.Styles() {
		// Look up in CTextStyles first
		if cts, ok := ctx.sheet.CTextStyles[sName]; ok {
			if cts.ValueFontSize > 0 {
				fontSize = cts.ValueFontSize
			} else if cts.CaptionFontSize > 0 {
				fontSize = cts.CaptionFontSize
			}
			if cts.ValueBold || cts.CaptionBold {
				fontStyleStr += "B"
			}
			if cts.ValueItalic || cts.CaptionItalic {
				fontStyleStr += "I"
			}
			if cts.ValueUnderline || cts.CaptionUnderline {
				fontStyleStr += "U"
			}
			if cts.ValueFont != "" {
				fontFamily = string(cts.ValueFont)
			} else if cts.CaptionFont != "" {
				fontFamily = string(cts.CaptionFont)
			}
			if cts.ValueColor != "" {
				fontColor = cts.ValueColor
			} else if cts.CaptionColor != "" {
				fontColor = cts.CaptionColor
			}
		} else if ts, ok := ctx.sheet.TextStyles[sName]; ok {
			// Fallback to normal TextStyles lookup
			if ts.FontSize > 0 {
				fontSize = ts.FontSize
			}
			if ts.Bold {
				fontStyleStr += "B"
			}
			if ts.Italic {
				fontStyleStr += "I"
			}
			if ts.Underline {
				fontStyleStr += "U"
			}
			if ts.Font != "" {
				fontFamily = string(ts.Font)
			}
			if ts.Color != "" {
				fontColor = ts.Color
			}
		}
	}

	// 2. Element-level direct font attribute takes final precedence
	if ctextObj.Font != "" {
		fontFamily = ctextObj.Font
	}

	// 3. Resolve using dynamic font loader/fallback
	family, styleFlags := ctx.resolveFont(fontFamily, fontStyleStr)
	return family, styleFlags, fontSize, fontColor
}

func (ctx *renderContext) resolveFont(family string, styleStr string) (string, string) {
	// Standardize style string (only keep B, I, U)
	var resolvedStyle string
	styleStr = strings.ToUpper(styleStr)
	if strings.Contains(styleStr, "B") {
		resolvedStyle += "B"
	}
	if strings.Contains(styleStr, "I") {
		resolvedStyle += "I"
	}
	if strings.Contains(styleStr, "U") {
		resolvedStyle += "U"
	}

	family = strings.TrimSpace(strings.ToLower(family))
	if family == "" || family == "arial" || family == "helvetica" || family == "sans-serif" {
		return "Arial", resolvedStyle
	}
	if family == "times" || family == "times-new-roman" || family == "serif" {
		return "Times", resolvedStyle
	}
	if family == "courier" || family == "courier-new" || family == "monospace" {
		return "Courier", resolvedStyle
	}

	// For custom fonts, check if already registered
	key := family + "_" + resolvedStyle
	if ctx.registeredFonts == nil {
		ctx.registeredFonts = make(map[string]bool)
	}

	if ctx.registeredFonts[key] {
		return family, resolvedStyle
	}

	// Look for a .ttf file in potential paths
	pathsToTry := []string{
		fmt.Sprintf("testdata/%s.ttf", family),
		fmt.Sprintf("pkg/renderer/testdata/%s.ttf", family),
		fmt.Sprintf("fonts/%s.ttf", family),
		fmt.Sprintf("%s.ttf", family),
	}

	// If style has B/I/U, we can also look for styled TTF files
	if resolvedStyle != "" {
		styleName := ""
		if resolvedStyle == "B" {
			styleName = "-bold"
		} else if resolvedStyle == "I" {
			styleName = "-italic"
		} else if resolvedStyle == "BI" || resolvedStyle == "IB" {
			styleName = "-bolditalic"
		}
		pathsToTry = append([]string{
			fmt.Sprintf("testdata/%s%s.ttf", family, styleName),
			fmt.Sprintf("pkg/renderer/testdata/%s%s.ttf", family, styleName),
			fmt.Sprintf("fonts/%s%s.ttf", family, styleName),
			fmt.Sprintf("%s%s.ttf", family, styleName),
		}, pathsToTry...)
	}

	for _, p := range pathsToTry {
		if _, err := os.Stat(p); err == nil {
			ctx.pdf.AddUTF8Font(family, resolvedStyle, p)
			ctx.registeredFonts[key] = true
			return family, resolvedStyle
		}
	}

	// Graceful fallback if custom font is not found
	if strings.Contains(family, "ocr") || strings.Contains(family, "code") || strings.Contains(family, "mono") {
		return "Courier", resolvedStyle
	}
	return "Arial", resolvedStyle
}

// resolveFont2 converts bool bold/italic flags to a style string and delegates to resolveFont.
func (ctx *renderContext) resolveFont2(family string, bold, italic bool) (string, string, float64) {
	styleStr := ""
	if bold {
		styleStr += "B"
	}
	if italic {
		styleStr += "I"
	}
	fam, flags := ctx.resolveFont(family, styleStr)
	return fam, flags, 0 // fontSize handled by caller
}

// ─── Style resolvers for typed elements ───────────────────────────────────────

func (ctx *renderContext) resolveNumberStyle(el *structure.Number) style.NumberStyle {
	s := style.NumberStyle{FontSize: 8.0, Alignment: properties.AlignmentCenter}
	// Apply default style if set
	if ctx.sheet.DefaultNumberStyle != nil {
		applyNumberStyle(&s, *ctx.sheet.DefaultNumberStyle)
	}
	// Apply named styles (override default)
	for _, sName := range el.Styles() {
		if ns, ok := ctx.sheet.NumberStyles[sName]; ok {
			applyNumberStyle(&s, ns)
		}
	}
	if s.FontSize == 0 {
		s.FontSize = 8.0
	}
	return s
}

func (ctx *renderContext) resolveMoneyStyle(el *structure.Money) style.MoneyStyle {
	s := style.MoneyStyle{FontSize: 8.0, Alignment: properties.AlignmentRight, Decimals: 0}
	if ctx.sheet.DefaultMoneyStyle != nil {
		applyMoneyStyle(&s, *ctx.sheet.DefaultMoneyStyle)
	}
	for _, sName := range el.Styles() {
		if ms, ok := ctx.sheet.MoneyStyles[sName]; ok {
			applyMoneyStyle(&s, ms)
		}
	}
	if s.FontSize == 0 {
		s.FontSize = 8.0
	}
	return s
}

func (ctx *renderContext) resolveDateStyle(el *structure.Date) style.DateStyle {
	s := style.DateStyle{FontSize: 8.0, Alignment: properties.AlignmentCenter, Format: "dd/MM/yyyy"}
	if ctx.sheet.DefaultDateStyle != nil {
		applyDateStyle(&s, *ctx.sheet.DefaultDateStyle)
	}
	for _, sName := range el.Styles() {
		if ds, ok := ctx.sheet.DateStyles[sName]; ok {
			applyDateStyle(&s, ds)
		}
	}
	if s.FontSize == 0 {
		s.FontSize = 8.0
	}
	return s
}

func (ctx *renderContext) resolveTimeStyle(el *structure.Time) style.TimeStyle {
	s := style.TimeStyle{FontSize: 8.0, Alignment: properties.AlignmentCenter, Format: "HH:mm"}
	if ctx.sheet.DefaultTimeStyle != nil {
		applyTimeStyle(&s, *ctx.sheet.DefaultTimeStyle)
	}
	for _, sName := range el.Styles() {
		if ts, ok := ctx.sheet.TimeStyles[sName]; ok {
			applyTimeStyle(&s, ts)
		}
	}
	if s.FontSize == 0 {
		s.FontSize = 8.0
	}
	return s
}

func (ctx *renderContext) resolveFloatStyle(el *structure.Float) style.FloatStyle {
	s := style.FloatStyle{FontSize: 8.0, Alignment: properties.AlignmentCenter, Decimals: 2}
	if ctx.sheet.DefaultFloatStyle != nil {
		applyFloatStyle(&s, *ctx.sheet.DefaultFloatStyle)
	}
	for _, sName := range el.Styles() {
		if fs, ok := ctx.sheet.FloatStyles[sName]; ok {
			applyFloatStyle(&s, fs)
		}
	}
	if s.FontSize == 0 {
		s.FontSize = 8.0
	}
	return s
}

// ─── Style applier helpers (merge non-zero fields onto a base) ────────────────

func applyNumberStyle(dst *style.NumberStyle, src style.NumberStyle) {
	if src.FontSize > 0 {
		dst.FontSize = src.FontSize
	}
	if src.Font != "" {
		dst.Font = src.Font
	}
	if src.Color != "" {
		dst.Color = src.Color
	}
	if src.Alignment != "" {
		dst.Alignment = src.Alignment
	}
	dst.Bold = src.Bold
	dst.Italic = src.Italic
	dst.UseComma = src.UseComma
}

func applyMoneyStyle(dst *style.MoneyStyle, src style.MoneyStyle) {
	if src.FontSize > 0 {
		dst.FontSize = src.FontSize
	}
	if src.Font != "" {
		dst.Font = src.Font
	}
	if src.Color != "" {
		dst.Color = src.Color
	}
	if src.Alignment != "" {
		dst.Alignment = src.Alignment
	}
	dst.Bold = src.Bold
	dst.Italic = src.Italic
	dst.UseComma = src.UseComma
	// decimals: always copy from src (0 is a valid value)
	dst.Decimals = src.Decimals
	if src.Symbol != "" {
		dst.Symbol = src.Symbol
	}
	if src.SymbolPosition != "" {
		dst.SymbolPosition = src.SymbolPosition
	}
}

func applyDateStyle(dst *style.DateStyle, src style.DateStyle) {
	if src.FontSize > 0 {
		dst.FontSize = src.FontSize
	}
	if src.Font != "" {
		dst.Font = src.Font
	}
	if src.Color != "" {
		dst.Color = src.Color
	}
	if src.Alignment != "" {
		dst.Alignment = src.Alignment
	}
	dst.Bold = src.Bold
	dst.Italic = src.Italic
	if src.Format != "" {
		dst.Format = src.Format
	}
}

func applyTimeStyle(dst *style.TimeStyle, src style.TimeStyle) {
	if src.FontSize > 0 {
		dst.FontSize = src.FontSize
	}
	if src.Font != "" {
		dst.Font = src.Font
	}
	if src.Color != "" {
		dst.Color = src.Color
	}
	if src.Alignment != "" {
		dst.Alignment = src.Alignment
	}
	dst.Bold = src.Bold
	dst.Italic = src.Italic
	if src.Format != "" {
		dst.Format = src.Format
	}
}

func applyFloatStyle(dst *style.FloatStyle, src style.FloatStyle) {
	if src.FontSize > 0 {
		dst.FontSize = src.FontSize
	}
	if src.Font != "" {
		dst.Font = src.Font
	}
	if src.Color != "" {
		dst.Color = src.Color
	}
	if src.Alignment != "" {
		dst.Alignment = src.Alignment
	}
	dst.Bold = src.Bold
	dst.Italic = src.Italic
	dst.UseComma = src.UseComma
	// decimals: always copy
	dst.Decimals = src.Decimals
}

// ─── Format functions ─────────────────────────────────────────────────────────

// formatNumber formats a numeric string as an integer, optionally with commas.
func formatNumber(val string, s style.NumberStyle) string {
	val = strings.TrimSpace(val)
	if val == "" {
		return val
	}
	// Parse as float first to handle "1234.0" inputs gracefully
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val
	}
	n := int64(math.Round(f))
	if s.UseComma {
		return formatWithCommas(strconv.FormatInt(n, 10), "")
	}
	return strconv.FormatInt(n, 10)
}

// formatMoney formats a monetary value with configurable decimals, commas, and symbol.
func formatMoney(val string, s style.MoneyStyle) string {
	val = strings.TrimSpace(val)
	if val == "" {
		return val
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val
	}
	// Format with required decimal places
	formatted := strconv.FormatFloat(f, 'f', s.Decimals, 64)
	// Split into integer and decimal parts
	parts := strings.SplitN(formatted, ".", 2)
	intPart := parts[0]
	decPart := ""
	if len(parts) == 2 {
		decPart = "." + parts[1]
	}
	if s.UseComma {
		intPart = formatWithCommas(intPart, "")
	}
	result := intPart + decPart
	// Attach symbol
	if s.Symbol != "" {
		if s.SymbolPosition == "suffix" {
			result = result + s.Symbol
		} else {
			result = s.Symbol + result
		}
	}
	return result
}

// formatDate parses an ISO 8601 date string and reformats it using a simple pattern.
// Supported output tokens: yyyy, yy, MM, M, dd, d
func formatDate(val string, s style.DateStyle) string {
	val = strings.TrimSpace(val)
	if val == "" {
		return val
	}
	// Try common input formats
	inputFormats := []string{"2006-01-02", "02/01/2006", "01/02/2006", "2006/01/02", "02-01-2006"}
	var parsed gotime.Time
	var err error
	for _, fmt_ := range inputFormats {
		parsed, err = gotime.Parse(fmt_, val)
		if err == nil {
			break
		}
	}
	if err != nil {
		return val
	} // can't parse — return as-is
	if s.Format == "" {
		return val
	}
	// Convert user-friendly tokens to Go layout
	goLayout := s.Format
	goLayout = strings.ReplaceAll(goLayout, "yyyy", "2006")
	goLayout = strings.ReplaceAll(goLayout, "yy", "06")
	goLayout = strings.ReplaceAll(goLayout, "MM", "01")
	goLayout = strings.ReplaceAll(goLayout, "M", "1")
	goLayout = strings.ReplaceAll(goLayout, "dd", "02")
	goLayout = strings.ReplaceAll(goLayout, "d", "2")
	return parsed.Format(goLayout)
}

// formatTime parses a time string and reformats using a simple pattern.
// Supported output tokens: HH, H, mm, m, ss, s
func formatTime(val string, s style.TimeStyle) string {
	val = strings.TrimSpace(val)
	if val == "" {
		return val
	}
	inputFormats := []string{"15:04:05", "15:04", "3:04 PM", "3:04PM"}
	var parsed gotime.Time
	var err error
	for _, fmt_ := range inputFormats {
		parsed, err = gotime.Parse(fmt_, val)
		if err == nil {
			break
		}
	}
	if err != nil {
		return val
	}
	if s.Format == "" {
		return val
	}
	goLayout := s.Format
	goLayout = strings.ReplaceAll(goLayout, "HH", "15")
	goLayout = strings.ReplaceAll(goLayout, "H", "15")
	goLayout = strings.ReplaceAll(goLayout, "mm", "04")
	goLayout = strings.ReplaceAll(goLayout, "m", "4")
	goLayout = strings.ReplaceAll(goLayout, "ss", "05")
	goLayout = strings.ReplaceAll(goLayout, "s", "5")
	return parsed.Format(goLayout)
}

// formatFloat formats a float string with a set number of decimal places and optional commas.
func formatFloat(val string, s style.FloatStyle) string {
	val = strings.TrimSpace(val)
	if val == "" {
		return val
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val
	}
	formatted := strconv.FormatFloat(f, 'f', s.Decimals, 64)
	if !s.UseComma {
		return formatted
	}
	parts := strings.SplitN(formatted, ".", 2)
	intPart := formatWithCommas(parts[0], "")
	if len(parts) == 2 {
		return intPart + "." + parts[1]
	}
	return intPart
}

// formatWithCommas inserts thousand-separator commas into an integer string.
func formatWithCommas(s, neg string) string {
	// Handle negative
	prefix := neg
	if strings.HasPrefix(s, "-") {
		prefix = "-"
		s = s[1:]
	}
	n := len(s)
	if n <= 3 {
		return prefix + s
	}
	var buf strings.Builder
	rem := n % 3
	if rem > 0 {
		buf.WriteString(s[:rem])
		if n > rem {
			buf.WriteByte(',')
		}
	}
	for i := rem; i < n; i += 3 {
		buf.WriteString(s[i : i+3])
		if i+3 < n {
			buf.WriteByte(',')
		}
	}
	return prefix + buf.String()
}

// ─── Render methods for typed elements ───────────────────────────────────────

func (ctx *renderContext) renderTypedText(display string, align properties.Alignment, fontFamily string, bold, italic bool, fontSize float64, color string, x, y, w, h float64) {
	styleStr := ""
	if bold {
		styleStr += "B"
	}
	if italic {
		styleStr += "I"
	}
	family, styleFlags := ctx.resolveFont(fontFamily, styleStr)
	ctx.pdf.SetFont(family, styleFlags, fontSize)
	r, g, b := parseHexColorRGB(color)
	ctx.pdf.SetTextColor(r, g, b)
	alignStr := "L"
	if align == properties.AlignmentCenter {
		alignStr = "C"
	} else if align == properties.AlignmentRight {
		alignStr = "R"
	}
	lineHeight := fontSize*0.3527 + 1.2
	ctx.pdf.SetXY(x, y)
	ctx.pdf.CellFormat(w, lineHeight, display, "", 0, alignStr, false, 0, "")
}

func (ctx *renderContext) renderNumber(el *structure.Number, x, y, w, h float64) {
	s := ctx.resolveNumberStyle(el)
	display := formatNumber(el.Value, s)
	color := s.Color
	if color == "" {
		color = "#000000"
	}
	ctx.renderTypedText(display, s.Alignment, string(s.Font), s.Bold, s.Italic, s.FontSize, color, x, y, w, h)
}

func (ctx *renderContext) renderMoney(el *structure.Money, x, y, w, h float64) {
	s := ctx.resolveMoneyStyle(el)
	display := formatMoney(el.Value, s)
	color := s.Color
	if color == "" {
		color = "#000000"
	}
	ctx.renderTypedText(display, s.Alignment, string(s.Font), s.Bold, s.Italic, s.FontSize, color, x, y, w, h)
}

func (ctx *renderContext) renderDate(el *structure.Date, x, y, w, h float64) {
	s := ctx.resolveDateStyle(el)
	display := formatDate(el.Value, s)
	color := s.Color
	if color == "" {
		color = "#000000"
	}
	ctx.renderTypedText(display, s.Alignment, string(s.Font), s.Bold, s.Italic, s.FontSize, color, x, y, w, h)
}

func (ctx *renderContext) renderTime(el *structure.Time, x, y, w, h float64) {
	s := ctx.resolveTimeStyle(el)
	display := formatTime(el.Value, s)
	color := s.Color
	if color == "" {
		color = "#000000"
	}
	ctx.renderTypedText(display, s.Alignment, string(s.Font), s.Bold, s.Italic, s.FontSize, color, x, y, w, h)
}

func (ctx *renderContext) renderFloat(el *structure.Float, x, y, w, h float64) {
	s := ctx.resolveFloatStyle(el)
	display := formatFloat(el.Value, s)
	color := s.Color
	if color == "" {
		color = "#000000"
	}
	ctx.renderTypedText(display, s.Alignment, string(s.Font), s.Bold, s.Italic, s.FontSize, color, x, y, w, h)
}
