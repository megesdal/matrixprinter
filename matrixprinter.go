package matrixprinter

import (
	"fmt"
	"io"
)

type Table struct {
	buf [][]string

	ncols     int
	colwidths []int
	colaligns []int

	col  int
	line []string

	value string
}

func NewTable() *Table {
	return &Table{
		buf:       make([][]string, 0),
		ncols:     0,
		colwidths: make([]int, 0),
		colaligns: make([]int, 0),
		col:       0,
		line:      make([]string, 0),
	}
}

// AppendInt adds an int to the next column
func (t *Table) AppendInt(i int) *Table {
	return t.Append(fmt.Sprintf("%d", i))
}

// Append adds a string to the next column
func (t *Table) Append(s string) *Table {

	swidth := len(s)
	if t.col == len(t.colwidths) {
		t.colwidths = append(t.colwidths, swidth)
	} else if t.colwidths[t.col] < swidth {
		t.colwidths[t.col] = swidth
	}

	if t.col == len(t.colaligns) {
		t.colaligns = append(t.colaligns, 1)
	}

	if t.col == len(t.line) {
		t.line = append(t.line, s)
	} else {
		t.line[t.col] = s
	}

	t.col++
	return t
}

// ColLeft makes the given column left-aligned
func (t *Table) ColLeft(col int) *Table {

	for col > len(t.colaligns) {
		t.colaligns = append(t.colaligns, 1)
	}

	if col == len(t.colaligns) {
		t.colaligns = append(t.colaligns, -1)
	} else {
		t.colaligns[col] = -1
	}
	return t
}

// EndRow completes the current row and starts a new one
func (t *Table) EndRow() *Table {
	t.buf = append(t.buf, t.line)
	t.col = 0
	t.line = make([]string, len(t.colwidths))
	return t
}

// Print writes out the contents of the matrix using the supplied Writer
func (t *Table) Print(w io.Writer) {

	for i := 0; i < len(t.buf); i++ {
		line := t.buf[i]
		t.printRow(w, line)
		io.WriteString(w, "\n")
	}
	t.printRow(w, t.line)
	io.WriteString(w, "\n") //Do we want a newline here?
}

func (out *Table) printRow(w io.Writer, line []string) {

	first := false
	for j := 0; j < len(line); j++ {
		width := out.colwidths[j] * out.colaligns[j]
		if width != 0 {
			if first {
				first = false
			} else {
				io.WriteString(w, " ")
			}
			format := fmt.Sprintf("%%%ds", width)
			io.WriteString(w, fmt.Sprintf(format, line[j]))
		}
	}
}

// TODO: port the java impl...
/*
func (out *Table) sortBy(final int col, final int fromRow) {

    	List<List<String>> subBuf = buf.subList(fromRow, buf.size());
    	buf = buf.subList(0, fromRow);
    	Collections.sort(subBuf, new Comparator<List<String>>() {

			public int compare(List<String> a, List<String> b) {
				String aStr = a.get(col);
				String bStr = b.get(col);
				try {
					return Integer.valueOf(aStr).compareTo(Integer.valueOf(bStr));
				} catch (NumberFormatException ex) {
					return aStr.compareTo(bStr);
				}
			}

    	});
    	buf.addAll(subBuf);
    }
}
*/
