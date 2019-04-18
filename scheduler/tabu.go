package scheduler

import (
	"math/rand"
	"time"
)

// TabuList collects a list of moves to not take in the next iterations
type TabuList struct {
	tabuMoves    []int
	size         int
	currentIndex int
}

func newTabuList(size int) *TabuList {
	tabuList := &TabuList{
		tabuMoves:    make([]int, size),
		size:         size,
		currentIndex: 0,
	}

	for index := range tabuList.tabuMoves {
		tabuList.tabuMoves[index] = -1
	}

	return tabuList
}

func (tabuList *TabuList) get(numberOfOperators int) int {
	var op int

	rand.Seed(time.Now().UTC().UnixNano())

	isTabu := true

	for isTabu == true {
		op = randInt(0, numberOfOperators)
		isTabu = false

		for _, tabuMove := range tabuList.tabuMoves {
			if tabuMove == op {
				isTabu = true
				break
			}
		}
	}

	return op
}

func (tabuList *TabuList) update(operator int) {
	tabuList.tabuMoves[tabuList.currentIndex] = operator
	tabuList.currentIndex = (tabuList.currentIndex + 1) % tabuList.size
}

func (tabuList *TabuList) clear() {
	for index := range tabuList.tabuMoves {
		tabuList.tabuMoves[index] = -1
	}
	tabuList.currentIndex = 0
}
