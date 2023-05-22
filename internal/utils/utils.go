package utils

import (
	"errors"
)

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxSum(arr []int, k int) (int, int, error) {
	if len(arr) < k {
		return 0, 0, errors.New("Array to small")
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
