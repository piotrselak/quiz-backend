package web

import (
	"encoding/json"
	"fmt"
	"github.com/piotrselak/back/db"
	"github.com/piotrselak/back/domain"
	"github.com/piotrselak/back/repository"
	"net/http"
)

// FetchAllQuizes All errors are 500 - everything should work even when database is empty
func FetchAllQuizes(w http.ResponseWriter, r *http.Request) {
	session := db.GetSessionFromContext(r)
	ctx := r.Context()
	quizes, err := repository.GetAllQuizes(ctx, session)
	if err == fmt.Errorf("could not find column") {
		jsonErr, _ := json.Marshal([]domain.Quiz{})
		w.WriteHeader(500)
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

func FetchQuizesCount(w http.ResponseWriter, r *http.Request) {
	session := db.GetSessionFromContext(r)
	ctx := r.Context()
	count, err := repository.CountQuizes(ctx, session)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	result, err := json.Marshal(count)
	if err != nil {
		http.Error(w, "Request failed due to error during parsing json", 500)
		return
	}
	w.Write(result)
}

// CreateNewQuiz Error 500 happens only if something is parsed badly by neo4j - server error
func CreateNewQuiz(w http.ResponseWriter, r *http.Request) {
	session := db.GetSessionFromContext(r)
	ctx := r.Context()

	var q domain.QuizForPost
	err := json.NewDecoder(r.Body).Decode(&q)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := repository.CreateQuiz(ctx, session, q)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func FilterByLikes() {

}

func FilterByMostPlays() {

}

func FetchNumberOfQuizes() {

}
