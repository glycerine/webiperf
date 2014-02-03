package main

import (
	"fmt"
	"log"
	"regexp"
)

// strict path only allows a single dot in the file extension,
// to avoid .. and injection attacks.

// use the perl class \w ===  [0-9A-Za-z_]
// specifically omitting intentionally for security: < and % and > and ( and ). (?: starts a non-capturing group,
// as per http://code.google.com/p/re2/wiki/Syntax
var strictPath string = `(?:\w|\.|-|/)+`

var validVerbs string = "edit|save|view|css|media|script|templates"

var validPathRegex = regexp.MustCompile(`^/(` + validVerbs + `)/(\w` + strictPath + `)$`)

var getPathRegex = regexp.MustCompile("^/(" + validVerbs + ")/(.+)$")

var dotdot = regexp.MustCompile("[.][.]")

func IsValidPath(s string) (ok bool, path string) {

	checkDots := dotdot.FindStringSubmatch(s)
	if checkDots != nil {
		// don't allow .. in paths
		log.Printf("invalid path detected in IsValidPath(): '%s' is bad because we detected dot-dot '..'\n", s)
		return false, ""
	}

	m := validPathRegex.FindStringSubmatch(s)
	if m == nil {
		log.Printf("invalid path detected in IsValidPath(): '%s' is bad because we it didn't pass the validPathRegex check. check validVerbs above too.'\n", s)
		return false, ""
	} else {
		pa := getPathRegex.FindStringSubmatch(s)
		if pa == nil {
			panic("")
		}
		return true, pa[2]
	}
}

func TestValidPath() {
	for _, s := range []string{"one", "two", "/edit/", "/view/a", "/media/img/tabs/blue.png", "/script/de-per/b/c/jquery-1.10.2.min.map", "..", "a/..", ".", "/edit/."} {

		ok, pa := IsValidPath(s)
		if ok {
			fmt.Printf("valid path: '%s'   with path-part: '%s'\n", s, pa)
		} else {
			fmt.Printf("Invalid Page Title: '%s'\n", s)
		}
	}
}

/*
func main() {
	TestValidPath()
}
*/
