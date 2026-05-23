package data

import (
	"encoding/json"
	"fmt"

	"github.com/niarbSoftHub/squash/pkg/interfaces"
	"github.com/niarbSoftHub/squash/pkg/properties"
)

type DataRoot struct {
	VariableValuePair map[string]string       `json:"variableValuePair,omitempty"`
	Details           []interfaces.DataDetail `json:"details,omitempty"`
}

type GroupData struct {
	DataRoot
}

func (g GroupData) DetailType() properties.DataDetailsType {
	return properties.DataGroup
}

// MarshalJSON sets the type field
func (g GroupData) MarshalJSON() ([]byte, error) {
	type Alias GroupData
	return json.Marshal(&struct {
		Type properties.DataDetailsType `json:"type"`
		Alias
	}{
		Type:  properties.DataGroup,
		Alias: (Alias)(g),
	})
}

type TableData struct {
	Columns  []string   `json:"columns,omitempty"`
	DataRows [][]string `json:"dataRows,omitempty"`
}

func (t TableData) DetailType() properties.DataDetailsType {
	return properties.DataTable
}

func (t TableData) MarshalJSON() ([]byte, error) {
	type Alias TableData
	return json.Marshal(&struct {
		Type properties.DataDetailsType `json:"type"`
		Alias
	}{
		Type:  properties.DataTable,
		Alias: (Alias)(t),
	})
}

// Custom UnmarshalJSON for DataRoot to handle polymorphic Details
func (dr *DataRoot) UnmarshalJSON(b []byte) error {
	type Alias DataRoot
	aux := &struct {
		Details []json.RawMessage `json:"details"`
		*Alias
	}{
		Alias: (*Alias)(dr),
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	dr.Details = make([]interfaces.DataDetail, len(aux.Details))
	for i, raw := range aux.Details {
		var typeExtract struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(raw, &typeExtract); err != nil {
			return err
		}

		switch properties.DataDetailsType(typeExtract.Type) {
		case properties.DataGroup:
			var g GroupData
			if err := json.Unmarshal(raw, &g); err != nil {
				return err
			}
			dr.Details[i] = g
		case properties.DataTable:
			var t TableData
			if err := json.Unmarshal(raw, &t); err != nil {
				return err
			}
			dr.Details[i] = t
		default:
			return fmt.Errorf("unknown DataDetail type: %s", typeExtract.Type)
		}
	}
	return nil
}
