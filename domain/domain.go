package domain

type QuizWithQuestions struct {
	Quiz      `json:"quiz"`
	Questions []Question `json:"questions"`
}
