package domain

type QuizWithQuestions struct {
	Quiz      `json:"quiz"`
	Questions []QuestionForFetch `json:"questions"`
}
