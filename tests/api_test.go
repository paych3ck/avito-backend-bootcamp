package tests

import (
	"avito-backend-bootcamp/database"
	"avito-backend-bootcamp/models"
	"avito-backend-bootcamp/routers"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func getToken(router *gin.Engine, userType string) (string, error) {
	payload := map[string]string{
		"user_type": userType,
	}

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("GET", "/dummyLogin", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return "", fmt.Errorf("failed to get token: %v", w.Body.String())
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return "", err
	}

	return response["token"], nil
}

func TestMain(m *testing.M) {
	os.Chdir("/app")

	err := database.InitTestDB()
	if err != nil {
		log.Fatalf("Failed to initialize test database: %v", err)
	}

	err = database.ClearTestDB()
	if err != nil {
		log.Fatalf("Failed to clear test database: %v", err)
	}

	code := m.Run()
	database.DB.Close()
	os.Exit(code)
}

func TestHouseCreatePostModerator(t *testing.T) {
	routes := routers.ApiHandleFunctions{}
	router := routers.NewRouter(routes)

	token, err := getToken(router, "moderator")
	assert.NoError(t, err)

	developer := "TestDeveloperModerator"
	payload := models.HouseCreatePostRequest{
		Address:   "TestAddressModerator",
		Year:      2024,
		Developer: &developer,
	}

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/house/create", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.House
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, payload.Address, response.Address)
	assert.Equal(t, payload.Year, response.Year)
	assert.Equal(t, *payload.Developer, *response.Developer)
	assert.NotZero(t, response.Id)
	assert.WithinDuration(t, time.Now(), response.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), response.UpdateAt, time.Second)
}

func TestHouseCreatePostClient(t *testing.T) {
	routes := routers.ApiHandleFunctions{}
	router := routers.NewRouter(routes)

	token, err := getToken(router, "client")
	assert.NoError(t, err)

	developer := "TestDeveloperClient"
	payload := models.HouseCreatePostRequest{
		Address:   "TestAddressClient",
		Year:      2024,
		Developer: &developer,
	}

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/house/create", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Only moderator can create house", response["error"])
}

func TestFlatCreatePostModerator(t *testing.T) {
	routes := routers.ApiHandleFunctions{}
	router := routers.NewRouter(routes)

	token, err := getToken(router, "moderator")
	assert.NoError(t, err)

	payload := models.FlatCreatePostRequest{
		HouseId:    1,
		FlatNumber: 101,
		Price:      10101,
		Rooms:      1,
	}

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/flat/create", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Flat
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, payload.HouseId, response.HouseId)
	assert.Equal(t, payload.FlatNumber, response.FlatNumber)
	assert.Equal(t, payload.Price, response.Price)
	assert.Equal(t, payload.Rooms, response.Rooms)
}

func TestFlatCreatePostClient(t *testing.T) {
	routes := routers.ApiHandleFunctions{}
	router := routers.NewRouter(routes)

	token, err := getToken(router, "client")
	assert.NoError(t, err)

	payload := models.FlatCreatePostRequest{
		HouseId:    1,
		FlatNumber: 202,
		Price:      20202,
		Rooms:      2,
	}

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/flat/create", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Flat
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, payload.HouseId, response.HouseId)
	assert.Equal(t, payload.FlatNumber, response.FlatNumber)
	assert.Equal(t, payload.Price, response.Price)
	assert.Equal(t, payload.Rooms, response.Rooms)
}

func TestFlatUpdatePostModerator(t *testing.T) {
	routes := routers.ApiHandleFunctions{}
	router := routers.NewRouter(routes)

	token, err := getToken(router, "moderator")
	assert.NoError(t, err)

	payload := models.FlatUpdatePostRequest{
		Id:     1,
		Status: models.APPROVED,
	}

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/flat/update", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Flat
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, payload.Id, response.Id)
	assert.Equal(t, payload.Status, response.Status)
}

func TestHouseIdGetModerator(t *testing.T) {
	routes := routers.ApiHandleFunctions{}
	router := routers.NewRouter(routes)

	token, err := getToken(router, "moderator")
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/house/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HouseIdGet200Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response.Flats))
}

func TestHouseIdGetClient(t *testing.T) {
	routes := routers.ApiHandleFunctions{}
	router := routers.NewRouter(routes)

	token, err := getToken(router, "client")
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/house/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HouseIdGet200Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.Flats))
}
