package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rustingoff/admin_panel_rep/internal/model"
	"github.com/rustingoff/admin_panel_rep/internal/service"
	"log"
	"net/http"
	"strconv"
)

type ClientController interface {
	CreateClient(c *gin.Context)
	UpdateClient(c *gin.Context)
	DeleteClient(c *gin.Context)

	GetAllClients(c *gin.Context)
	GetClient(c *gin.Context)
}

type clientController struct {
	cService service.ClientService
}

func GetClientController(s service.ClientService) ClientController {
	return &clientController{cService: s}
}

func (cc *clientController) CreateClient(c *gin.Context) {
	var client model.Client

	err := c.BindJSON(&client)
	if err != nil {
		log.Printf("FAILED bind json to client structure with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid request")
		return
	}

	err = cc.cService.CreateClient(client)
	if err != nil {
		log.Printf("FAILED create client with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "server can't create a client")
		return
	}

	c.JSON(http.StatusCreated, "created")
}

func (cc *clientController) UpdateClient(c *gin.Context) {
	var client model.ClientUpdate

	clientID, err := strconv.Atoi(c.Param("clientID"))
	if err != nil || clientID < 1 {
		log.Printf("FAILED convert param to int with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid url")
		return
	}

	err = c.BindJSON(&client)
	if err != nil {
		log.Printf("FAILED bind json to client update structure with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid request")
	}

	err = cc.cService.UpdateClient(client, uint(clientID))
	if err != nil {
		log.Printf("FAILED update client with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "server can't update a client")
	}

	c.JSON(http.StatusOK, "updated")
}

func (cc *clientController) DeleteClient(c *gin.Context) {
	clientID, err := strconv.Atoi(c.Param("clientID"))
	if err != nil || clientID < 1 {
		log.Printf("FAILED convert param to int with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid url")
		return
	}

	err = cc.cService.DeleteClient(uint(clientID))
	if err != nil {
		log.Printf("FAILED delete client with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "server can't delete a client")
	}

	c.JSON(http.StatusOK, "updated")
}

func (cc *clientController) GetAllClients(c *gin.Context) {

	response, err := cc.cService.GetAllClients()
	if err != nil {
		log.Printf("FAILED to get clients with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "invalid server response")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (cc *clientController) GetClient(c *gin.Context) {
	clientID, err := strconv.Atoi(c.Param("clientID"))
	if err != nil || clientID < 1 {
		log.Printf("FAILED convert param to int with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid url")
		return
	}

	response, err := cc.cService.GetClient(uint(clientID))
	if err != nil {
		log.Printf("FAILED to get client with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "invalid server response")
		return
	}

	c.JSON(http.StatusOK, response)
}
