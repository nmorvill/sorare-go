package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func MaxSum(arr []int, k int) (int, int, error) {
	if len(arr) < k {
		return 0, 0, errors.New("Array too small")
	}

	res := 0
	for i := 0; i < k; i++ {
		res += arr[i]
	}

	startingIndex := 0
	currSum := res
	for i := k; i < len(arr); i++ {
		currSum += arr[i] - arr[i-k]
		if currSum > res {
			startingIndex = i - k + 1
			res = currSum
		}
	}

	return res, startingIndex, nil
}

func GetColorCodeOfRank(rank int, maxRank int) string {
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

func GetGameweekFromDate(date time.Time) int {
	gw374 := time.Date(2023, time.May, 23, 14, 0, 0, 0, time.Local)
	daysDiff := -int(gw374.Sub(date).Hours() / 24)
	gwDiff := (daysDiff / 7) * 2
	if daysDiff%7 >= 3 {
		gwDiff += 1
	}
	return gwDiff + 374
}

func GetGameweekFromString(str string) int {
	date, err := time.Parse("2006-01-02T15:04:05Z", str)
	if err != nil {
		return 0
	}
	gw374 := time.Date(2023, time.May, 23, 14, 0, 0, 0, time.Local)
	daysDiff := -int(gw374.Sub(date).Hours() / 24)
	gwDiff := (daysDiff / 7) * 2
	if daysDiff%7 >= 3 {
		gwDiff += 1
	}
	return gwDiff + 374
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func ClearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(strings.Join(strings.Fields(str), ""), "")
}
