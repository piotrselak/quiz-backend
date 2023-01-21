package domain

import (
	"encoding/json"
	"fmt"
)

type Question struct {
	Index        int64    `json:"index"`
	QuestionText string   `json:"questionText"`
	Answers      []string `json:"answers"`
	ValidAnswers []string `json:"validAnswers"`
	Type         string   `json:"type"`
}

func (question Question) ToCypher(char string) Cypher {
	answers, _ := json.Marshal(question.Answers)
	validAnswers, _ := json.Marshal(question.ValidAnswers)
	properties := fmt.Sprintf("{index: %d, questionText: '%s', answers: %s, validAnswers: %s, type: '%s'}",
		question.Index, question.QuestionText, string(answers), string(validAnswers), question.Type)
	return fmt.Sprintf("(%s:Question %s)", char, properties)
}

type QuestionForPost struct {
	Data     []Question `json:"data"`
	EditHash string     `json:"editHash"`
}

// QuestionForFetch is used for fetching questions when frontend should not get valid answers to question
type QuestionForFetch struct {
	Index        int64    `json:"index"`
	QuestionText string   `json:"questionText"`
	Answers      []string `json:"answers"`
	Type         string   `json:"type"`
}
