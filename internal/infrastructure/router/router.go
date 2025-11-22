package router

import (
	"net/http"

	"github.com/Cursus-platforms/cursus-server-go/internal/features/course"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Cursus-platforms/cursus-server-go/internal/features/auth"
)

type HandleDependencies struct {
	AuthHandler   *auth.Handler
	CourseHandler *course.Handler
}

func NewRouter(deps HandleDependencies) http.Handler {
	r := chi.NewRouter()

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Tối đa 5 phút
	})

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(corsOptions.Handler)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api/v1", func(r chi.Router) {
		//Auth routes
		r.Post("/auth/login", deps.AuthHandler.HandleLogin)
		r.Post("/auth/refresh", deps.AuthHandler.HandleRefresh)
		r.Post("/auth/register", deps.AuthHandler.HandleRequestVerification)

		//User routes
		// r.Group(func(r chi.Router) {
		// 	r.Use(customMiddleware.RequireRole(model.ROLE_ADMIN, model.ROLE_INSTRUCTOR, model.ROLE_STUDENT))

		// 	r.Get("/users/me", deps.UserHandler.GetMyProfile)
		// })

	})

	return r
}
