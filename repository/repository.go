package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"reflect"

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
				fmt.Sprintf("CREATE %s", quiz.ToQuiz(uuid.New().String(), 0).ToCypher("q")),
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
	res := checkIfQuizExists(ctx, session, id)
	if !res {
		return errors.New("not found")
	}

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

func RemoveQuiz(ctx context.Context, session neo4j.SessionWithContext, id string) error {
	res := checkIfQuizExists(ctx, session, id)
	if !res {
		return errors.New("not found")
	}

	cypherScript := fmt.Sprintf(
		"MATCH (player:Player)-[played:Played]->(q:Quiz {id: '%s'})"+
			"-[has:Has]->(question:Question) DELETE player, played, q, has, question", id)
	_, _ = session.Run(ctx, cypherScript, map[string]any{})

	cypherScript = fmt.Sprintf("MATCH (q2:Quiz {id: '%s'})-[has2:Has]->"+
		"(question2:Question) DELETE q2, has2, question2", id)
	_, _ = session.Run(ctx, cypherScript, map[string]any{})

	return nil
}

func FetchSpecificQuiz(ctx context.Context, session neo4j.SessionWithContext, id string) (domain.QuizWithQuestions, error) {
	result, err := neo4j.ExecuteRead[*domain.QuizWithQuestions](ctx, session,
		func(transaction neo4j.ManagedTransaction) (*domain.QuizWithQuestions, error) {
			cypherScript := fmt.Sprintf("MATCH (quiz:Quiz {id: '%s'})"+
				"-[:Has]->(question: Question) RETURN quiz, question", id)
			neoRecords, err := transaction.Run(ctx,
				cypherScript,
				map[string]any{})

			if err != nil {
				return nil, err
			}

			records, err := neoRecords.Collect(ctx)
			if err != nil {
				return nil, err
			}

			quiz, err := toQuiz(records[0])
			if err != nil {
				return nil, err
			}

			var resultRecords []domain.QuestionForFetch

			for _, record := range records {

				questionForFetch, err := toQuestionForFetch(record)
				if err != nil {
					return nil, err
				}

				resultRecords = append(resultRecords, *questionForFetch)
			}

			return &domain.QuizWithQuestions{
				Quiz:      *quiz,
				Questions: resultRecords,
			}, nil
		})
	if err != nil {
		return domain.QuizWithQuestions{}, err
	}
	return *result, nil
}

func FetchRecordsForQuiz() {

}

func checkIfQuizExists(ctx context.Context, session neo4j.SessionWithContext, id string) bool {
	result, err := neo4j.ExecuteRead[domain.Quiz](ctx, session,
		func(transaction neo4j.ManagedTransaction) (domain.Quiz, error) {
			neoRecords, err := transaction.Run(ctx,
				fmt.Sprintf("MATCH (quiz:Quiz {id: '%s'}) RETURN quiz", id),
				map[string]any{})

			if err != nil {
				return domain.Quiz{}, err
			}

			records, err := neoRecords.Collect(ctx)
			if err != nil {
				return domain.Quiz{}, err
			}

			var resultRecord domain.Quiz

			for _, record := range records {
				quiz, err := toQuiz(record)
				if err != nil {
					return domain.Quiz{}, err
				}

				resultRecord = *quiz
			}

			return resultRecord, nil
		})

	if err != nil || reflect.DeepEqual(domain.Quiz{}, result) {
		return false
	}
	return true
}
