package table

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"main/wrap"
)

func Send(headers []string, rows [][]string) {
	columns := make([]table.Column, len(headers))
	maxColWidth := 60

	for index, header := range headers {
		width := len(header)
		for _, row := range rows {
			if index < len(row) && len(row[index]) > width {
				width = len(row[index])
			}
		}
		if width > maxColWidth {
			width = maxColWidth
		}
		columns[index] = table.Column{Title: header, Width: width + 2}
	}

	wrappedRows := make([]table.Row, 0)
	for rowIdx, row := range rows {
		maxLines := 1
		wrappedCells := make([][]string, len(row))

		for i, cell := range row {
			if i < len(columns) {
				cellWidth := columns[i].Width - 2
				wrappedCells[i] = wrap.Send(cell, cellWidth)
				if len(wrappedCells[i]) > maxLines {
					maxLines = len(wrappedCells[i])
				}
			}
		}

		for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
			newRow := make([]string, len(row))
			for cellIdx, wrappedCell := range wrappedCells {
				if lineIdx < len(wrappedCell) {
					newRow[cellIdx] = wrappedCell[lineIdx]
				} else {
					newRow[cellIdx] = ""
				}
			}
			wrappedRows = append(wrappedRows, newRow)
		}

		// Add empty row after each logical row (except the last one)
		if rowIdx < len(rows)-1 {
			emptyRow := make([]string, len(row))
			for index := range emptyRow {
				emptyRow[index] = ""
			}
			wrappedRows = append(wrappedRows, emptyRow)
		}
	}

	tableLocal := table.New(
		table.WithColumns(columns),
		table.WithRows(wrappedRows),
		table.WithFocused(false),
		table.WithHeight(len(wrappedRows)+2),
	)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	// Remove selection highlight
	styles.Selected = lipgloss.NewStyle()
	tableLocal.SetStyles(styles)

	fmt.Println(tableLocal.View())
}
