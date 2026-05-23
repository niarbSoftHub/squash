package structure

import (
	"encoding/xml"
	"fmt"

	"github.com/niarbSoftHub/squash/pkg/interfaces"
)

// ElementWrapper wraps an interfaces.Element so we can use xml:",any"
// to parse ordered polymorphic children.
type ElementWrapper struct {
	Element interfaces.Element
}

func (ew *ElementWrapper) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var e interfaces.Element
	switch start.Name.Local {
	case "space":
		var spacer Spacer
		if err := d.DecodeElement(&spacer, &start); err != nil {
			return err
		}
		e = &spacer
	case "text":
		var text Text
		if err := d.DecodeElement(&text, &start); err != nil {
			return err
		}
		e = &text
	case "number":
		var num Number
		if err := d.DecodeElement(&num, &start); err != nil {
			return err
		}
		e = &num
	case "money":
		var money Money
		if err := d.DecodeElement(&money, &start); err != nil {
			return err
		}
		e = &money
	case "date":
		var date Date
		if err := d.DecodeElement(&date, &start); err != nil {
			return err
		}
		e = &date
	case "time":
		var t Time
		if err := d.DecodeElement(&t, &start); err != nil {
			return err
		}
		e = &t
	case "float":
		var fl Float
		if err := d.DecodeElement(&fl, &start); err != nil {
			return err
		}
		e = &fl
	case "c-text":
		var ctext CText
		if err := d.DecodeElement(&ctext, &start); err != nil {
			return err
		}
		e = &ctext
	case "image":
		var img Image
		if err := d.DecodeElement(&img, &start); err != nil {
			return err
		}
		e = &img
	case "cell":
		var cell Cell
		if err := d.DecodeElement(&cell, &start); err != nil {
			return err
		}
		e = &cell
	case "group":
		var g Group
		if err := d.DecodeElement(&g, &start); err != nil {
			return err
		}
		e = &g
	case "layout-h":
		var hl HorizontalLayout
		if err := d.DecodeElement(&hl, &start); err != nil {
			return err
		}
		e = &hl
	case "layout-v":
		var vl VerticalLayout
		if err := d.DecodeElement(&vl, &start); err != nil {
			return err
		}
		e = &vl
	case "signatures":
		var sigs Signatures
		if err := d.DecodeElement(&sigs, &start); err != nil {
			return err
		}
		e = &sigs
	default:
		return fmt.Errorf("unknown xml tag for element: %s", start.Name.Local)
	}
	ew.Element = e
	return nil
}

func (ew ElementWrapper) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if ew.Element == nil {
		return nil
	}
	return e.EncodeElement(ew.Element, xml.StartElement{Name: xml.Name{Local: string(ew.Element.Type())}})
}

// DetailWrapper handles polymorphic parsing of the interfaces.Detail marker.
type DetailWrapper struct {
	Detail interfaces.Detail
}

func (dw *DetailWrapper) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// The start element is <details>
	// We need to decode the child inside it.
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch se := t.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "group":
				var g Group
				if err := d.DecodeElement(&g, &se); err != nil {
					return err
				}
				dw.Detail = &g
			case "layout-h":
				var hl HorizontalLayout
				if err := d.DecodeElement(&hl, &se); err != nil {
					return err
				}
				dw.Detail = &hl
			default:
				return fmt.Errorf("invalid detail tag (must be group or layout-h): %s", se.Name.Local)
			}
		case xml.EndElement:
			if se.Name == start.Name {
				return nil
			}
		}
	}
}

func (dw DetailWrapper) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if dw.Detail == nil {
		return nil
	}
	// start element is <details>. We write it.
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	// write the actual detail object
	if err := e.EncodeElement(dw.Detail, xml.StartElement{Name: xml.Name{Local: string(dw.Detail.Type())}}); err != nil {
		return err
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}
