package http

import (
	"encoding/json"
	"fmt"
	"github.com/piotrselak/back/domain"
	"net/http"

	"github.com/piotrselak/back/db"
	"github.com/piotrselak/back/repository"
)

func FetchAllQuizes(w http.ResponseWriter, r *http.Request) {
	session := db.GetSessionFromContext(r)
	ctx := r.Context()
	quizes, err := repository.GetAllQuizes(ctx, session)
	if err == fmt.Errorf("could not find column") {
		jsonErr, _ := json.Marshal([]domain.Quiz{})
		w.Write(jsonErr)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	result, err := json.Marshal(quizes)
	if err != nil {
		http.Error(w, "Request failed due to error during parsing json", 500)
		return
	}
	w.Write(result)
}

func CreateNewQuiz(w http.ResponseWriter, r *http.Request) {
	session := db.GetSessionFromContext(r)
	ctx := r.Context()

	var q domain.QuizForPost
	err := json.NewDecoder(r.Body).Decode(&q)
	fmt.Println(r.Body, q, err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = repository.CreateQuiz(ctx, session, q)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}
