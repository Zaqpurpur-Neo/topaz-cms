package controllers

import (
	"github.com/gin-gonic/gin"
)

type KritaArtworkController struct {
	RoutePath string
}

func NewKritaArtworkController() *KritaArtworkController {
	return &KritaArtworkController{
		RoutePath: "",
	}
}

func (c *KritaArtworkController) HandleRequest(ctx *gin.Context) {

}
