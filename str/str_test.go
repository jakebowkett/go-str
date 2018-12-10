package str

import (
	"sort"
	"testing"
)

func TestNth(t *testing.T) {

	cases := []struct {
		n    int
		s    string
		sub  string
		want int
	}{
		//   0123456789  12
		{5, "hi hi hi hi hi", "hi", 12},
		{5, "hi hi hi    hi", "hi", -1},
		{1, "世界世界世界世界", "世", 0},
		{3, "世界世界世界世界", "世", 4},
		{-1, "世界世界世界世界", "世", 6},
		{-2, "世界世界世界世界", "世", 4},
		{2, "💩💩💩", "💩", 1}, // poop emoji
		{1, "hi", "hi", 0},
		{0, "hi", "hi", -1},
	}

	for _, c := range cases {
		if got := Nth(c.s, c.sub, c.n); got != c.want {
			t.Errorf("Nth(%q, %q, %d) return %d, wanted %d.", c.s, c.sub, c.n, got, c.want)
		}
	}
}

func TestPadLeft(t *testing.T) {

	cases := []struct {
		n    int
		pad  rune
		s    string
		want string
	}{
		{5, ' ', "Hello", "Hello"},
		{10, ' ', "Hello", "     Hello"},
		{3, ' ', "Hello", "Hello"},
		{-1, ' ', "Hello", "Hello"},
		{3, ' ', "", "   "},
		{5, ' ', "世界", "   世界"},
		{5, ' ', "💩💩💩", "  💩💩💩"}, // poop emoji
		{5, '世', "hi", "世世世hi"},
	}

	for _, c := range cases {
		if got := PadLeft(c.s, c.pad, c.n); got != c.want {
			t.Errorf("PadLeft(%q, %q, %d) return %q, wanted %q.", c.s, c.pad, c.n, got, c.want)
		}
	}
}

func TestPadRight(t *testing.T) {

	cases := []struct {
		n    int
		pad  rune
		s    string
		want string
	}{
		{5, ' ', "Hello", "Hello"},
		{10, ' ', "Hello", "Hello     "},
		{3, ' ', "Hello", "Hello"},
		{-1, ' ', "Hello", "Hello"},
		{3, ' ', "", "   "},
		{5, ' ', "世界", "世界   "},
		{5, ' ', "💩💩💩", "💩💩💩  "}, // poop emoji
		{5, '世', "hi", "hi世世世"},
	}

	for _, c := range cases {
		if got := PadRight(c.s, c.pad, c.n); got != c.want {
			t.Errorf("PadRight(%q, %q, %d) return %q, wanted %q.", c.s, c.pad, c.n, got, c.want)
		}
	}
}

func TestPadToLongest(t *testing.T) {

	cases := []struct {
		pad  rune
		ss   []string
		want []string
	}{
		{
			' ',
			[]string{
				"hi",
				"   ",
				"世界",
				"💩💩💩",
				"世界世界世界",
				"\n",
			},
			[]string{
				"hi    ",
				"      ",
				"世界    ",
				"💩💩💩   ",
				"世界世界世界",
				"\n     ",
			},
		},
		{
			'界',
			[]string{
				"hi",
				"   ",
				"世界",
				"💩💩💩",
				"世界世界世界",
				"\n",
			},
			[]string{
				"hi界界界界",
				"   界界界",
				"世界界界界界",
				"💩💩💩界界界",
				"世界世界世界",
				"\n界界界界界",
			},
		},
	}

	for _, c := range cases {
		if got := PadToLongest(c.ss, c.pad); !strSliceEqual(got, c.want) {
			t.Errorf(
				"PadToLongest(ss, %q)\n"+
					"    return %q\n"+
					"    wanted %q.",
				c.pad, got, c.want)
		}
	}
}

func TestLen(t *testing.T) {

	cases := []struct {
		want int
		s    string
	}{
		{5, "Hello"},
		{2, "世界"},
		{1, "💩"},            // poop emoji
		{3, "💩💩💩"},          // poop emoji
		{10, "💩💩世💩 💩💩l,\n"}, // poop emoji
		{0, ""},
	}

	for _, c := range cases {
		if got := Len(c.s); got != c.want {
			t.Errorf("Len(%q) return %d, wanted %d.", c.s, got, c.want)
		}
	}
}

