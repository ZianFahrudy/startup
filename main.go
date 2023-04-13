package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/payment"
	"bwastartup/transaction"
	"bwastartup/user"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/example"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var s snap.Client

func main() {
	// fmt.Println("================ Request with Snap Client ================")
	// initializeSnapClient()
	// createTransaction()

	// fmt.Println("================ Request Snap token ================")
	// createTokenTransactionWithGateway()

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:root@tcp(localhost:8080)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}
	campaignRepository := campaign.NewRepository(db)
	usersRepository := user.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(usersRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService, paymentService)

	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/images", "./images")

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email-check", userHandler.CheckEmailAvaibility)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	router.Run(":9888")

}
func initializeSnapClient() {
	s.New("SB-Mid-server-aZHMJQxWrIDZRMQaJhb2ogVr", midtrans.Sandbox)
}
func createTransaction() {
	// Optional : here is how if you want to set append payment notification for this request
	s.Options.SetPaymentAppendNotification("https://example.com/append")

	// Optional : here is how if you want to set override payment notification for this request
	s.Options.SetPaymentOverrideNotification("https://example.com/override")
	// Send request to Midtrans Snap API

	resp, err := s.CreateTransaction(GenerateSnapReq())
	if err != nil {
		fmt.Println("Error :", err.GetMessage())
	}
	fmt.Println("Response : ", resp)
}

func GenerateSnapReq() *snap.Request {

	// Initiate Customer address
	custAddress := &midtrans.CustomerAddress{
		FName:       "John",
		LName:       "Doe",
		Phone:       "081234567890",
		Address:     "Baker Street 97th",
		City:        "Jakarta",
		Postcode:    "16000",
		CountryCode: "IDN",
	}

	// Initiate Snap Request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  "MID-GO-ID-" + example.Random(),
			GrossAmt: 200000,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:    "John",
			LName:    "Doe",
			Email:    "john@doe.com",
			Phone:    "081234567890",
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "ITEM1",
				Price: 200000,
				Qty:   1,
				Name:  "Someitem",
			},
		},
	}
	return snapReq
}

func createTokenTransactionWithGateway() {
	s.Options.SetPaymentOverrideNotification("https://example.com/url2")

	resp, err := s.CreateTransactionUrl(GenerateSnapReq())
	if err != nil {
		fmt.Println("Error :", err.GetMessage())
	}
	fmt.Println("Responsess : ", resp)
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return

		}

		tokenString := "xx"

		arrayToken := strings.Split(authHeader, " ")

		fmt.Println(len(arrayToken))

		if len(tokenString) == 2 {
			tokenString = arrayToken[1]

		}
		fmt.Println(tokenString)

		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)

		if err != nil {
			response := helper.APIResponse("Unauthotized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return

		}

		c.Set("currentUser", user)

	}

}
