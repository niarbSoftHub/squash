package structure_test

import (
	"encoding/xml"
	"testing"

	"github.com/niarbSoftHub/squash/pkg/models/structure"
)

func TestXMLDeserialization(t *testing.T) {
	xmlData := `
	<report orientation="0" pageSize="0">
		<reportHeader>
			<layout-v>
				<text styles="companyNameStyle">{{companyName}}</text>
				<text styles="reportNameStyle">{{reportName}}</text>
			</layout-v>
		</reportHeader>
		<details>
			<group>
				<header>
					<layout-v>
						<c-text styles="groupCtextStyle1" caption="Sub company">{{subCompany}}</c-text>
					</layout-v>
				</header>
				<details>
					<layout-h>
						<cell><text>{{employeeId}}</text></cell>
						<cell ratio="2"><text>{{employeeName}}</text></cell>
					</layout-h>
				</details>
			</group>
		</details>
	</report>
	`

	var report structure.Report
	if err := xml.Unmarshal([]byte(xmlData), &report); err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// Verify data
	if report.ReportHeader == nil {
		t.Fatalf("ReportHeader is nil")
	}

	if len(report.ReportHeader.Elements) != 2 {
		t.Fatalf("Expected 2 elements in ReportHeader, got %d", len(report.ReportHeader.Elements))
	}

	text1, ok := report.ReportHeader.Elements[0].Element.(*structure.Text)
	if !ok {
		t.Fatalf("Expected first element to be Text")
	}
	if text1.Value != "{{companyName}}" {
		t.Errorf("Expected {{companyName}}, got %s", text1.Value)
	}
	if text1.Styles()[0] != "companyNameStyle" {
		t.Errorf("Expected companyNameStyle, got %s", text1.Styles()[0])
	}

	// Verify details
	if report.Details == nil || report.Details.Detail == nil {
		t.Fatalf("Details is nil")
	}

	group, ok := report.Details.Detail.(*structure.Group)
	if !ok {
		t.Fatalf("Expected detail to be Group")
	}

	if group.Header == nil || len(group.Header.Elements) != 1 {
		t.Fatalf("Group header incorrect")
	}

	ctext, ok := group.Header.Elements[0].Element.(*structure.CText)
	if !ok {
		t.Fatalf("Expected c-text in group header")
	}
	if ctext.Caption != "Sub company" || ctext.Value != "{{subCompany}}" {
		t.Errorf("CText incorrect: %v", ctext)
	}

	// Re-serialize back to XML
	b, err := xml.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatalf("Failed to remarshal Report: %v", err)
	}

	if len(b) == 0 {
		t.Errorf("Remarshaled XML is empty")
	}
}
