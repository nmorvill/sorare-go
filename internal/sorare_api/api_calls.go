package sorare_api

import "sort"

type ClubExport struct {
	Abbreviation string       `json:"abbr"`
	Name         string       `json:"name"`
	LogoURL      string       `json:"logoURL"`
	NbTeams      int          `json:"nbTeams"`
	Games        []GameExport `json:"games"`
}

type GameExport struct {
	OpponentRank int    `json:"oppRank"`
	LogoURL      string `json:"logoURL"`
	Location     string `json:"location"`
}

func GetCalendars() []ClubExport {
	var ret []ClubExport
	for _, league := range getAllDomesticLeaguesSlugs() {
		ret = append(ret, getAllClubsFromLeague(league)...)
	}
	sort.Slice(ret, func(i, j int) bool {
		return (getMuStrengthOfClub(ret[i]) > getMuStrengthOfClub(ret[j]))
	})
	return ret
}

func getAllClubsFromLeague(league string) []ClubExport {
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
			if game.Competition.Format == "DOMESTIC_LEAGUE" && nbGames < 5 {
				var g GameExport
				if club.Team.Slug == game.HomeTeam.Slug {
					g.Location = "HOME"
					g.OpponentRank = ranks[game.AwayTeam.Slug]
					g.LogoURL = game.AwayTeam.PictureUrl
				} else {
					g.Location = "AWAY"
					g.OpponentRank = ranks[game.HomeTeam.Slug]
					g.LogoURL = game.HomeTeam.PictureUrl
				}
				c.Games = append(c.Games, g)
				nbGames++
			}
		}
		ret = append(ret, c)
	}
	return ret
}

func getMuStrengthOfClub(club ClubExport) float32 {
	var sum float32 = 0
	l := len(club.Games)
	if l == 0 {
		return 0
	}
	for _, game := range club.Games {
		sum += (float32(game.OpponentRank) / float32(club.NbTeams))
	}
	return sum / float32(l)
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
