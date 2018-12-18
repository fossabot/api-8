package server

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/devlover-id/api/pkg/api"
	"github.com/devlover-id/api/pkg/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var originCheck = regexp.MustCompile("https?://(.+.)?pinterkode.id")

// buildRouter construct and return http router
func buildRouter(prod bool) http.Handler {
	if prod {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx := context.Background()
	router := gin.New()

	//add cors to router
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "HEAD"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		AllowOriginFunc: func(origin string) bool {
			return originCheck.MatchString(origin)
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/", api.WrapGin(ctx, index))
	router.GET("/ping", api.WrapGin(ctx, ping))

	v1 := router.Group("/v1")
	{
		v1.POST("/users", api.WrapGin(ctx, user.V1PostUser))
	}
	return router
}

func index(ctx context.Context, req api.Request) api.Response {
	return api.OKResp(map[string]string{
		"hello": "world",
	})
}

func ping(ctx context.Context, req api.Request) api.Response {
	return api.OKResp(map[string]string{
		"message": "pong",
	})
}
