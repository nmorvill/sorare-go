package sorare_api

import (
	"fmt"
	footballapi "sorare-mu/internal/football_api"
	"sorare-mu/internal/utils"
	"sync"
)

type Division int64

const (
	Champ   Division = 0
	Chall   Division = 1
	Asia    Division = 2
	America Division = 3
	Div2    Division = 4
)

func GetCalendars() []ClubExport {
	var ret []ClubExport
	wg := sync.WaitGroup{}
	for _, league := range getAllDomesticLeagues() {
		wg.Add(1)
		go func(league string) {
			ret = append(ret, getAllClubsFromLeague(league)...)
			wg.Done()
		}(league.Slug)
	}
	wg.Wait()
	return ret
}

func GetLeagues() []LeagueExport {
	return getAllDomesticLeagues()
}

func getAllClubsFromLeague(league string) []ClubExport {
	clubsFromLeague := getClubsOfLeague(league)
	ranks, streaks := getRanksAndStreaks(clubsFromLeague, league)
	divisions, colors := getDivisions(league)
	division, exists := divisions[clubsFromLeague.Competition.DisplayName]
	var color string
	if !exists {
		color = "#ffffff"
		fmt.Println(clubsFromLeague.Competition.DisplayName)
	} else {
		color = colors[division]
	}
	nbTeams := len(clubsFromLeague.Competition.Contestants)
	if league == "mlspa" {
		nbTeams = 15
	}
	var ret []ClubExport
	for _, club := range clubsFromLeague.Competition.Contestants {
		var c ClubExport
		c.Abbreviation = club.Team.Code
		c.Name = club.Team.Name
		c.NbTeams = nbTeams
		c.LogoURL = club.Team.PictureUrl
		c.Rank = ranks[club.Team.Slug]
		c.Color = color
		c.Slug = club.Team.Slug
		c.League = league

		for _, game := range club.Team.UpcomingGames {
			if game.Competition.Format == "DOMESTIC_LEAGUE" {
				var g GameExport
				if club.Team.Slug == game.HomeTeam.Slug {
					g.IsHome = true
					g.OpponentRank = ranks[game.AwayTeam.Slug]
					g.Streak = streaks[game.AwayTeam.Slug]
					g.LogoURL = game.AwayTeam.PictureUrl
					g.OpponentName = game.AwayTeam.Name
				} else {
					g.IsHome = false
					g.OpponentRank = ranks[game.HomeTeam.Slug]
					g.Streak = streaks[game.HomeTeam.Slug]
					g.LogoURL = game.HomeTeam.PictureUrl
					g.OpponentName = game.HomeTeam.Name
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

func getRanksAndStreaks(clubs COMPS, league string) (map[string]int, map[string][5]int) {
	var ranks map[string]int = make(map[string]int)
	var streaks map[string][5]int = make(map[string][5]int)
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

func getStreak(club TEAM) [5]int {
	var streak [5]int
	i := 0
	for _, game := range club.LastFiveGames {
		if game.Status == "played" {
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

func getAllDomesticLeagues() []LeagueExport {
	leagues := getAllLeagues()
	var ret []LeagueExport
	for _, league := range leagues.LeaguesOpenForGameStats {
		if league.Format == "DOMESTIC_LEAGUE" {
			var l LeagueExport
			l.DisplayName = league.DisplayName
			l.Slug = league.Slug
			ret = append(ret, l)
		}
	}
	return ret
}

func getDivisions(displayName string) (map[string]Division, map[Division]string) {
	divisions := map[string]Division{
		"French Ligue 1":                Champ,
		"Bundesliga":                    Champ,
		"Premier League":                Champ,
		"Serie A":                       Champ,
		"LaLiga Santander":              Champ,
		"Super League":                  Chall,
		"Jupiler Pro League":            Chall,
		"Russian Premier League":        Chall,
		"Primeira Liga":                 Chall,
		"Premiership":                   Chall,
		"Spor Toto Süper Lig":           Chall,
		"Eredivisie":                    Chall,
		"Austrian Bundesliga":           Chall,
		"Superliga":                     Chall,
		"SuperSport HNL":                Chall,
		"Eliteserien":                   Chall,
		"Major League Soccer":           America,
		"Liga MX":                       America,
		"Superliga Argentina de Fútbol": America,
		"Primera A":                     America,
		"Campeonato Brasileiro Série A": America,
		"Primera División del Perú":     America,
		"Primera División de Chile":     America,
		"Liga Pro":                      America,
		"Ligue 2":                       Div2,
		"2. Bundesliga":                 Div2,
		"Football League Championship":  Div2,
		"LaLiga Smartbank":              Div2,
		"Serie B":                       Div2,
		"J1 League":                     Asia,
		"K League 1":                    Asia,
		"Chinese Super League":          Asia,
	}
	colors := map[Division]string{
		Champ:   "#293C9B",
		America: "#D73939",
		Chall:   "#FDE617",
		Div2:    "#22C545",
		Asia:    "#8A5ED1",
	}
	return divisions, colors
}
