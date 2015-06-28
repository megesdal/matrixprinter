package matrixprinter 

import (
	"fmt"
	"io"
)

type MatrixPrinter struct {
	buf [][]string

	ncols     int
	colwidths []int
	colaligns []int

	col  int
	line []string

	value string
}

func New() *MatrixPrinter {
	return &MatrixPrinter{
		buf:       make([][]string, 0),
		ncols:     0,
		colwidths: make([]int, 0),
		colaligns: make([]int, 0),
		col:       0,
		line:      make([]string, 0),
	}
}

// AppendInt adds an int to the next column
func (out *MatrixPrinter) AppendInt(i int) *MatrixPrinter {
	return out.Append(fmt.Sprintf("%d", i))
}

// Append adds a string to the next column
func (out *MatrixPrinter) Append(s string) *MatrixPrinter {

	swidth := len(s)
	if out.col == len(out.colwidths) {
		out.colwidths = append(out.colwidths, swidth)
	} else if out.colwidths[out.col] < swidth {
		out.colwidths[out.col] = swidth
	}

	if out.col == len(out.colaligns) {
		out.colaligns = append(out.colaligns, 1)
	}

	if out.col == len(out.line) {
		out.line = append(out.line, s)
	} else {
		out.line[out.col] = s
	}

	out.col++
	return out
}

// ColLeft makes the given column left-aligned
func (out *MatrixPrinter) ColLeft(col int) *MatrixPrinter {

	for col > len(out.colaligns) {
		out.colaligns = append(out.colaligns, 1)
	}

	if col == len(out.colaligns) {
		out.colaligns = append(out.colaligns, -1)
	} else {
		out.colaligns[col] = -1
	}
	return out
}

// EndRow completes the current row and starts a new one
func (out *MatrixPrinter) EndRow() *MatrixPrinter {
	out.buf = append(out.buf, out.line)
	out.col = 0
	out.line = make([]string, len(out.colwidths))
	return out
}

// Print writes out the contents of the matrix using the supplied Writer
func (out *MatrixPrinter) Print(w io.Writer) {

	for i := 0; i < len(out.buf); i++ {
		line := out.buf[i]
		out.printRow(w, line)
		io.WriteString(w, "\n")
	}
	out.printRow(w, out.line)
	io.WriteString(w, "\n") //Do we want a newline here?
}

func (out *MatrixPrinter) printRow(w io.Writer, line []string) {

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
func (out *MatrixPrinter) sortBy(final int col, final int fromRow) {

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
