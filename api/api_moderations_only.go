package api

import (
	"avito-backend-bootcamp/auth"
	"avito-backend-bootcamp/database"
	"avito-backend-bootcamp/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ModerationsOnlyAPI struct {
}

func (api *ModerationsOnlyAPI) FlatUpdatePost(c *gin.Context) {
	jwtTokenStr := c.GetHeader("Authorization")

	if jwtTokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	claims, err := auth.ValidateJwtToken(jwtTokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		return
	}

	if claims.UserType != string(models.MODERATOR) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only moderator can update flat status"})
		return
	}

	var updateFlatRequest models.FlatUpdatePostRequest
	if err := c.ShouldBindJSON(&updateFlatRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentFlat, err := database.GetFlatByID(updateFlatRequest.Id)
	if err != nil {
		log.Printf("Error fetching flat: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch flat"})
		return
	}

	if currentFlat.Status == models.ON_MODERATION && updateFlatRequest.Status == models.ON_MODERATION {
		c.JSON(http.StatusConflict, gin.H{"error": "Flat is already under moderation"})
		return
	}

	flat, err := database.UpdateFlatStatus(updateFlatRequest.Id, string(updateFlatRequest.Status))
	if err != nil {
		log.Printf("Error updating flat status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update flat status"})
		return
	}

	c.JSON(http.StatusOK, flat)
}

func (api *ModerationsOnlyAPI) HouseCreatePost(c *gin.Context) {
	jwtTokenStr := c.GetHeader("Authorization")

	if jwtTokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	claims, err := auth.ValidateJwtToken(jwtTokenStr)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		return
	}

	if claims.UserType != string(models.MODERATOR) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only moderator can create house"})
		return
	}

	var createHouseRequest models.HouseCreatePostRequest

	if err := c.ShouldBindJSON(&createHouseRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	house := models.House{
		Address:   createHouseRequest.Address,
		Year:      createHouseRequest.Year,
		Developer: createHouseRequest.Developer,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
	}

	if err := database.CreateHouse(&house); err != nil {
		log.Printf("Error creating house: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create house"})
		return
	}

	c.JSON(http.StatusOK, house)
}
