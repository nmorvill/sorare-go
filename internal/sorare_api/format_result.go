package sorare_api

import (
	"regexp"
	"sorare-mu/internal/utils"
	"sort"
	"strings"
	"time"
)

const RESULTS_PAGE = 50

type ClubExport struct {
	Abbreviation string       `json:"abbr"`
	Slug         string       `json:"string"`
	Name         string       `json:"name"`
	LogoURL      string       `json:"logoURL"`
	NbTeams      int          `json:"nbTeams"`
	Games        []GameExport `json:"games"`
	Rank         int          `json:"rank"`
	Color        string       `json:"color"`
	Division     Division     `json:"division"`
	League       string       `json:"league"`
}

type LeagueExport struct {
	Slug        string `json:"slug"`
	DisplayName string `json:"displayName"`
}

type GameExport struct {
	OpponentName string `json:"oppName"`
	OpponentRank int    `json:"oppRank"`
	LogoURL      string `json:"logoURL"`
	IsHome       bool   `json:"location"`
	Color        string `json:"color"`
	Existing     bool   `json:"existing"`
	IsInSequence bool   `json:"isInSequence"`
	Gameweek     int    `json:"gameweek"`
	Streak       [5]int `json:"streak"`
}

type GraphExport struct {
	Clubs       []GraphClubExport `json:"clubs"`
	GraphWidth  int               `json:"graphWidth"`
	GraphHeight int               `json:"graphHeight"`
}

type GraphClubExport struct {
	X    int        `json:"x"`
	Y    int        `json:"y"`
	Club ClubExport `json:"club"`
	Id   string     `json:"id"`
}

type TableExport struct {
	Clubs    []ClubExport `json:"clubs"`
	LastClub ClubExport   `json:"lastClub"`
	URL      string       `json:"url"`
	Page     int          `json:"page"`
	HasNext  bool         `json:"hasNext"`
}

func ArrangeTable(results []ClubExport, nbGames int, minGames int, sequence int, allGameweeks bool, search string, league string, url string, page int) TableExport {
	var res []ClubExport
	if len(search) > 0 {
		results = filterSearch(results, search)
	}
	if league != "all" && league != "" {
		results = filterLeague(results, league)
	}

	if allGameweeks {
		res = getGamesByGW(results, minGames, nbGames)
	} else {
		res = getGamesByOrder(results, minGames, nbGames)
	}

	if sequence == 0 {
		res = sortByOverallCalendar(res)
	} else {
		res = sortByBestSequence(res, sequence)
	}

	r1 := regexp.MustCompile(`page=[0-9]+`)
	start, end := page*RESULTS_PAGE, (page+1)*RESULTS_PAGE
	hasNext := true
	if end > len(res) {
		end = len(res) - 1
		hasNext = false
	}
	newURL := r1.ReplaceAllString(url, "")
	ret := TableExport{Clubs: res[start:end], LastClub: res[end], URL: newURL[1 : len(newURL)-1], Page: page + 1, HasNext: hasNext}

	return ret
}

func ArrangeGraph(results []ClubExport, nbGames int, minGames int, search string, graphWidth int, graphHeight int, league string) []GraphClubExport {
	var ret []GraphClubExport
	if len(search) > 0 {
		results = filterSearch(results, search)
	}
	if league != "all" && league != "" {
		results = filterLeague(results, league)
	}

	results = getGamesByOrder(results, minGames, nbGames)
	ret = getGraphPoints(results, graphWidth, graphHeight)

	return ret
}

func sortByBestSequence(clubs []ClubExport, maxSequence int) []ClubExport {
	sort.Slice(clubs, func(i, j int) bool {
		valueI, startI := getBestSequence(clubs[i], maxSequence)
		for k := startI; k < startI+maxSequence; k++ {
			clubs[i].Games[k].IsInSequence = true
		}
		valueJ, startJ := getBestSequence(clubs[j], maxSequence)
		for k := startJ; k < startJ+maxSequence; k++ {
			clubs[j].Games[k].IsInSequence = true
		}
		return valueI > valueJ
	})
	return clubs
}

func sortByOverallCalendar(clubs []ClubExport) []ClubExport {
	sort.Slice(clubs, func(i, j int) bool {
		return getMuStrengthOfClub(clubs[i]) > getMuStrengthOfClub(clubs[j])
	})
	return clubs
}

func getGamesByGW(clubs []ClubExport, minGames int, maxGames int) []ClubExport {
	var ret []ClubExport
	gw := utils.GetGameweekFromDate(time.Now())
	for _, club := range clubs {
		g := make([]GameExport, maxGames)
		j := 0
		for i := 0; i < maxGames; i++ {
			if j < len(club.Games) && gw+i == club.Games[j].Gameweek {
				g[i] = club.Games[j]
				j++
			} else {
				g[i] = GameExport{Existing: false}
			}
		}
		club.Games = g
		if j >= minGames {
			ret = append(ret, club)
		}
	}
	return ret
}

func getGamesByOrder(clubs []ClubExport, minGames int, maxGames int) []ClubExport {
	var ret []ClubExport
	for _, club := range clubs {
		if len(club.Games) >= minGames {
			for len(club.Games) < maxGames {
				club.Games = append(club.Games, GameExport{Existing: false})
			}
			club.Games = club.Games[:maxGames]
			ret = append(ret, club)
		}
	}
	return ret
}

func getMuStrengthOfClub(club ClubExport) float32 {
	var sum float32 = 0
	var count float32 = 0
	for _, game := range club.Games {
		if game.Existing {
			sum += (float32(game.OpponentRank) / float32(club.NbTeams))
			count++
		}
	}
	return sum / count
}

func filterSearch(clubs []ClubExport, search string) []ClubExport {
	var ret []ClubExport
	search = strings.ToLower(search)
	for _, club := range clubs {
		if strings.Contains(strings.ToLower(club.Name), search) {
			ret = append(ret, club)
		}
	}
	return ret
}

func filterLeague(clubs []ClubExport, league string) []ClubExport {
	var ret []ClubExport
	for _, club := range clubs {
		if club.League == league {
			ret = append(ret, club)
		}
	}
	return ret
}

func getBestSequence(club ClubExport, maxSequence int) (float32, int) {
	var res []int
	for _, game := range club.Games {
		if game.Existing {
			res = append(res, game.OpponentRank)
		}
	}
	ret, start, err := utils.MaxSum(res, maxSequence)
	if err != nil {
		//fmt.Println(err.Error())
	}
	return float32(ret) / float32(club.NbTeams), start
}

func getGraphPoints(clubs []ClubExport, graphWidth int, graphHeight int) []GraphClubExport {
	var ret []GraphClubExport
	for _, club := range clubs {
		var x int = int(getMuStrengthOfClub(club) * float32(graphWidth))
		var rank float32 = float32(club.Rank) / float32(club.NbTeams)
		var y int = int(float32(graphHeight) * rank)
		if len(club.Color) == 0 {
			club.Color = "white"
		}

		ret = append(ret, GraphClubExport{Club: club, X: x, Y: y, Id: utils.ClearString(club.Slug)})
	}
	return ret
}
