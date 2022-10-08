package domain

import (
	"errors"
)

var ErrEmptyAuthHeader = errors.New("empty auth header")
var ErrInvalidAuthHeader = errors.New("invalid auth header")

var ErrEmptyToken = errors.New("token is empty")
var ErrAccessTokenParse = errors.New("access parse token error")
var ErrAccessTokenExpired = errors.New("access token expired")
var ErrRefreshTokenExpired = errors.New("refresh token expired")
var ErrRefreshTokenParse = errors.New("parse refresh token error")
var ErrRefreshToken = errors.New("refresh tokens error")

var ErrUserCredNotFound = errors.New("user with such credentials not found")
var ErrUserInputParam = errors.New("invalid user input param")
var ErrCantCreateUser = errors.New("can't create user")
var ErrSearchUserError = errors.New("search user error")

var ErrorFileType = errors.New("file type is not supported")
var ErrorCreateTempFile = errors.New("failed to create temp file")
var ErrorWriteTempFile = errors.New("failed to write chunk to temp file")
