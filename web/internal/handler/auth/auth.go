package auth

import (
	"context"
	"net/http"
	"quiz_platform/internal/handler/repository"
	"quiz_platform/internal/misc/config"
	"quiz_platform/internal/utility"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

type registrationForm struct {
	Username string `json:"username" binding:"required,min=5,max=100"`
	Email    string `json:"email" binding:"required,min=4,max=200"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type loginForm struct {
	Email    string `json:"email" binding:"required,min=4,max=200"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func RegisterGetHandler(c *gin.Context) {
	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "register.html", utility.MergeMaps(*baseH, gin.H{
		"title": "Register"}))
}

func RegisterPostHandler(c *gin.Context) {
	ctx := context.Background()
	var form registrationForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usernameLength := utf8.RuneCountInString(form.Username)
	emailLength := utf8.RuneCountInString(form.Email)
	passwordLength := utf8.RuneCountInString(form.Password)

	if usernameLength < 5 || usernameLength > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be between 5 and 100 characters."})
		return
	}
	if emailLength < 4 || emailLength > 200 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email must be between 4 and 200 characters."})
		return
	}
	if passwordLength < 8 || passwordLength > 72 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be between 8 and 72 characters."})
		return
	}

	hash, err := utility.HashPassword(form.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = repository.UserRepositoryInstance.AddUser(
		ctx, form.Username, form.Email, hash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func LoginGetHandler(c *gin.Context) {
	baseHInterface, _ := c.Get("BaseH")
	baseH, _ := baseHInterface.(*gin.H)
	c.HTML(http.StatusOK, "login.html", utility.MergeMaps(*baseH, gin.H{
		"title": "Login"}))
}

func LoginPostHandler(c *gin.Context) {
	ctx := context.Background()
	var form loginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := repository.UserRepositoryInstance.GetUserByEmail(ctx, form.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = utility.ValidatePassword([]byte(user.PasswordHash), []byte(form.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utility.CreateToken(
		user.Id,
		time.Second*time.Duration(config.GlobalConfig.TokenInfo.ExpiresIn),
		config.GlobalConfig.TokenInfo.PrivateKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"token": token.Token})
}
