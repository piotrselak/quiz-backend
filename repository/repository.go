package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/piotrselak/back/domain"
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
	quiz domain.QuizForPost) error {
	_, err := neo4j.ExecuteWrite[*domain.QuizForPost](ctx, session,
		func(transaction neo4j.ManagedTransaction) (*domain.QuizForPost, error) {

			_, err := transaction.Run(ctx,
				fmt.Sprintf("CREATE %s", quiz.ToQuiz(uuid.New().String(), 0, 0).ToCypher("q")),
				map[string]any{})

			if err != nil {
				return nil, err
			}

			return &quiz, nil
		})
	if err != nil {
		return err
	}
	return nil
}

func AddQuestions(ctx context.Context, session neo4j.SessionWithContext,
	id string, questions domain.QuestionForPost) error {
	_, err := neo4j.ExecuteWrite[*domain.QuestionForPost](ctx, session,
		func(transaction neo4j.ManagedTransaction) (*domain.QuestionForPost, error) {

			var cypherQuestions = ""

			for index, question := range questions.Data {
				newQuestion := fmt.Sprintf("(q)%s%s",
					domain.Has{}.ToCypherRight(""), question.ToCypher(""))
				var split string
				if index == len(questions.Data)-1 {
					split = ""
				} else {
					split = ", "
				}
				cypherQuestions = cypherQuestions + newQuestion + split
			}

			_, err := transaction.Run(ctx,
				fmt.Sprintf("MATCH (q:Quiz {id: '%s'}) CREATE %s", id, cypherQuestions),
				map[string]any{})
			fmt.Println(cypherQuestions)
			fmt.Println(err)

			if err != nil {
				return nil, err
			}

			return &questions, nil
		})
	if err != nil {
		return err
	}
	return nil
}
