package verification

import (
	"github.com/piotrselak/back/domain"
	"sort"
)

// ParseAnswers returns number of valid answers
func ParseAnswers(userQuestions []domain.QuestionForFetch, validQuestions []domain.Question) int {
	sort.Slice(userQuestions, func(i, j int) bool {
		return userQuestions[i].Index < userQuestions[j].Index
	})
	sort.Slice(validQuestions, func(i, j int) bool {
		return validQuestions[i].Index < validQuestions[j].Index
	})
	var counter int = 0
	for i := 0; i < len(userQuestions); i++ {
		if equal(userQuestions[i].Answers, validQuestions[i].Answers) { // here might be an error
			counter += 1
		}
	}
	return counter
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
