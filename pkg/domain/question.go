package domain

import (
	"encoding/json"
	"fmt"
)

type Question struct {
	QuestionText string        `json:"questionText"`
	Answers      []interface{} `json:"answers"`
}

func (question Question) ToCypher(char string) Cypher {
	//q, _ := json.Marshal(question)
	//properties := string(q)
	//return fmt.Sprintf("(%s:Question %s)", char, properties)
	properties := fmt.Sprintf("{questionText: '%s', answers: '%s'}",
		question.QuestionText, question.Answers)
	return fmt.Sprintf("(%s:Question %s)", char, properties)
}

type Answer struct {
	Answer []interface{} `json:"answer"`
}

func (answer Answer) ToCypher(char string) Cypher {
	q, _ := json.Marshal(answer)
	properties := string(q)
	// maybe remove "" in properties
	return fmt.Sprintf("(%s:Answer %s)", char, properties)
}
