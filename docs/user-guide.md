# Squash User Guide

Comprehensive documentation for the Squash PDF report generation library.

[Full user guide content - this file contains all the detailed documentation from the comprehensive README]

# Squash — PDF Report Generation Library

[![Go Reference](https://pkg.go.dev/badge/github.com/niarbSoftHub/squash.svg)](https://pkg.go.dev/github.com/niarbSoftHub/squash)
[![Go Report Card](https://goreportcard.com/badge/github.com/niarbSoftHub/squash)](https://goreportcard.com/report/github.com/niarbSoftHub/squash)

**Squash** is a Go library for generating professional PDF reports with complex layouts, dynamic data rendering, and declarative styling. Define your report structure once as XML, provide data as JSON, and specify styles as JSON—Squash handles the rest.

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [Architecture](#architecture)
5. [Page Configuration](#page-configuration)
6. [Report Sections](#report-sections)
7. [Element Reference](#element-reference)
8. [Style System](#style-system)
9. [Data Model](#data-model)
10. [Font Handling](#font-handling)
11. [Automatic Behaviors](#automatic-behaviors)
12. [Known Limitations](#known-limitations)
13. [Complete Example](#complete-example)
14. [Troubleshooting](#troubleshooting)

---

## Introduction

### What is Squash?

Squash solves the problem of generating complex PDF reports without writing low-level layout code. Instead of programmatically building each page, you define your report declaratively using three separate inputs:

| Input | Format | Purpose |
|---|---|---|
| **Template** | XML | Describes the report layout (headers, groups, rows, sections) |
| **Data** | JSON | Provides the actual values to populate the template (hierarchical key-value pairs and tables) |
| **Styles** | JSON | Defines visual properties (fonts, colors, borders, alignment) that can be named and reused |

Squash merges these three inputs to produce a **single PDF file**.

### Key Features

- **Declarative Input**: Define structure, data, and styling separately for maximum maintainability
- **Type-Safe Element Rendering**: Dedicated components for numbers, money, dates, times, floats with automatic formatting
- **Flexible Layouts**: Horizontal and vertical layout containers with proportional column ratios
- **Professional Styling**: Borders, padding, background colors, fonts, alignment, and vertical alignment
- **Smart Style Hierarchy**: Three-level cascading defaults (renderer → stylesheet → element)
- **Data-Driven Repetition**: Automatically repeat structural groups for each data entry
- **Variable Substitution**: Use `{{variableName}}` placeholders to bind data to templates
- **PDF Output**: High-quality PDF generation using `gofpdf`

---

## Installation

```bash
go get github.com/niarbSoftHub/squash
```

### Dependencies

- `github.com/phpdave11/gofpdf v1.4.3` — PDF rendering engine
- `golang.org/x/image v0.40.0` — Required for WebP image support (indirect dependency)

---

## Quick Start

The typical workflow involves six steps:

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
	// 1. Load and unmarshal JSON data
	dataBytes, _ := os.ReadFile("data.json")
	var dataRoot data.DataRoot
	json.Unmarshal(dataBytes, &dataRoot)

	// 2. Load and unmarshal XML template
	templateBytes, _ := os.ReadFile("template.xml")
	var report structure.Report
	xml.Unmarshal(templateBytes, &report)

	// 3. Load and unmarshal JSON stylesheet
	styleBytes, _ := os.ReadFile("styles.json")
	var styleSheet style.StyleSheet
	json.Unmarshal(styleBytes, &styleSheet)

	// 4. Expand the template tree with data (resolve variables and repeat sections)
	expanded := engine.ExpandTree(&report, &dataRoot)

	// 5. Render the expanded tree to a PDF document
	generatedDoc, _ := renderer.Generate(expanded, &styleSheet)
	doc, _ := generatedDoc.Generate()

	// 6. Save the PDF to disk
	doc.Save("output.pdf")
}
```

---

## Architecture

Squash operates as a three-stage pipeline:

### Stage 1: Parse

The user loads three JSON/XML files into Go structs:
- XML template → `structure.Report`
- JSON data → `data.DataRoot`
- JSON styles → `style.StyleSheet`

### Stage 2: Expand

The `engine.ExpandTree()` function walks the template tree and binds data to it:
- Resolves all `{{variableName}}` placeholders with values from the data
- Repeats structural groups (like table rows) once for each data entry
- Produces an expanded tree ready for rendering

### Stage 3: Render

The `renderer.Generate()` function converts the expanded tree into a PDF:
- Traverses the element tree
- Looks up styles in the stylesheet and applies them
- Calls `gofpdf` APIs to draw text, cells, borders, images, etc.
- Manages pages, margins, headers, footers, and page breaks

### Style Hierarchy

Squash resolves styles at three levels, with each level overriding the previous:

```
┌─────────────────────────────────────┐
│  Renderer Defaults (font size 8pt,  │
│  text: left, money: right, etc.)    │
└──────────┬──────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│  Stylesheet Defaults                │
│  ("default-money-style", etc.)      │
└──────────┬──────────────────────────┘
           ↓
┌─────────────────────────────────────┐
│  Element Styles                     │
│  (styles="myCustomStyle", etc.)     │
└──────────┴──────────────────────────┘
           ↓
       Final Style
```

---

## Page Configuration

The `<report>` element is the root of your template and controls page-level settings.

### Attributes

| Attribute | Type | Default | Unit | Description |
|---|---|---|---|---|
| `orientation` | int | `0` | — | `0` = Portrait, `1` = Landscape |
| `pageSize` | int | `4` | — | Page size constant (see table below) |
| `marginLeft` | float32 | `10` | mm | Left margin |
| `marginRight` | float32 | `10` | mm | Right margin |
| `marginTop` | float32 | `10` | mm | Top margin |
| `marginBottom` | float32 | `10` | mm | Bottom margin |

### Page Size Reference

The `pageSize` attribute accepts integer constants representing ISO and North American paper sizes.

| Value | Name | Size (mm) | Value | Name | Size (mm) |
|---|---|---|---|---|---|
| 0 | A0 | 841×1189 | 14 | C0 | 1297×917 |
| 1 | A1 | 594×841 | 15 | C1 | 917×649 |
| 2 | A2 | 420×594 | 16 | C2 | 649×458 |
| 3 | A3 | 297×420 | 17 | C3 | 458×324 |
| **4** | **A4 (default)** | **210×297** | 18 | C4 | 324×229 |
| 5 | A5 | 148×210 | 19 | C5 | 229×162 |
| 6 | A6 | 105×148 | 20 | C6 | 162×114 |
| 7 | A7 | 74×105 | 21 | C7 | 114×81 |
| 8 | A8 | 52×74 | 22 | C8 | 81×57 |
| 9 | A9 | 37×52 | 23 | C9 | 57×40 |
| 10 | A10 | 26×37 | 24 | C10 | 40×28 |
| 11 | Letter | 216×279 | 25 | B0 | 1414×1000 |
| 12 | Legal | 216×356 | 26–35 | B1–B10 | (decreasing) |
| 13 | Tabloid | 279×432 | 36 | DL | 110×220 |

**Note:** When `pageSize="0"`, it defaults to A4 (not A0) due to Go's zero-value semantics. To use A0, `pageSize` must be explicitly set; this is a known quirk of the library.

### Example

```xml
<report orientation="1" pageSize="3" marginLeft="4" marginTop="4" marginRight="4" marginBottom="4">
  <!-- Content here -->
</report>
```

This creates a landscape A3 page with 4mm margins on all sides.

---

## Report Sections

A report can contain five top-level sections, each with distinct rendering behavior.

| Section | XML Element | Rendered | Purpose |
|---|---|---|---|
| Report Header | `<reportHeader>` | Once, at the top of page 1 | Title, company info, report metadata |
| Page Header | `<pageHeader>` | Top of every page (re-printed after page breaks) | Column headers, section titles |
| Details | `<details>` | Main body of the report | Data rows and groups |
| Report Footer | `<reportFooter>` | Once, at the end of the last page | Grand totals, final signatures |
| Page Footer | `<pageFooter>` | Bottom of every page | Signatures, page numbers, disclaimers |

### Basic Structure

```xml
<report orientation="0" pageSize="4" marginLeft="10" marginTop="10" marginRight="10" marginBottom="10">
  <reportHeader>
    <layout-v>
      <!-- Content rendered once at the start -->
    </layout-v>
  </reportHeader>

  <pageHeader>
    <layout-v>
      <!-- Content repeated at the top of every page -->
    </layout-v>
  </pageHeader>

  <details>
    <!-- Main data body: groups and/or rows -->
  </details>

  <reportFooter>
    <layout-v>
      <!-- Content rendered once at the end -->
    </layout-v>
  </reportFooter>

  <pageFooter>
    <layout-v>
      <!-- Content rendered at the bottom of every page -->
    </layout-v>
  </pageFooter>
</report>
```

---

## Element Reference

All elements support the optional attributes `styles="styleName"` and `font="fontName"` to apply styling and override the font.

### Text Elements

#### `<text>`

Renders plain text with automatic line wrapping.

**Attributes:**
- `styles` — Comma-separated style names (optional)
- `font` — Font family name (optional, overrides style font)

**Content:** Text with `{{variableName}}` placeholders

**Style Type:** `TextStyle`

**Example:**
```xml
<text styles="myTextStyle">Hello {{userName}}, welcome to {{companyName}}!</text>
```

---

#### `<c-text>` (Caption-Text)

Renders a label-value pair with a separator (e.g., "Department: HR").

**Attributes:**
- `caption` — The label text
- `separator` — Text between caption and value (default: `:`)
- `styles` — Style names (optional)
- `font` — Font family (optional)

**Content:** The value text (supports `{{variables}}`)

**Style Type:** `CTextStyle`

**Example:**
```xml
<c-text caption="Department" separator=":">{{department}}</c-text>
```

Renders as: `Department: Finance`

---

#### `<number>`

Renders an integer value (no decimal point).

**Attributes:**
- `styles` — Style names (optional)

**Content:** Numeric value (supports `{{variables}}`)

**Style Type:** `NumberStyle`

**Example:**
```xml
<number styles="countStyle">{{employee_count}}</number>
```

---

#### `<money>`

Renders a monetary value with currency symbol, decimal places, and thousands separator.

**Attributes:**
- `styles` — Style names (optional)

**Content:** Numeric value (supports `{{variables}}`)

**Style Type:** `MoneyStyle`

**Example:**
```xml
<money>{{salary}}</money>
```

With appropriate styling (`symbol: "৳"`, `decimals: 2`), renders as: `৳ 50,000.00`

---

#### `<date>`

Renders a date in a configurable format.

**Attributes:**
- `styles` — Style names (optional)

**Content:** Date value (supports `{{variables}}`)

**Style Type:** `DateStyle`

**Accepted Input Formats:**
- ISO 8601: `2026-05-21`
- Slash-separated: `21/05/2026`, `05/21/2026`, `2026/05/21`
- Hyphenated: `21-05-2026`

**Output Format Tokens:**
- `yyyy` — 4-digit year (e.g., 2026)
- `yy` — 2-digit year (e.g., 26)
- `MM` — 2-digit month (e.g., 05)
- `M` — Month without leading zero (e.g., 5)
- `dd` — 2-digit day (e.g., 21)
- `d` — Day without leading zero (e.g., 21)

**Default Format:** `dd/MM/yyyy`

**Example:**
```xml
<date styles="dateStyle">{{joining_date}}</date>
```

With format `"dd/MM/yyyy"`, input `2026-05-21` renders as: `21/05/2026`

---

#### `<time>`

Renders a time value in a configurable format.

**Attributes:**
- `styles` — Style names (optional)

**Content:** Time value (supports `{{variables}}`)

**Style Type:** `TimeStyle`

**Accepted Input Formats:**
- 24-hour: `14:30`, `14:30:45`
- 12-hour: `2:30 PM`, `2:30PM`

**Output Format Tokens:**
- `HH` — 2-digit hour (00–23)
- `H` — Hour without leading zero (0–23)
- `mm` — 2-digit minute (00–59)
- `m` — Minute without leading zero (0–59)
- `ss` — 2-digit second (00–59)
- `s` — Second without leading zero (0–59)

**Default Format:** `HH:mm`

**Example:**
```xml
<time styles="timeStyle">{{check_in_time}}</time>
```

With format `"HH:mm:ss"`, input `14:30:45` renders as: `14:30:45`

---

#### `<float>`

Renders a decimal number with configurable precision.

**Attributes:**
- `styles` — Style names (optional)

**Content:** Numeric value (supports `{{variables}}`)

**Style Type:** `FloatStyle`

**Example:**
```xml
<float styles="precisionStyle">{{overtime_hours}}</float>
```

With `decimals: 2`, value `7.5` renders as: `7.50`

---

### Media Elements

#### `<image>`

Renders an image (JPEG, PNG, GIF, or WebP).

**Attributes:**
- `width` — Width in millimeters (optional)
- `height` — Height in millimeters (optional)

**Content:** File path to the image (supports `{{variables}}`)

**Supported Formats:** JPEG, PNG, GIF, WebP

**Aspect Ratio:** If only `width` or `height` is specified, the other dimension is calculated to maintain aspect ratio. If neither is specified, the image is sized to fill the available width.

**Example:**
```xml
<image width="12" height="12">{{photo_path}}</image>
```

---

#### `<signatures>`

Renders a row of signature lines (space for handwritten signatures).

**Content:** One or more `<signature>` elements

**Example:**
```xml
<signatures>
  <signature>Prepared by</signature>
  <signature>Checked by</signature>
  <signature>Approved by</signature>
</signatures>
```

---

### Layout Elements

#### `<layout-v>` (Vertical Layout)

Stacks child elements vertically.

**Attributes:**
- `styles` — Style names (optional, applies VerticalLayoutStyle)

**Content:** Any elements (text, numbers, cells, nested layouts, etc.)

**Style Type:** `VerticalLayoutStyle`

**Example:**
```xml
<layout-v>
  <text>Header 1</text>
  <text>Header 2</text>
  <text>Header 3</text>
</layout-v>
```

---

#### `<layout-h>` (Horizontal Layout)

Arranges child cells side-by-side with proportional widths.

**Attributes:**
- `height` — Fixed height in millimeters; `0` = auto-measure (optional)
- `styles` — Style names (optional, applies HorizontalLayoutStyle)

**Content:** One or more `<cell>` elements

**Style Type:** `HorizontalLayoutStyle`

**Example:**
```xml
<layout-h>
  <cell ratio="25"><text>Column 1</text></cell>
  <cell ratio="75"><text>Column 2</text></cell>
</layout-h>
```

The ratios do not need to sum to 100; they are relative proportions. Here, column 1 is 25% and column 2 is 75% of the available width.

---

#### `<cell>`

A column container within a `<layout-h>`. Each cell can contain any single element.

**Attributes:**
- `ratio` — Proportional width relative to other cells in the same row (required)
- `styles` — Style names (optional, applies CellStyle)

**Content:** Any single element (text, number, money, image, or nested layout-v)

**Style Type:** `CellStyle`

**Example:**
```xml
<cell ratio="30" styles="dataCell">
  <text>Cell content</text>
</cell>
```

---

#### `<space>`

Renders vertical whitespace.

**Attributes:**
- `height` — Height in millimeters (default: 10)

**Example:**
```xml
<space height="5"/>
```

---

### Data-Driven Elements

#### `<group>`

A repeating section that is rendered once for each entry in a `GroupData` within the data model.

**Attributes:**
- `pageBreak` — `"true"` to force a page break before this group (only for groups after the first)

**Content:**
- `<header>` — Rendered at the top of each group iteration
- `<details>` — Contains either a nested `<group>` or a `<layout-h>` (the row template)
- `<footer>` — Rendered at the end of each group iteration

**Example:**
```xml
<group pageBreak="true">
  <header>
    <layout-v>
      <text>{{sectionName}}</text>
    </layout-v>
  </header>
  <details>
    <layout-h>
      <!-- Row template: repeated for each row in the group's table data -->
      <cell ratio="50"><text>{{col1}}</text></cell>
      <cell ratio="50"><text>{{col2}}</text></cell>
    </layout-h>
  </details>
  <footer>
    <layout-v>
      <text>Total: {{groupTotal}}</text>
    </layout-v>
  </footer>
</group>
```

---

## Style System

### Overview

Styles are JSON objects that define visual properties for elements. Styles are named and registered in a stylesheet, then applied to elements by referencing the style name.

### Style Declaration Forms

**Form 1: Nested under typed maps** (recommended for organization)

```json
{
  "textStyles": {
    "myTextStyle": {
      "styleType": "text-style",
      "fontSize": 12,
      "bold": true,
      "color": "#333333"
    }
  },
  "moneyStyles": {
    "boldMoney": {
      "styleType": "money-style",
      "bold": true,
      "color": "#FF0000"
    }
  }
}
```

**Form 2: Root-level with `styleType` discriminator** (shorthand)

```json
{
  "myTextStyle": {
    "styleType": "text-style",
    "fontSize": 12,
    "bold": true
  }
}
```

### Default Styles

Apply globally using the `"default-<type>-style"` naming pattern:

```json
{
  "default-text-style": {
    "styleType": "text-style",
    "fontSize": 8,
    "alignment": "left",
    "color": "#000000"
  },
  "default-money-style": {
    "styleType": "money-style",
    "fontSize": 7.5,
    "alignment": "right",
    "decimals": 2,
    "useComma": true,
    "symbol": "৳",
    "symbolPosition": "suffix"
  }
}
```

Default styles override renderer defaults but can be overridden by element-specific styles.

### Style Types and Properties

Please refer to the full user guide for complete documentation of all style types (TextStyle, NumberStyle, MoneyStyle, DateStyle, TimeStyle, FloatStyle, CellStyle, etc.).

---

## Data Model

Data is provided as JSON and unmarshaled into `data.DataRoot`.

### Variable Scoping

Variables are resolved in order of specificity:
1. **Row-level variables** (from `TableData.columns` matched to row values)
2. **Group-level variables** (from `GroupData.variableValuePair`)
3. **Report-level variables** (from `DataRoot.variableValuePair`)

---

## Font Handling

The library supports three built-in font families:
- `arial`, `helvetica`, `sans-serif` → Arial
- `times`, `times-new-roman`, `serif` → Times New Roman
- `courier`, `courier-new`, `monospace` → Courier

Custom TrueType fonts can be placed in the `fonts/` directory with `.ttf` extension.

---

## Automatic Behaviors

- **Auto page breaks**: Triggered automatically when content would overflow
- **Page numbers**: Automatically printed bottom-right of every page as `N / Total`
- **Group header repeat**: Group headers are re-printed at the top of new pages when a page break occurs mid-group
- **Footer safe zone**: Page footer height is measured before rendering to prevent overlap

---

## Known Limitations

Properties that are parsed but **not rendered**:
- `strikeThrough`, `superscript`, `subscript`, `rotation` on TextStyle
- Border width on CellStyle (all borders hardcoded to 0.2mm)
- Borders/background on HorizontalLayoutStyle and VerticalLayoutStyle
- `separatorBehavior` on CTextStyle

---

## Troubleshooting

### Variables Not Appearing
- Check variable name spelling and case
- Verify variable is in the correct scope
- Use a text editor to search in your JSON file

### Styles Not Applying
- Verify style name matches between XML and JSON
- Ensure `styleType` matches the element type
- Validate JSON syntax

### Fonts Not Found
- Place font files in the `fonts/` directory
- Use correct naming: `fontname.ttf`, `fontname-bold.ttf`, etc.
- Reference by base name in stylesheet

---

## License

Licensed under GNU AFFERO GENERAL PUBLIC LICENSE.

Although this library is licensed under GNU AGPL, custom licenses without charge are available for startups and growing companies upon request. The owner holds the right to deny such requests.

---

For the complete user guide with all detailed examples, property tables, and advanced topics, please refer to the full documentation.
