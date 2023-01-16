package domain

import (
	"fmt"
)

// !!! question for fetch shouldnt have valid answers in it

type Question struct {
	Index        int64      `json:"index"`
	QuestionText string     `json:"questionText"`
	Answers      [][]string `json:"answers"`
	Type         string     `json:"type"`
}

func (question Question) ToCypher(char string) Cypher {
	properties := fmt.Sprintf("{index: '%d', questionText: '%s', answers: '%s', type: '%s'}",
		question.Index, question.QuestionText, question.Answers, question.Type)
	return fmt.Sprintf("(%s:Question %s)", char, properties)
}

type QuestionForPost struct {
	Data []Question `json:"data"`
}
