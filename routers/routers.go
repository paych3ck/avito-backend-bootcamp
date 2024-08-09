package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"avito-backend-bootcamp/api"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
}

func NewRouter(handleFunctions ApiHandleFunctions) *gin.Engine {
	return NewRouterWithGinEngine(gin.Default(), handleFunctions)
}

func NewRouterWithGinEngine(router *gin.Engine, handleFunctions ApiHandleFunctions) *gin.Engine {
	for _, route := range getRoutes(handleFunctions) {
		if route.HandlerFunc == nil {
			route.HandlerFunc = DefaultHandleFunc
		}
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

func DefaultHandleFunc(c *gin.Context) {
	c.String(http.StatusNotImplemented, "501 not implemented")
}

type ApiHandleFunctions struct {
	AuthOnlyAPI        api.AuthOnlyAPI
	ModerationsOnlyAPI api.ModerationsOnlyAPI
	NoAuthAPI          api.NoAuthAPI
}

func getRoutes(handleFunctions ApiHandleFunctions) []Route {
	return []Route{
		{
			"FlatCreatePost",
			http.MethodPost,
			"/flat/create",
			handleFunctions.AuthOnlyAPI.FlatCreatePost,
		},
		{
			"HouseIdGet",
			http.MethodGet,
			"/house/:id",
			handleFunctions.AuthOnlyAPI.HouseIdGet,
		},
		{
			"HouseIdSubscribePost",
			http.MethodPost,
			"/house/:id/subscribe",
			handleFunctions.AuthOnlyAPI.HouseIdSubscribePost,
		},
		{
			"FlatUpdatePost",
			http.MethodPost,
			"/flat/update",
			handleFunctions.ModerationsOnlyAPI.FlatUpdatePost,
		},
		{
			"HouseCreatePost",
			http.MethodPost,
			"/house/create",
			handleFunctions.ModerationsOnlyAPI.HouseCreatePost,
		},
		{
			"DummyLoginGet",
			http.MethodGet,
			"/dummyLogin",
			handleFunctions.NoAuthAPI.DummyLoginGet,
		},
		{
			"LoginPost",
			http.MethodPost,
			"/login",
			handleFunctions.NoAuthAPI.LoginPost,
		},
		{
			"RegisterPost",
			http.MethodPost,
			"/register",
			handleFunctions.NoAuthAPI.RegisterPost,
		},
	}
}
