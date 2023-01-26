package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/piotrselak/back/domain"
)

func toQuizForFetch(record *neo4j.Record) (*domain.QuizForFetch, error) {
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

	modifiersRaw, err := neo4j.GetProperty[[]interface{}](itemNode, "modifiers")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var modifiers []string
	for _, modifier := range modifiersRaw {
		modifiers = append(modifiers, modifier.(string))
	}
	return &domain.QuizForFetch{Id: id, Name: name, Rating: rating, Modifiers: modifiers}, nil
}

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

	modifiersRaw, err := neo4j.GetProperty[[]interface{}](itemNode, "modifiers")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var modifiers []string
	for _, modifier := range modifiersRaw {
		modifiers = append(modifiers, modifier.(string))
	}

	return &domain.Quiz{Id: id, Name: name, Rating: rating, EditHash: editHash, Modifiers: modifiers}, nil
}

func toQuestionForFetch(record *neo4j.Record) (*domain.QuestionForFetch, error) {
	rawItemNode, found := record.Get("question")
	if !found {
		return nil, fmt.Errorf("could not find column")
	}
	itemNode := rawItemNode.(neo4j.Node)

	index, err := neo4j.GetProperty[int64](itemNode, "index")
	if err != nil {
		return nil, err
	}

	questionText, err := neo4j.GetProperty[string](itemNode, "questionText")
	if err != nil {
		return nil, err
	}

	answers, err := neo4j.GetProperty[[]any](itemNode, "answers") //idk if any is all right
	if err != nil {
		return nil, err
	}

	var answersFinal map[int]string
	for i, ans := range answers {
		answersFinal[i] = ans.(string)
	}

	qType, err := neo4j.GetProperty[string](itemNode, "type")
	if err != nil {
		return nil, err
	}

	return &domain.QuestionForFetch{
		Index:        index,
		QuestionText: questionText,
		Answers:      answersFinal,
		Type:         qType,
	}, nil
}

func toQuestion(record *neo4j.Record) (*domain.Question, error) {
	rawItemNode, found := record.Get("question")
	if !found {
		return nil, fmt.Errorf("could not find column")
	}
	itemNode := rawItemNode.(neo4j.Node)

	index, err := neo4j.GetProperty[int64](itemNode, "index")
	if err != nil {
		return nil, err
	}

	questionText, err := neo4j.GetProperty[string](itemNode, "questionText")
	if err != nil {
		return nil, err
	}

	answers, err := neo4j.GetProperty[[]any](itemNode, "answers")
	if err != nil {
		return nil, err
	}

	var answersFinal map[int]string
	for i, ans := range answers {
		answersFinal[i] = ans.(string)
	}

	validAnswers, err := neo4j.GetProperty[[]any](itemNode, "validAnswers")
	if err != nil {
		return nil, err
	}

	var validAnswersFinal map[int]string
	for i, ans := range validAnswers {
		validAnswersFinal[i] = ans.(string)
	}

	qType, err := neo4j.GetProperty[string](itemNode, "type")
	if err != nil {
		return nil, err
	}

	return &domain.Question{
		Index:        index,
		QuestionText: questionText,
		Answers:      answersFinal,
		ValidAnswers: validAnswersFinal,
		Type:         qType,
	}, nil
}

func toUser(record *neo4j.Record) (*domain.User, error) {
	rawItemNode, found := record.Get("p")
	if !found {
		return nil, fmt.Errorf("could not find column")
	}
	itemNode := rawItemNode.(neo4j.Node)

	name, err := neo4j.GetProperty[string](itemNode, "name")
	if err != nil {
		return nil, err
	}

	return &domain.User{
		Name: name,
	}, nil
}

func toPlayed(record *neo4j.Record) (*domain.Played, error) {
	rawItemNode, found := record.Get("r")
	if !found {
		return nil, fmt.Errorf("could not find column")
	}
	itemNode := rawItemNode.(neo4j.Relationship)

	score, err := neo4j.GetProperty[int64](itemNode, "score")
	if err != nil {
		return nil, err
	}

	return &domain.Played{
		Score: score,
	}, nil
}
