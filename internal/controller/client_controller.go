package controller

import (
	"fmt"
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
	idnp, err := strconv.ParseUint(c.PostForm("idnp"), 10, 64)
	if err != nil {
		log.Printf("invalid idnp with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid idnp")
		return
	}

	sum, err := strconv.ParseUint(c.PostForm("sum"), 10, 32)
	if err != nil {
		log.Printf("invalid sum with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid sum")
		return
	}

	period, err := strconv.ParseUint(c.PostForm("time"), 10, 32)
	if err != nil {
		log.Printf("invalid time with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid time")
		return
	}

	monthRate, err := strconv.ParseFloat(c.PostForm("monthly_rate"), 32)
	if err != nil {
		log.Printf("invalid monthly rate with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid monthly rate")
		return
	}

	var client = model.Client{
		NrContract:  c.PostForm("nr_contract"),
		FirstName:   c.PostForm("first_name"),
		LastName:    c.PostForm("last_name"),
		IDNP:        idnp,
		Phone:       c.PostForm("phone"),
		Sum:         sum,
		Time:        period,
		SignDate:    c.PostForm("sign_date"),
		MonthlyRate: fmt.Sprintf("%.2f", monthRate),
	}

	err = cc.cService.CreateClient(client)
	if err != nil {
		log.Printf("FAILED create client with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "can't create a client")
		return
	}

	c.HTML(http.StatusCreated, "index.html", nil)
}

func (cc *clientController) UpdateClient(c *gin.Context) {

	//todo return client for update
	idnp, err := strconv.ParseUint(c.PostForm("idnp"), 10, 64)
	if err != nil {
		idnp = 0
	}

	sum, err := strconv.ParseUint(c.PostForm("sum"), 10, 32)
	if err != nil {
		sum = 0
	}

	period, err := strconv.ParseUint(c.PostForm("time"), 10, 32)
	if err != nil {
		period = 0
	}

	var client = model.ClientUpdate{
		NrContract: c.PostForm("nr_contract"),
		FirstName:  c.PostForm("first_name"),
		LastName:   c.PostForm("last_name"),
		IDNP:       idnp,
		Phone:      c.PostForm("phone"),
		Sum:        sum,
		Time:       period,
		SignDate:   c.PostForm("sign_date"),
	}

	clientID, err := strconv.ParseUint(c.PostForm("ID"), 10, 8)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid data")
		return
	}

	err = cc.cService.UpdateClient(client, uint(clientID))
	if err != nil {
		log.Printf("FAILED update client with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "server can't update a client")
		return
	}

	c.HTML(http.StatusOK, "client_update.html", clientID)
}

func (cc *clientController) DeleteClient(c *gin.Context) {

	clientID, err := strconv.Atoi(c.Query("id"))
	if err != nil || clientID < 1 {
		log.Printf("FAILED convert query to int with error: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid url")
		return
	}

	err = cc.cService.DeleteClient(uint(clientID))
	if err != nil {
		log.Printf("FAILED delete client with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "server can't delete a client")
		return
	}

	c.Redirect(http.StatusPermanentRedirect, "/api/clients")
}

func (cc *clientController) GetAllClients(c *gin.Context) {

	response, err := cc.cService.GetAllClients()
	if err != nil {
		log.Printf("FAILED to get clients with error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, "invalid server response")
		return
	}

	c.HTML(http.StatusOK, "client.html", response)
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
