package cd

import (
	"testing"
	"unix-shell-utility/pkg/cmds"

	"github.com/stretchr/testify/assert"
)

func TestCd(t *testing.T) {
	wd := `C:\Users\maria\Desktop\Изучение go\wb\l2\test\exercises\eighth`
	cd, err := New()
	if err != nil {
		t.Fatal(err)
	}

	m := map[string]string{
		"main":         `C:\Users\maria\Desktop\Изучение go\wb\l2\test\exercises\eighth\main`,
		"main/main.go": ErrIsFile.Error(),
		"..":           `C:\Users\maria\Desktop\Изучение go\wb\l2\test\exercises`,
		"../..":        `C:\Users\maria\Desktop\Изучение go\wb\l2\test`,
		"":             cmds.ErrWrongArgs.Error(),
	}

	for k, v := range m {
		cd.changeWd(wd)
		res, _ := cd.Do(k)
		if res != "" {
			assert.Equal(t, res, v)
			continue
		}
		assert.Equal(t, cd.WorkDir(), v)
	}
}
