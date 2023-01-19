package domain

import (
	"fmt"
)

// Quiz Neo4j Node
type Quiz struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Rating    float64  `json:"rating"`
	EditHash  string   `json:"editHash"`
	Modifiers []string `json:"modifiers"`
}

type QuizForPost struct {
	Name      string   `json:"name"`
	EditHash  string   `json:"editHash"`
	Modifiers []string `json:"modifiers"`
}

// Has Neo4j Relationship
type Has struct{}

func (quiz Quiz) ToCypher(char string) Cypher {
	properties := fmt.Sprintf("{id: '%s', name: '%s', rating: %.1f, editHash: '%s', modifiers: %s}",
		quiz.Id, quiz.Name, quiz.Rating, quiz.EditHash, quiz.Modifiers)
	return fmt.Sprintf("(%s:Quiz %s)", char, properties)
}

func (r Has) ToCypherRight(char string) Cypher {
	return fmt.Sprintf("-[%s:Has]->", char)
}

func (r Has) ToCypherLeft(char string) Cypher {
	return fmt.Sprintf("<-[%s:Has]-", char)
}

func (quiz QuizForPost) ToQuiz(id string, rating float64) Quiz {
	return Quiz{Id: id, Name: quiz.Name, Rating: rating,
		EditHash: quiz.EditHash, Modifiers: quiz.Modifiers}
}
