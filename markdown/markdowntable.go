package markdown

// Creates a markdown table from a list of lists

import (
	"fmt"
	"strings"
)

func MarkdownTable(tabledata [][]string) string {

	numColumns := len(tabledata[0])
	reduceColNum := numColumns - 1 //reduce width of this col by -1
	border := "|"
	for range numColumns {
		border += " %s |"
	}

	maxWidths := make([]int, numColumns)
	separator := make([]string, numColumns)
	for i := range numColumns {
		separator[i] = "---"
	}

	tabledata = append(tabledata, []string{})
	copy(tabledata[1:], tabledata)
	tabledata[1] = separator

	for _, row := range tabledata {
		for cIdx, col := range row {
			if len(col) > maxWidths[cIdx] {
				maxWidths[cIdx] = len(col)
			}
		}
	}
	maxWidths[reduceColNum]--

	for rIdx, row := range tabledata {
		fill := " "
		if rIdx == 1 {
			fill = "-"
		}
		for cIdx, col := range row {
			mw := maxWidths[cIdx]
			if rIdx == 1 && cIdx == reduceColNum {
				mw++
			}
			tabledata[rIdx][cIdx] = fmt.Sprintf("%-"+fmt.Sprint(mw)+"s", col)
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
