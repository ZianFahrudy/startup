package main

import (
	"bwastartup/auth"
	"bwastartup/handler"
	"bwastartup/user"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:root@tcp(localhost:8080)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	usersRepository := user.NewRepository(db)
	userService := user.NewService(usersRepository)
	authService := auth.NewService()

	fmt.Println(authService.GenerateToken(1001))

	userService.SaveAvatar(3, "images/1-profile.png")

	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email-check", userHandler.CheckEmailAvaibility)
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run(":9888")

}
