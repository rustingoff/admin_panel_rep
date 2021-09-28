package controller

import "github.com/gin-gonic/gin"

type ClientController interface {
	CreateClient(c *gin.Context)
}

type clientController struct {
}

func GetClientController() ClientController {
	return &clientController{}
}

func (cc *clientController) CreateClient(c *gin.Context) {

}
