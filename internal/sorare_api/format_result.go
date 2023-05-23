package sorare_api

import (
	"sorare-mu/internal/utils"
	"sort"
	"time"
)

func ArrangeResults(results []ClubExport, mode string, nbGames int, minGames int, sequence int, allGameweeks bool) []ClubExport {
	var ret []ClubExport
	if allGameweeks {
		ret = getGamesByGW(results, minGames, nbGames)
	} else {
		ret = getGamesByOrder(results, minGames, nbGames)
	}

	if mode == "Calendar" {
		ret = sortByOverallCalendar(ret)
	} else if mode == "Sequence" {
		ret = sortByBestSequence(ret, sequence)
	}
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
