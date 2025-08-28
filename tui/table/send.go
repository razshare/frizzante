package table

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/razshare/frizzante/tui/wrap"
)

func Send(headers []string, rows [][]string, options ...Options) {
	if len(headers) == 0 || len(rows) == 0 {
		return
	}

	defaultOptions := DefaultOptions()
	if len(options) > 0 {
		defaultOptions = options[0]
	}

	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				if len(cell) > colWidths[i] {
					colWidths[i] = len(cell)
				}
			}
		}
	}

	for i := range colWidths {
		if colWidths[i] > defaultOptions.MaxColumnWidth {
			colWidths[i] = defaultOptions.MaxColumnWidth
		}
	}

	var processedRows [][]string
	logicalRowIndices := make([]int, 0)

	for rowIdx, row := range rows {
		maxLines := 1
		wrappedCells := make([][]string, len(headers))

		for i := 0; i < len(headers); i++ {
			cellContent := ""
			if i < len(row) {
				cellContent = row[i]
			}

			wrappedCells[i] = wrap.Send(cellContent, colWidths[i])
			if len(wrappedCells[i]) > maxLines {
				maxLines = len(wrappedCells[i])
			}
		}

		for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
			newRow := make([]string, len(headers))
			for cellIdx := 0; cellIdx < len(headers); cellIdx++ {
				if cellIdx < len(wrappedCells) && lineIdx < len(wrappedCells[cellIdx]) {
					newRow[cellIdx] = wrappedCells[cellIdx][lineIdx]
				}
			}
			processedRows = append(processedRows, newRow)
			logicalRowIndices = append(logicalRowIndices, rowIdx)
		}

		if rowIdx < len(rows)-1 {
			emptyRow := make([]string, len(headers))
			processedRows = append(processedRows, emptyRow)
			logicalRowIndices = append(logicalRowIndices, -1) // -1 indicates separator
		}
	}

	tbl := table.New().
		Headers(headers...).
		Rows(processedRows...).
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return defaultOptions.HeaderStyle
			}

			if row < len(logicalRowIndices) {
				logicalRow := logicalRowIndices[row]
				if logicalRow == -1 {
					return lipgloss.NewStyle()
				}

				if logicalRow%2 == 0 {
					return defaultOptions.RowStyle
				} else {
					return defaultOptions.AltRowStyle
				}
			}

			return lipgloss.NewStyle()
		})

	fmt.Println(tbl.Render())
}
