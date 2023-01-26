package domain

type UserAnswers struct {
	Answers        []QuestionForFetch `json:"answers"`
	Name           string             `json:"name"`
	NegativePoints bool               `json:"negativePoints"`
}
