/*
Package str provides string functions that handle
multi-byte characters correctly.
*/
package str

import (
	"fmt"
	"golang.org/x/text/unicode/norm"
	"strings"
)

/*
Len returns the number of characters in a string rather
than the number of bytes.
*/
func Len(s string) int {
	var it norm.Iter
	var n int
	it.InitString(norm.NFC, s)
	for !it.Done() {
		n++
		it.Next()
	}
	return n
}

/*
Char returns character n (rather than byte n) in s. An error will
be returned if n is a negative number or if it exceeds the number
of characters in s.
*/
func Char(s string, i int) (string, error) {
	return Slice(s, i, i+1)
}

func charsUntil(s string, maxBytePos int) []string {

	var it norm.Iter
	var n int

	it.InitString(norm.NFC, s)
	cc := make([]string, 0, maxBytePos)

	for !it.Done() {
		if n > maxBytePos {
			break
		}
		n++
		cc = append(cc, string(it.Next()))
	}

	return cc
}

/*
Chars returns a slice of all the characters (rather than bytes) in s.
If s is an empty string the slice will be non-nil and zero length.
*/
func Chars(s string) []string {
	return charsUntil(s, len(s))
}

/*
CharSet returns a slice of all the characters (rather than bytes) in
s, excluding duplicates. If s is an empty string the slice will be
non-nil and zero length. Characters appear in the same order they do
is s.
*/
func CharSet(s string) []string {

	cc := charsUntil(s, len(s))
	set := make([]string, 0, len(cc))
	seen := make(map[string]bool, len(cc))

	for _, c := range cc {
		if seen[c] {
			continue
		}
		seen[c] = true
		set = append(set, c)
	}

	return set
}

/*
Slice returns a substring of s. The start and end parameters refer
to character indices rather than byte indices.
*/
func Slice(s string, start, end int) (string, error) {

	if start < 0 || end < 0 {
		return "", fmt.Errorf("Negative indices are not allowed. Start was %d, end was %d.", start, end)
	}
	if start > end {
		return "", fmt.Errorf("Start index cannot be greater than end index. Start was %d, end was %d.", start, end)
	}

	cc := charsUntil(s, end)

	if end > len(cc) {
		return "", fmt.Errorf("End index is out of bounds. String contains %d characters, end index was %d.", len(cc), end)
	}

	return strings.Join(cc[start:end], ""), nil
}

/*
Capitalise returns a copy of s with its first character
converted to upper case if possible.
*/
func Capitalise(s string) string {

	if s == "" {
		return s
	}

	var it norm.Iter
	it.InitString(norm.NFC, s)

	first := string(it.Next())
	if it.Done() {
		return first
	}

	return strings.ToUpper(first) + s[it.Pos():]
}

/*
Words returns the words in s as a slice of strings. Word
boundaries include any space character (as defined by Unicode),
endashes, and emdashes. In addition, grammatical marks
adjacent to word boundaries are omitted.

For example, a call to Words with the string "hi\n\nthere"
would produce the slice []string{"hi", "there"}.
*/
func Words(s string) []string {

	cc := charsUntil(s, len(s))
	words := []string{}

	precededByBoundary := true
	idx := -1

	for i, c := range cc {
		if grammarOnBoundary(cc, i, precededByBoundary) {
			continue
		}
		if isBoundaryChar(c) {
			precededByBoundary = true
			continue
		}
		if precededByBoundary {
			words = append(words, c)
			precededByBoundary = false
			idx++
			continue
		}
		words[idx] += c
		precededByBoundary = false
	}

	return words
}

func grammarOnBoundary(cc []string, i int, precededByBoundary bool) bool {
	if !isGrammar(cc[i]) {
		return false
	}
	if !precededByBoundary && !boundaryNext(cc, i) && !grammarNext(cc, i) {
		return false
	}
	return true
}

func isGrammar(c string) bool {
	grammar := `!?,.'"[]()*~{}-<>`
	return strings.Contains(grammar, c)
}

func grammarNext(cc []string, i int) bool {
	i++
	if i == len(cc) {
		return false
	}
	return isGrammar(cc[i])
}

func isBoundaryChar(c string) bool {
	splitters := "–—" // endash and emdash.
	if strings.Contains(splitters, c) {
		return true
	}
	if strings.TrimSpace(c) == "" {
		return true
	}
	return false
}

func boundaryNext(cc []string, i int) bool {
	i++
	if i == len(cc) {
		return true
	}
	return isBoundaryChar(cc[i])
}

/*
WordCount returns the number of words in s. See Words for
more information on how it determines what is a word.
*/
func WordCount(s string) int {

	cc := charsUntil(s, len(s))
	precededByBoundary := true
	var count int

	for _, c := range cc {
		if isGrammar(c) {
			continue
		}
		if isBoundaryChar(c) {
			precededByBoundary = true
			continue
		}
		if precededByBoundary {
			count++
			precededByBoundary = false
		}
	}

	return count
}
