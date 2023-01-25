package id

import (
	"encoding/json"
	"fmt"
	"github.com/piotrselak/back/domain"
	"github.com/piotrselak/back/verification"
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
		fmt.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
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

func FetchSpecificQuizWithAnswers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	session := db.GetSessionFromContext(r)
	ctx := r.Context()
	id := ctx.Value("quizID").(string)
	hash := r.Header.Get("editHash")

	validHash, err := repository.FetchQuizHash(ctx, session, id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprint(err), 404)
		return
	}

	if hash != validHash {
		http.Error(w, http.StatusText(401), http.StatusUnauthorized)
		return
	}

	q, err := repository.FetchSpecificQuizWithAnswers(ctx, session, id)
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

// saves record as well
// 1 good answer = 10 points
func VerifyAnswers(w http.ResponseWriter, r *http.Request) {
	session := db.GetSessionFromContext(r)
	ctx := r.Context()
	id := ctx.Value("quizID").(string)

	var answers domain.UserAnswers
	err := json.NewDecoder(r.Body).Decode(&answers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	quizWithAnswers, err := repository.FetchSpecificQuizWithAnswers(ctx, session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	questions := quizWithAnswers.Questions

	validCounter := verification.ParseAnswers(answers.Answers, questions)
	fmt.Println(validCounter)

	score := int64(validCounter * 10)

	playerRecord := domain.RecordUnit{
		User:   domain.User{Name: answers.Name},
		Played: domain.Played{Score: score},
	}

	err = repository.SaveRecord(ctx, session, id, playerRecord)
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	jsonRecord, err := json.Marshal(playerRecord)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Write(jsonRecord)
}

func FetchRecordsForQuiz(w http.ResponseWriter, r *http.Request) {
	session := db.GetSessionFromContext(r)
	ctx := r.Context()
	id := ctx.Value("quizID").(string)

	records, err := repository.FetchRecordsForQuiz(ctx, session, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	recordsParsed := *records

	jsonRecords, err := json.Marshal(recordsParsed)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(jsonRecords)
}

func LikeQuiz() {}
