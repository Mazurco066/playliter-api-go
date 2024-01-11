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
	Invite(*gin.Context)
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
	user := ctl.validateTokenData(c)
	if user == nil {
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

	bandObj := band.Band{
		Title:       newBand.Title,
		Description: newBand.Description,
		OwnerID:     user.ID,
		Owner:       *user,
	}
	if newBand.Logo != nil {
		bandObj.Logo = newBand.Logo
	}

	if persistErr := ctl.BandUC.Create(&bandObj); persistErr != nil {
		helpers.HTTPRes(c, http.StatusInternalServerError, "Error persisting band!", persistErr.Error())
		return
	}

	bandOutput := ctl.mapToBandOutput(&bandObj)
	helpers.HTTPRes(c, http.StatusOK, "Band successfully created!", bandOutput)
}

// @Summary Get band
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/bands/:id [get]
func (ctl *bandController) Get(c *gin.Context) {
	user := ctl.validateTokenData(c)
	if user == nil {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	id, err := ctl.stringToUint(c.Param(("id")))
	if err != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	bandResult, err := ctl.BandUC.FindById(id)
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
	if bandResult.OwnerID != user.ID && !ctl.isBandMember(bandResult.Members, user.ID) {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	bandOutput := ctl.mapToBandOutput(bandResult)
	helpers.HTTPRes(c, http.StatusOK, "Band retrieved!", bandOutput)
}

// @Summary Invite account to join the band
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/bands/:id/invite/:account_id [post]
func (ctl *bandController) Invite(c *gin.Context) {
	user := ctl.validateTokenData(c)
	if user == nil {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	id, err := ctl.stringToUint(c.Param(("id")))
	if err != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	invitedId, err := ctl.stringToUint(c.Param(("account_id")))
	if err != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	bandResult, err := ctl.BandUC.FindById(id)
	if err != nil {
		es := err.Error()
		if strings.Contains(es, "not found") {
			helpers.HTTPRes(c, http.StatusNotFound, "Band not found", nil)
			return
		}
		helpers.HTTPRes(c, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	// Validate if user is a current band admin
	if bandResult.OwnerID != user.ID && !ctl.isBandAdmin(bandResult.Members, user.ID) {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	// Verify if desired user is already a band member
	if invitedId == bandResult.OwnerID || ctl.isBandMember(bandResult.Members, invitedId) {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invited account is already a band member", nil)
		return
	}

	invitedUser, err := ctl.AccountUc.GetAccountById(invitedId)
	if err != nil {
		es := err.Error()
		if strings.Contains(es, "not found") {
			helpers.HTTPRes(c, http.StatusNotFound, "Account not found", nil)
			return
		}
		helpers.HTTPRes(c, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	// Just double checking if an invite for this account already exists
	inviteExists := ctl.BandUC.InviteExists(invitedUser, bandResult)
	if inviteExists {
		helpers.HTTPRes(c, http.StatusBadRequest, "Account was already invited. Please wait for a response from given account.", nil)
		return
	}

	bandRequestObj := band.BandRequest{
		BandID:    bandResult.ID,
		Band:      *bandResult,
		InvitedID: invitedUser.ID,
		Invited:   *invitedUser,
		Status:    "pending",
	}

	if persistErr := ctl.BandUC.CreateInvite(&bandRequestObj); persistErr != nil {
		helpers.HTTPRes(c, http.StatusInternalServerError, "Error persisting band request!", persistErr.Error())
		return
	}

	requestOutput := ctl.mapToBandRequestOutput(&bandRequestObj)
	helpers.HTTPRes(c, http.StatusOK, "Account successfully invited", requestOutput)
}

// @Summary List account bands
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/bands [get]
func (ctl *bandController) List(c *gin.Context) {
	user := ctl.validateTokenData(c)
	if user == nil {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	var paging commoninputs.PagingParams
	if err := c.BindQuery(&paging); err != nil {
		paging.Limit = 100
		paging.Offset = 0
	}

	results, err := ctl.BandUC.FindByAccount(user, &paging)
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

	user, err := ctl.AccountUc.GetAccountByEmail(id.(string))
	if err != nil {
		return nil
	}

	return user
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

func (ctl *bandController) mapToBandRequestOutput(b *band.BandRequest) *bandoutputs.BandRequestOutput {
	return &bandoutputs.BandRequestOutput{
		ID: b.ID,
		Band: &bandoutputs.BandOutput{
			ID:          b.Band.ID,
			Logo:        *b.Band.Logo,
			Title:       b.Band.Title,
			Description: b.Band.Description,
			Owner: &accountoutputs.AccountOutput{
				ID:           b.Band.Owner.ID,
				Name:         b.Band.Owner.Name,
				Username:     b.Band.Owner.Username,
				Email:        b.Band.Owner.Email,
				Avatar:       *b.Band.Owner.Avatar,
				IsEmailValid: b.Band.Owner.IsEmailValid,
				Role:         b.Band.Owner.Role,
				IsActive:     b.Band.Owner.IsActive,
			},
		},
		Invited: &accountoutputs.AccountOutput{
			ID:           b.Invited.ID,
			Name:         b.Invited.Name,
			Username:     b.Invited.Username,
			Email:        b.Invited.Email,
			Avatar:       *b.Invited.Avatar,
			IsEmailValid: b.Invited.IsEmailValid,
			Role:         b.Invited.Role,
			IsActive:     b.Invited.IsActive,
		},
		Status: b.Status,
	}
}
