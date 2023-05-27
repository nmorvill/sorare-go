package views

import (
	"fmt"
	"html/template"
)

func GetTableTemplate() *template.Template {
	t, err := template.ParseFiles("./web/templates/upcomingTable.gohtml")
	if err != nil {
		fmt.Println(err.Error())
	}
	return t
}

func GetGraphTemplate() *template.Template {
	t, err := template.ParseFiles("./web/templates/upcomingGraph.gohtml")
	if err != nil {
		fmt.Println(err.Error())
	}
	return t
}

func GetDetailsTemplate() *template.Template {
	t, err := template.ParseFiles("./web/templates/upcomingDetails.gohtml")
	if err != nil {
		fmt.Println(err.Error())
	}
	return t
}
