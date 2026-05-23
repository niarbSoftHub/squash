package interfaces

import "github.com/niarbSoftHub/squash/pkg/properties"

// Element is the base interface for all structural items.
type Element interface {
	Type() properties.ElementType
	Parent() Parent
	Styles() []string
}

// Parent is an Element that can contain other Elements.
type Parent interface {
	Element
	AddElement(Element)
}

// Detail is a marker interface for repeating structural elements in the report body.
type Detail interface {
	Element
	// Marker interface, but technically Details are Elements, so we embed Element.
}

// DataDetail is an interface for data model items that correspond to details.
type DataDetail interface {
	DetailType() properties.DataDetailsType
}
