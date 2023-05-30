package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sorare-mu/internal/cache"
	"sorare-mu/internal/sorare_api"
	"sorare-mu/web/views"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	r.Static("/css", "./static/style")
	r.Static("/assets", "./static/assets")

	r.GET("/", func(c *gin.Context) {
		res := getIndexTemplate()
		c.Data(http.StatusOK, "text/html; charset=utf-8", res.Bytes())
	})

	r.GET("/matchups", func(c *gin.Context) {
		presentation := c.DefaultQuery("presentation", "Table")
		mode := c.DefaultQuery("mode", "Calendar")
		nbGames, _ := strconv.Atoi(c.DefaultQuery("nbGames", "5"))
		minGames, _ := strconv.Atoi(c.DefaultQuery("minGames", "3"))
		sequence, _ := strconv.Atoi(c.DefaultQuery("sequence", "3"))
		allGameweeks := c.DefaultQuery("allGameweeks", "off") == "on"
		league := c.DefaultQuery("league", "all")
		search := c.DefaultQuery("search", "")

		var res bytes.Buffer
		if presentation == "Table" {
			res = getTableResult(mode, nbGames, minGames, sequence, allGameweeks, search, league)
		} else {
			res = getGraphResult(nbGames, minGames, search, league)
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", res.Bytes())
	})

	r.Run()
}

func getTableResult(mode string, nbGames int, minGames int, sequence int, allGameweeks bool, search string, league string) bytes.Buffer {
	calendars := cache.GetData("calendars", sorare_api.GetCalendars)
	calendars = sorare_api.ArrangeResults(calendars, mode, nbGames, minGames, sequence, allGameweeks, search, league)

	ret := getTableTemplate(calendars)
	return ret
}

func getGraphResult(nbGames int, minGames int, search string, league string) bytes.Buffer {
	calendars := cache.GetData("calendars", sorare_api.GetCalendars)
	ret := getGraphTemplate(calendars, nbGames, minGames, 1000, 600, search, league)
	return ret
}

func getTableTemplate(clubs []sorare_api.ClubExport) bytes.Buffer {
	t := views.GetTableTemplate()
	var out bytes.Buffer
	err := t.Execute(&out, clubs)
	if err != nil {
		fmt.Println(err.Error())
	}
	return out
}

func getGraphTemplate(clubs []sorare_api.ClubExport, nbGames int, minGames int, graphWidth int, graphHeight int, search string, league string) bytes.Buffer {
	c := sorare_api.ArrangeGraph(clubs, nbGames, minGames, search, graphWidth, graphHeight, league)
	g := sorare_api.GraphExport{Clubs: c, GraphWidth: graphWidth, GraphHeight: graphHeight}
	t := views.GetGraphTemplate()
	var out bytes.Buffer
	err := t.Execute(&out, g)
	if err != nil {
		fmt.Println(err.Error())
	}
	return out
}

func getIndexTemplate() bytes.Buffer {
	l := cache.GetData("leagues", sorare_api.GetLeagues)
	t := views.GetIndexTemplate()
	var out bytes.Buffer
	err := t.Execute(&out, l)
	if err != nil {
		fmt.Println(err.Error())
	}
	return out
}
