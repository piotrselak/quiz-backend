package domain

import (
	"encoding/json"
	"fmt"
)

type Question struct {
	Index        int64      `json:"index"`
	QuestionText string     `json:"questionText"`
	Answers      [][]string `json:"answers"`
	Type         string     `json:"type"`
}

func (question Question) ToCypher(char string) Cypher {
	answers, _ := json.Marshal(question.Answers)
	properties := fmt.Sprintf("{index: %d, questionText: '%s', answers: %s, type: '%s'}",
		question.Index, question.QuestionText, string(answers), question.Type)
	return fmt.Sprintf("(%s:Question %s)", char, properties)
}

type QuestionForPost struct {
	Data []Question `json:"data"`
}

// QuestionForFetch is used for fetching questions when frontend should not get valid answers to question
type QuestionForFetch struct {
	Index        int64    `json:"index"`
	QuestionText string   `json:"questionText"`
	Answers      []string `json:"answers"`
	Type         string   `json:"type"`
}
