package renderer_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/niarbSoftHub/squash/pkg/engine"
	"github.com/niarbSoftHub/squash/pkg/models/data"
	"github.com/niarbSoftHub/squash/pkg/models/structure"
	"github.com/niarbSoftHub/squash/pkg/models/style"
	"github.com/niarbSoftHub/squash/pkg/renderer"
)

func TestPayrollRender(t *testing.T) {
	// 1. Load static JSON data
	dataBytes, err := os.ReadFile("testdata/payroll_data.json")
	if err != nil {
		t.Fatalf("Failed to read payroll_data.json: %v", err)
	}
	var finalDataRoot data.DataRoot
	if err := json.Unmarshal(dataBytes, &finalDataRoot); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// 2. Load XML Template from file
	xmlDataBytes, err := os.ReadFile("testdata/payroll_template.xml")
	if err != nil {
		t.Fatalf("Failed to read XML template: %v", err)
	}
	var report structure.Report
	if err := xml.Unmarshal(xmlDataBytes, &report); err != nil {
		t.Fatalf("Failed to parse xml: %v", err)
	}

	// 3. Load Styles from JSON file
	styleBytes, err := os.ReadFile("testdata/payroll_styles.json")
	if err != nil {
		t.Fatalf("Failed to read styles: %v", err)
	}
	var styleSheet style.StyleSheet
	if err := json.Unmarshal(styleBytes, &styleSheet); err != nil {
		t.Fatalf("Failed to unmarshal styles: %v", err)
	}

	// 4. Expand Tree
	expanded := engine.ExpandTree(&report, &finalDataRoot)

	// 5. Render
	m, err := renderer.Generate(expanded, &styleSheet)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	doc, err := m.Generate()
	if err != nil {
		t.Fatalf("Maroto PDF generation failed: %v", err)
	}

	// 6. Write PDF
	err = doc.Save("payroll_output.pdf")
	if err != nil {
		t.Fatalf("Failed to save PDF: %v", err)
	}
}

func BenchmarkPayrollRender(b *testing.B) {
	dataBytes, _ := os.ReadFile("testdata/payroll_data.json")
	var finalDataRoot data.DataRoot
	_ = json.Unmarshal(dataBytes, &finalDataRoot)

	xmlDataBytes, _ := os.ReadFile("testdata/payroll_template.xml")
	var report structure.Report
	_ = xml.Unmarshal(xmlDataBytes, &report)

	styleBytes, _ := os.ReadFile("testdata/payroll_styles.json")
	var styleSheet style.StyleSheet
	_ = json.Unmarshal(styleBytes, &styleSheet)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expanded := engine.ExpandTree(&report, &finalDataRoot)
		m, _ := renderer.Generate(expanded, &styleSheet)
		_, _ = m.Generate()
	}
}

func TestPayrollMemoryUsage(t *testing.T) {
	// Force GC to get a clean baseline
	runtime.GC()
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	// Run full generation
	dataBytes, _ := os.ReadFile("testdata/payroll_data.json")
	var finalDataRoot data.DataRoot
	_ = json.Unmarshal(dataBytes, &finalDataRoot)

	xmlDataBytes, _ := os.ReadFile("testdata/payroll_template.xml")
	var report structure.Report
	_ = xml.Unmarshal(xmlDataBytes, &report)

	styleBytes, _ := os.ReadFile("testdata/payroll_styles.json")
	var styleSheet style.StyleSheet
	_ = json.Unmarshal(styleBytes, &styleSheet)

	expanded := engine.ExpandTree(&report, &finalDataRoot)
	m, _ := renderer.Generate(expanded, &styleSheet)
	_, _ = m.Generate()

	// Read peak stats immediately
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	fmt.Println("=============================================")
	fmt.Println("         ACTUAL RAM MEMORY STATS             ")
	fmt.Println("=============================================")
	fmt.Printf("Active Heap Memory:    %.2f MB (RAM actively in use)\n", float64(memAfter.Alloc)/1024/1024)
	fmt.Printf("OS Allocated Memory:   %.2f MB (Virtual memory footprint)\n", float64(memAfter.Sys)/1024/1024)
	fmt.Printf("Total Garbage Swept:   %.2f MB (Cleaned up automatically)\n", float64(memAfter.TotalAlloc-memBefore.TotalAlloc)/1024/1024)
	fmt.Println("=============================================")
}
