package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	blackfriday "github.com/russross/blackfriday"
)

// Currently *no* template caching, so they can be edited without restarting the app.
// The 'var templates' perf optimization we skip for now so that we can edit templates
//   and see the result immediately.
// global init of templates:
// var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

var validTitleRegex = regexp.MustCompile("^/(edit|save|view|css|media|script)/([a-zA-Z0-9/]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validTitleRegex.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}

	return m[2], nil // The title is the second subexpression.
}

type Page struct {
	Title       string
	Body        []byte
	MdProcessed template.HTML
	IperfCmd    string
}

func (page *Page) Save() error {
	filename := "pagemd/" + page.Title + ".md"
	err := ioutil.WriteFile(filename, page.Body, 0600)
	if err != nil {
		panic(err)
	}

	// save rendered to html version
	out := blackfriday.MarkdownBasic(page.Body)
	page.MdProcessed = template.HTML(out)

	htmlname := "pagehtml/" + page.Title + ".html"
	return ioutil.WriteFile(htmlname, []byte(page.MdProcessed), 0600)
}

func LoadPage(title string) (*Page, error) {
	filename := "pagemd/" + title + ".md"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// render md to html
	out := blackfriday.MarkdownBasic(body)
	MdProcessed := template.HTML(out)

	return &Page{Title: title, Body: body, MdProcessed: MdProcessed}, nil
}

func main() {

	StartSimpleWebServer()

}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	//log.Printf("RootHandler called with RequestURI: %#+v", r.RequestURI)
	IperfHandler(w, r)
	//ViewHandler(w, r, "index")
	//http.ServeFile(w, r, "index.html")
	//http.NotFound(w, r)
}

// media, css, script, and any new verb that goes
//  through to a filesystem path that should be security checked first should
//  use makeVerbDirHandler(). It rejects invalid paths (e.g. with '..' in them)
//  and logs attempt + succes/failure.
func makeVerbDirHandler(verbdir string) http.HandlerFunc {

	h := http.StripPrefix("/"+verbdir+"/", http.FileServer(http.Dir(verbdir)))
	loggingWrapper := WrapHTTPHandler{m: h}

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("verbDirHandler for '%s' called with URL.Path: '%s'", verbdir, r.URL.Path)

		// Validate
		// Here we will extract the path after the verb
		// and call the provided handler 'fn' with it
		valid, pa := IsValidPath(r.URL.Path)
		if !valid {
			log.Printf("verbDirHandler for '%s' rejecting invalid path '%s'", verbdir, r.URL.Path)
			http.NotFound(w, r)
			return
		}
		log.Printf("verbDirHandler for '%s' accepting valid path '%s'", verbdir, pa)
		loggingWrapper.ServeHTTP(w, r)
	}
}

// the file <tmpl>.html on disk is rendered using Page p
func renderTemplate(w http.ResponseWriter, fnTmpl string, p *Page) {
	out := blackfriday.MarkdownBasic(p.Body)
	p.MdProcessed = template.HTML(out)
	//log.Printf("renderTemplate file: %s.html: from .Body: '%v'   ---->  .MdProcessed: '%v'\n", fnTmpl, string(p.Body), p.MdProcessed)

	myT, err := template.ParseFiles("templates/" + fnTmpl + ".html")
	if err != nil {
		log.Printf("renderTemplate could not load file: %s.html", fnTmpl)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// using the var templates cache:
	//myT := templates.Lookup(fnTmpl + ".html")
	err = myT.Execute(w, p)

	//	err := templates.ExecuteTemplate(w, fnTmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// factor out common error handling that was the call to getTitle() + err
//  handling by generating a closure as a handler.
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Here we will extract the page title from the Request,
		// and call the provided handler 'fn'
		// NB used to be call to getTitle()
		m := validPathRegex.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := LoadPage(title)
	log.Printf("ViewHandler() called with title: '%s'\n", title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := LoadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)

}

func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	//log.Printf("SaveHandler called for title: %v  with body: %v\n", title, body)
	p := &Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func CssHandler(w http.ResponseWriter, r *http.Request, path string) {
	title := r.URL.Path[len("/css/"):]
	if true {
		log.Printf("CssHandler called with %#+v\n", title)
	}
	h := http.StripPrefix("/css/", http.FileServer(http.Dir("css")))
	h.ServeHTTP(w, r)
}

func StartSimpleWebServer() {

	http.HandleFunc("/", RootHandler)

	http.HandleFunc("/css/", makeVerbDirHandler("css"))
	http.HandleFunc("/script/", makeVerbDirHandler("script"))
	http.HandleFunc("/media/", makeVerbDirHandler("media"))

	http.HandleFunc("/view/", makeHandler(ViewHandler))
	http.HandleFunc("/edit/", makeHandler(EditHandler))
	http.HandleFunc("/save/", makeHandler(SaveHandler))

	http.HandleFunc("/configure-iperf", IperfHandler)
	http.HandleFunc("/simple", SimpleHandler)
	http.HandleFunc("/upload", handleFileUpload)
	http.HandleFunc("/templates/", makeVerbDirHandler("templates"))

	port := "8090"
	fmt.Printf("listening on localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("done with ListenAndServe on localhost:%s\n", port)

}
