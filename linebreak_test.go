package linebreak_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dgryski/go-linebreak"
)

func ExampleWrap() {
	text := "a b c d e f g h i j k l m n o p qqqqqqqqq"
	width := 9
	textWrapped := linebreak.Wrap(text, width, width)

	fmt.Println(textWrapped)

	// Output:
	// a b c d
	// e f g h
	// i j k l
	// m n o p
	// qqqqqqqqq
}

type testCase struct {
	text     string
	width    int
	maxwidth int
	want     string
}

var coreTestCases = []testCase{
	{
		text:  "",
		width: 80,
		want:  "",
	},
	{
		text:  "hello, world.",
		width: 80,
		want:  "hello, world.",
	},
	{
		text:  "hello,\nworld.",
		width: 80,
		want:  "hello, world.",
	},
	{
		text:  "hello, world.",
		width: 6,
		want:  "hello,\nworld.",
	},
	{
		text:  "aaaaaa b cccccc d",
		width: 4,
		want:  "aaaaaa\nb\ncccccc\nd",
	},
	{
		text:  "a b c d e f g h i j k l m n o p qqqqqqqqq",
		width: 16,
		want: strings.Join([]string{
			"a b c d e f g h",
			"i j k l m n o p",
			"qqqqqqqqq",
		}, "\n"),
	},
}

func TestWrap(t *testing.T) {
	cases := []testCase{
		{
			text:  "a b c d e f g h i j k l m n o p qqqqqqqqq",
			width: 9,
			want: strings.Join([]string{
				"a b c d",
				"e f g h",
				"i j k l",
				"m n o p",
				"qqqqqqqqq",
			}, "\n"),
		},
	}
	for _, c := range append(cases, coreTestCases...) {
		if c.maxwidth == 0 {
			c.maxwidth = c.width
		}
		t.Run(fmt.Sprintf("wrap %q width %d %d", c.text, c.width, c.maxwidth), func(t *testing.T) {
			if got := linebreak.Wrap(c.text, c.width, c.maxwidth); got != c.want {
				t.Errorf("got:\n%s\nwant:\n%s\n", got, c.want)
			}
		})
	}
}

func TestGreedy(t *testing.T) {
	cases := []testCase{
		{
			text:  "a b c d e f g h i j k l m n o p qqqqqqqqq",
			width: 9,
			want: strings.Join([]string{
				"a b c d e",
				"f g h i j",
				"k l m n o",
				"p",
				"qqqqqqqqq",
			}, "\n"),
		},
	}
	for _, c := range append(cases, coreTestCases...) {
		if c.maxwidth == 0 {
			c.maxwidth = c.width
		}
		t.Run(fmt.Sprintf("wrap %q width %d %d", c.text, c.width, c.maxwidth), func(t *testing.T) {
			if got := linebreak.Greedy(c.text, c.width, c.maxwidth); got != c.want {
				t.Errorf("got:\n%s\nwant:\n%s\n", got, c.want)
			}
		})
	}
}
