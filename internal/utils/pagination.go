package utils

import (
	"strconv"
)

func CalculatePagination(nbElements int, sizePage int, requestedPage string) (nbPage int, currentPage int, indexStart int, indexEnd int) {
	nbPage = nbElements/sizePage + MinIntegerOf(nbElements%sizePage, 1)

	currentPage, err := strconv.Atoi(requestedPage)
	if err != nil || currentPage < 1 {
		currentPage = 1
	}
	currentPage = MinIntegerOf(currentPage, nbPage)

	indexStart = MaxIntegerOf(currentPage-1, 0) * sizePage
	indexEnd = MinIntegerOf(indexStart+sizePage, nbElements)
	return nbPage, currentPage, indexStart, indexEnd
}
