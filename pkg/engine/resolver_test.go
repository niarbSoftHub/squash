package engine_test

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/niarbSoftHub/squash/pkg/engine"
	"github.com/niarbSoftHub/squash/pkg/models/data"
	"github.com/niarbSoftHub/squash/pkg/models/structure"
)

func TestResolveVariables(t *testing.T) {
	vars := map[string]string{
		"companyName": "Test Co",
		"year":        "2023",
	}

	res := engine.ResolveString("Welcome to {{companyName}} in {{year}}!", vars)
	if res != "Welcome to Test Co in 2023!" {
		t.Errorf("Expected 'Welcome to Test Co in 2023!', got '%s'", res)
	}

	res2 := engine.ResolveString("Missing {{varName}} is ignored.", vars)
	if res2 != "Missing {{varName}} is ignored." {
		t.Errorf("Expected 'Missing {{varName}} is ignored.', got '%s'", res2)
	}
}

func TestExpandTree(t *testing.T) {
	// 1. Setup Data Model
	dataRoot := data.DataRoot{
		VariableValuePair: map[string]string{
			"reportName": "Employee List",
		},
	}
	tableData := data.TableData{
		Columns: []string{"id", "name"},
		DataRows: [][]string{
			{"1", "Alice"},
			{"2", "Bob"},
			{"3", "Charlie"},
		},
	}
	dataRoot.Details = append(dataRoot.Details, tableData)

	// Serialize/Deserialize trick to set types correctly on data model
	b, _ := json.Marshal(dataRoot)
	var finalDataRoot data.DataRoot
	json.Unmarshal(b, &finalDataRoot)

	// 2. Setup Structure Model
	xmlData := `
	<report>
		<reportHeader>
			<layout-v>
				<text>{{reportName}}</text>
			</layout-v>
		</reportHeader>
		<details>
			<layout-h>
				<cell><text>{{id}}</text></cell>
				<cell><text>{{name}}</text></cell>
			</layout-h>
		</details>
	</report>
	`
	var report structure.Report
	if err := xml.Unmarshal([]byte(xmlData), &report); err != nil {
		t.Fatalf("Failed to parse xml: %v", err)
	}

	// 3. Expand!
	expanded := engine.ExpandTree(&report, &finalDataRoot)

	// 4. Verify
	// Check report header is resolved
	text, _ := expanded.ReportHeader.Elements[0].Element.(*structure.Text)
	if text.Value != "Employee List" {
		t.Errorf("Expected 'Employee List', got '%s'", text.Value)
	}

	// Details should now be a VerticalLayout containing 3 HorizontalLayout rows
	vLayout, ok := expanded.Details.Detail.(*structure.VerticalLayout)
	if !ok {
		t.Fatalf("Expected Details to be expanded into VerticalLayout, got %T", expanded.Details.Detail)
	}

	if len(vLayout.Elements) != 3 {
		t.Fatalf("Expected 3 rows in expanded table, got %d", len(vLayout.Elements))
	}

	// Check row 2 (Bob)
	row2, ok := vLayout.Elements[1].Element.(*structure.HorizontalLayout)
	if !ok {
		t.Fatalf("Expected child to be HorizontalLayout")
	}

	cell1, _ := row2.Elements[0].Element.(*structure.Cell)
	textID, _ := cell1.Element.Element.(*structure.Text)
	if textID.Value != "2" {
		t.Errorf("Expected '2', got '%s'", textID.Value)
	}

	cell2, _ := row2.Elements[1].Element.(*structure.Cell)
	textName, _ := cell2.Element.Element.(*structure.Text)
	if textName.Value != "Bob" {
		t.Errorf("Expected 'Bob', got '%s'", textName.Value)
	}
}
