package api

import (
	"avito-backend-bootcamp/auth"
	"avito-backend-bootcamp/database"
	"avito-backend-bootcamp/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type NoAuthAPI struct {
}

func (api *NoAuthAPI) DummyLoginGet(c *gin.Context) {
	var dummyLoginRequest models.DummyLoginRequest

	if err := c.ShouldBindJSON(&dummyLoginRequest); err != nil {
		response := models.DummyLoginGet500Response{
			Message:   "Invalid data",
			RequestId: uuid.New().String(),
			Code:      http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if dummyLoginRequest.UserType != models.CLIENT && dummyLoginRequest.UserType != models.MODERATOR {
		response := models.DummyLoginGet500Response{
			Message:   "Invalid user type",
			RequestId: uuid.New().String(),
			Code:      http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	jwtToken, err := auth.GenerateJwtToken("dummylogin@example.com", string(dummyLoginRequest.UserType))
	if err != nil {
		response := models.DummyLoginGet500Response{
			Message:   "Failed to generate token",
			RequestId: uuid.New().String(),
			Code:      http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := models.DummyLoginGet200Response{
		Token: jwtToken,
	}

	c.JSON(http.StatusOK, response)
}

func (api *NoAuthAPI) LoginPost(c *gin.Context) {
	var loginRequest models.LoginPostRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := database.GetUserByEmail(loginRequest.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	jwtToken, err := auth.GenerateJwtToken(user.Email, user.UserType)
	if err != nil {
		response := models.DummyLoginGet500Response{
			Message:   "Failed to generate token",
			RequestId: uuid.New().String(),
			Code:      http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := models.DummyLoginGet200Response{
		Token: jwtToken,
	}

	c.JSON(http.StatusOK, response)
}

func (api *NoAuthAPI) RegisterPost(c *gin.Context) {
	var registerRequest models.RegisterPostRequest

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if registerRequest.UserType != models.CLIENT && registerRequest.UserType != models.MODERATOR {
		log.Printf("Invalid user type: %v", registerRequest.UserType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Email:    registerRequest.Email,
		Password: string(hashedPassword),
		UserType: string(registerRequest.UserType),
	}

	if err := database.CreateUser(&user); err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(200, gin.H{"status": "OK"})
}
