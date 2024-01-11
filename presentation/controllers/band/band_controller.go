package bandcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	accountusecase "github.com/mazurco066/playliter-api-go/data/usecases/account"
	bandusecase "github.com/mazurco066/playliter-api-go/data/usecases/band"
	bandinputs "github.com/mazurco066/playliter-api-go/domain/inputs/band"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"github.com/mazurco066/playliter-api-go/domain/models/band"
	accountoutputs "github.com/mazurco066/playliter-api-go/domain/outputs/account"
	bandoutputs "github.com/mazurco066/playliter-api-go/domain/outputs/band"
	"github.com/mazurco066/playliter-api-go/presentation/helpers"
)

type BandController interface {
	Create(*gin.Context)
}

type bandController struct {
	AccountUc accountusecase.AccountUseCase
	BandUC    bandusecase.BandUseCase
}

func NewBandController(
	accountUc accountusecase.AccountUseCase,
	bandUc bandusecase.BandUseCase,
) BandController {
	return &bandController{
		AccountUc: accountUc,
		BandUC:    bandUc,
	}
}

// @Summary Register a new band under an account owner
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/bands [post]
func (ctl *bandController) Create(c *gin.Context) {
	account := ctl.validateTokenData(c)
	if account == nil {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	var newBand bandinputs.RegisterInput
	if err := c.BindJSON(&newBand); err != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", nil)
		return
	}

	validate := validator.New()
	if validationErr := validate.Struct(newBand); validationErr != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", validationErr.Error())
		return
	}

	band := band.Band{
		Title:       newBand.Title,
		Description: newBand.Description,
		OwnerID:     account.ID,
		Owner:       *account,
	}
	if newBand.Logo != nil {
		band.Logo = newBand.Logo
	}

	if persistErr := ctl.BandUC.Create(&band); persistErr != nil {
		helpers.HTTPRes(c, http.StatusInternalServerError, "Error persisting band!", persistErr.Error())
		return
	}

	bandOutput := ctl.mapToBandOutput(&band)
	helpers.HTTPRes(c, http.StatusOK, "Band successfully created!", bandOutput)
}

/* =========== PRIVATE METHODS =========== */

func (ctl *bandController) validateTokenData(c *gin.Context) *account.Account {
	id, exists := c.Get("user_email")
	if exists == false {
		return nil
	}

	account, err := ctl.AccountUc.GetAccountByEmail(id.(string))
	if err != nil {
		return nil
	}

	return account
}

// Map band to struct aux function
func (ctl *bandController) mapToBandOutput(b *band.Band) *bandoutputs.BandOutput {
	return &bandoutputs.BandOutput{
		ID:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		Logo:        *b.Logo,
		Owner: &accountoutputs.AccountOutput{
			ID:           b.Owner.ID,
			Name:         b.Owner.Name,
			Username:     b.Owner.Username,
			Email:        b.Owner.Email,
			Avatar:       *b.Owner.Avatar,
			IsEmailValid: b.Owner.IsEmailValid,
			Role:         b.Owner.Role,
			IsActive:     b.Owner.IsActive,
		},
	}
}
