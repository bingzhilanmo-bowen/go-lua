package router

import (
	"bingzhilanmo/go-lua/config"
	"bingzhilanmo/go-lua/controller"
	_ "bingzhilanmo/go-lua/docs"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	healthURL = "health"
	baseURL   = "/api/v1/"
)

func RegisterRouter(router *gin.Engine) {

	conf := config.GetGlobalConfig()
	if conf.Base.InitSwagger {
		url := ginSwagger.URL(fmt.Sprintf("swagger/doc.json")) // The url pointing to API definition
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	v1 := router.Group(baseURL)

	keyword := v1.Group("/keyword")
	{
		keyword.POST("/created", controller.CreateNewKeyword)
		keyword.PUT("/updated", controller.UpdateKeyword)
		keyword.GET("/page", controller.QueryKeywordPage)
		keyword.GET("/detail", controller.QueryKeywordDetail)
		keyword.POST("/run", controller.RunKeywordScript)
	}
}