package rest

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tarkovskynik/Golang-ninja-project/internal/domain"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tarkovskynik/Golang-ninja-project/docs"
)

type Users interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error)
	ParseToken(token string) (int, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
	GetRefreshTokenTTL() time.Duration
}

type FilesServece interface {
	Upload(ctx context.Context, input domain.File) (string, error)
	GetFiles(ctx context.Context, id int) ([]domain.File, error)
	StoreFileInfo(ctx context.Context, input domain.File) error
}

type Handler struct {
	usersService      Users
	filesService      FilesServece
	maxUploadFileSize int64
	fileTypes         map[string]interface{}
}

func NewHandler(users Users, files FilesServece, maxUploadFileSize int64, fileTypes map[string]interface{}) *Handler {
	return &Handler{
		usersService:      users,
		filesService:      files,
		maxUploadFileSize: maxUploadFileSize,
		fileTypes:         fileTypes,
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

	filesApi := router.Group("/s3")
	filesApi.Use(h.authMiddleware())
	{
		filesApi.POST("/upload", h.fileUploadS3)
		filesApi.GET("/files", h.getFilesS3)
	}

	return router
}
