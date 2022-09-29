package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"time"

	"github.com/tarkovskynik/Golang-ninja-project/internal/domain"

	_ "github.com/tarkovskynik/Golang-ninja-project/docs"
)

type Users interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error)
	ParseToken(token string) (int, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
	GetRefreshTokenTTL() time.Duration
}

type Handler struct {
	usersService Users
}

func NewHandler(users Users) *Handler {
	return &Handler{
		usersService: users,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	usersApi := router.Group("/auth")
	{
		usersApi.POST("/sign-up", h.signUp)
		usersApi.POST("/sign-in", h.signIn)
		usersApi.GET("/refresh", h.refresh)
	}

	//filesApi := router.Group("")
	//filesApi.Use(h.authMiddleware())
	//{
	//	filesApi.POST("/upload")
	//	filesApi.GET("/files")
	//}
	return router
}
