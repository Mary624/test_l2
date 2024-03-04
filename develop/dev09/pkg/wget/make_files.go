package wget

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func (t *TreeHtml) AddFiles(pathRes string, withImgs bool) error {
	dir := getDirName(pathRes)
	err := os.Mkdir(dir, os.ModeDir)
	if err != nil {
		return err
	}

	pathRes = path.Join(pathRes, dir)
	err = makeFile(t.Start, pathRes, withImgs)
	return err
}

func makeFile(node *Node, pathRes string, withPic bool) error {
	file := fmt.Sprintf("h%d.html", node.Num)
	f, err := os.Create(path.Join(pathRes, file))
	if err != nil {
		return err
	}

	// заменяем ссылки в html
	name := fmt.Sprintf("h%d", node.Num)
	l := mapLinks(node)
	if len(l) > 0 {
		replaceLinks(node.Data, l, true, "")
	}
	if len(node.Imgs) > 0 && withPic {
		replaceLinks(node.Data, node.Imgs, false, name)
	}
	if len(node.Csss) > 0 {
		replaceCsss(node.Data, node.Csss, name)
	}

	// сохраняем src в файлы
	err = saveSrcToDir(node, pathRes, name, false)
	if err != nil {
		return err
	}
	if withPic {
		err = saveSrcToDir(node, pathRes, name, true)
		if err != nil {
			return err
		}
	}

	var b bytes.Buffer
	err = html.Render(&b, node.Data)
	if err != nil {
		return err
	}

	_, err = f.Write(b.Bytes())
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var m sync.Mutex
	for _, n := range node.Htmls {
		wg.Add(1)
		go func(wg *sync.WaitGroup, m *sync.Mutex, n *Node, pathRes string, withPic bool) {
			defer wg.Done()
			e := makeFile(n, pathRes, withPic)
			if e != nil {
				m.Lock()
				err = e
				m.Unlock()
			}
		}(&wg, &m, n, pathRes, withPic)
	}
	wg.Wait()
	return err
}

func getLinks(n *html.Node) []string {
	var links []string
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, getLinks(c)...)
	}
	return links
}

func getImgs(n *html.Node) []string {
	var imgs []string
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, a := range n.Attr {
			if a.Key == "src" {
				imgs = append(imgs, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		imgs = append(imgs, getImgs(c)...)
	}
	return imgs
}

func replaceCsss(n *html.Node, l map[string]int, nameParent string) {
	b := false
	sheet := -1
	if n.Type == html.ElementNode && n.Data == "link" {
		for i, a := range n.Attr {
			if a.Key == "rel" && a.Val == "stylesheet" {
				b = true
			} else if a.Key == "href" {
				sheet = i
			}
			if sheet >= 0 && b {
				break
			}
		}
	}
	if sheet >= 0 && b {
		replaceLink(n, l, sheet, fmt.Sprintf("src/%sc", nameParent), "css")
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		replaceCsss(c, l, nameParent)
	}
}

func getCsss(n *html.Node) []string {
	b := false
	href := ""
	var csss []string
	if n.Type == html.ElementNode && n.Data == "link" {
		for _, a := range n.Attr {
			if a.Key == "rel" && a.Val == "stylesheet" {
				b = true
			} else if a.Key == "href" {
				href = a.Val
			}
			if href != "" && b {
				break
			}
		}
	}
	if href != "" && b {
		csss = append(csss, href)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		csss = append(csss, getCsss(c)...)
	}
	return csss
}

func replaceLinks(n *html.Node, l map[string]int, isLink bool, nameParent string) {
	data, key := "a", "href"
	if !isLink {
		data, key = "img", "src"
	}
	if n.Type == html.ElementNode && n.Data == data {
		for i := 0; i < len(n.Attr); i++ {
			if n.Attr[i].Key == key {
				name, t := fmt.Sprintf("src/%sp", nameParent), getType(n.Attr[i].Val)
				if isLink {
					name = "h"
					t = "html"
				}
				replaceLink(n, l, i, name, t)
			} else if n.Attr[i].Key == "srcset" && !isLink {
				n.Attr[i].Val = ""
				n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		replaceLinks(c, l, isLink, nameParent)
	}
}

func replaceLink(n *html.Node, l map[string]int, i int, name, t string) {
	val := n.Attr[i].Val
	if len(val) > 2 && val[0] == '/' {
		for h, v := range l {
			if val[:2] == "//" {
				if "http://"+val[2:] == h {
					n.Attr[i].Val = fmt.Sprintf("%s%d.%s", name, v, t)
					return
				}
			}
			u, _ := url.Parse(h)
			f := fmt.Sprintf("http://%s", u.Host)
			r, _ := url.JoinPath(f, val)
			if h == r {
				n.Attr[i].Val = fmt.Sprintf("%s%d.%s", name, v, t)
			}
		}
		return
	}
	if v, ok := l[val]; ok {
		n.Attr[i].Val = fmt.Sprintf("%s%d.%s", name, v, t)
	}
}

func getType(link string) string {
	l := strings.Split(link, ".")
	res := strings.ToLower(l[len(l)-1])
	if strings.Contains(res, "/") {
		res = strings.Split(res, "/")[0]
	}
	if strings.Contains(res, "?") {
		res = strings.Split(res, "?")[0]
	}
	return res
}

func mapLinks(n *Node) map[string]int {
	htmls := n.Htmls
	res := make(map[string]int, len(htmls))

	for _, h := range htmls {
		res[h.Link] = h.Num
	}
	return res
}

func getDirName(pathRes string) string {
	dir := "result"
	num := 1
	for {
		if _, err := os.Stat(path.Join(pathRes, dir)); os.IsNotExist(err) {
			break
		} else {
			dir = fmt.Sprintf("result%d", num)
			num++
		}
	}
	return dir
}

func saveSrcToDir(node *Node, pathRes, nameParent string, isImg bool) error {
	pathSrc := path.Join(pathRes, "src")
	if _, err := os.Stat(pathSrc); os.IsNotExist(err) {
		err := os.MkdirAll(pathSrc, os.ModePerm)
		if err != nil {
			return err
		}
	}
	m := node.Csss
	if isImg {
		m = node.Imgs
	}
	for k, v := range m {
		err := saveSrc(k, pathSrc, nameParent, v, isImg)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveSrc(link, pathRes, nameParent string, num int, isImg bool) error {
	r, err := http.Get(link)
	if err != nil {
		return nil
	}
	if r.StatusCode == http.StatusForbidden {
		return nil
	}
	defer r.Body.Close()

	name := nameParent + "c"
	if isImg {
		name = nameParent + "p"
	}
	t := "css"
	if isImg {
		t = getType(link)
	}
	file, err := os.Create(path.Join(pathRes, fmt.Sprintf("%s%d.%s", name, num, t)))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	return err
}
