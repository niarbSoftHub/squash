# Squash — PDF Report Generation Library

[![Go Reference](https://pkg.go.dev/badge/github.com/niarbSoftHub/squash.svg)](https://pkg.go.dev/github.com/niarbSoftHub/squash)
[![Go Report Card](https://goreportcard.com/badge/github.com/niarbSoftHub/squash)](https://goreportcard.com/report/github.com/niarbSoftHub/squash)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

**Squash** is a Go library for generating professional PDF reports with complex layouts, dynamic data rendering, and declarative styling. Define your report structure once as XML, provide data as JSON, and specify styles as JSON—Squash handles the rest.

## Quick Overview

Generate PDF reports from three declarative inputs:

| Input | Format | Purpose |
|---|---|---|
| **Template** | XML | Report layout (headers, groups, rows, sections) |
| **Data** | JSON | Actual values to populate the template |
| **Styles** | JSON | Visual properties (fonts, colors, borders, alignment) |

```go
package main

import (
	"encoding/json"
	"encoding/xml"
	"os"

	"github.com/niarbSoftHub/squash/pkg/engine"
	"github.com/niarbSoftHub/squash/pkg/models/data"
	"github.com/niarbSoftHub/squash/pkg/models/structure"
	"github.com/niarbSoftHub/squash/pkg/models/style"
	"github.com/niarbSoftHub/squash/pkg/renderer"
)

func main() {
	// Load data
	dataBytes, _ := os.ReadFile("data.json")
	var dataRoot data.DataRoot
	json.Unmarshal(dataBytes, &dataRoot)

	// Load template
	templateBytes, _ := os.ReadFile("template.xml")
	var report structure.Report
	xml.Unmarshal(templateBytes, &report)

	// Load styles
	styleBytes, _ := os.ReadFile("styles.json")
	var styleSheet style.StyleSheet
	json.Unmarshal(styleBytes, &styleSheet)

	// Generate PDF
	expanded := engine.ExpandTree(&report, &dataRoot)
	doc, _ := renderer.Generate(expanded, &styleSheet)
	result, _ := doc.Generate()
	result.Save("output.pdf")
}
```

## Features

- ✅ **Declarative XML templates** for report structure
- ✅ **Type-safe elements** (text, number, money, date, time, float, image)
- ✅ **Flexible layouts** with proportional column widths
- ✅ **Professional styling** (borders, padding, colors, fonts, alignment)
- ✅ **Data-driven repetition** (groups and rows)
- ✅ **Variable substitution** with `{{variableName}}` syntax
- ✅ **Auto page breaks** and group header repetition
- ✅ **37 page sizes** (A, B, C series + Letter, Legal, Tabloid, DL)
- ✅ **Custom fonts** (TTF files)

## Installation

```bash
go get github.com/niarbSoftHub/squash
```

## Documentation

📖 **[Full User Guide](./docs/user-guide.md)**

The user guide covers:
- Element reference (14 element types)
- Style system (10 style types with complete property tables)
- Data model and variable scoping
- Page configuration and all 37 page sizes
- Font handling (built-in and custom fonts)
- Automatic behaviors (page breaks, page numbers, etc.)
- Known limitations
- Complete working examples
- Troubleshooting guide

## Support & Feedback

### Found a Bug?

Please create an issue on the project repository:

1. **Go to**: [GitHub Issues](https://github.com/niarbSoftHub/squash/issues)
2. **Click**: "New Issue"
3. **Describe**:
   - What you were trying to do
   - What happened
   - What you expected
   - Steps to reproduce (with code if possible)
   - Your Go version (`go version`)

### Have a Feature Request?

Open an issue with the tag `[Feature Request]` and describe:
- What functionality you need
- Why you need it
- A concrete use case

## License

**GNU AFFERO GENERAL PUBLIC LICENSE v3**

See the [LICENSE](./LICENSE) file for details.

### Custom Licensing

Although this library is released under GNU AGPL, **custom licenses without charge are available for startups and growing companies** upon request.

**To request a custom license:**
- Contact the project maintainer
- Provide details about your company (name, industry, stage)
- Describe your use case

**Note:** The maintainer reserves the right to deny custom license requests.

## Dependencies

- `github.com/phpdave11/gofpdf v1.4.3` — PDF rendering engine
- `golang.org/x/image v0.40.0` — WebP image support

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on [GitHub](https://github.com/niarbSoftHub/squash).

---

**Repository:** [github.com/niarbSoftHub/squash](https://github.com/niarbSoftHub/squash)
