package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/piotrselak/back/internal/repository"
	"github.com/piotrselak/back/pkg/db"
	"github.com/piotrselak/back/pkg/domain"
)

func FetchAllQuizes(w http.ResponseWriter, r *http.Request) {
	session := db.GetSessionFromContext(r)
	ctx := r.Context()
	quizes, err := repository.GetAllQuizes(ctx, session)
	if err == fmt.Errorf("could not find column") {
		json, _ := json.Marshal([]domain.Quiz{})
		w.Write(json)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	result, err := json.Marshal(quizes)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Write(result)
}
