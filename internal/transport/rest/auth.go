package rest

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tarkovskynik/Golang-ninja-project/pkg/logger"

	"github.com/tarkovskynik/Golang-ninja-project/internal/domain"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce json
// @Param   input body     domain.SignUpInput          true "account info"
// @Success 201   {object} domain.Response
// @Failure 400   {object} domain.Response
// @Failure 500   {object} domain.Response
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input domain.SignUpInput
	if err := c.BindJSON(&input); err != nil {
		logger.LogError("signUp", err)
		c.JSON(http.StatusBadRequest, domain.Response{Error: domain.ErrUserInputParam.Error()})
		return
	}

	err := h.usersService.SignUp(c.Request.Context(), input)
	if err != nil {
		logger.LogError("signUp", err)
		c.JSON(http.StatusInternalServerError, domain.Response{Error: domain.ErrCantCreateUser.Error()})
		return
	}

	c.JSON(http.StatusCreated, domain.Response{Status: "ok"})
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce json
// @Param   input body     domain.SignInInput          true "account credentials"
// @Success 201   {object} domain.Response
// @Failure 400   {object} domain.Response
// @Failure 500   {object} domain.Response
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input domain.SignInInput
	if err := c.BindJSON(&input); err != nil {
		logger.LogError("signIn", err)
		c.JSON(http.StatusBadRequest, domain.Response{Error: domain.ErrUserInputParam.Error()})
		return
	}

	accessToken, refreshToken, err := h.usersService.SignIn(c.Request.Context(), input)
	if err != nil {
		logger.LogError("signIn", err)
		if errors.Is(err, domain.ErrUserCredNotFound) {
			c.JSON(http.StatusBadRequest, domain.Response{Error: domain.ErrSearchUserError.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, domain.Response{Error: domain.ErrSearchUserError.Error()})
		return
	}
	refreshTokenTTL := h.usersService.GetRefreshTokenTTL().Seconds()
	c.SetCookie("refresh-token", refreshToken, int(refreshTokenTTL), "/", "localhost", false, true)
	c.JSON(http.StatusOK, domain.Response{Status: "ok", Token: "Bearer " + accessToken})
}

// @Summary     Refresh
// @Description Refresh tokens
// @Tags        auth
// @Produce     json
// @Success     200     {object} domain.Response
// @Failure     400,500 {object} domain.Response
// @Router      /auth/refresh [get]
func (h *Handler) refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		logger.LogError("refresh", err)
		c.JSON(http.StatusBadRequest, domain.Response{Error: domain.ErrRefreshTokenParse.Error()})
		return
	}

	accessToken, refreshToken, err := h.usersService.RefreshTokens(c.Request.Context(), cookie)
	if err != nil {
		logger.LogError("refresh", err)
		c.JSON(http.StatusInternalServerError, domain.Response{Error: domain.ErrRefreshToken.Error()})
		return
	}
	refreshTokenTTL := h.usersService.GetRefreshTokenTTL().Seconds()
	c.SetCookie("refresh-token", refreshToken, int(refreshTokenTTL), "/", "localhost", false, true)
	c.JSON(http.StatusOK, domain.Response{Status: "ok", Token: accessToken})
}
