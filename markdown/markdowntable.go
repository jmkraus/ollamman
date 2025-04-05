package markdown

// Creates a markdown table from a list of lists

import (
	"fmt"
	"strings"
)

func MarkdownTable(tabledata [][]string) string {

	numColumns := len(tabledata[0])
	border := "|"
	for range numColumns {
		border += " %s |"
	}

	// define variables
	maxWidths := make([]int, numColumns)
	fixWidths := make([]int, numColumns)
	separator := make([]string, numColumns)
	for i := range numColumns {
		separator[i] = "---"
	}

	// create heading separator row
	tabledata = append(tabledata, []string{})
	copy(tabledata[1:], tabledata)
	tabledata[1] = separator

	// identify fix widths (heading is prefixed with a "+")
	for idx, col := range tabledata[0] {
		if strings.HasPrefix(col, "+") {
			fixWidths[idx] = len(col) - 1
			tabledata[0][idx] = col[1:]
		} else {
			fixWidths[idx] = 0
		}
	}

	// identify max widths
	for _, row := range tabledata {
		for cIdx, col := range row {
			if len(col) > maxWidths[cIdx] {
				maxWidths[cIdx] = len(col)
			}
			if fixWidths[cIdx] != 0 {
				maxWidths[cIdx] = fixWidths[cIdx]
			}
		}
	}

	for rIdx, row := range tabledata {
		fill := " "
		if rIdx == 1 {
			fill = "-"
		}
		for cIdx, col := range row {
			tabledata[rIdx][cIdx] = fmt.Sprintf("%-"+fmt.Sprint(maxWidths[cIdx])+"s", col)
			tabledata[rIdx][cIdx] = strings.Replace(tabledata[rIdx][cIdx], " ", fill, -1)
		}
	}

	rows := make([]string, len(tabledata))
	for i, row := range tabledata {
		rowValues := make([]any, len(row))
		for j, v := range row {
			rowValues[j] = v
		}
		rows[i] = fmt.Sprintf(border, rowValues...)
	}

	table := strings.Join(rows, "\n")
	return table
}
