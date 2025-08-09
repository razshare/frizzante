package table

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"main/wrap"
)

const (
	MaxColumnWidth = 60
	ColumnPadding  = 2
	HeaderHeight   = 2
)

func Send(headers []string, rows [][]string) {
	if len(headers) == 0 || len(rows) == 0 {
		return
	}

	columns := make([]table.Column, len(headers))
	
	for index, header := range headers {
		width := len(header)
		
		for _, row := range rows {
			if index < len(row) {
				cellLen := len(row[index])
				if cellLen > width {
					width = cellLen
				}
			}
		}
		
		if width > MaxColumnWidth {
			width = MaxColumnWidth
		}
		
		columns[index] = table.Column{
			Title: header,
			Width: width + ColumnPadding,
		}
	}

	estimatedCapacity := len(rows) * 3
	wrappedRows := make([]table.Row, 0, estimatedCapacity)
	
	emptyRow := make([]string, len(headers))
	
	for rowIdx, row := range rows {
		if len(row) > len(columns) {
			row = row[:len(columns)]
		}
		
		maxLines := 1
		wrappedCells := make([][]string, len(columns))
		
		for i := 0; i < len(columns); i++ {
			cellContent := ""
			if i < len(row) {
				cellContent = row[i]
			}
			
			cellWidth := columns[i].Width - ColumnPadding
			wrappedCells[i] = wrap.Send(cellContent, cellWidth)
			
			if len(wrappedCells[i]) > maxLines {
				maxLines = len(wrappedCells[i])
			}
		}
		
		for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
			newRow := make([]string, len(columns))
			for cellIdx := 0; cellIdx < len(columns); cellIdx++ {
				if cellIdx < len(wrappedCells) && lineIdx < len(wrappedCells[cellIdx]) {
					newRow[cellIdx] = wrappedCells[cellIdx][lineIdx]
				}
			}
			wrappedRows = append(wrappedRows, newRow)
		}
		
		if rowIdx < len(rows)-1 {
			wrappedRows = append(wrappedRows, emptyRow)
		}
	}

	tableLocal := table.New(
		table.WithColumns(columns),
		table.WithRows(wrappedRows),
		table.WithFocused(false),
		table.WithHeight(len(wrappedRows)+HeaderHeight),
	)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	styles.Selected = lipgloss.NewStyle()
	tableLocal.SetStyles(styles)

	fmt.Println(tableLocal.View())
}
