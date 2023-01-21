package id

import (
	"encoding/json"
	"fmt"
	"github.com/piotrselak/back/domain"
	"net/http"

	"github.com/piotrselak/back/db"
	"github.com/piotrselak/back/repository"
)

func AddQuestions(w http.ResponseWriter, r *http.Request) {
	session := db.GetSessionFromContext(r)
	ctx := r.Context()
	id := ctx.Value("quizID").(string)

	var q domain.QuestionForPost
	err := json.NewDecoder(r.Body).Decode(&q)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validHash, err := repository.FetchQuizHash(ctx, session, id)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	if q.EditHash != validHash {
		http.Error(w, http.StatusText(401), http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = repository.AddQuestions(ctx, session, id, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func RemoveQuiz(w http.ResponseWriter, r *http.Request) {
	session := db.GetSessionFromContext(r)
	ctx := r.Context()
	id := ctx.Value("quizID").(string)

	var hash domain.Hash
	err := json.NewDecoder(r.Body).Decode(&hash)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validHash, err := repository.FetchQuizHash(ctx, session, id)

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	if hash.EditHash != validHash {
		http.Error(w, http.StatusText(401), http.StatusUnauthorized)
		return
	}

	err = repository.RemoveQuiz(ctx, session, id)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func FetchSpecificQuiz(w http.ResponseWriter, r *http.Request) {
	session := db.GetSessionFromContext(r)
	ctx := r.Context()
	id := ctx.Value("quizID").(string)

	q, err := repository.FetchSpecificQuiz(ctx, session, id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprint(err), 404)
		return
	}
	marshalled, err := json.Marshal(q)
	if err != nil {
		panic(err)
	}
	w.Write(marshalled)
}

func ModifyQuiz() {

}

func LikeQuiz() {}

// saves record as well
func VerifyAnswers() {

}

func FetchScoreStatisticsForQuiz() {

}

//maybe some apoc
