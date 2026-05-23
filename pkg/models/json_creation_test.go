package models_test

import (
	"encoding/json"
	"testing"

	"github.com/niarbSoftHub/squash/pkg/interfaces"
	"github.com/niarbSoftHub/squash/pkg/models/data"
)

func TestDataModelJson(t *testing.T) {
	dataRoot := data.DataRoot{
		VariableValuePair: map[string]string{
			"companyName": "The Boring Company",
			"reportName":  "Financial Report-2019",
			"dateRange":   "1 June 2018 – 31 May 2019",
		},
		Details: []interfaces.DataDetail{},
	}

	company1 := data.GroupData{
		DataRoot: data.DataRoot{
			VariableValuePair: map[string]string{
				"subCompany": "Company A",
				"region":     "Europe",
			},
		},
	}

	companyDepartment11 := data.GroupData{
		DataRoot: data.DataRoot{
			VariableValuePair: map[string]string{
				"department":    "HR",
				"employeeCount": "512",
			},
		},
	}

	departmentTable11 := data.TableData{
		Columns: []string{"employeeId", "employeeName", "joiningDate", "salary"},
		DataRows: [][]string{
			{"1", "Mr. Baker", "1/11/2018", "50000"},
			{"2", "Mr. Lipton", "5/07/2018", "40000"},
		},
	}

	companyDepartment11.Details = append(companyDepartment11.Details, departmentTable11)
	company1.Details = append(company1.Details, companyDepartment11)
	dataRoot.Details = append(dataRoot.Details, company1)

	// Serialize
	b, err := json.Marshal(dataRoot)
	if err != nil {
		t.Fatalf("Failed to marshal DataRoot: %v", err)
	}

	// Deserialize
	var decoded data.DataRoot
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal DataRoot: %v", err)
	}

	b2, err := json.Marshal(decoded)
	if err != nil {
		t.Fatalf("Failed to remarshal DataRoot: %v", err)
	}

	if string(b) != string(b2) {
		t.Errorf("Serialization mismatch.\nExpected: %s\nGot: %s", string(b), string(b2))
	}
}
