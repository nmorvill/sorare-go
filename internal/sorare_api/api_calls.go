package sorare_api

import (
	footballapi "sorare-mu/internal/football_api"
	"sorare-mu/internal/utils"
	"sync"
)

type ClubExport struct {
	Abbreviation string       `json:"abbr"`
	Name         string       `json:"name"`
	LogoURL      string       `json:"logoURL"`
	NbTeams      int          `json:"nbTeams"`
	Games        []GameExport `json:"games"`
	Rank         int          `json:"rank"`
}

type GameExport struct {
	OpponentRank int    `json:"oppRank"`
	LogoURL      string `json:"logoURL"`
	IsHome       bool   `json:"location"`
	Color        string `json:"color"`
	Existing     bool   `json:"existing"`
	IsInSequence bool   `json:"isInSequence"`
	Gameweek     int    `json:"gameweek"`
	Streak       [3]int `json:"streak"`
}

func GetCalendars() []ClubExport {
	var ret []ClubExport
	wg := sync.WaitGroup{}
	for _, league := range getAllDomesticLeaguesSlugs() {
		wg.Add(1)
		go func(league string) {
			ret = append(ret, getAllClubsFromLeague(league)...)
			wg.Done()
		}(league)
	}
	wg.Wait()
	return ret
}

func getAllClubsFromLeague(league string) []ClubExport {
	clubsFromLeague := getClubsOfLeague(league)
	ranks, streaks := getRanksAndStreaks(clubsFromLeague, league)
	nbTeams := len(clubsFromLeague.Competition.Contestants)
	if league == "mlspa" {
		nbTeams = 14
	}
	var ret []ClubExport
	for _, club := range clubsFromLeague.Competition.Contestants {
		var c ClubExport
		c.Abbreviation = club.Team.Code
		c.Name = club.Team.Name
		c.NbTeams = nbTeams
		c.LogoURL = club.Team.PictureUrl
		c.Rank = ranks[club.Team.Slug]

		for _, game := range club.Team.UpcomingGames {
			if game.Competition.Format == "DOMESTIC_LEAGUE" {
				var g GameExport
				if club.Team.Slug == game.HomeTeam.Slug {
					g.IsHome = true
					g.OpponentRank = ranks[game.AwayTeam.Slug]
					g.Streak = streaks[game.AwayTeam.Slug]
					g.LogoURL = game.AwayTeam.PictureUrl
				} else {
					g.IsHome = false
					g.OpponentRank = ranks[game.HomeTeam.Slug]
					g.Streak = streaks[game.HomeTeam.Slug]
					g.LogoURL = game.HomeTeam.PictureUrl
				}
				g.Color = utils.GetColorCodeOfRank(g.OpponentRank, c.NbTeams)
				g.Existing = true
				g.IsInSequence = false
				g.Gameweek = utils.GetGameweekFromString(game.Date)
				c.Games = append(c.Games, g)
			}
		}

		ret = append(ret, c)
	}
	return ret
}

func getRanksAndStreaks(clubs COMPS, league string) (map[string]int, map[string][3]int) {
	var ranks map[string]int = make(map[string]int)
	var streaks map[string][3]int = make(map[string][3]int)
	var westernConferenceRanks map[string]int
	if league == "mlspa" {
		westernConferenceRanks = footballapi.GetWesternConferenceRanking()
	}
	for _, club := range clubs.Competition.Contestants {
		if club.Rank == 0 && westernConferenceRanks != nil {
			club.Rank = westernConferenceRanks[club.Team.Name]
		}
		ranks[club.Team.Slug] = club.Rank
		streaks[club.Team.Slug] = getStreak(club.Team)
	}
	return ranks, streaks
}

func getStreak(club TEAM) [3]int {
	var streak [3]int
	i := 0
	for _, game := range club.LastFiveGames {
		if game.Status == "played" && i < 3 {
			var result int
			if game.HomeGoals > game.AwayGoals {
				result = 1
			} else if game.AwayGoals > game.HomeGoals {
				result = 2
			} else {
				result = 0
			}
			if game.AwayTeam.Slug == club.Slug {
				if result == 1 {
					result = 2
				} else if result == 2 {
					result = 1
				}
			}
			streak[i] = result
			i++
		}
	}
	return streak
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
