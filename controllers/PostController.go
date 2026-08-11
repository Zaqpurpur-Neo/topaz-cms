package controllers

import "github.com/gin-gonic/gin"

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

func (pc *PostController) CreatePost(c *gin.Context) {

}

func (pc *PostController) GetSinglePost(c *gin.Context) {
}

func (pc *PostController) GetPosts(c *gin.Context) {

}
