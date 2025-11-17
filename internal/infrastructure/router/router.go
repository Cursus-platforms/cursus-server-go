package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Cursus-platforms/cursus-server-go/internal/features/auth"
)

type HandleDependencies struct {
	AuthHandler *auth.Handler
}

func NewRouter(deps HandleDependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api/v1", func(r chi.Router) {
		//Auth routes
		r.Post("/auth/register", deps.AuthHandler.HandleRegister)
		// r.Post("/auth/login", deps.AuthHandler.HandleLogin)

		//User routes
		// r.Group(func(r chi.Router) {
		// 	r.Use(customMiddleware.RequireRole(model.ROLE_ADMIN, model.ROLE_INSTRUCTOR, model.ROLE_STUDENT))

		// 	r.Get("/users/me", deps.UserHandler.GetMyProfile)
		// })

	})

	return r
}
