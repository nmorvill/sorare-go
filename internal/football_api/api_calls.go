package footballapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RANKING struct {
	Response []struct {
		League struct {
			Standings [][]struct {
				Rank int `json:"rank"`
				Team struct {
					Name string `json:"name"`
				} `json:"team"`
			} `json:"standings"`
		} `json:"league"`
	} `json:"response"`
}

func requestFootballAPI[K interface{}](parameters string) K {
	url := "https://v3.football.api-sports.io/" + parameters
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "b69607fed21713f819d2adb72e6b0012")
	req.Header.Add("X-RapidAPI-Host", "api-football-v1.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var result K
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error unmarshalling ", err.Error())
	}
	return result
}

func GetWesternConferenceRanking() map[string]int {
	rankings := make(map[string]int)
	nameCorrespondance := map[string]string{
		"Austin":               "Austin FC",
		"Minnesota United FC":  "Minnesota United",
		"Vancouver Whitecaps":  "Vancouver Whitecaps FC",
		"St. Louis City":       "St. Louis City SC",
		"Los Angeles Galaxy":   "LA Galaxy",
		"Houston Dynamo":       "Houston Dynamo FC",
		"Seattle Sounders":     "Seattle Sounders FC",
		"Real Salt Lake":       "Real Salt Lake",
		"San Jose Earthquakes": "San Jose Earthquakes",
		"Sporting Kansas City": "Sporting Kansas City ",
		"Los Angeles FC":       "Los Angeles FC",
		"FC Dallas":            "FC Dallas",
		"Portland Timbers":     "Portland Timbers",
		"Colorado Rapids":      "Colorado Rapids",
	}
	res := requestFootballAPI[RANKING]("standings?season=2023&league=253")
	for _, club := range res.Response[0].League.Standings[1] {
		rankings[nameCorrespondance[club.Team.Name]] = club.Rank
	}
	return rankings
}
