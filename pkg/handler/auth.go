package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kazakovichna/todoListPrjct"
	"net/http"
	"time"
)

func (h *Handler) SignUp(c *gin.Context) {
	var input todoListPrjct.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	//UserEmail string `json:"user_email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) SignIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	tokens, err := h.services.Authorization.CreateSession(input.Username, input.Password)
	//token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:       "Access",
		Value:      tokens.AccessToken,
		Path:       "/",
		Domain:     "http://localhost:8080/",
		Expires:    time.Now().Add(15 * time.Minute),
	})

	c.JSON(http.StatusOK, map[string]interface{}{
		"AccessToken": tokens.AccessToken,
		"RefreshToken": tokens.RefreshToken,
	})
}

type refreshInput struct {
	RefreshToken string `json:"RefreshToken" binding:"required"`
}

func (h *Handler) refreshToken(c *gin.Context) {
	var input refreshInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.services.Authorization.RefreshTokenServices(input.RefreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"AccessToken": res.AccessToken,
		"RefreshToken": res.RefreshToken,
	})
}