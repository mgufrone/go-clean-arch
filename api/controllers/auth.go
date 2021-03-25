package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"recipes/api/constants"
	"recipes/api/http"
	"recipes/api/requests"
	"recipes/app/common"
	"recipes/app/config"
	"recipes/domains/shared"
	"recipes/domains/user"
	"strconv"
	"strings"
)


type AuthController struct {
	handler user.IUserUseCase
	cfg *config.Config
}

func (a *AuthController) Logout(ctx *gin.Context) {
	panic("implement me")
}

func (a *AuthController) IsLoggedIn(ctx *gin.Context) {
	token := strings.Replace(ctx.GetHeader("authorization"), "Bearer ", "", 1)
	if token == "" {
		http.Unauthorized(ctx, nil)
		return
	}
	tkn, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.cfg.JwtKey, nil
	})
	if err != nil {
		http.Unauthorized(ctx, err)
	}
	claims := tkn.Claims.(*jwt.StandardClaims)
	aud, _ := strconv.Atoi(claims.Audience)
	usr := &user.User{
		Model: shared.Model{ID: uint(aud)},
	}
	sess := &user.Session{
		User: usr,
		SessionID: claims.Id,
	}
	if err = a.handler.CheckSession(ctx, usr, sess); err != nil {
		http.ServerError(ctx, err)
		return
	}
	ctx.Set(common.UserKey, aud)
	ctx.Set(common.UserSessionKey, claims.Id)
	ctx.Next()
}

func (a *AuthController) generateToken(sess *user.Session) (token string, refresh string) {
	usrID := strconv.Itoa(int(sess.User.ID.(uint64)))
	claims := jwt.StandardClaims{
		Audience:  usrID,
		ExpiresAt: sess.ExpiresAt.Unix(),
		Id:        sess.SessionID,
		IssuedAt:  sess.UpdatedAt.Unix(),
		Issuer:    a.cfg.AppName,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	claims.ExpiresAt = sess.RefreshExpiresAt.Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err := jwtToken.SignedString(a.cfg.JwtKey)
	if err != nil {
		return "", ""
	}
	jwtRefresh, _ := refreshToken.SignedString(a.cfg.JwtKey)
	return jwtString, jwtRefresh
}
func (a *AuthController) Login(ctx *gin.Context) {
	req, ok := ctx.Get(constants.RequestKey)
	if !ok {
		http.BadRequest(ctx, nil)
		return
	}
	r := req.(requests.Login)
	usr, err := a.handler.FindByToken(ctx, r.Token)
	if err != nil {
		http.ServerError(ctx, err)
		return
	}
	ses, err := a.handler.Login(ctx, usr)
	if err != nil {
		http.ServerError(ctx, err)
		return
	}
	token, refresh := a.generateToken(ses)
	http.Ok(ctx, gin.H{"token": token, "refresh_token": refresh})
}
func NewAuthController(handler user.IUserUseCase, cfg *config.Config) *AuthController {
	return &AuthController{handler, cfg}
}
