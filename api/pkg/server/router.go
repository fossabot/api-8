package server

import (
	"net/http"
	"regexp"
	"time"

	"github.com/devlover-id/api/pkg/api"
	"github.com/devlover-id/api/pkg/auth"
	"github.com/devlover-id/api/pkg/oauth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var originCheck = regexp.MustCompile("https://(.+.)?devlover.id")

// buildRouter construct and return http router
func buildRouter(prod bool) http.Handler {
	if prod {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodGet},
		AllowHeaders:     []string{"Origin", "Authorization"},
		ExposeHeaders:    []string{},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return originCheck.MatchString(origin)
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/", api.WrapGin(index))
	router.GET("/ping", api.WrapGin(ping))

	v1 := router.Group("/v1")
	{
		v1.POST("/auth/register", api.WrapGin(auth.V1PostRegister))
		v1.POST("/auth/login", api.WrapGin(auth.V1PostLogin))
		v1.GET("/oauth/github", api.WrapGin(oauth.V1GetGithub))
	}
	return router
}

func index(req api.Request) api.Response {
	return api.OKResp(map[string]string{
		"hello": "world",
	})
}

func ping(req api.Request) api.Response {
	return api.OKResp(map[string]string{
		"message": "pong",
	})
}
