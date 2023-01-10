package repository

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/piotrselak/back/domain"
)

// ToQuiz - converts neo4j record to quiz
func toQuiz(record *neo4j.Record) (*domain.Quiz, error) {
	rawItemNode, found := record.Get("quiz")
	if !found {
		return nil, fmt.Errorf("could not find column")
	}
	itemNode := rawItemNode.(neo4j.Node)

	id, err := neo4j.GetProperty[string](itemNode, "id")
	if err != nil {
		return nil, err
	}

	name, err := neo4j.GetProperty[string](itemNode, "name")
	if err != nil {
		return nil, err
	}

	rating, err := neo4j.GetProperty[float64](itemNode, "rating")
	if err != nil {
		return nil, err
	}

	editHash, err := neo4j.GetProperty[string](itemNode, "editHash")
	if err != nil {
		return nil, err
	}

	return &domain.Quiz{Id: id, Name: name, Rating: rating, EditHash: editHash}, nil
}
