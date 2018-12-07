package str

import (
	"testing"
)

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
		{-1, "Hello", "", true},
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
		s    string
		want []string
	}{
		{"Hello", []string{"H", "e", "l", "o"}},
		{"世界世界", []string{"世", "界"}},
		{"💩💩💩", []string{"💩"}}, // poop emoji
		{"\n\n", []string{"\n"}},
		{"", []string{}},
	}

	for _, c := range cases {
		if got := CharSet(c.s); !strSliceEqual(got, c.want) {
			t.Errorf("CharSet(%q) return %v, wanted %v.", c.s, got, c.want)
		}
	}
}

func strSliceEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, _ := range s1 {
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
		{1, 0, "Hello", "", true},
		{0, 1, "Hello", "H", false},
		{1, 4, "Hello", "ell", false},
		{0, 5, "Hello", "Hello", false},
		{0, 6, "Hello", "", true},
		{0, -1, "Hello", "", true},
		{-1, 0, "Hello", "", true},
		{-1, -1, "Hello", "", true},
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
					"    return %q, %s\n"+
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

func TestWordCount(t *testing.T) {

	cases := []struct {
		want int
		s    string
	}{
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
		if got := WordCount(c.s); got != c.want {
			t.Errorf("WordCount(%q) return %d, wanted %d.", c.s, got, c.want)
		}
	}
}
