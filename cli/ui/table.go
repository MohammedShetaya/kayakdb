package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

// Table represents a table with columns and rows
type Table struct {
	columns []string
	rows    [][]string
	title   string
}

// NewTable creates a new table with the specified columns
func NewTable(columns ...string) *Table {
	return &Table{
		columns: columns,
		rows:    make([][]string, 0),
	}
}

// WithTitle sets a title for the table
func (t *Table) WithTitle(title string) *Table {
	t.title = title
	return t
}

// AddRow adds a row to the table
func (t *Table) AddRow(values ...string) *Table {
	// Ensure we have the right number of columns
	if len(values) != len(t.columns) {
		// Pad with empty strings if too few values
		for len(values) < len(t.columns) {
			values = append(values, "")
		}
		// Truncate if too many values
		if len(values) > len(t.columns) {
			values = values[:len(t.columns)]
		}
	}
	t.rows = append(t.rows, values)
	return t
}

// AddRows adds multiple rows to the table
func (t *Table) AddRows(rows [][]string) *Table {
	for _, row := range rows {
		t.AddRow(row...)
	}
	return t
}

// calculateColumnWidths determines the maximum width needed for each column
func (t *Table) calculateColumnWidths() []int {
	if len(t.columns) == 0 {
		return []int{}
	}

	widths := make([]int, len(t.columns))

	// Start with header widths
	for i, header := range t.columns {
		widths[i] = lipgloss.Width(header)
	}

	// Check each row for wider content
	for _, row := range t.rows {
		for i, cell := range row {
			if i < len(widths) {
				cellWidth := lipgloss.Width(cell)
				if cellWidth > widths[i] {
					widths[i] = cellWidth
				}
			}
		}
	}

	return widths
}

// renderRow renders a single row with proper alignment and styling
func (t *Table) renderRow(values []string, widths []int, isHeader bool) string {
	if len(values) == 0 {
		return ""
	}

	cells := make([]string, len(values))

	for i, value := range values {
		width := widths[i]
		if i < len(widths) {
			// Pad the content to the column width
			paddedValue := value + strings.Repeat(" ", width-lipgloss.Width(value))

			if isHeader {
				cells[i] = tableHeaderStyle.Width(width).Render(paddedValue)
			} else {
				cells[i] = tableCellStyle.Width(width).Render(paddedValue)
			}
		}
	}

	return strings.Join(cells, "")
}

// Print renders and prints the table
func (t *Table) Print() {
	if len(t.columns) == 0 {
		fmt.Println(Error("Table Error", "No columns defined").WithDetails("Call NewTable() with column names"))
		return
	}

	fmt.Println()

	// Print title if provided
	if t.title != "" {
		fmt.Println(titleStyle.Render(fmt.Sprintf("📊 %s", t.title)))
		fmt.Println()
	}

	widths := t.calculateColumnWidths()

	var content strings.Builder

	// Render header
	headerRow := t.renderRow(t.columns, widths, true)
	content.WriteString(headerRow)
	content.WriteString("\n")

	// Add separator line
	separator := make([]string, len(t.columns))
	for i, width := range widths {
		separator[i] = strings.Repeat("─", width+2) // +2 for padding
	}
	content.WriteString(lipgloss.NewStyle().Foreground(colorTable).Render(strings.Join(separator, "")))
	content.WriteString("\n")

	// Render data rows
	for _, row := range t.rows {
		dataRow := t.renderRow(row, widths, false)
		content.WriteString(dataRow)
		content.WriteString("\n")
	}

	// Apply table styling and print
	fmt.Println(tableStyle.Render(strings.TrimSpace(content.String())))
	fmt.Println()
}

// String returns the string representation of the table
func (t *Table) String() string {
	if len(t.columns) == 0 {
		return "Empty table - no columns defined"
	}

	widths := t.calculateColumnWidths()
	var result strings.Builder

	// Header
	headerRow := t.renderRow(t.columns, widths, true)
	result.WriteString(headerRow)
	result.WriteString("\n")

	// Separator
	separator := make([]string, len(t.columns))
	for i, width := range widths {
		separator[i] = strings.Repeat("─", width+2)
	}
	result.WriteString(strings.Join(separator, ""))
	result.WriteString("\n")

	// Data rows
	for _, row := range t.rows {
		dataRow := t.renderRow(row, widths, false)
		result.WriteString(dataRow)
		result.WriteString("\n")
	}

	return tableStyle.Render(strings.TrimSpace(result.String()))
}

// PrintSimple prints a simple table without styling (useful for plain text output)
func PrintSimpleTable(columns []string, rows [][]string) {
	table := NewTable(columns...)
	table.AddRows(rows)

	// Simple text-only version
	widths := table.calculateColumnWidths()

	fmt.Println()

	// Print header
	for i, col := range columns {
		fmt.Printf("%-*s", widths[i]+2, col)
	}
	fmt.Println()

	// Print separator
	for _, width := range widths {
		fmt.Print(strings.Repeat("-", width+2))
	}
	fmt.Println()

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				fmt.Printf("%-*s", widths[i]+2, cell)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
