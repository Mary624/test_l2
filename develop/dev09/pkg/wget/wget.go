package wget

import (
	"errors"
	"net/url"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

var (
	ErrCannotParseUrl = errors.New("can't parse url")
)

func Wget() error {
	// получаем аргументы
	var max int
	flag.IntVarP(&max, "level", "l", 1, "choose depth")

	var withImg bool
	flag.BoolVarP(&withImg, "page-requisites", "p", false, "with img")

	flag.Parse()
	link := flag.Arg(0)
	if !strings.Contains(link, ".") {
		return ErrCannotParseUrl
	}

	// если некорретный url, завершаем программу
	u, err := url.Parse(link)
	if err != nil {
		return ErrCannotParseUrl
	}
	if u.Host == "" || u.Scheme == "" {
		return ErrCannotParseUrl
	}
	// строим дерево из ссылок и src
	tree, err := NewTree(link, max)
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	// добавляем файлы
	err = tree.AddFiles(dir, withImg)
	if err != nil {
		return err
	}
	return nil
}
