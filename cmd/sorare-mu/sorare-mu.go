package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sorare-mu/internal/sorare_api"
	"sorare-mu/web/pages"
	"time"
)

func main() {
	fmt.Println("Starting !")
	start := time.Now()
	calendars := sorare_api.GetBestSequence()
	out := t(calendars)
	e := ioutil.WriteFile("./index.html", out.Bytes(), 0644)
	if e != nil {
		fmt.Println(e)
		return
	}
	elapsed := time.Since(start)
	fmt.Printf("Done in %s!", elapsed)
}

func t(clubs []sorare_api.ClubExport) bytes.Buffer {
	t := pages.GetTemplate()
	var out bytes.Buffer
	err := t.Execute(&out, clubs)
	if err != nil {
		fmt.Println(err.Error())
	}
	return out
}
