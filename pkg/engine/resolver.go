package engine

import (
	"regexp"
	"strings"

	"github.com/niarbSoftHub/squash/pkg/interfaces"
	"github.com/niarbSoftHub/squash/pkg/models/data"
	"github.com/niarbSoftHub/squash/pkg/models/structure"
)

var variableRegex = regexp.MustCompile(`\{\{\s*([^}\s]+)\s*\}\}`)

// ResolveString replaces {{varName}} with values from the vars map.
func ResolveString(text string, vars map[string]string) string {
	if vars == nil {
		return text
	}
	return variableRegex.ReplaceAllStringFunc(text, func(match string) string {
		varName := strings.TrimSpace(match[2 : len(match)-2])
		if val, ok := vars[varName]; ok {
			return val
		}
		return match
	})
}

// CloneElement deeply clones a structural element.
func CloneElement(e interfaces.Element) interfaces.Element {
	if e == nil {
		return nil
	}
	switch v := e.(type) {
	case *structure.Spacer:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		return &clone
	case *structure.Text:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		return &clone
	case *structure.CText:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		return &clone
	case *structure.Number:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		return &clone
	case *structure.Money:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		return &clone
	case *structure.Date:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		return &clone
	case *structure.Time:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		return &clone
	case *structure.Float:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		return &clone
	case *structure.Image:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		return &clone
	case *structure.Cell:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		if v.Element.Element != nil {
			clone.Element = structure.ElementWrapper{Element: CloneElement(v.Element.Element)}
		}
		return &clone
	case *structure.VerticalLayout:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		clone.Elements = make([]structure.ElementWrapper, len(v.Elements))
		for i, ew := range v.Elements {
			clone.Elements[i] = structure.ElementWrapper{Element: CloneElement(ew.Element)}
		}
		return &clone
	case *structure.HorizontalLayout:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		clone.Elements = make([]structure.ElementWrapper, len(v.Elements))
		for i, ew := range v.Elements {
			clone.Elements[i] = structure.ElementWrapper{Element: CloneElement(ew.Element)}
		}
		return &clone
	case *structure.Group:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		clone.PageBreak = v.PageBreak
		if v.Header != nil {
			clone.Header = CloneElement(v.Header).(*structure.VerticalLayout)
		}
		if v.Footer != nil {
			clone.Footer = CloneElement(v.Footer).(*structure.VerticalLayout)
		}
		if v.Details != nil && v.Details.Detail != nil {
			clone.Details = &structure.DetailWrapper{Detail: CloneElement(v.Details.Detail).(interfaces.Detail)}
		}
		return &clone
	case *structure.Report:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		if v.ReportHeader != nil {
			clone.ReportHeader = CloneElement(v.ReportHeader).(*structure.VerticalLayout)
		}
		if v.PageHeader != nil {
			clone.PageHeader = CloneElement(v.PageHeader).(*structure.VerticalLayout)
		}
		if v.ReportFooter != nil {
			clone.ReportFooter = CloneElement(v.ReportFooter).(*structure.VerticalLayout)
		}
		if v.PageFooter != nil {
			clone.PageFooter = CloneElement(v.PageFooter).(*structure.VerticalLayout)
		}
		if v.Details != nil && v.Details.Detail != nil {
			clone.Details = &structure.DetailWrapper{Detail: CloneElement(v.Details.Detail).(interfaces.Detail)}
		}
		return &clone
	case *structure.Signatures:
		if v == nil {
			return nil
		}
		clone := *v
		clone.StyleRaw = v.StyleRaw
		clone.Font = v.Font
		clone.Items = make([]string, len(v.Items))
		copy(clone.Items, v.Items)
		return &clone
	}
	return nil
}

// ExpandTree fully resolves all variables and expands data details into vertical layouts of repeated structural elements.
func ExpandTree(report *structure.Report, dataRoot *data.DataRoot) *structure.Report {
	clonedReport := CloneElement(report).(*structure.Report)

	// Resolve variables in Report Headers and Footers
	resolveVariablesInElement(clonedReport.ReportHeader, dataRoot.VariableValuePair)
	resolveVariablesInElement(clonedReport.PageHeader, dataRoot.VariableValuePair)
	resolveVariablesInElement(clonedReport.ReportFooter, dataRoot.VariableValuePair)
	resolveVariablesInElement(clonedReport.PageFooter, dataRoot.VariableValuePair)

	// Now expand the details section recursively
	if clonedReport.Details != nil && clonedReport.Details.Detail != nil {
		expandedDetails := expandDetail(clonedReport.Details.Detail, dataRoot)
		clonedReport.Details.Detail = expandedDetails.(interfaces.Detail)
	}

	return clonedReport
}

