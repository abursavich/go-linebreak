// Package linebreak wraps text at a given width.
package linebreak

import "strings"

// This code is a translation of `linear` from http://xxyxyz.org/line-breaking/
// https://en.wikipedia.org/wiki/SMAWK_algorithm
// A. Aggarwal, T. Tokuyama. Consecutive interval query and dynamic programming on intervals. Discrete Applied Mathematics 85, 1998.

// Wrap formats text at the given width in linear time.
func Wrap(text string, width, maxwidth int) string {
	words := strings.Fields(text)
	count := len(words)
	if count == 0 {
		return ""
	}
	offsets := make([]int, count+1)
	for i, w := range words {
		offsets[i+1] = offsets[i] + len(w)
	}

	minima := make([]int64, count+1)
	for i := 1; i < len(minima); i++ {
		minima[i] = 1000000000000000000
	}
	breaks := make([]int, count+1)

	// closes over offsets, minima
	cost := func(i, j int) int64 {
		w := offsets[j] - offsets[i] + j - i - 1
		if w > maxwidth {
			return 10000000000 * int64(w-width)
		}
		d := abs(width - w)
		// last line has smaller extra space penalty
		if j == count {
			return minima[i] + int64(d*d)
		}
		return minima[i] + int64(d*d*d)
	}

	var smawk func([]int, []int)
	// smawk closes over cost, minima, breaks
	smawk = func(rows, columns []int) {
		var stack []int
		i := 0
		for i < len(rows) {
			if len(stack) > 0 {
				c := columns[len(stack)-1]
				if cost(peek(stack), c) < cost(rows[i], c) {
					if len(stack) < len(columns) {
						stack = push(stack, rows[i])
					}
					i++
				} else {
					stack = pop(stack)
				}
			} else {
				stack = push(stack, rows[i])
				i++
			}
		}
		rows = stack

		if len(columns) > 1 {
			smawk(rows, step(columns[1:], 2))
		}

		i = 0
		var j int
		for j < len(columns) {
			var end int
			if j+1 < len(columns) {
				end = breaks[columns[j+1]]
			} else {
				end = rows[len(rows)-1]
			}
			c := cost(rows[i], columns[j])
			if c < minima[columns[j]] {
				minima[columns[j]] = c
				breaks[columns[j]] = rows[i]
			}
			if rows[i] < end {
				i++
			} else {
				j += 2
			}
		}
	}

	n := count + 1
	i := 0
	offset := 0
	var r1 []int
	var r2 []int
	for {
		r := min(n, 1<<uint(i+1))
		edge := (1 << uint(i)) + offset
		r1 = genrange(r1, 0+offset, edge)
		r2 = genrange(r2, edge, r+offset)
		smawk(r1, r2)
		x := minima[r-1+offset]
		// because python code has 'for ... else'
		var terminatedFor bool
		for j := 1 << uint(i); j < r-1; j++ {
			y := cost(j+offset, r-1+offset)
			if y <= x {
				n -= j
				i = 0
				offset += j
				terminatedFor = true
				break
			}
		}
		if !terminatedFor {
			if r == n {
				break
			}
			i++
		}
	}

	// Work backwards finding line breaks.
	for x, j := 0, count; ; x++ {
		breaks[count-x], j = j, breaks[j]
		if j <= 0 {
			breaks = breaks[count-x:]
			break
		}
	}

	// Work forwards writing lines.
	var b strings.Builder
	b.Grow(offsets[count] + count - 1)
	i = 0
	for _, j := range breaks {
		appendLine(&b, words[i:j])
		i = j
	}
	return b.String()
}

// Greedy formats text at the given width greedily.
func Greedy(text string, width, maxwidth int) string {
	fields := strings.Fields(text)
	count := len(fields)
	if count == 0 {
		return ""
	}

	var b strings.Builder
	b.Grow(len(text))

	i := 0
	lineWidth := len(fields[i])
	for k := 1; k < count; k++ {
		fieldWidth := len(fields[k])
		nextWidth := lineWidth + 1 + fieldWidth

		switch {
		// Append the existing line if adding the field will go
		// over the max width OR farther from the target width.
		case nextWidth > maxwidth || abs(width-lineWidth) < abs(width-nextWidth):
			appendLine(&b, fields[i:k])
			i = k
			lineWidth = fieldWidth
		default:
			lineWidth = nextWidth
		}
	}
	appendLine(&b, fields[i:])

	return b.String()
}

func appendLine(b *strings.Builder, fields []string) {
	if len(fields) == 0 {
		return
	}
	if b.Len() > 0 {
		b.WriteRune('\n')
	}
	b.WriteString(fields[0])
	for _, field := range fields[1:] {
		b.WriteRune(' ')
		b.WriteString(field)
	}
}

// trivial int stack
func push(s []int, i int) []int { return append(s, i) }
func pop(s []int) []int         { return s[:len(s)-1] }
func peek(s []int) int          { return s[len(s)-1] }

// python list[a::b]
func step(ints []int, step int) []int {
	r := make([]int, 0, 1+(len(ints)/step))
	for i := 0; i < len(ints); i += step {
		r = append(r, ints[i])
	}
	return r
}

// python range(a,b)
func genrange(r []int, start, stop int) []int {
	if r != nil {
		r = r[:0]
	}
	if stop <= start {
		return r
	}
	for i := start; i < stop; i++ {
		r = append(r, i)
	}
	return r
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
