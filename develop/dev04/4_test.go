package four

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFour(t *testing.T) {
	input := []string{"мама", "Амам", "олег", "аамм", "Мама", "пятак", "пятка", "тяпка",
		"тяпка", "листок", "слиток", "столик"}
	res := AnagramSearch(input)
	actual := map[string][]string{
		"мама":   {"аамм", "амам", "мама"},
		"пятак":  {"пятак", "пятка", "тяпка"},
		"листок": {"листок", "слиток", "столик"},
	}
	assert.Equal(t, res, actual)
}
