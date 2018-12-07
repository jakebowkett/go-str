/*
Package str provides string functions that handle
multi-byte characters correctly.
*/
package str

import (
	"fmt"
	"strings"
)

/*
Len returns the number of characters in a string rather
than the number of bytes.
*/
func Len(s string) int {
	return strings.Count(s, "") - 1
}

/*
Char returns character n (rather than byte n) in s. An error will
be returned if n is a negative number or if it exceeds the number
of characters in s.
*/
func Char(s string, i int) (string, error) {
	return Slice(s, i, i+1)
}

/*
Chars returns a slice of all the characters (rather than bytes) in s.
If s is an empty string the slice will be non-nil and zero length.
*/
func Chars(s string) []string {
	return strings.Split(s, "")
}

/*
CharSet returns a slice of all the characters (rather than bytes) in
s, excluding duplicates. If s is an empty string the slice will be
non-nil and zero length. Characters appear in the same order as they
do in s.
*/
func CharSet(s string) []string {
	return makeSet(strings.Split(s, ""))
}

func makeSet(ss []string) []string {

	set := make([]string, 0, len(ss))
	seen := make(map[string]bool, len(ss))

	for _, s := range ss {
		if seen[s] {
			continue
		}
		seen[s] = true
		set = append(set, s)
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

	cc := strings.SplitN(s, "", end+1)

	if end > len(cc) {
		return "", fmt.Errorf("End index is out of bounds. String contains %d characters, end was %d.", len(cc), end)
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
	ss := strings.SplitN(s, "", 2)
	ss[0] = strings.ToUpper(ss[0])
	return strings.Join(ss, "")
}

/*
Words returns the words in s as a slice of strings. Word
boundaries include any space character (as defined by Unicode),
forward slashes, endashes, and emdashes. In addition, grammatical
marks adjacent to word boundaries are omitted.

	ww := Words(`"Here's a sentence," said the narrator/programmer.`)
	// ww is []string{"Here's", "a", "sentence", "said", "the", "narrator", "programmer"}

*/
func Words(s string) []string {

	cc := strings.Split(s, "")

	// Approximate how long our words slice will need to be
	// to avoid repeated expansions.
	avgWordLen := 5.5
	words := make([]string, 0, int(float64(len(cc))/avgWordLen))

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

/*
WordSet is the same as Words but removes duplicates from
its results. See Words for more information on how it
determines what is a word.
*/
func WordSet(s string) []string {
	return makeSet(Words(s))
}

/*
WordOccurence returns a map whose keys are words in s and
whose values are the number of appearances in s. See Words
for more information on how it determines what is a word.
*/
func WordOccurence(s string) map[string]int {

	ww := Words(s)
	occurence := make(map[string]int, len(ww))

	for _, w := range ww {
		occurence[w]++
	}

	return occurence
}

/*
WordCount returns the number of words in s. See Words for
more information on how it determines what is a word.
*/
func WordCount(s string) int {

	cc := strings.Split(s, "")
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
	splitters := "–—/" // endash, emdash, and forward slash. Didn't mean to rhyme.
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
