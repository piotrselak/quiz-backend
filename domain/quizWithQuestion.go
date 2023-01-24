package domain

type QuizWithQuestions struct {
	QuizForFetch `json:"quiz"`
	Questions    []QuestionForFetch `json:"questions"`
}

type QuizWithQuestionsAndAnswers struct {
	QuizForFetch `json:"quiz"`
	Questions    []Question `json:"questions"`
}
