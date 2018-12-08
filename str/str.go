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
	return strings.Repeat(string(padChar), diff) + s
}

/*
PadRight suffixes s with padChar until s contains length number
of characters.
*/
func PadRight(s string, padChar rune, length int) string {
	diff := length - len([]rune(s))
	if diff <= 0 {
		return s
	}
	return s + strings.Repeat(string(padChar), diff)
}

/*
PadToLongest suffixes each string in ss with padChar until it
contains as many characters as the longest string in ss.
*/
func PadToLongest(ss []string, padChar rune) []string {
	var longest int
	for i := range ss {
		length := len([]rune(ss[i]))
		if length > longest {
			longest = length
		}
	}
	for i := range ss {
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

If fold is set to true characters of different cases will be considered
duplicates and condensed into a single lowercase character in the
resulting slice. For example, if s were "Hi, Hello" the output with
fold set to true would be []string{"h", "i", ",", " ", "e", "l", "o"}
*/
func CharSet(s string, fold bool) []string {
	return makeSet(strings.Split(s, ""), fold)
}

/*
CharsByOccurrence returns a slice of Occurrence where each index
represents a single character and the number of times it appears in s.
The slice is ordered from most frequent to least frequent but no
other guarantees on order are made beyond that.

If fold is set to true characters of different cases will be considered
equal and all entries in the resulting slice will be in lowercase.
For example, "h" and "H" will each count towards an occurrence of the
single character, "h".
*/
func CharsByOccurrence(s string, fold bool) []Occurrence {
	return occurrences(strings.Split(s, ""), fold)
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

	ww := Words(`"Here's a sentence," said the narrator/programmer.`)
	// ww is []string{
	// 	"Here's",
	// 	"a",
	// 	"sentence",
	// 	"said",
	// 	"the",
	// 	"narrator",
	// 	"programmer",
	// }

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

/*
WordSet is the same as Words but removes duplicates from its results.
Words will appear in order of their first appearance in s.

If fold is set to true words of different cases will be considered
duplicates and condensed into a single lowercase word in the resulting
slice. For example, if s were "hello, Hello, hELlo there!" the output
with fold set to true would be []string{"hello", "there"}

See Words for what a word is in this context.
*/
func WordSet(s string, fold bool) []string {
	return makeSet(Words(s), fold)
}

/*
Occurrence contains two fields: SubStr is a substring that occurred
in a string while N is the number of times it appeared in said string.
*/
type Occurrence struct {
	SubStr string
	N      int
}

type occMap []Occurrence

func (om occMap) Len() int {
	return len(om)
}

func (om occMap) Less(i, j int) bool {
	return om[i].N > om[j].N
}

func (om occMap) Swap(i, j int) {
	om[i], om[j] = om[j], om[i]
}

/*
WordsByOccurrence returns a slice of Occurence where each index
represents a single word and the number of times it appears in s. The
slice is ordered from most frequent to least frequent but beyond that
words are not guaranteed to be in the order they appeared in s.

If fold is set to true words of different cases will be considered
equal and all entries in the resulting slice will be in lowercase.
For example, "hello" and "Hello" will each count towards an occurrence
of the single word, "hello".

See Words for what a word is in this context.
*/
func WordsByOccurrence(s string, fold bool) []Occurrence {
	return occurrences(Words(s), fold)
}

func occurrences(ss []string, fold bool) []Occurrence {

	occurrences := make(map[string]int, len(ss))
	for _, s := range ss {
		if fold {
			s = strings.ToLower(s)
		}
		occurrences[s]++
	}

	om := make(occMap, 0, len(occurrences))
	for s, o := range occurrences {
		om = append(om, Occurrence{SubStr: s, N: o})
	}

	sort.Sort(om)

	return om
}

func grammarOnBoundary(cc []string, i int, precededByBoundary bool) bool {
	for {
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
