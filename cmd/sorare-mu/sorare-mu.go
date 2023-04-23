package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sorare-mu/internal/sorare_api"
)

func main() {
	fmt.Println("Starting !")
	r, err := json.Marshal(sorare_api.GetCalendars())
	if err != nil {
		fmt.Println(err)
		return
	}
	e := ioutil.WriteFile("./export.json", r, 0644)
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println("Done !")
}