// In-place variable resolution for a single tree segment
func resolveVariablesInElement(e interfaces.Element, vars map[string]string) {
	if e == nil {
		return
	}
	switch v := e.(type) {
	case *structure.Spacer:
		// Spacer elements do not contain text/variables to resolve
	case *structure.Text:
		if v == nil {
			return
		}
		v.Value = ResolveString(v.Value, vars)
	case *structure.Number:
		if v == nil {
			return
		}
		v.Value = ResolveString(v.Value, vars)
	case *structure.Money:
		if v == nil {
			return
		}
		v.Value = ResolveString(v.Value, vars)
	case *structure.Date:
		if v == nil {
			return
		}
		v.Value = ResolveString(v.Value, vars)
	case *structure.Time:
		if v == nil {
			return
		}
		v.Value = ResolveString(v.Value, vars)
	case *structure.Float:
		if v == nil {
			return
		}
		v.Value = ResolveString(v.Value, vars)
	case *structure.CText:
		if v == nil {
			return
		}
		v.Caption = ResolveString(v.Caption, vars)
		v.Value = ResolveString(v.Value, vars)
	case *structure.Image:
		if v == nil {
			return
		}
		v.Path = ResolveString(v.Path, vars)
	case *structure.Cell:
		if v == nil {
			return
		}
		if v.Element.Element != nil {
			resolveVariablesInElement(v.Element.Element, vars)
		}
	case *structure.VerticalLayout:
		if v == nil {
			return
		}
		for _, ew := range v.Elements {
			resolveVariablesInElement(ew.Element, vars)
		}
	case *structure.HorizontalLayout:
		if v == nil {
			return
		}
		for _, ew := range v.Elements {
			resolveVariablesInElement(ew.Element, vars)
		}
	case *structure.Group:
		if v == nil {
			return
		}
		resolveVariablesInElement(v.Header, vars)
		resolveVariablesInElement(v.Footer, vars)
		// We do not resolve Details here because they need per-row data scoping!
	case *structure.Signatures:
		if v == nil {
			return
		}
		for i, item := range v.Items {
			v.Items[i] = ResolveString(item, vars)
		}
	}
}

// expandDetail processes a single Detail element (Group or HorizontalLayout) against a DataRoot (which contains arrays of data).
func expandDetail(detail interfaces.Detail, dRoot *data.DataRoot) interfaces.Element {
	// The result of expanding a detail is a vertical stack of that detail repeated for each data item.
	expandedLayout := &structure.VerticalLayout{}

	switch v := detail.(type) {
	case *structure.Group:
		// A structural Group iterates over DataRoot's DataDetails (which should be GroupDatas)
		for _, dataDetail := range dRoot.Details {
			if gData, ok := dataDetail.(data.GroupData); ok {
				// Clone the group structure for this iteration
				groupClone := CloneElement(v).(*structure.Group)

				// Resolve header/footer with this specific group's variable scope
				resolveVariablesInElement(groupClone.Header, gData.VariableValuePair)
				resolveVariablesInElement(groupClone.Footer, gData.VariableValuePair)

				// Recursively expand the inner detail with this specific group's data scope
				if groupClone.Details != nil && groupClone.Details.Detail != nil {
					innerExpanded := expandDetail(groupClone.Details.Detail, &gData.DataRoot)
					groupClone.Details.Detail = innerExpanded.(interfaces.Detail)
				}

				expandedLayout.AddElement(groupClone)
			}
		}

	case *structure.HorizontalLayout:
		// A structural HorizontalLayout (row) iterates over DataRoot's DataDetails (which should be TableData)
		// We expect exactly ONE TableData at this level
		if len(dRoot.Details) == 1 {
			if tData, ok := dRoot.Details[0].(data.TableData); ok {
				// Iterate over table rows
				for _, rowValues := range tData.DataRows {
					rowClone := CloneElement(v).(*structure.HorizontalLayout)

					// Build a temporary map for this row to resolve variables
					rowVars := make(map[string]string)
					for i, colName := range tData.Columns {
						if i < len(rowValues) {
							rowVars[colName] = rowValues[i]
						}
					}

					// Resolve variables in this specific cloned row
					resolveVariablesInElement(rowClone, rowVars)

					expandedLayout.AddElement(rowClone)
				}
			}
		}
	}

	return expandedLayout
}
