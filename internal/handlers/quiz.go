package handlers

import (
	"fmt"
	"net/http"

	"github.com/piotrselak/back/internal/repository"
	"github.com/piotrselak/back/pkg/db"
)

func FetchAllQuizes(w http.ResponseWriter, r *http.Request) {
	session := db.GetSessionFromContext(r)
	ctx := r.Context()
	quizes, error := repository.GetAllQuizes(ctx, session)
	if err != nil {
		http.Error
	}

	fmt.Print(result)
}
