package controller

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nescool101/rentManager/service"
)

func StartHTTPServer() error {
	router := gin.Default()
	router.GET("/payers", getPayers)
	router.GET("/validate_email", validateEmail)
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	return router.Run(":" + port)
}

func getPayers(c *gin.Context) {
	//service.LoadPayers()
	//c.JSON(http.StatusOK, gin.H{"payers": service.GetAllPayers()})
}

func validateEmail(c *gin.Context) {
	service.NotifyAll()
	c.JSON(http.StatusOK, gin.H{"status": "Email notifications triggered."})
}
