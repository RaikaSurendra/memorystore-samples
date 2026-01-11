package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type HomeController struct {
	APIURL string
}

func NewHomeController() *HomeController {
	url := os.Getenv("API_URL")
	if url == "" {
		url = "localhost:8080"
	}
	return &HomeController{
		APIURL: url,
	}
}

func (hc *HomeController) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"apiUrl": hc.APIURL,
	})
}
