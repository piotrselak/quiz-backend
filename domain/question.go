package domain

import (
	"encoding/json"
	"fmt"

	"github.com/piotrselak/back/modules/db"
)

type Question struct {
	QuestionText string        `json:"questionText"`
	Answers      []interface{} `json:"answers"`
}

func (question Question) ToCypher(char string) db.Cypher {
	q, _ := json.Marshal(question)
	properties := string(q)
	return fmt.Sprintf("(%s:Question %s)", char, properties)
}

type Answer struct {
	Answer []interface{} `json:"answer"`
}

func (answer Answer) ToCypher(char string) db.Cypher {
	q, _ := json.Marshal(answer)
	properties := string(q)
	return fmt.Sprintf("(%s:Answer %s)", char, properties)
}
