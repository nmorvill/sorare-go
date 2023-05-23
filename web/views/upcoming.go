package views

import (
	"fmt"
	"html/template"
)

func GetTemplate() *template.Template {
	t, err := template.ParseFiles("./web/templates/upcoming.gohtml")
	if err != nil {
		fmt.Println(err.Error())
	}
	return t
}
