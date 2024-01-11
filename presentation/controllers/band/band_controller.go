package bandcontroller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	accountusecase "github.com/mazurco066/playliter-api-go/data/usecases/account"
	bandusecase "github.com/mazurco066/playliter-api-go/data/usecases/band"
	bandinputs "github.com/mazurco066/playliter-api-go/domain/inputs/band"
	commoninputs "github.com/mazurco066/playliter-api-go/domain/inputs/common"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"github.com/mazurco066/playliter-api-go/domain/models/band"
	accountoutputs "github.com/mazurco066/playliter-api-go/domain/outputs/account"
	bandoutputs "github.com/mazurco066/playliter-api-go/domain/outputs/band"
	"github.com/mazurco066/playliter-api-go/presentation/helpers"
)

type BandController interface {
	Create(*gin.Context)
	Get(*gin.Context)
	List(*gin.Context)
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

// @Summary Get band
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/bands/:id [get]
func (ctl *bandController) Get(c *gin.Context) {
	account := ctl.validateTokenData(c)
	if account == nil {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	id, err := ctl.stringToUint(c.Param(("id")))
	if err != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	band, err := ctl.BandUC.FindById(id)
	if err != nil {
		es := err.Error()
		if strings.Contains(es, "not found") {
			helpers.HTTPRes(c, http.StatusNotFound, "Band not found", nil)
			return
		}
		helpers.HTTPRes(c, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	// Validate if user is a current band member
	if band.OwnerID != account.ID && !ctl.isBandMember(band.Members, account.ID) {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	bandOutput := ctl.mapToBandOutput(band)
	helpers.HTTPRes(c, http.StatusOK, "Band retrieved!", bandOutput)
}

// @Summary List account bands
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/bands [get]
func (ctl *bandController) List(c *gin.Context) {
	account := ctl.validateTokenData(c)
	if account == nil {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	var paging commoninputs.PagingParams
	if err := c.BindQuery(&paging); err != nil {
		paging.Limit = 100
		paging.Offset = 0
	}

	results, err := ctl.BandUC.FindByAccount(account, &paging)
	if err != nil {
		helpers.HTTPRes(c, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	var resultOutput []*bandoutputs.BandOutput
	for _, b := range results {
		output := ctl.mapToBandOutput(b)
		resultOutput = append(resultOutput, output)
	}

	// Empty array if no results
	if resultOutput == nil {
		helpers.HTTPRes(c, http.StatusOK, "Bands successfully listed!", []string{})
		return
	}

	// Formatted array
	helpers.HTTPRes(c, http.StatusOK, "Bands successfully listed!", resultOutput)
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

func (ctl *bandController) isBandMember(members []band.Member, accountID uint) bool {
	for _, member := range members {
		if member.AccountID == accountID {
			return true
		}
	}
	return false
}

func (ctl *bandController) isBandAdmin(members []band.Member, accountID uint) bool {
	for _, member := range members {
		if member.AccountID == accountID && member.Role == "admin" {
			return true
		}
	}
	return false
}

func (ctl *bandController) stringToUint(IDParam string) (uint, error) {
	userID, err := strconv.Atoi(IDParam)
	if err != nil {
		return 0, errors.New("id should be a number")
	}
	return uint(userID), nil
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
