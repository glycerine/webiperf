package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func SimpleFormToIperfMap(form *map[string][]string) *map[string]string {
	// get a list of every value the map should contain,
	// to avoid rendering <no value> for missing values.
	values, _ := GrepNamesFromTemplateFile("templates/simple.html")

	m := make(map[string]string)

	// first extract from form
	for k, v := range *form {
		m[k] = strings.Join(v, " ")
	}

	fmt.Printf("after form extraction, m is %#+v\n", m)

	// next add empty strings for any missing values...to
	// avoid producing "<no value>"
	for _, v := range values {
		_, ok := m[v]
		if !ok {
			m[v] = ""
		}
	}
	return &m
}

func SimpleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("SimpleHandler() called\n")
	fmt.Printf("request r is %+#v\n", r)

	err := r.ParseForm()
	if err != nil {
		log.Printf("SampleHandler() could not ParseForm on request '%#+v': returned err: %v\n", r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawFormInputMap := (map[string][]string)(r.Form)

	finishedOutputMap := SimpleFormToIperfMap(&rawFormInputMap)

	fmt.Printf("=================================\n")
	fmt.Printf("input rawFormInputMap is %#+v\n", rawFormInputMap)
	fmt.Printf("---------------------------------\n")
	fmt.Printf("output finishedOutputMap is %#+v\n", *finishedOutputMap)
	fmt.Printf("=================================\n")

	fnTmpl := "simple"
	myT, err := template.ParseFiles("templates/" + fnTmpl + ".html")
	if err != nil {
		log.Printf("renderTemplate could not load file: %s.html", fnTmpl)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = myT.Execute(w, finishedOutputMap)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
