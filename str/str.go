/*
Package str provides functions that handle single-
and multi-byte character strings in a convenient way.
*/
package str

import (
	"errors"
	"strings"
)

/*
In returns true if s is in ss.
*/
func In(ss []string, s string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}
	return false
}

/*
Len returns the number of characters in a string rather
than the number of bytes.
*/
func Len(s string) int {
	return len([]rune(s))
}

/*
SplitBefore mirrors the SplitAfter function in the standard
library's strings package. It slices s into substrings before
each instance of sep and returns a slice of those substrings.

It is equivalent to SplitBeforeN where n is -1.
*/
func SplitBefore(s, sep string) []string {
	return SplitBeforeN(s, sep, -1)
}

func SplitBeforeN(s, sep string, n int) []string {

	var subStr []string

	// Standard library returns a nil slice when n is 0.
	if n == 0 {
		return subStr
	}

	if sep == "" {
		return splitBeforeEmptySep(s, n)
	}

	for i := len(s) - len(sep); i >= 0; i-- {

		if len(subStr) == n-1 {
			break
		}

		// sep can't be inside s if sep is longer.
		if len(s) < len(sep) {
			break
		}

		// If index from the end of the string
		// is less than sep, sep can't be there.
		if len(s)-(i+1) < len(sep) {
			continue
		}

		if s[i:i+len(sep)] == sep {
			subStr = append(subStr, s[i:])
			s = s[:i]
		}
	}

	subStr = append(subStr, s)
	ReverseSlice(subStr)

	return subStr
}

func splitBeforeEmptySep(s string, n int) []string {

	var subStr []string

	// The standard library strings package returns
	// a zero-length non-nil slice when s is empty.
	if s == "" {
		return subStr
	}

	if n < 0 {
		return strings.Split(s, "")
	}

	rr := []rune(s)
	if n > len(rr) {
		return strings.Split(s, "")
	}

	n--

	subStr = append(subStr, string(rr[0:len(rr)-n]))
	for _, r := range rr[len(rr)-n:] {
		subStr = append(subStr, string(r))
	}

	return subStr
}

/*
ReverseSlice reverses the strings in ss in place.

	ss := []string{"hi", "how", "are", "you?"}
	str.ReverseSlice(ss)
	// ss is now []string{"you?", "are", "how", "hi"}

*/
func ReverseSlice(ss []string) {

	L := 0
	R := len(ss) - 1

	for {

		// Base case.
		if L >= R {
			return
		}

		// Swap elements.
		tmp := ss[L]
		ss[L] = ss[R]
		ss[R] = tmp

		// Increment inwards.
		L++
		R--
	}
}

/*
Reverse returns a new string with its characters in the reverse order.

	s := str.Reverse("hello") // "olleh"

*/
func Reverse(s string) string {

	rr := []rune(s)

	L := 0
	R := len(rr) - 1

	for {

		if L >= R {
			return string(rr)
		}

		tmp := rr[L]
		rr[L] = rr[R]
		rr[R] = tmp

		L++
		R--
	}

}

/*
Nth returns the character index of the nth instance of
subStr in s. If n is negative it will search from the end
of the string. Nth will return -1 if the nth instance
of subStr cannot be found or if n is 0.

Note that for consistency with several functions in the
standard library "strings" package, Nth considers the
empty substring "" to exist between characters as well
as at the start and end of a string. For example:

	strings.LastIndex("hi", "") // 2
	str.Nth("hi", "", -1) // 2

	strings.Index("", "") // 0
	strings.Count("", "") // 1
	str.Nth("", "", 1) // 0

	strings.Count("hi", "") // 3
	str.Nth("hi", "", 3) // 2

*/
func Nth(s, subStr string, n int) int {

	if n == 0 {
		return -1
	}
	if Len(subStr) > Len(s) {
		return -1
	}

	// For consistency with the standard library's
	// strings.Index and strings.Count we treat an
	// empty substring as valid.
	if subStr == "" {
		return nthEmptyString(s, n)
	}

	if n < 0 {
		return nthLast(s, subStr, -n)
	}
	return nthFirst(s, subStr, n)
}

func nthEmptyString(s string, n int) int {

	rr := []rune(s)

	if abs(n) > len(rr)+1 {
		return -1
	}

	if n < 0 {
		return (len(rr) - -n) + 1
	}

	return n - 1
}

