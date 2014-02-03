package main

import (
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
)

func TestSaveLoad(t *testing.T) {

	cv.Convey("When we save a test page, we should load it back and get the same page", t, func() {
		p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
		p1.Save()
		p2, err := LoadPage("TestPage")
		if err != nil {
			panic(err)
		}
		cv.So(p1.Title, cv.ShouldEqual, p2.Title)
		cv.So(p1.Body, cv.ShouldResemble, p2.Body)
	})
}
