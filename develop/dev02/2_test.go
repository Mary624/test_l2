package two

import (
	"strings"
	"testing"
)

func TestTwo(t *testing.T) {
	m := map[string]string{
		"a":        "a",
		"aa":       "aa",
		"aa10":     "a" + strings.Repeat("a", 10),
		"aa10a":    "a" + strings.Repeat("a", 10) + "a",
		"aa010a":   "",
		`qwe\`:     "",
		`qwe\\`:    `qwe\`,
		`qwe\45`:   `qwe44444`,
		`qwe\4\5`:  `qwe45`,
		`qwe\\5`:   `qwe\\\\\`,
		`a4bc2d5e`: `aaaabccddddde`,
		`abcd`:     `abcd`,
		`45a`:      ``,
		"":         "",
	}
	for k, v := range m {
		res, _ := StringUnpacking(k)
		if res != v {
			t.Fatalf("%s->%s != %s", k, res, v)
		}
	}
}