func TestChar(t *testing.T) {

	cases := []struct {
		i       int
		s       string
		want    string
		wantErr bool
	}{
		{10, "Hello", "", true},
		{-1, "Hello", "o", false},
		{0, "Hello", "H", false},
		{3, "Hello", "l", false},
		{2, "世界世界", "世", false},
		{1, "💩💩💩💩", "💩", false}, // poop emoji
		{1, "", "", true},
	}

	for _, c := range cases {
		if got, err := Char(c.s, c.i); got != c.want || c.wantErr && err == nil {
			errStr := "nil"
			if c.wantErr {
				errStr = "error"
			}
			t.Errorf(
				"Char(%q, %d)\n"+
					"    return %q, %s\n"+
					"    wanted %q, %s.\n",
				c.s, c.i, got, err, c.want, errStr)
		}
	}
}

func TestChars(t *testing.T) {

	cases := []struct {
		s    string
		want []string
	}{
		{"Hello", []string{"H", "e", "l", "l", "o"}},
		{"世界", []string{"世", "界"}},
		{"💩💩💩", []string{"💩", "💩", "💩"}}, // poop emoji
		{"\n\n", []string{"\n", "\n"}},
		{"", []string{}},
	}

	for _, c := range cases {
		if got := Chars(c.s); !strSliceEqual(got, c.want) {
			t.Errorf("Chars(%q) return %v, wanted %v.", c.s, got, c.want)
		}
	}
}

func TestCharSet(t *testing.T) {

	cases := []struct {
		fold bool
		s    string
		want []string
	}{
		// Unfolded.
		{
			false,
			"Hello",
			[]string{"H", "e", "l", "o"},
		},
		{
			false,
			"世界世界",
			[]string{"世", "界"},
		},
		{
			false,
			"💩💩💩",
			[]string{"💩"}, // poop emoji
		},
		{
			false,
			"\n\n",
			[]string{"\n"},
		},
		{
			false,
			"",
			[]string{},
		},

		// Folded.
		{
			true,
			"Hello hi",
			[]string{"h", "e", "l", "o", " ", "i"},
		},
		{
			true,
			"世界世界",
			[]string{"世", "界"},
		},
		{
			true,
			"世a界ASG世f界",
			[]string{"世", "a", "界", "s", "g", "f"},
		},
		{
			true,
			"💩💩💩",
			[]string{"💩"}, // poop emoji
		},
		{
			true,
			"\n\n",
			[]string{"\n"},
		},
		{
			true,
			"",
			[]string{},
		},
	}

	for _, c := range cases {
		if got := CharSet(c.s, c.fold); !strSliceEqual(got, c.want) {
			t.Errorf("CharSet(%q, %v) return %v, wanted %v.", c.s, c.fold, got, c.want)
		}
	}
}

func strSliceEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func TestSlice(t *testing.T) {

	cases := []struct {
		n1      int
		n2      int
		s       string
		want    string
		wantErr bool
	}{
		{1, 0, "Hello", "ello", false},
		{0, 1, "Hello", "H", false},
		{1, 4, "Hello", "ell", false},
		{0, 5, "Hello", "Hello", false},
		{0, 6, "Hello", "", true},
		{0, -1, "Hello", "Hell", false},
		{-1, 0, "Hello", "o", false},
		{-1, -1, "Hello", "", false},
		{1, 2, "世界", "界", false},
		{2, 3, "💩💩💩💩💩", "💩", false}, // poop emoji
		{0, 0, "", "", false},
		{0, 1, "", "", true},
		{0, 2, "", "", true},
	}

	for _, c := range cases {
		if got, err := Slice(c.s, c.n1, c.n2); got != c.want || c.wantErr && err == nil {
			errStr := "nil"
			if c.wantErr {
				errStr = "error"
			}
			t.Errorf(
				"Slice(%q, %d, %d)\n"+
					"    return %q, %v\n"+
					"    wanted %q, %s.\n",
				c.s, c.n1, c.n2, got, err, c.want, errStr)
		}
	}
}

func TestCapitalise(t *testing.T) {

	cases := []struct {
		s    string
		want string
	}{
		{"Hello", "Hello"},
		{"hello", "Hello"},
		{"世界", "世界"},
		{"💩", "💩"}, // poop emoji
		{"", ""},
		{"\n", "\n"},
	}

	for _, c := range cases {
		if got := Capitalise(c.s); got != c.want {
			t.Errorf("Capitalise(%q) return %q, wanted %q.", c.s, got, c.want)
		}
	}
}

