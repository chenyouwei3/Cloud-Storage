package utils

import "strconv"

func GetPage(currPage, pageSize string) (int, int, error) {
	//当前页码,每页大小
	curr, err := strconv.Atoi(currPage)
	if err != nil {
		return 0, 0, err
	}
	size, err := strconv.Atoi(pageSize)
	if err != nil {
		return 0, 0, err
	}
	skip := (curr - 1) * size
	return skip, size, nil
}
