package router

import (
	"client/internal/auth"
	"client/internal/handler"
	"os"

	"github.com/labstack/echo/v4"
)

type Router struct {
	*handler.Handler
}

func NewRouter(h *handler.Handler) *Router {
	return &Router{
		Handler: h,
	}
}

func (r *Router) Register(e *echo.Echo) {
	jwtMiddleware := auth.JWT()

	apiGroup := e.Group("/api")

	// register auth routes
	cake := apiGroup.Group("/cakes", jwtMiddleware)
	cake.GET("/:id", r.Handler.GetByID)
	cake.POST("", r.Handler.Create)
	cake.GET("/search", r.Handler.Search)
	cake.PUT("/:id", r.Handler.UpdateByID)

	// register auth routes
	auth := apiGroup.Group("/auth")
	auth.POST("/register", r.Handler.RegisterUser)
	auth.POST("/login", r.Handler.LoginUser)

	// register user routes
	user := apiGroup.Group("/users", jwtMiddleware)
	user.GET("/current", r.Handler.CurrentUser)

	// Serve static files (profile pictures) from the 'picture' directory.
	e.Static("/picture", os.Getenv("PATH_TO_UPLOAD"))
}
