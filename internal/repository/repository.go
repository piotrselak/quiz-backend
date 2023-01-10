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
				quiz, err := toQuiz(record)
				if err != nil {
					return nil, err
				}

				resultRecords = append(resultRecords, quiz)
			}

			return resultRecords, nil
		})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CreateQuiz(ctx context.Context, session neo4j.SessionWithContext,
	quizWithQuestions domain.QuizWithQuestions) error {
	_, err := neo4j.ExecuteWrite[*domain.QuizWithQuestions](ctx, session,
		func(transaction neo4j.ManagedTransaction) (*domain.QuizWithQuestions, error) {
			questions := quizWithQuestions.Questions
			var cypherScript = ""

			for ind, question := range questions {
				var splitChar string
				if ind == 0 {
					splitChar = ""
				} else {
					splitChar = ", "
				}
				questionChar := fmt.Sprintf("q%s", ind)
				cypherScript = cypherScript + fmt.Sprintf("%s-[:Has]->%s%s",
					question.ToCypher(questionChar), splitChar)

			}

			_, err := transaction.Run(ctx,
				fmt.Sprintf("CREATE %s", cypherScript),
				map[string]any{})

			if err != nil {
				return nil, err
			}

			return &quizWithQuestions, nil
		})
	if err != nil {
		return err
	}
	return nil
}
