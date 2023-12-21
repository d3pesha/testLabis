package handlers

import (
	"github.com/d3pesha/testLabis/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

// ClientService struct represents a service for handling clients.
type ClientService struct {
	db *sqlx.DB
}

// NewClientService creates a new instance of ClientService.
func NewClientService(db *sqlx.DB) *ClientService {
	return &ClientService{db: db}
}

// @Summary Get all clients
// @Description Get a list of all clients
// @Produce json
// @Success 200 {array} models.Client
// @Router /clients [get]
func (cs *ClientService) GetClients(ctx *gin.Context) {
	var clients []models.Client

	err := cs.db.Select(&clients, "SELECT * FROM users")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range clients {
		var objects []models.Object
		err = cs.db.Select(&objects, "SELECT * FROM objects WHERE id_user = $1", clients[i].ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		clients[i].Objects = objects
	}

	ctx.JSON(http.StatusOK, clients)
}

// @Summary Get a client by ID
// @Description Get a client by its ID
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {object} models.Client
// @Failure 400 {object} models.ErrorResponse
// @Router /clients/{id} [get]
func (cs *ClientService) GetClientByID(ctx *gin.Context) {
	var client models.Client

	id := ctx.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	err = cs.db.Get(&client, "SELECT * FROM users WHERE id = $1", clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var objects []models.Object
	err = cs.db.Select(&objects, "SELECT * FROM objects WHERE id_user = $1", clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	client.Objects = objects

	ctx.JSON(http.StatusOK, client)
}

// @Summary Create a new client
// @Description Create a new client
// @Accept json
// @Produce json
// @Param client body models.Client true "Client information"
// @Success 201 {object} models.Client
// @Failure 400 {object} models.ErrorResponse
// @Router /clients [post]
func (cs *ClientService) CreateClient(ctx *gin.Context) {
	var client models.Client
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var count int
	err := cs.db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1", client.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email is already registered"})
		return
	}

	err = cs.db.QueryRowx("INSERT INTO users (fio, email, password) VALUES ($1, $2, $3) RETURNING id, fio, email", client.FIO, client.Email, client.Password).StructScan(&client)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, client)

}

// @Summary Delete a client by ID
// @Description Delete a client by its ID
// @Produce json
// @Param id path int true "Client ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Router /clients/{id} [delete]
func (cs *ClientService) DeleteClient(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	var objectCount int
	err = cs.db.Get(&objectCount, "SELECT COUNT(*) FROM objects WHERE id_user = $1", clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if objectCount > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "can't delete client with objects"})
		return
	}

	_, err = cs.db.Exec("DELETE FROM users WHERE id = $1", clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusNoContent, nil)

}

// @Summary Update a client by ID
// @Description Update a client by its ID
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param client body models.Client true "Updated client information"
// @Success 200 {object} models.Client
// @Failure 400 {object} models.ErrorResponse
// @Router /clients/{id} [PUT]
func (cs *ClientService) UpdateClient(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	var updatedClient models.Client
	if err = ctx.ShouldBindJSON(&updatedClient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = cs.db.QueryRowx("UPDATE users SET fio = $1, email = $2, password = $3 WHERE id = $4	RETURNING id, fio, email",
		updatedClient.FIO, updatedClient.Email, updatedClient.Password, clientID).StructScan(&updatedClient)

	ctx.JSON(http.StatusOK, updatedClient)

}
