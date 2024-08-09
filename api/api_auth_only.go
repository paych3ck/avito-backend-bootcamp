package api

import (
	"avito-backend-bootcamp/auth"
	"avito-backend-bootcamp/database"
	"avito-backend-bootcamp/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthOnlyAPI struct {
}

func (api *AuthOnlyAPI) FlatCreatePost(c *gin.Context) {
	jwtTokenStr := c.GetHeader("Authorization")
	if jwtTokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	_, err := auth.ValidateJwtToken(jwtTokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		return
	}

	var createFlatRequest models.FlatCreatePostRequest
	if err := c.ShouldBindJSON(&createFlatRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flat := models.Flat{
		HouseId:    createFlatRequest.HouseId,
		FlatNumber: createFlatRequest.FlatNumber,
		Price:      createFlatRequest.Price,
		Rooms:      createFlatRequest.Rooms,
		Status:     models.CREATED,
	}

	if err := database.CreateFlat(&flat); err != nil {
		log.Printf("Error creating flat: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create flat"})
		return
	}

	if err := database.UpdateHouse(createFlatRequest.HouseId); err != nil {
		log.Printf("Error updating house: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update house"})
		return
	}

	c.JSON(http.StatusOK, flat)
}

func (api *AuthOnlyAPI) HouseIdGet(c *gin.Context) {
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

	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid house ID"})
		return
	}

	var flats []models.Flat
	if claims.UserType == string(models.MODERATOR) {
		flats, err = database.GetFlatsByHouseID(houseID, "all")
	} else {
		flats, err = database.GetFlatsByHouseID(houseID, string(models.APPROVED))
	}
	if err != nil {
		log.Printf("Error getting flats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get flats"})
		return
	}

	response := models.HouseIdGet200Response{
		Flats: flats,
	}
	c.JSON(http.StatusOK, response)
}

func (api *AuthOnlyAPI) HouseIdSubscribePost(c *gin.Context) {
	// Your handler implementation
	c.JSON(200, gin.H{"status": "OK"})
}