func TestWords(t *testing.T) {

	cases := []struct {
		s    string
		want []string
	}{
		{"grammar at end,,)", []string{"grammar", "at", "end"}},
		{"    Status: happy", []string{"Status", "happy"}},
		{"Status: happy", []string{"Status", "happy"}},
		{"Status::(happy)", []string{"Status::(happy"}},
		{"ei\nther/or", []string{"ei", "ther", "or"}},
		{"either/or", []string{"either", "or"}},
		{"either/\nor", []string{"either", "or"}},
		{`"here's an em—dash"`, []string{"here's", "an", "em", "dash"}},
		{`"here's some dialogue!"`, []string{"here's", "some", "dialogue"}},
		{"it's grammar!", []string{"it's", "grammar"}},
		{"Hello there, friend!", []string{"Hello", "there", "friend"}},
		{"hi,,    my name is thing", []string{"hi", "my", "name", "is", "thing"}},
		{"世界 世界世界", []string{"世界", "世界世界"}},
		{"hello I am poop 💩 hi", []string{"hello", "I", "am", "poop", "💩", "hi"}}, // poop emoji
		{"", []string{}},
		{"\n\n\n  \n\n\n", []string{}},
		{"interrupted\n\n\nstring.\n\n\n", []string{"interrupted", "string"}},
	}

	for _, c := range cases {
		if got := Words(c.s); !strSliceEqual(got, c.want) {
			t.Errorf("Words(%q) return %v, wanted %v.", c.s, got, c.want)
		}
	}
}

func TestWordSet(t *testing.T) {

	cases := []struct {
		fold bool
		s    string
		want []string
	}{

		// Unfolded
		{
			false,
			"I'm really, really tired of thinking of ways to test shit.",
			[]string{
				"I'm",
				"really",
				"tired",
				"of",
				"thinking",
				"ways",
				"to",
				"test",
				"shit",
			},
		},
		{
			false,
			"REALLY, Really, really... tired.",
			[]string{
				"REALLY",
				"Really",
				"really",
				"tired",
			},
		},

		// Folded
		{
			true,
			"I'm really, really tired of thinking of ways to test shit.",
			[]string{
				"i'm",
				"really",
				"tired",
				"of",
				"thinking",
				"ways",
				"to",
				"test",
				"shit",
			},
		},
		{
			true,
			"REALLY, Really, really... tired.",
			[]string{
				"really",
				"tired",
			},
		},
	}

	for _, c := range cases {
		if got := WordSet(c.s, c.fold); !strSliceEqual(got, c.want) {
			t.Errorf("WordSet(%q, %v) return %v, wanted %v.", c.s, c.fold, got, c.want)
		}
	}
}

func TestWordCount(t *testing.T) {

	cases := []struct {
		want int
		s    string
	}{
		{4, `"here's a forward/slash!"`},
		{4, `"here's an em—dash"`},
		{3, `"here's some dialogue!"`},
		{3, "Hello there, friend!"},
		{5, "hi,,    my name is thing"},
		{2, "世界 世界世界"},
		{8, "hello I am poop 💩 that's my face"}, // poop emoji
		{0, ""},
		{0, "\n\n\n  \n\n\n"},
		{9, "This is a good string\n\n\nthat continues after newlines.\n\n\n"},
	}

	for _, c := range cases {
		wordsLen := len(Words(c.s))
		got := WordCount(c.s)
		if got != wordsLen {
			t.Errorf("WordCount(%q) return %d, len(Words(%q)) is %d.", c.s, got, c.s, wordsLen)
		}
		if got != c.want {
			t.Errorf("WordCount(%q) return %d, wanted %d.", c.s, got, c.want)
		}
	}
}

