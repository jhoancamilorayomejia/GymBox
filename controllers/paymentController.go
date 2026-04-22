package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

func CreatePreference(c *gin.Context) {
	var requestData struct {
		Amount    float64 `json:"amount"`
		Username  string  `json:"username"`
		PlanType  string  `json:"planType"`
		Quantity  int     `json:"quantity"`
		DateStart string  `json:"dateStart"`
		DateEnd   string  `json:"dateEnd"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	accessToken := os.Getenv("MP_ACCESS_TOKEN")
	if accessToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MP_ACCESS_TOKEN not set"})
		return
	}

	cfg, err := config.New(accessToken)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration error"})
		return
	}

	client := preference.NewClient(cfg)

	// Título enriquecido con fechas
	titulo := fmt.Sprintf(
		"RayoBox · %s x%d · %s · Del %s al %s",
		requestData.PlanType,
		requestData.Quantity,
		requestData.Username,
		requestData.DateStart,
		requestData.DateEnd,
	)

	request := preference.Request{
		Items: []preference.ItemRequest{
			{
				Title:     titulo,
				Quantity:  1,
				UnitPrice: requestData.Amount,
			},
		},
		BackURLs: &preference.BackURLsRequest{
			Success: "http://localhost:8080/success",
			Failure: "http://localhost:8080/failure",
			Pending: "http://localhost:8080/pending",
		},
	}

	resource, err := client.Create(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating preference"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"preferenceId": resource.ID})
}
