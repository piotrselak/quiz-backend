package domain

import (
	"fmt"
)

// Quiz Neo4j Node
type Quiz struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Rating   float64 `json:"rating"`
	EditHash string  `json:"editHash"`
}

func (quiz Quiz) ToCypher(char string) Cypher {
	//q, _ := json.Marshal(quiz)
	//properties := string(q)
	//return fmt.Sprintf("(%s:Quiz %s)", char, properties)
	properties := fmt.Sprintf("{id: '%s', name: '%s', rating: '%.1f', editHash: '%s'}",
		quiz.Id, quiz.Name, quiz.Rating, quiz.EditHash)
	return fmt.Sprintf("(%s:Quiz %s)", char, properties)
}

type QuizForPost struct {
	Name     string `json:"name"`
	EditHash string `json:"editHash"`
}

func (quiz QuizForPost) ToQuiz(id string, rating float64) Quiz {
	return Quiz{Id: id, Name: quiz.Name, Rating: rating, EditHash: quiz.EditHash}
}

// Has Neo4j Relationship
type Has struct{}

func (r *Has) ToCypherRight(char string) Cypher {
	return fmt.Sprintf("-[%s:Has]->")
}

func (r *Has) ToCypherLeft(char string) Cypher {
	return fmt.Sprintf("<-[%s:Has]-")
}
