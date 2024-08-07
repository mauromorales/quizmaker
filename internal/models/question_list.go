package models

import (
	"math/rand/v2"
	"sort"
)

type QuestionList []Question

func (ql QuestionList) InDifficultyRange(min, max int) QuestionList {
	result := QuestionList{}
	for _, q := range ql {
		if q.Difficulty <= max && q.Difficulty >= min {
			result = append(result, q)
		}
	}

	return result
}

func (ql QuestionList) Suffled() QuestionList {
	dest := make(QuestionList, len(ql))
	perm := rand.Perm(len(ql))
	for i, v := range perm {
		dest[v] = ql[i]
	}

	return dest
}

func (ql QuestionList) OrderedByDifficulty() QuestionList {
	sort.Slice(ql, func(i, j int) bool {
		return ql[i].Difficulty < ql[j].Difficulty
	})

	return ql
}

// Limit return exactly `limit` number of questions. If the total questions
// are more than the requests, it will remove the additional ones, with preference
// to the questions with difficulty that is already represented. For example
// If there are 5 total questions like this:
// - 2 with difficulty 1
// - 1 with difficulty 2
// - 2 with difficutly 3
// when limit is "3" the method will return one questions from each difficulty
// (instead for example: 2 with difficulty 1 and 1 with difficulty 2).
// In other words, it tries to have questions from all difficulties if possible.
// Desired questions are picked starting from lower to higher levels
// but when a difficulty level runs out of questions, one from a higher level is
// picked instead.
func (ql QuestionList) Limit(limit int) QuestionList {
	difficultiesMap := map[int]QuestionList{}

	for _, q := range ql {
		if difficultiesMap[q.Difficulty] == nil { // initialize map
			difficultiesMap[q.Difficulty] = QuestionList{q}
		} else {
			difficultiesMap[q.Difficulty] = append(difficultiesMap[q.Difficulty], q)
		}
	}

	uniqueDifficulties := []int{}
	for k := range difficultiesMap {
		uniqueDifficulties = append(uniqueDifficulties, k)
	}
	sort.Ints(uniqueDifficulties)

	result := QuestionList{}
	totalAvailableQuestions := len(ql)
	pickedQuestions := 0
	// loops the difficutlies over and over again until the desired number of
	// questions are selected.
	for pickedQuestions < limit && pickedQuestions <= totalAvailableQuestions {
		for _, d := range uniqueDifficulties {
			var q Question
			if len(difficultiesMap[d]) == 0 {
				continue // no more questions in this difficulty
			}

			// pop a question from this difficulty
			q, difficultiesMap[d] = difficultiesMap[d][0], difficultiesMap[d][1:]
			result = append(result, q)

			// stop here if we are done
			pickedQuestions++
			if pickedQuestions == limit || pickedQuestions == totalAvailableQuestions {
				break
			}
		}
	}

	return result
}
