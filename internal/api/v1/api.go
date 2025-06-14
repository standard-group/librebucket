package v1

import (
    "net/http"

    "github.com/go-chi/chi/v5"
)

// RegisterAPIRoutes registers all API v1 routes to the router.
func RegisterAPIRoutes(r chi.Router) {
    // User authentication
    r.Post("/register", RegisterUserHandler)
    r.Post("/login", LoginUserHandler)
    r.With(AuthMiddleware).Get("/user", GetCurrentUserHandler)

    // Repository management
    r.With(AuthMiddleware).Get("/repos", ListRepositoriesHandler)
    r.With(AuthMiddleware).Post("/repos", CreateRepositoryHandler)
    r.Get("/repos/{owner}/{repo}", GetRepositoryHandler)
    r.With(AuthMiddleware).Delete("/repos/{owner}/{repo}", DeleteRepositoryHandler)
    r.With(AuthMiddleware).Patch("/repos/{owner}/{repo}", UpdateRepositoryHandler)

    // Git objects
    r.Get("/repos/{owner}/{repo}/commits", ListCommitsHandler)
    r.Get("/repos/{owner}/{repo}/commits/{sha}", GetCommitHandler)
    r.Get("/repos/{owner}/{repo}/blob/{ref}/{path:.*}", GetBlobHandler)
    r.Get("/repos/{owner}/{repo}/tree/{ref}/{path:.*}", GetTreeHandler)

    // User management (admin)
    r.With(AuthMiddleware, AdminMiddleware).Get("/users", ListUsersHandler)
    r.With(AuthMiddleware, AdminMiddleware).Get("/users/{username}", GetUserHandler)
    r.With(AuthMiddleware, AdminMiddleware).Delete("/users/{username}", DeleteUserHandler)
    r.With(AuthMiddleware, AdminMiddleware).Patch("/users/{username}", UpdateUserHandler)

    // Health check
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    })
}