func nthFirst(s, subStr string, n int) int {

	// i below is the byte position so we record
	// what character we're on.
	var charPos int
	var seen int

	for i, _ := range s {
		if i+len(subStr) > len(s) {
			return -1
		}
		if string(s[i:i+len(subStr)]) == subStr {
			seen++
			if seen == n {
				return charPos
			}
		}
		charPos++
	}
	return -1
}

func nthLast(s, subStr string, n int) int {

	rr := []rune(s)
	seen := 0
	subLen := Len(subStr) // char len not byte len

	for i := len(rr) - subLen; i >= 0; i-- {
		if string(rr[i:i+subLen]) == subStr {
			seen++
			if seen == n {
				return i
			}
		}
	}
	return -1
}

/*
Char returns character n (rather than byte n) in s. Negative values
for n are treated as an offset from the end of the string. An error
will be returned if the absolute value of n is greater than the number
of characters in s.

	s, _ := Char("Hello", 0)   // "H"
	s, _ := Char("Hello", 1)   // "e"
	s, _ := Char("Hello", -1)  // "o"
	s, _ := Char("世界", -1)   // "界"
	s, _ := Char("Hello", 8)   // Error; out of bounds.
	s, _ := Char("Hello", -7)  // Error; out of bounds.

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
CharsByOccurrence returns an unordered OccMap where each index
represents a single character and the number of times it appears in s.
OccMap implements sort.Interface; see OccMap for more details.

If fold is set to true characters of different cases will be considered
equal and all entries in the resulting slice will be in lowercase.
For example, "h" and "H" will each count towards an occurrence of the
single character, "h".
*/
func CharsByOccurrence(s string, fold bool) OccMap {
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
of s and continue until end.

Returns an error if start or end have an absolute value greater
than the number of characters in s.

	s, _ := Slice("Hello", 1, 2)     // "e"
	s, _ := Slice("Hello", -4, -1)   // "ell"
	s, _ := Slice("Hello", -1, 0)    // "o"
	s, _ := Slice("Hello", -1, 2)    // "oHe"
	s, _ := Slice("世界地球風", 1, 3) // "界地"
	s, _ := Slice("Hello", 2, 8)     // Error; out of bounds.
	s, _ := Slice("Hello", -6, 2)    // Error; out of bounds.

*/
func Slice(s string, start, end int) (string, error) {
	cc := []rune(s)
	if abs(start) > len(cc) || abs(end) > len(cc) {
		return "", errors.New("index out of bounds")
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
Words returns the words in s as a slice of strings in order of
their appearance. Word boundaries include any space character
(as defined by Unicode), forward slashes, endashes, and emdashes.
In addition, grammatical marks adjacent to word boundaries are
omitted. Grammatical marks are defined as one of the following:

	!?,.'"[]()*~{}:;-<>+=|%&@#$^\`

Grammar within a word, such as apostrophes indicating contractions,
are retained.

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
Occurrences is a pairing of a subtring and an int representing how
often that substring occurs in the string it was sourced from.
*/
type Occurrences struct {
	SubStr string
	N      int
}

/*
OccMap is an unordered sequence of Occurrences. It implements
sort.Interface and will sort the substrings from most frequent to
least, though no order is guaranteed among substrings with the
same number of occurrences.
*/
type OccMap []Occurrences

func (om OccMap) Len() int {
	return len(om)
}

func (om OccMap) Less(i, j int) bool {
	return om[i].N > om[j].N
}

func (om OccMap) Swap(i, j int) {
	om[i], om[j] = om[j], om[i]
}

/*
WordsByOccurrence returns an unordered OccMap where each index
represents a single word and the number of times it appears in s.
OccMap implements sort.Interface; see OccMap for more details.

If fold is set to true words of different cases will be considered
equal and all entries in the resulting slice will be in lowercase.
For example, "hello" and "Hello" will each count towards an occurrence
of the single word, "hello".

See Words for what a word is in this context.
*/
func WordsByOccurrence(s string, fold bool) OccMap {
	return occurrences(Words(s), fold)
}

func occurrences(ss []string, fold bool) OccMap {

	occs := make(map[string]int, len(ss))
	for _, s := range ss {
		if fold {
			s = strings.ToLower(s)
		}
		occs[s]++
	}

	om := make(OccMap, 0, len(occs))
	for s, o := range occs {
		om = append(om, Occurrences{SubStr: s, N: o})
	}

	return om
}
