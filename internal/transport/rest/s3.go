package rest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tarkovskynik/Golang-ninja-project/pkg/logger"

	"github.com/tarkovskynik/Golang-ninja-project/internal/domain"
)

func newErrorResponse(ctx *gin.Context, statusCode int, tag string, err error) {
	logger.LogError(tag, err)
	ctx.AbortWithStatusJSON(statusCode, domain.Response{Status: "error", Error: err.Error()})
}

// @Summary Upload picture to S3
// @Security JWT
// @Tags files
// @Description upload picture
// @ModuleID fileUploadImage
// @Accept mpfd
// @Produce json
// @Param file formData file true "file"
// @Success 200 	{object} domain.Response
// @Failure 400,404 {object} domain.Response
// @Failure 500 	{object} domain.Response
// @Failure default {object} domain.Response
// @Router /s3/upload [post]
func (h *Handler) fileUploadS3(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "fileUploadS3", err)
		return
	}

	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, h.maxUploadFileSize) // 10 megabytes

	iFile, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "fileUploadS3", err)
		return
	}
	defer iFile.Close()

	buffer := make([]byte, fileHeader.Size)

	if _, err := iFile.Read(buffer); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "fileUploadS3", err)
		return
	}

	contentType := http.DetectContentType(buffer)

	// Validate File Type
	if _, ex := h.fileTypes[contentType]; !ex {
		newErrorResponse(ctx, http.StatusBadRequest, "fileUploadS3", domain.ErrorFileType)
		return
	}

	tempFilename := fmt.Sprintf("temp.file.%d-%s", userId, fileHeader.Filename)

	f, err := os.OpenFile(tempFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "fileUploadS3", domain.ErrorCreateTempFile)
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, bytes.NewReader(buffer)); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "fileUploadS3", domain.ErrorWriteTempFile)
		return
	}

	file := domain.File{
		UserID:      int64(userId),
		Name:        fileHeader.Filename,
		Size:        fileHeader.Size,
		Type:        domain.Image,
		ContentType: contentType,
	}

	file.URL, err = h.filesService.Upload(ctx.Request.Context(), file)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "fileUploadS3", err)
		return
	}

	// размер, дата загрузки, айди пользователя, ссылка на внешнее хранилище
	file.UploadAt = time.Now()
	err = h.filesService.StoreFileInfo(ctx.Request.Context(), file)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "fileUploadS3", err)
		return
	}

	ctx.JSON(http.StatusOK, &domain.Response{Status: "ok", URL: file.URL})
}

// @Summary Get list of files from S3
// @Security JWT
// @Tags files
// @Description Get list of files from S3
// @ModuleID    fileUploadImage
// @Produce     json
// @Success     200       {object} domain.Response
// @Failure     400,500   {object} domain.Response
// @Failure     default   {object} domain.Response
// @Router      /s3/files [get]
func (h *Handler) getFilesS3(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "getFilesS3", err)
		return
	}

	files, err := h.filesService.GetFiles(ctx, userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "getFilesS3", err)
		return
	}

	ctx.JSON(http.StatusOK, &domain.Response{Status: "ok", Files: files})
}
