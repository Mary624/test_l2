package wget

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/net/html"
)

type TreeHtml struct {
	Start *Node
}

type Node struct {
	Data  *html.Node
	Link  string
	Imgs  map[string]int
	Csss  map[string]int
	Num   int
	Htmls []*Node
}

type Img struct {
	Link string
	Num  int
}

type SafeInt struct {
	val int64
}

func NewSaveInt(v int64) *SafeInt {
	return &SafeInt{
		val: v,
	}
}

func (v *SafeInt) Get() int {
	return int(atomic.LoadInt64(&v.val))
}

func (v *SafeInt) Increment() {
	atomic.AddInt64(&v.val, 1)
}

func (v *SafeInt) Set(newV int) {
	atomic.StoreInt64(&v.val, int64(newV))
}

func NewTree(link string, max int) (*TreeHtml, error) {
	if max < 1 {
		return nil, nil
	}

	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	num := NewSaveInt(0)
	start, err := buildTree(u.Host, link, 1, max, num)
	if err != nil {
		return nil, err
	}

	res := TreeHtml{}
	res.Start = start

	return &res, nil
}

func buildTree(host, link string, level, max int, num *SafeInt) (*Node, error) {
	node, err := getNodeFromLink(link, host, num)
	if err != nil {
		return nil, err
	}
	num.Increment()

	// получаем src
	getImgsFromNode(node, host)
	getCsssFromNode(node, host)
	// движемся до последнего уровня
	if level < max {
		links := getLinks(node.Data)
		if len(links) == 0 {
			return node, nil
		}
		level++

		// рекурсивно добавляем новые ссылки
		getLinksFromNode(node, links, host, link, level, max, num)
	}

	return node, nil
}

func getCsssFromNode(node *Node, host string) {
	csss := getCsss(node.Data)
	mCsss := make(map[string]int, len(csss))

	var wg sync.WaitGroup
	var m sync.Mutex
	for i, css := range csss {
		if len(css) == 0 {
			continue
		}
		wg.Add(1)
		go addMap(&wg, &m, i, css, host, mCsss)

	}
	wg.Wait()
	node.Csss = mCsss
}

func getImgsFromNode(node *Node, host string) {
	imgs := getImgs(node.Data)
	mImgs := make(map[string]int, len(imgs))
	var wg sync.WaitGroup
	var m sync.Mutex

	for i, img := range imgs {
		if len(img) == 0 {
			continue
		}
		wg.Add(1)
		go addMap(&wg, &m, i, img, host, mImgs)

	}
	wg.Wait()
	node.Imgs = mImgs
}

func getLinksFromNode(node *Node, links []string, host, link string, level, max int, num *SafeInt) {
	var wg sync.WaitGroup
	var m sync.Mutex

	nodes := make([]*Node, 0, len(links))

	for _, l := range links {
		if len(l) == 0 {
			continue
		}
		wg.Add(1)
		go addLinks(&wg, &m, host, l, link, level, max, num, &nodes)
	}
	wg.Wait()
	node.Htmls = nodes
}

func addMap(wg *sync.WaitGroup, m *sync.Mutex, i int, v, host string, mapRes map[string]int) {
	defer wg.Done()
	res := v
	if len(v) > 2 && v[:2] == "//" {
		res = "http://" + v[2:]
	} else {
		u, err := url.Parse(v)
		if err != nil {
			return
		}
		if u.Scheme == "" && u.Host == "" {
			if u.Path != "" {
				res, err = url.JoinPath("http://"+host, v)
				if err != nil {
					return
				}
			}
		}
	}
	m.Lock()
	mapRes[res] = i
	m.Unlock()
}

func addLinks(wg *sync.WaitGroup, m *sync.Mutex, host, l, link string, level, max int, num *SafeInt, nodes *[]*Node) {
	defer wg.Done()
	if l[0] != '/' {
		u, err := url.Parse(link)
		if err != nil {
			return
		}
		if u.Host != host {
			host = u.Host
		}
	}
	newNode, err := buildTree(host, l, level, max, num)
	if err != nil {
		return
	}
	m.Lock()
	(*nodes) = append(*nodes, newNode)
	m.Unlock()
}

func getNodeFromLink(link, host string, num *SafeInt) (*Node, error) {
	var err error
	if link[0] == '/' {
		link, err = url.JoinPath("http://"+host, link)
		if err != nil {
			return nil, err
		}
	}

	code, err := onPage(link)
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(strings.NewReader(code))
	if err != nil {
		return nil, err
	}

	var node Node
	node.Data = doc
	node.Link = link
	node.Num = num.Get()
	return &node, nil
}

func onPage(link string) (string, error) {
	res, err := http.Get(link)
	if err != nil {
		return "", err
	}
	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(content), nil
}
