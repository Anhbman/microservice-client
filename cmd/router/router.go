package router

import (
	"os"

	"github.com/labstack/echo/v4"
	"client/cmd/handler"
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
	// jwtMiddleware := jwt.JWT(utils.JWTSecret)
	// create api group
	// apiGroup := e.Group("/api")
	// auth := apiGroup.Group("/auth")
	// auth.POST("/register", h.SignUp)
	// auth.POST("/login", h.Login)

	// create user group
	// user := apiGroup.Group("/users", jwtMiddleware)
	// user.GET("/current", h.CurrentUser)

	// // create cake group
	// cake := apiGroup.Group("/cakes", jwtMiddleware)
	// // cake routes
	// cake.POST("", h.Create)
	// // cake.GET("", h.Get)
	// cake.GET("/search", h.Search)
	// cake.GET("/:id", h.GetByID)
	// cake.DELETE("/:id", h.DeleteByID)
	// cake.PUT("/:id", h.UpdateByID)


	apiGroup := e.Group("/api")

	// register auth routes
	cake := apiGroup.Group("/cakes")
	cake.GET("/:id", r.Handler.GetByID)
	cake.POST("", r.Handler.Create)
	cake.GET("/search", r.Handler.Search)
	cake.PUT("/:id", r.Handler.UpdateByID)


	// Serve static files (profile pictures) from the 'picture' directory.
	e.Static("/picture", os.Getenv("PATH_TO_UPLOAD"))
}
