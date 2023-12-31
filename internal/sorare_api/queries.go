package sorare_api

import (
	"context"
	"os"

	"github.com/machinebox/graphql"
)

type LOFGS struct {
	LeaguesOpenForGameStats []struct {
		Slug        string `json:"slug"`
		Format      string `json:"format"`
		DisplayName string `json:"displayName"`
	} `json:"leaguesOpenForGameStats"`
}

type TEAM struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	PictureUrl    string `json:"pictureUrl"`
	Slug          string `json:"slug"`
	LastFiveGames []struct {
		Status    string `json:"status"`
		HomeGoals int    `json:"homeGoals"`
		AwayGoals int    `json:"awayGoals"`
		HomeTeam  struct {
			Slug string `json:"slug"`
		} `json:"homeTeam"`
		AwayTeam struct {
			Slug string `json:"slug"`
		} `json:"awayTeam"`
	} `json:"lastFiveGames"`
	UpcomingGames []struct {
		Date        string `json:"date"`
		Competition struct {
			Format string `json:"format"`
		} `json:"competition"`
		HomeTeam struct {
			Slug       string `json:"slug"`
			PictureUrl string `json:"pictureUrl"`
			Name       string `json:"name"`
		} `json:"homeTeam"`
		AwayTeam struct {
			Slug       string `json:"slug"`
			PictureUrl string `json:"pictureUrl"`
			Name       string `json:"name"`
		} `json:"awayTeam"`
	} `json:"upcomingGames"`
}

type COMPS struct {
	Competition struct {
		DisplayName string `json:"displayName"`
		Contestants []struct {
			Rank int  `json:"rank"`
			Team TEAM `json:"team"`
		} `json:"contestants"`
	} `json:"competition"`
}

func callSorareApi[K interface{}](req *graphql.Request) K {
	api_key := os.Getenv("SORARE_API_KEY")
	client := graphql.NewClient("https://api.sorare.com/graphql")
	req.Header.Set("APIKEY", api_key)
	var ret K
	if err := client.Run(context.Background(), req, &ret); err != nil {
		panic(err)
	}
	return ret
}

func getAllLeagues() LOFGS {
	q := `
	{
		leaguesOpenForGameStats
		{
		  slug
		  format
		  displayName
		}
	  }
	`
	res := callSorareApi[LOFGS](graphql.NewRequest(q))
	return res
}

func getClubsOfLeague(league string) COMPS {
	q := graphql.NewRequest(`
	query($slug: String!, $start: Int!) {
		competition(slug:$slug) {
			displayName
			contestants(seasonStartYear:$start) {
				rank
				team {
					... on Club {
						code
						name
						pictureUrl
						slug
						lastFiveGames {
							status
							homeGoals
							awayGoals
							homeTeam {
								... on Club {
								slug
								}
							}
							awayTeam {
								... on Club {
								slug
								}
							}
						}
						upcomingGames(first:15) {
							date
							competition {
								format
							}
							homeTeam {
								... on Club {
									slug
									pictureUrl
									name
								}
							}
							awayTeam {
								... on Club {
									slug
									pictureUrl
									name
								}
							}
						}
					}
				}
			}
		}
	}
	`)
	q.Var("slug", league)
	q.Var("start", getStartYearOfLeague(league))
	res := callSorareApi[COMPS](q)
	return res
}

func getStartYearOfLeague(league string) int {
	leagues := map[string]int{
		"mlspa":                         2023,
		"j1-league":                     2023,
		"campeonato-brasileiro-serie-a": 2023,
		"superliga-argentina-de-futbol": 2023,
		"eliteserien":                   2023,
		"k-league-1":                    2023,
		"primera-division-cl":           2023,
		"liga-pro-ec":                   2023,
		"chinese-super-league":          2023,
	}
	value, is := leagues[league]
	if is {
		return value
	} else {
		return 2022
	}
}
