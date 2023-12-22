package handlers

import (
	"database/sql"
	"fmt"
	"github.com/d3pesha/testLabis/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"strconv"
	"time"
)

var contractCounter = 1

// @tags Contracts
type ContractService struct {
	db *sqlx.DB
}

func NewContractService(db *sqlx.DB) *ContractService {
	return &ContractService{db: db}
}

func generateContractNumber(db *sqlx.DB) (string, error) {
	currentYear := time.Now().Year() % 100
	contractNumber := fmt.Sprintf("%06d/%02d", contractCounter, currentYear)
	contractCounter++
	var count int

	err := db.Get(&count, "SELECT COUNT(*) FROM contracts WHERE number = $1", contractNumber)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return generateContractNumber(db)
	}

	return contractNumber, nil
}

// @Summary Get a contract by ID
// @Description Get a contract by its ID
// @Produce json
// @Param id path int true "Contract ID"
// @Success 200 {object} models.Contract
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /contracts/{id} [get]
// @tags contracts
func (os *ContractService) GetContract(ctx *gin.Context) {
	var contract models.Contract

	id := ctx.Param("id")
	contractID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid contract ID"})
		return
	}

	err = os.db.Get(&contract, "SELECT * FROM contracts WHERE id = $1", contractID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, contract)
}

func (os *ContractService) CreateContractMethod(ctx *gin.Context, objectID int) {
	var contract models.Contract

	currentDateTime := time.Now()

	contractNumber, err := generateContractNumber(os.db)
	if err != nil {
		return
	}

	_, err = os.db.Exec("UPDATE contracts SET status = false WHERE id_object = $1", contract.ObjectID)
	_, err = os.db.Exec("UPDATE objects SET is_visible = false WHERE id = $1", contract.ObjectID)

	err = os.db.QueryRowx("INSERT INTO contracts (id_object, data, number, status) VALUES ($1, $2, $3, $4) RETURNING id, id_object, data, number, status",
		objectID, currentDateTime, contractNumber, true).Scan(&contract.ID, &contract.ObjectID, &contract.Data, &contract.Number, &contract.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	/*var count int

	err = os.db.Get(&count, "SELECT COUNT(*) FROM objects WHERE id = $1", objectID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}


	if count > 1 {
		_, err = os.db.Exec("UPDATE contracts SET status = false WHERE id_object = $1", contract.ObjectID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = os.db.Exec("UPDATE objects SET is_visible = false WHERE id = $1", contract.ObjectID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}*/

	ctx.JSON(http.StatusCreated, contract)

}

// @Summary Create a new contract
// @Description Create a new contract
// @Accept json
// @Produce json
// @Param contract body models.Contract true "Contract information"
// @Success 201 {object} models.Contract
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /contracts [post]
// @tags contracts
func (os *ContractService) CreateContract(ctx *gin.Context) {
	var contract models.Contract
	if err := ctx.ShouldBindJSON(&contract); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var object models.Object

	err := os.db.Get(&object, "SELECT * FROM objects WHERE id = $1", contract.ObjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Object not found"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	currentDateTime := time.Now()

	contractNumber, err := generateContractNumber(os.db)
	if err != nil {
		return
	}

	_, err = os.db.Exec("UPDATE contracts SET status = false WHERE id_object = $1", contract.ObjectID)
	//_, err = os.db.Exec("UPDATE objects SET is_visible = false WHERE id = $1", contract.ObjectID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = os.db.QueryRowx("INSERT INTO contracts (id_object, 	data, number, status) VALUES ($1, $2, $3, $4) RETURNING id, id_object, data, number, status",
		contract.ObjectID, currentDateTime, contractNumber, true).Scan(&contract.ID, &contract.ObjectID, &contract.Data, &contract.Number, &contract.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var statusObject bool

	err = os.db.Get(&statusObject, "SELECT is_visible FROM objects WHERE id = $1", contract.ObjectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if statusObject != true {
		_, err = os.db.Exec("UPDATE objects SET is_visible = true WHERE id = $1", contract.ObjectID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusCreated, contract)

}

// @Summary Delete a contract by ID
// @Description Delete a contract by its ID
// @Produce json
// @Param id path int true "Contract ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /contracts/{id} [delete]
// @tags contracts
func (os *ContractService) DeleteContract(ctx *gin.Context) {
	id := ctx.Param("id")
	contractID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid contract ID"})
		return
	}

	_, err = os.db.Exec("DELETE FROM contracts WHERE id = $1", contractID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// @Summary Update a contract by ID
// @Description Update a contract by its ID
// @Accept json
// @Produce json
// @Param id path int true "Contract ID"
// @Param contract body models.Contract true "Updated contract information"
// @Success 200 {object} models.Contract
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /contracts/{id} [PUT]
// @tags contracts
func (os *ContractService) UpdateContract(ctx *gin.Context) {
	id := ctx.Param("id")
	contractID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid contract ID"})
		return
	}

	var updatedContract models.Contract
	if err = ctx.ShouldBindJSON(&updatedContract); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = os.db.QueryRowx("UPDATE contracts SET status = $1 WHERE id = $2	RETURNING id, id_object, data, number, status",
		updatedContract.Status, contractID).StructScan(&updatedContract)

	ctx.JSON(http.StatusOK, updatedContract)
}
