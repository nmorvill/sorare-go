package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sorare-mu/internal/cache"
	"sorare-mu/internal/sorare_api"
	"sorare-mu/web/views"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.Static("/css", "./static/style")
	r.Static("/assets", "./static/assets")

	r.LoadHTMLGlob("./static/pages/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/calendars", func(c *gin.Context) {
		mode := c.DefaultQuery("mode", "Calendar")
		nbGames, _ := strconv.Atoi(c.DefaultQuery("nbGames", "5"))
		minGames, _ := strconv.Atoi(c.DefaultQuery("minGames", "3"))
		sequence, _ := strconv.Atoi(c.DefaultQuery("sequence", "3"))
		allGameweeks := c.DefaultQuery("allGameweeks", "off") == "on"
		search := c.DefaultQuery("search", "")

		res := getResult(mode, nbGames, minGames, sequence, allGameweeks, search)

		c.Data(http.StatusOK, "text/html; charset=utf-8", res.Bytes())
	})

	r.Run()
}

func getResult(mode string, nbGames int, minGames int, sequence int, allGameweeks bool, search string) bytes.Buffer {
	fmt.Println("Starting !")
	start := time.Now()

	calendars := cache.GetData("calendars", sorare_api.GetCalendars)
	calendars = sorare_api.ArrangeResults(calendars, mode, nbGames, minGames, sequence, allGameweeks, search)

	ret := getTemplate(calendars)
	elapsed := time.Since(start)
	fmt.Printf("Done in %s!", elapsed)

	return ret
}

func getTemplate(clubs []sorare_api.ClubExport) bytes.Buffer {
	t := views.GetTemplate()
	var out bytes.Buffer
	err := t.Execute(&out, clubs)
	if err != nil {
		fmt.Println(err.Error())
	}
	return out
}
