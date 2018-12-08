/*
Package str provides functions that handle single-
and multi-byte character strings in a convenient way.
*/
package str

import (
	"fmt"
	"sort"
	"strings"
)

/*
PadLeft prefixes s with padChar until s contains length number
of characters.
*/
func PadLeft(s string, padChar rune, length int) string {
	diff := length - len([]rune(s))
	if diff <= 0 {
		return s
	}
	return s + strings.Repeat(string(padChar), diff)
}

/*
PadLeft suffixes s with padChar until s contains length number
of characters.
*/
func PadRight(s string, padChar rune, length int) string {
	diff := length - len([]rune(s))
	if diff <= 0 {
		return s
	}
	return strings.Repeat(string(padChar), diff) + s
}

/*
PadToLongest suffixes each string in ss with padChar until it
contains as many characters as the longest string in ss.
*/
func PadToLongest(ss []string, padChar rune) {
	var longest int
	for i, _ := range ss {
		length := len([]rune(s[i]))
		if length > longest {
			longest = length
		}
	}
	for i, _ := range ss {
		diff := longest - len([]rune(ss[i]))
		ss[i] += strings.Repeat(string(padChar), diff)
	}
	return ss
}

/*
Len returns the number of characters in a string rather
than the number of bytes.
*/
func Len(s string) int {
	return len([]rune(s))
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
	return makeSet(strings.Split(s, ""), false)
}

/*
CharSetFold is the same as CharSet except it considers characters of
different cases equivalent. For example, "a" == "A". The strings in
the resulting slice are all lowercase.
*/
func CharSetFold(s string) []string {
	return makeSet(strings.Split(s, ""), true)
}

func makeSet(ss []string, fold bool) []string {

	set := make([]string, 0, len(ss))
	seen := make(map[string]bool, len(ss))

	for _, s := range ss {
		if fold {
			s = strings.ToLower(s)
		}
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
to character indices rather than byte indices. Both start and end
may be negative, in which case they refer to the offset from the
end of s. If start is greater than end it will wrap to the beginning
of s.

Returns an error if start or end have an absolute value greater
than the number of characters in s.
*/
func Slice(s string, start, end int) (string, error) {
	cc := []rune(s)
	if abs(start) > len(cc) || abs(end) > len(cc) {
		return "", fmt.Errorf("index out of bounds")
	}
	if start < 0 {
		start = len(cc) + start
	}
	if end < 0 {
		end = len(cc) + end
	}
	if start > end {
		return string(cc[start:]) + string(cc[0:end]), nil
	}
	return string(cc[start:end]), nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

/*
Capitalise returns a copy of s with its first character
converted to upper case if possible.
*/
func Capitalise(s string) string {
	rr := []rune(s)
	if len(rr) < 2 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(string(rr[0])) + string(rr[1:])
}

/*
Words returns the words in s as a slice of strings in order of
their appearance. Word boundaries include any space character
(as defined by Unicode), forward slashes, endashes, and emdashes.
In addition, grammatical marks adjacent to word boundaries are
omitted. Grammatical marks are defined as one of the following:
!?,.'"[]()*~{}:;-<>+=|%&@#$^\`

	// ww is []string{"Here's", "a", "sentence", "said", "the", "narrator", "programmer"}
	ww := Words(`"Here's a sentence," said the narrator/programmer.`)

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
WordSet is the same as Words but removes duplicates from its results.
Words will appear in order of their first appearance in s.

See Words for what a word is in this context.
*/
func WordSet(s string) []string {
	return makeSet(Words(s), false)
}

/*
WordSetFold is the same as WordSet except it considers characters
of different cases to be duplicates. For example, the words "hello",
"Hello", and "hELlo" would all be considered equal. WordSetFold will
return all words as lowercase. Words will appear in order of their
first appearance in s.

See Words for what a word is in this context.
*/
func WordSetFold(s string) []string {
	return makeSet(Words(s), true)
}

type Word struct {
	Occurence int
	Word      string
}

type wordMap []Word

func (wm wordMap) Len() int {
	return len(wm)
}

func (wm wordMap) Less(i, j int) bool {
	return wm[i].Occurence > wm[j].Occurence
}

func (wm wordMap) Swap(i, j int) {
	wm[i], wm[j] = wm[j], wm[i]
}

/*
WordOccurence returns a slice of Word where each index represents
a single word and the number of times it appears in s. The slice
is ordered from most frequent to least frequent but beyond that
words are not guaranteed to be in the order they appeared in s.

See Words for what a word is in this context.
*/
func WordOccurence(s string) []Word {

	ww := Words(s)

	occurence := make(map[string]int, len(ww))
	for _, w := range ww {
		occurence[w]++
	}

	wm := make(wordMap, 0, len(ww))
	for w, o := range occurence {
		wm = append(wm, Word{Occurence: o, Word: w})
	}

	sort.Sort(wm)

	return wm
}

/*
WordCount returns the number of words in s.

See Words for what a word is in this context.
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
	for {
		if i == len(cc) {
			return true
		}
		if !isGrammar(cc[i]) {
			return false
		}
		if precededByBoundary {
			return true
		}
		if boundaryNext(cc, i) {
			return true
		}
		i++
	}
}

func isGrammar(c string) bool {
	grammar := `!?,.'"[]()*~{}:;-<>+=|%&@#$^\` + "`"
	return strings.Contains(grammar, c)
}

func isBoundaryChar(c string) bool {
	splitters := "–—/" // endash, emdash, and forward slash
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
