package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/piotrselak/back/domain"
	"strconv"
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

	rating, err := neo4j.GetProperty[string](itemNode, "rating")
	if err != nil {
		return nil, err
	}

	editHash, err := neo4j.GetProperty[string](itemNode, "editHash")
	if err != nil {
		return nil, err
	}

	ratingFloat, err := strconv.ParseFloat(rating, 64)
	if err != nil {
		return nil, err
	}
	return &domain.Quiz{Id: id, Name: name, Rating: ratingFloat, EditHash: editHash}, nil
}

func toQuestionForFetch(record *neo4j.Record) (*domain.QuestionForFetch, error) {
	rawItemNode, found := record.Get("question")
	if !found {
		return nil, fmt.Errorf("could not find column")
	}
	itemNode := rawItemNode.(neo4j.Node)

	index, err := neo4j.GetProperty[string](itemNode, "index")
	if err != nil {
		fmt.Println("Expected fuckup")
		return nil, err
	}

	questionText, err := neo4j.GetProperty[string](itemNode, "questionText")
	if err != nil {
		return nil, err
	}

	_, err = neo4j.GetProperty[string](itemNode, "answers") //idk if any is all right
	if err != nil {
		return nil, err
	}
	//answersParsed, err := formatToJSON(answers)
	if err != nil {
		return nil, err
	}
	//fmt.Println(answersParsed)
	var answersFinal []string
	//err = json.Unmarshal([]byte(answersParsed), &answersFinal)
	//if err != nil {
	//	fmt.Println(answersParsed)
	//	fmt.Println(err)
	//}

	qType, err := neo4j.GetProperty[string](itemNode, "type")
	if err != nil {
		return nil, err
	}

	intIndex, err := strconv.ParseInt(index, 10, 64)
	if err != nil {
		return nil, err
	}

	return &domain.QuestionForFetch{
		Index:        intIndex,
		QuestionText: questionText,
		Answers:      answersFinal,
		Type:         qType,
	}, nil
}
