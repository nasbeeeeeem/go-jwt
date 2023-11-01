package server

import (
	"log"
	"net/http"

	"go-jwt/pkg/infrastructure/database"
	"go-jwt/pkg/infrastructure/repositoryimpl"

	"go-jwt/pkg/interfaces/api/handler"
	"go-jwt/pkg/interfaces/api/middleware"
	"go-jwt/pkg/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// サーバー起動処理
func Server(addr string) {
	// DI
	db, err := database.NewDBClient()
	if err != nil {
		panic(err)
	}
	userRepoImp := repositoryimpl.NewUserRepo(db)
	userUseCase := usecase.NewUseCase(userRepoImp)
	userHandler := handler.NewHandler(userUseCase)

	r = gin.Default()

	r.Use(cors.New(cors.Config{
		// 許可するアクセス元
		AllowOrigins: []string{
			"http://192.168.10.114:3000",
			"http://localhost:3000",
		},
		// 許可するメソッド
		AllowMethods: []string{
			"POST",
			"GET",
		},

		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		// Cookieの許可
		AllowCredentials: true,
	}))

	// アクセスポイント
	r.POST("/signup", userHandler.HandlerSing)
	r.POST("login", userHandler.HandlerLogin)
	r.GET("/logout", userHandler.HandlerLogout)

	secured := r.Group("/secured").Use(middleware.Auth())
	secured.GET("/ping", Ping)

	log.Println("Server running...")
	if err := r.Run(addr); err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}
