package handlers

import (
	"database/sql"
	"github.com/d3pesha/testLabis/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
)

// @tags Objects
type ObjectService struct {
	db *sqlx.DB
}

func NewObjectService(db *sqlx.DB) *ObjectService {
	return &ObjectService{db: db}
}

// @Summary Get objects by object ID
// @Description Get a list of objects by object ID
// @Produce json
// @Param id path int true "object ID"
// @Success 200 {array} models.Object
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /objects/{id} [get]
// @tags objects
func (os *ObjectService) GetObjects(ctx *gin.Context) {
	var objects []models.Object

	id := ctx.Param("id")
	objectID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid object ID"})
		return
	}

	err = os.db.Select(&objects, "SELECT * FROM objects WHERE id = $1", objectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range objects {
		var contracts []models.Contract
		err = os.db.Select(&contracts, "SELECT * FROM contracts WHERE id_object = $1", objects[i].ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		objects[i].Contracts = contracts
		ctx.JSON(http.StatusOK, objects)
	}
}

// @Summary Create a new object
// @Description Create a new object
// @Accept json
// @Produce json
// @Param object body models.Object true "Object information"
// @Success 201 {object} models.Object
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /objects [post]
// @tags objects
func (os *ObjectService) CreateObjects(ctx *gin.Context) {
	var object models.Object
	if err := ctx.ShouldBindJSON(&object); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var client models.Client

	err := os.db.Get(&client, "SELECT * FROM users WHERE id = $1", object.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	err = os.db.QueryRowx("INSERT INTO objects (id_user, address, is_visible) VALUES ($1, $2, $3) RETURNING id, id_user, address, is_visible",
		object.UserID, object.Address, object.IsVisible).Scan(&object.ID, &object.UserID, &object.Address, &object.IsVisible)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cs := NewContractService(os.db)
	cs.CreateContractMethod(ctx, object.ID)

	ctx.JSON(http.StatusCreated, object)
}

// @Summary Delete an object by ID
// @Description Delete an object by its ID
// @Produce json
// @Param id path int true "Object ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /objects/{id} [delete]
// @tags objects
func (os *ObjectService) DeleteObjects(ctx *gin.Context) {
	id := ctx.Param("id")
	objectID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid object ID"})
		return
	}

	var contractStatus bool
	err = os.db.Get(&contractStatus, "SELECT is_visible FROM objects WHERE id = $1", objectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if contractStatus == true {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "can't delete object with contract"})
		return
	}

	_, err = os.db.Exec("DELETE FROM objects WHERE id = $1", objectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// @Summary Update an object by ID
// @Description Update an object by its ID
// @Accept json
// @Produce json
// @Param id path int true "Object ID"
// @Param object body models.Object true "Updated object information"
// @Success 200 {object} models.Object
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /objects/{id} [PUT]
// @tags objects
func (os *ObjectService) UpdateObject(ctx *gin.Context) {
	id := ctx.Param("id")
	objectID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid object ID"})
		return
	}

	var updatedObject models.Object
	if err = ctx.ShouldBindJSON(&updatedObject); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = os.db.QueryRowx("UPDATE objects SET address = $1, is_visible = $2 WHERE id = $3	RETURNING id, id_user, address, is_visible",
		updatedObject.Address, updatedObject.IsVisible, objectID).StructScan(&updatedObject)

	ctx.JSON(http.StatusOK, updatedObject)
}
