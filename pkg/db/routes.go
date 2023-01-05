package db

import (
	"context"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func OpenSession(driver neo4j.DriverWithContext) func(http.Handler) http.Handler {
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
