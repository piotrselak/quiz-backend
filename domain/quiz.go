package domain

import (
	"encoding/json"
	"fmt"

	"github.com/piotrselak/back/modules/db"
)

// Quiz Neo4j Node
type Quiz struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Rating   float64 `json:"rating"`
	EditHash string  `json:"editHash"`
}

func (quiz Quiz) ToCypher(char string) db.Cypher {
	q, _ := json.Marshal(quiz)
	properties := string(q)
	return fmt.Sprintf("(%s:Quiz %s)", char, properties)
}

// Has Neo4j Relationship
type Has struct{}

func (r *Has) ToCypherRight(char string) db.Cypher {
	return fmt.Sprintf("-[%s:Has]->")
}

func (r *Has) ToCypherLeft(char string) db.Cypher {
	return fmt.Sprintf("<-[%s:Has]-")
}
