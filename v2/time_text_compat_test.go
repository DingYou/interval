package interval

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestParseTimeEndpointAliases(t *testing.T) {
	const stamp = "2026-11-15T03:00:00.123456789Z"
	for _, tc := range []struct{ in, want string }{
		{"empty", "empty"},
		{"[ NuLl , NULL ]", "(,)"}, {"(null,Null)", "(,)"},
		{"[null,NULL)", "(,)"}, {"(null,NULL]", "(,)"},
		{"(nil,nil)", "(,)"}, {"[null,nil]", "(,)"}, {"[NiL,NuLl)", "(,)"},
		{"(NULL,inf)", "(,)"}, {"(-inf,NULL)", "(,)"}, {"(,NULL]", "(,)"}, {"[NULL,)", "(,)"},
		{"(nil,inf)", "(,)"}, {"(-inf,nil)", "(,)"}, {"(,nil]", "(,)"}, {"[nil,)", "(,)"},
		{"[NULL," + stamp + "]", "(," + stamp + "]"},
		{"[" + stamp + ",Null]", "[" + stamp + ",)"},
		{"[NiL," + stamp + "]", "(," + stamp + "]"},
		{"[" + stamp + ",nil]", "[" + stamp + ",)"},
		{"[" + stamp + ",Inf)", "[" + stamp + ",)"},
		{"(-InF,+INF)", "(,)"},
		{"(,+inf)", "(,)"}, {"(,)", "(,)"},
		{`["` + stamp + `","` + stamp + `"]`, "[" + stamp + "," + stamp + "]"},
	} {
		t.Run(tc.in, func(t *testing.T) {
			got, err := ParseTime(tc.in)
			if err != nil || FormatTime(got) != tc.want {
				t.Fatalf("ParseTime(%q) = %q, %v; want %q", tc.in, FormatTime(got), err, tc.want)
			}
		})
	}
}

func TestParseTimeRejectsInvalidAliases(t *testing.T) {
	const stamp = "2026-11-15T03:00:00Z"
	for _, in := range []string{
		"(inf,nil)", "(+inf,nil)", "(nil,-inf)",
		"(-inf,inf]", "[-inf,inf)",
		"[," + stamp + "]", "[" + stamp + ",]",
	} {
		t.Run(in, func(t *testing.T) {
			if _, err := ParseTime(in); !errors.Is(err, ErrSyntax) {
				t.Fatalf("got %v; want ErrSyntax", err)
			}
		})
	}
	for _, in := range []string{`("NULL",)`, `("nil",)`, `("-inf",)`, `(,"+inf")`} {
		var parseErr *time.ParseError
		if _, err := ParseTime(in); !errors.As(err, &parseErr) {
			t.Errorf("ParseTime(%q) got %v; want timestamp parse error", in, err)
		}
	}
}

func TestParseGenericAliasLikeStringsRemainFinite(t *testing.T) {
	decode := func(s string) (string, error) { return s, nil }
	for _, text := range []string{"[NULL,nil]", "[+inf,-inf]", "[inf,inf]"} {
		got, err := Parse(text, decode, strings.Compare)
		if err != nil || got.String() != text {
			t.Errorf("Parse(%q) = %q, %v; want unchanged finite endpoints", text, got.String(), err)
		}
	}
}
