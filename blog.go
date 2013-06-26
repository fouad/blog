package main

import (
	"flag"
	"github.com/drone/routes"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Page struct {
	Title string
	Key   string
	Md    string
	Html  template.HTML
}

var port = flag.Int("port", 3000, "port to listen on")

var pages map[string]*Page
var index, page *template.Template

func PageHandler(rw http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()

	page_key := params.Get(":page_key")

	log.Println(page_key)

	p := pages[page_key]

	page.Execute(rw, p)
}

func IndexHandler(rw http.ResponseWriter, req *http.Request) {
	index.Execute(rw, pages)
}

func main() {
	// load all of the pages from `pages/`
	pdir, _ := ioutil.ReadDir("pages")

	pages = make(map[string]*Page)

	for _, v := range pdir {
		key := strings.TrimSuffix(v.Name(), ".md")

		file, _ := ioutil.ReadFile("pages/" + v.Name())

		// grab first line of the file
		re, _ := regexp.Compile(`.*`)

		token := re.FindString(string(file))

		title := strings.TrimPrefix(token, "## ")

		// Parse Markdown to HTML
		html := blackfriday.MarkdownCommon(file)

		page := Page{title, key, string(file), template.HTML(string(html))}

		pages[key] = &page
	}

	index, _ = template.ParseFiles("tmpl/index.html")

	page, _ = template.ParseFiles("tmpl/page.html")

	mux := routes.New()

	mux.Get("/", IndexHandler)
	mux.Get("/:page_key/", PageHandler)

	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./public/css/"))))
	http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./public/img/"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./public/js/"))))
	http.Handle("/", mux)

	flag.Parse()

	log.Println(":" + strconv.Itoa(*port))

	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
