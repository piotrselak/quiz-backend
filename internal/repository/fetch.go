package repository

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/piotrselak/back/pkg/domain"
)

func GetAllQuizes(ctx context.Context, session neo4j.SessionWithContext) ([]*domain.Quiz, error) {
	result, err := neo4j.ExecuteRead[[]*domain.Quiz](ctx, session,
		func(transaction neo4j.ManagedTransaction) ([]*domain.Quiz, error) {
			neoRecords, err := transaction.Run(ctx,
				"MATCH (quiz:Quiz) RETURN quiz",
				map[string]any{})

			if err != nil {
				return nil, err
			}

			records, err := neoRecords.Collect(ctx)
			if err != nil {
				return nil, err
			}

			var resultRecords []*domain.Quiz

			for _, record := range records {
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

				quiz := &domain.Quiz{Id: id, Name: name, Rating: rating, EditHash: editHash}

				resultRecords = append(resultRecords, quiz)
			}

			return resultRecords, nil
		})
	if err != nil {
		return nil, err
	}
	return result, nil
}
