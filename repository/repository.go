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

func GetAllQuizes(ctx context.Context, session neo4j.SessionWithContext) ([]*domain.QuizForFetch, error) {
	result, err := neo4j.ExecuteRead[[]*domain.QuizForFetch](ctx, session,
		func(transaction neo4j.ManagedTransaction) ([]*domain.QuizForFetch, error) {
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

			var resultRecords []*domain.QuizForFetch

			for _, record := range records {
				quiz, err := toQuizForFetch(record)
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
	quiz domain.QuizForPost) (string, error) {
	result, err := neo4j.ExecuteWrite[domain.Quiz](ctx, session,
		func(transaction neo4j.ManagedTransaction) (domain.Quiz, error) {
			newQuiz := quiz.ToQuiz(uuid.New().String(), 0)
			_, err := transaction.Run(ctx,
				fmt.Sprintf("CREATE %s", newQuiz.ToCypher("q")),
				map[string]any{})

			if err != nil {
				return domain.Quiz{}, err
			}

			return newQuiz, nil
		})
	if err != nil {
		return "", err
	}
	return result.Id, nil
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

	cypherScript = fmt.Sprintf("MATCH (q2:Quiz {id: '%s'}) DELETE q2", id)
	_, _ = session.Run(ctx, cypherScript, map[string]any{})

	return nil
}

func FetchSpecificQuizWithAnswers(ctx context.Context, session neo4j.SessionWithContext, id string) (domain.QuizWithQuestionsAndAnswers, error) {
	res := checkIfQuizExists(ctx, session, id)
	if !res {
		return domain.QuizWithQuestionsAndAnswers{}, errors.New("not found")
	}

	result, err := neo4j.ExecuteRead[*domain.QuizWithQuestionsAndAnswers](ctx, session,
		func(transaction neo4j.ManagedTransaction) (*domain.QuizWithQuestionsAndAnswers, error) {
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

			if len(records) == 0 {
				return nil, errors.New("quiz has no questions")
			}

			quiz, err := toQuizForFetch(records[0])
			if err != nil {
				return nil, err
			}

			var resultRecords []domain.Question

			for _, record := range records {

				question, err := toQuestion(record)
				if err != nil {
					return nil, err
				}

				resultRecords = append(resultRecords, *question)
			}

			return &domain.QuizWithQuestionsAndAnswers{
				QuizForFetch: *quiz,
				Questions:    resultRecords,
			}, nil
		})
	if err != nil {
		return domain.QuizWithQuestionsAndAnswers{}, err
	}
	return *result, nil
}

func FetchSpecificQuiz(ctx context.Context, session neo4j.SessionWithContext, id string) (domain.QuizWithQuestions, error) {
	res := checkIfQuizExists(ctx, session, id)
	if !res {
		return domain.QuizWithQuestions{}, errors.New("not found")
	}

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

			if len(records) == 0 {
				return nil, errors.New("quiz has no questions")
			}

			quiz, err := toQuizForFetch(records[0])
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
				QuizForFetch: *quiz,
				Questions:    resultRecords,
			}, nil
		})
	if err != nil {
		return domain.QuizWithQuestions{}, err
	}
	return *result, nil
}

func FetchQuizHash(ctx context.Context, session neo4j.SessionWithContext, id string) (string, error) {
	res := checkIfQuizExists(ctx, session, id)
	if !res {
		return "", errors.New("not found")
	}

	result, err := neo4j.ExecuteRead[string](ctx, session,
		func(transaction neo4j.ManagedTransaction) (string, error) {
			cypherScript := fmt.Sprintf("MATCH (quiz:Quiz {id: '%s'}) RETURN quiz.editHash", id)
			neoRecords, err := transaction.Run(ctx,
				cypherScript,
				map[string]any{})

			if err != nil {
				return "", err
			}

			record, err := neoRecords.Single(ctx)
			if err != nil {
				return "", err
			}

			hash, _, err := neo4j.GetRecordValue[string](record, "quiz.editHash")
			if err != nil {
				return "", err
			}
			return hash, nil
		})
	if err != nil {
		return "", err
	}
	return result, nil
}

func UpdateQuiz() {

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