func TestWordsByOccurrence(t *testing.T) {

	cases := []struct {
		s    string
		fold bool
		sort bool
		want OccMap
	}{
		{
			s:    "grammar grammar at end,,)",
			fold: false,
			sort: true,
			want: OccMap{
				{SubStr: "grammar", N: 2},
				{SubStr: "at", N: 1},
				{SubStr: "end", N: 1},
			},
		},
		{
			s:    `"Here's the dialogue," said the narrator/programmer to the listener!! And here's this.`,
			fold: false,
			sort: true,
			want: OccMap{
				{SubStr: "the", N: 3},
				{SubStr: "Here's", N: 1},
				{SubStr: "dialogue", N: 1},
				{SubStr: "said", N: 1},
				{SubStr: "narrator", N: 1},
				{SubStr: "programmer", N: 1},
				{SubStr: "to", N: 1},
				{SubStr: "listener", N: 1},
				{SubStr: "And", N: 1},
				{SubStr: "here's", N: 1},
				{SubStr: "this", N: 1},
			},
		},
		{
			s:    `"Here's the dialogue," said the narrator/programmer to the listener!! And here's this.`,
			fold: true,
			sort: true,
			want: OccMap{
				{SubStr: "the", N: 3},
				{SubStr: "here's", N: 2},
				{SubStr: "dialogue", N: 1},
				{SubStr: "said", N: 1},
				{SubStr: "narrator", N: 1},
				{SubStr: "programmer", N: 1},
				{SubStr: "to", N: 1},
				{SubStr: "listener", N: 1},
				{SubStr: "and", N: 1},
				{SubStr: "this", N: 1},
			},
		},
		{
			s:    "thing, Thing, and THING",
			fold: true,
			sort: true,
			want: OccMap{
				{SubStr: "thing", N: 3},
				{SubStr: "and", N: 1},
			},
		},
	}

	for _, c := range cases {
		got := WordsByOccurrence(c.s, c.fold)
		if c.sort {
			sort.Sort(got)
		}
		if !occSliceCorrect(got, c.want) {
			t.Errorf(
				"WordsByOccurrence(%q)\n"+
					"    return %v\n"+
					"    wanted %v",
				c.s, got, c.want)
		}
	}
}

func TestCharsByOccurrence(t *testing.T) {

	cases := []struct {
		s    string
		fold bool
		sort bool
		want OccMap
	}{
		{
			s:    "Hello there!",
			fold: false,
			sort: true,
			want: OccMap{
				{SubStr: "e", N: 3},
				{SubStr: "l", N: 2},
				{SubStr: "H", N: 1},
				{SubStr: "o", N: 1},
				{SubStr: " ", N: 1},
				{SubStr: "t", N: 1},
				{SubStr: "h", N: 1},
				{SubStr: "r", N: 1},
				{SubStr: "!", N: 1},
			},
		},
		{
			s:    "Hello there!",
			fold: true,
			sort: true,
			want: OccMap{
				{SubStr: "e", N: 3},
				{SubStr: "l", N: 2},
				{SubStr: "h", N: 2},
				{SubStr: "o", N: 1},
				{SubStr: " ", N: 1},
				{SubStr: "t", N: 1},
				{SubStr: "r", N: 1},
				{SubStr: "!", N: 1},
			},
		},
		{
			s:    "Hello, 世界! Small 世界.",
			fold: true,
			sort: true,
			want: OccMap{
				{SubStr: "l", N: 4},
				{SubStr: " ", N: 3},
				{SubStr: "世", N: 2},
				{SubStr: "界", N: 2},
				{SubStr: "h", N: 1},
				{SubStr: "e", N: 1},
				{SubStr: "o", N: 1},
				{SubStr: ",", N: 1},
				{SubStr: "!", N: 1},
				{SubStr: "s", N: 1},
				{SubStr: "m", N: 1},
				{SubStr: "a", N: 1},
				{SubStr: ".", N: 1},
			},
		},
	}

	for _, c := range cases {
		got := CharsByOccurrence(c.s, c.fold)
		if c.sort {
			sort.Sort(got)
		}
		if !occSliceCorrect(got, c.want) {
			t.Errorf(
				"CharsByOccurrence(%q)\n"+
					"    return %v\n"+
					"    wanted %v",
				c.s, got, c.want)
		}
	}
}

func occSliceCorrect(gotOcc, wantOcc OccMap) bool {

	if len(gotOcc) != len(wantOcc) {
		return false
	}
	if len(gotOcc) == 0 {
		return true
	}

	seen := make(map[string]bool, len(gotOcc))
	prevN := gotOcc[0].N

	for _, got := range gotOcc {

		// There shouldn't be duplicates.
		if seen[got.SubStr] {
			return false
		}
		seen[got.SubStr] = true

		// The substring we got should be one
		// we're expecting.
		if !inOccSlice(wantOcc, got.SubStr) {
			return false
		}

		// Later occurrences shouldn't be more
		// numerous than previous ones.
		if got.N > prevN {
			return false
		}
		prevN = got.N
	}

	return true
}

func inOccSlice(om OccMap, s string) bool {
	for _, o := range om {
		if o.SubStr == s {
			return true
		}
	}
	return false
}
