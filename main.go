package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	db "github.com/piotrselak/back/db"
	"github.com/piotrselak/back/web/quiz"
	quizId "github.com/piotrselak/back/web/quiz/id"
)

func main() {
	driver := db.InitNeo4j()
	ctx := context.Background()
	defer func(driver neo4j.DriverWithContext, ctx context.Context) {
		err := driver.Close(ctx)
		if err != nil {
			panic(err)
		}
	}(driver, ctx)

	r := chi.NewRouter()
	corsConfig := cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "editHash"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(corsConfig))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("api"))
	})

	r.Mount("/quiz", quizRouter(driver))
	r.Mount("/quiz/{id}", specificQuizRouter(driver))
	err := http.ListenAndServe(":3333", r)
	if err != nil {
		panic(err)
	}
}

func quizRouter(driver neo4j.DriverWithContext) http.Handler {
	r := chi.NewRouter()
	r.Use(openSession(driver))
	r.Get("/", quiz.FetchAllQuizes)
	r.Post("/", quiz.CreateNewQuiz)
	return r
}

func specificQuizRouter(driver neo4j.DriverWithContext) http.Handler {
	r := chi.NewRouter()
	r.Use(openSession(driver))
	r.Use(QuizIDCtx)
	r.Put("/", quizId.AddQuestions)
	r.Delete("/", quizId.RemoveQuiz)
	r.Get("/", quizId.FetchSpecificQuiz)
	r.Post("/", quizId.VerifyAnswers)
	r.Get("/answers", quizId.FetchSpecificQuizWithAnswers)
	r.Get("/score", quizId.FetchRecordsForQuiz)
	return r
}

func QuizIDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		quizID := chi.URLParam(r, "id")
		ctx := context.WithValue(r.Context(), "quizID", quizID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func openSession(driver neo4j.DriverWithContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := driver.NewSession(r.Context(), neo4j.SessionConfig{})
			defer func() {
				err := session.Close(r.Context())
				if err != nil {
					return
				}
			}()

			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
