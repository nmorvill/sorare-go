package sorare_api

import (
	"context"
	"os"

	"github.com/machinebox/graphql"
)

type LOFGS struct {
	LeaguesOpenForGameStats []struct {
		Slug   string `json:"slug"`
		Format string `json:"format"`
	} `json:"leaguesOpenForGameStats"`
}

type COMPS struct {
	Competition struct {
		Contestants []struct {
			Rank int `json:"rank"`
			Team struct {
				Code          string `json:"code"`
				Name          string `json:"name"`
				PictureUrl    string `json:"pictureUrl"`
				Slug          string `json:"slug"`
				UpcomingGames []struct {
					Date        string `json:"date"`
					Competition struct {
						Format string `json:"format"`
					} `json:"competition"`
					HomeTeam struct {
						Slug       string `json:"slug"`
						PictureUrl string `json:"pictureUrl"`
					} `json:"homeTeam"`
					AwayTeam struct {
						Slug       string `json:"slug"`
						PictureUrl string `json:"pictureUrl"`
					} `json:"awayTeam"`
				} `json:"upcomingGames"`
			} `json:"team"`
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
			contestants(seasonStartYear:$start) {
				rank
				team {
					... on Club {
						code
						name
						pictureUrl
						slug
						upcomingGames(first:10) {
							date
							competition {
								format
							}
							homeTeam {
								... on Club {
									slug
									pictureUrl
								}
							}
							awayTeam {
								... on Club {
									slug
									pictureUrl
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
	if league == "mlspa" {
		return 2023
	} else {
		return 2022
	}
}
