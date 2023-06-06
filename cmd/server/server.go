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
		nbGames, _ := strconv.Atoi(c.DefaultQuery("nbGames", "5"))
		minGames, _ := strconv.Atoi(c.DefaultQuery("minGames", "0"))
		sequence, _ := strconv.Atoi(c.DefaultQuery("sequence", "0"))
		allGameweeks := c.DefaultQuery("allGameweeks", "off") == "on"
		league := c.DefaultQuery("league", "all")
		search := c.DefaultQuery("search", "")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

		var res bytes.Buffer
		if presentation == "Table" {
			res = getTableResult(nbGames, minGames, sequence, allGameweeks, search, league, c.Request.URL.String(), page)
		} else {
			res = getGraphResult(nbGames, minGames, search, league)
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", res.Bytes())
	})

	r.Run()
}

func getTableResult(nbGames int, minGames int, sequence int, allGameweeks bool, search string, league string, url string, page int) bytes.Buffer {
	calendars := cache.GetData("calendars", sorare_api.GetCalendars)
	export := sorare_api.ArrangeTable(calendars, nbGames, minGames, sequence, allGameweeks, search, league, url, page)

	ret := getTableTemplate(export)
	return ret
}

func getGraphResult(nbGames int, minGames int, search string, league string) bytes.Buffer {
	calendars := cache.GetData("calendars", sorare_api.GetCalendars)
	ret := getGraphTemplate(calendars, nbGames, minGames, 1000, 600, search, league)
	return ret
}

func getTableTemplate(export sorare_api.TableExport) bytes.Buffer {
	t := views.GetTableTemplate()
	var out bytes.Buffer
	err := t.Execute(&out, export)
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
