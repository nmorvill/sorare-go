package sorare_api

import (
	"fmt"
	"sorare-mu/internal/utils"
	"sort"
	"strconv"
	"sync"
)

const MinGames = 3
const MaxSequence = 5

type ClubExport struct {
	Abbreviation string       `json:"abbr"`
	Name         string       `json:"name"`
	LogoURL      string       `json:"logoURL"`
	NbTeams      int          `json:"nbTeams"`
	Games        []GameExport `json:"games"`
	BestSequence float32      `json:"bestSequence"`
}

type GameExport struct {
	OpponentRank int    `json:"oppRank"`
	LogoURL      string `json:"logoURL"`
	IsHome       bool   `json:"location"`
	Color        string `json:"color"`
	Existing     bool   `json:"existing"`
	IsInSequence bool   `json:"isInSequence"`
}

func GetBestSequence() []ClubExport {
	var ret []ClubExport
	wg := sync.WaitGroup{}
	for _, league := range getAllDomesticLeaguesSlugs() {
		wg.Add(1)
		go func(league string) {
			ret = append(ret, getAllClubsFromLeague(league, 5)...)
			wg.Done()
		}(league)
	}
	wg.Wait()
	sort.Slice(ret, func(i, j int) bool {
		return float32(ret[i].BestSequence) > float32(ret[j].BestSequence)
	})
	return ret
}

func GetCalendars() []ClubExport {
	var ret []ClubExport
	wg := sync.WaitGroup{}
	for _, league := range getAllDomesticLeaguesSlugs() {
		wg.Add(1)
		go func(league string) {
			ret = append(ret, getAllClubsFromLeague(league, 5)...)
			wg.Done()
		}(league)
	}
	wg.Wait()
	sort.Slice(ret, func(i, j int) bool {
		return getMuStrengthOfClub(ret[i]) > getMuStrengthOfClub(ret[j])
	})
	return ret
}

func getAllClubsFromLeague(league string, minGames int) []ClubExport {
	clubsFromLeague := getClubsOfLeague(league)
	ranks := getAllRanks(clubsFromLeague)
	nbTeams := len(clubsFromLeague.Competition.Contestants)
	var ret []ClubExport
	for _, club := range clubsFromLeague.Competition.Contestants {
		var c ClubExport
		c.Abbreviation = club.Team.Code
		c.Name = club.Team.Name
		c.NbTeams = nbTeams
		c.LogoURL = club.Team.PictureUrl
		nbGames := 0

		for _, game := range club.Team.UpcomingGames {
			if game.Competition.Format == "DOMESTIC_LEAGUE" && nbGames < 10 {
				var g GameExport
				if club.Team.Slug == game.HomeTeam.Slug {
					g.IsHome = true
					g.OpponentRank = ranks[game.AwayTeam.Slug]
					g.LogoURL = game.AwayTeam.PictureUrl
				} else {
					g.IsHome = false
					g.OpponentRank = ranks[game.HomeTeam.Slug]
					g.LogoURL = game.HomeTeam.PictureUrl
				}
				g.Color = getColorCodeOfRank(g.OpponentRank, c.NbTeams)
				g.Existing = true
				g.IsInSequence = false
				c.Games = append(c.Games, g)
				nbGames++
			}
		}

		for len(c.Games) < 10 {
			c.Games = append(c.Games, GameExport{Existing: false})
		}

		if nbGames >= minGames {
			best, start := getBestSequence(c)
			c.BestSequence = best
			for i := start; i < start+MaxSequence; i++ {
				c.Games[i].IsInSequence = true
			}
			ret = append(ret, c)
		}
	}
	return ret
}

func getColorCodeOfRank(rank int, maxRank int) string {
	if rank == 0 {
		return "#000000"
	}
	var green [3]int = [3]int{87, 223, 72}
	var red [3]int = [3]int{255, 67, 57}
	per := float32(rank) / float32(maxRank)
	r := red[0] + int(per*(float32(green[0]-red[0])))
	g := red[1] + int(per*(float32(green[1]-red[1])))
	b := red[2] + int(per*(float32(green[2]-red[2])))
	return "#" + strconv.FormatInt(int64(r), 16) + strconv.FormatInt(int64(g), 16) + strconv.FormatInt(int64(b), 16)
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

func getBestSequence(club ClubExport) (float32, int) {
	var res []int
	for _, game := range club.Games {
		if game.Existing {
			res = append(res, game.OpponentRank)
		}
	}
	ret, start, err := utils.MaxSum(res, MaxSequence)
	if err != nil {
		fmt.Println(err.Error())
	}
	return float32(ret) / float32(club.NbTeams), start
}

func getAllRanks(clubs COMPS) map[string]int {
	var ret map[string]int = make(map[string]int)
	for _, club := range clubs.Competition.Contestants {
		ret[club.Team.Slug] = club.Rank
	}
	return ret
}

func getAllDomesticLeaguesSlugs() []string {
	leagues := getAllLeagues()
	var ret []string
	for _, league := range leagues.LeaguesOpenForGameStats {
		if league.Format == "DOMESTIC_LEAGUE" {
			ret = append(ret, league.Slug)
		}
	}
	return ret
}
