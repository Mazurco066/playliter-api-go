package bandcontroller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	ExpelMember(*gin.Context)
	Get(*gin.Context)
	Invite(*gin.Context)
	List(*gin.Context)
	Remove(*gin.Context)
	RespondInvite(*gin.Context)
	Update(*gin.Context)
	UpdateMember(*gin.Context)
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

// @Summary Delete band endpoint
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/bands/:id [delete]
func (ctl *bandController) Remove(c *gin.Context) {
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

	// Verify if user is the band owner
	if user.ID != bandResult.OwnerID {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	// TODO: Remove band concerts, songs and mambers, also delete invites
	if persistErr := ctl.BandUC.Remove(bandResult); persistErr != nil {
		helpers.HTTPRes(c, http.StatusInternalServerError, "Error deleting band!", persistErr.Error())
		return
	}

	helpers.HTTPRes(c, http.StatusNoContent, "Band successfully deleted!", nil)
}

// @Summary Respond band invite
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/bands/:id/invite/:invite_id [patch]
func (ctl *bandController) RespondInvite(c *gin.Context) {
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

	inviteId, err := ctl.stringToUint(c.Param(("invite_id")))
	if err != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	inviteResult, err := ctl.BandUC.FindInviteById(inviteId)
	if err != nil {
		es := err.Error()
		if strings.Contains(es, "not found") {
			helpers.HTTPRes(c, http.StatusNotFound, "Band invitation not found", nil)
			return
		}
		helpers.HTTPRes(c, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	if id != inviteResult.BandID || user.ID != inviteResult.InvitedID {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid band invite!", nil)
		return
	}

	var updateInput bandinputs.UpdateInviteInput
	if err := c.BindJSON(&updateInput); err != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", nil)
		return
	}

	validate := validator.New()
	if validationErr := validate.Struct(updateInput); validationErr != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", validationErr.Error())
		return
	}

	if updateInput.Status != "accepted" && updateInput.Status != "denied" {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", nil)
		return
	}

	// It means that user accepted the invitation so now he will become a member
	if updateInput.Status == "accepted" {
		memberObj := band.Member{
			BandID:    inviteResult.BandID,
			Band:      inviteResult.Band,
			AccountID: inviteResult.InvitedID,
			Account:   inviteResult.Invited,
			Role:      "member",
			JoinedAt:  time.Now(),
		}

		if persistErr := ctl.BandUC.CreateMember(&memberObj); persistErr != nil {
			helpers.HTTPRes(c, http.StatusInternalServerError, "Error persisting the new band member!", persistErr.Error())
			return
		}
	}

	inviteResult.Status = updateInput.Status
	if persistErr := ctl.BandUC.UpdateInvite(inviteResult); persistErr != nil {
		helpers.HTTPRes(c, http.StatusInternalServerError, "Error persisting the band invite data!", persistErr.Error())
		return
	}

	helpers.HTTPRes(c, http.StatusOK, "Invite successfully responded", nil)
}

// @Summary Updated band data endpoint
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/bands/:id [patch]
func (ctl *bandController) Update(c *gin.Context) {
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

	// Verify if user is a band admin
	if user.ID != bandResult.OwnerID && !ctl.isBandAdmin(bandResult.Members, user.ID) {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	var updateInput bandinputs.UpdateInput
	if err := c.BindJSON(&updateInput); err != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", nil)
		return
	}

	validate := validator.New()
	if validationErr := validate.Struct(updateInput); validationErr != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", validationErr.Error())
		return
	}

	// Update band and persist
	if updateInput.Title != "" {
		bandResult.Title = updateInput.Title
	}
	if updateInput.Description != "" {
		bandResult.Description = updateInput.Description
	}
	if updateInput.Logo != nil {
		bandResult.Logo = updateInput.Logo
	}

	if persistErr := ctl.BandUC.Update(bandResult); persistErr != nil {
		helpers.HTTPRes(c, http.StatusInternalServerError, "Error persisting band!", persistErr.Error())
		return
	}

	bandOutput := ctl.mapToBandOutput(bandResult)
	helpers.HTTPRes(c, http.StatusOK, "Band successfully updated", bandOutput)
}

// @Summary Promote or Demote band member into admin
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/bands/:id/member/:member_id [patch]
func (ctl *bandController) UpdateMember(c *gin.Context) {
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

	memberId, err := ctl.stringToUint(c.Param(("member_id")))
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

	memberResult, err := ctl.BandUC.FindMemberById(memberId)
	if err != nil {
		es := err.Error()
		if strings.Contains(es, "not found") {
			helpers.HTTPRes(c, http.StatusNotFound, "Band Member not found", nil)
			return
		}
		helpers.HTTPRes(c, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	// Verify if desired user is already a band admin
	if user.ID != bandResult.OwnerID && !ctl.isBandAdmin(bandResult.Members, user.ID) {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	var updateInput bandinputs.UpdateMemberInput
	if err := c.BindJSON(&updateInput); err != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", nil)
		return
	}

	validate := validator.New()
	if validationErr := validate.Struct(updateInput); validationErr != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", validationErr.Error())
		return
	}

	if updateInput.Role != "admin" && updateInput.Role != "member" {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", nil)
		return
	}

	memberOutput := ctl.mapToMemberOutput(memberResult)
	if memberResult.Role == updateInput.Role {
		helpers.HTTPRes(c, http.StatusOK, "No need to update this member!", memberOutput)
		return
	}

	// Updating member reference
	memberResult.Role = updateInput.Role
	if persistErr := ctl.BandUC.UpdateMember(memberResult); persistErr != nil {
		helpers.HTTPRes(c, http.StatusInternalServerError, "Error persisting band member!", persistErr.Error())
		return
	}

	memberOutput = ctl.mapToMemberOutput(memberResult)
	helpers.HTTPRes(c, http.StatusOK, "Member successfully updated!", memberOutput)
}

// @Summary Expel members from band
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/bands/:id/member/:member_id [delete]
func (ctl *bandController) ExpelMember(c *gin.Context) {
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

	memberId, err := ctl.stringToUint(c.Param(("member_id")))
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

	memberResult, err := ctl.BandUC.FindMemberById(memberId)
	if err != nil {
		es := err.Error()
		if strings.Contains(es, "not found") {
			helpers.HTTPRes(c, http.StatusNotFound, "Band Member not found", nil)
			return
		}
		helpers.HTTPRes(c, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	// Verify if desired user is already a band admin
	if user.ID != bandResult.OwnerID && !ctl.isBandAdmin(bandResult.Members, user.ID) {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	// Remove member from band
	if persistErr := ctl.BandUC.RemoveMember(memberResult); persistErr != nil {
		helpers.HTTPRes(c, http.StatusInternalServerError, "Error deleting band member!", persistErr.Error())
		return
	}

	helpers.HTTPRes(c, http.StatusNoContent, "Band member successfully expeled!", nil)
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

func (ctl *bandController) mapToMemberOutput(b *band.Member) *bandoutputs.MemberOutput {
	return &bandoutputs.MemberOutput{
		ID: b.ID,
		Band: &bandoutputs.BandOutput{
			ID:          b.Band.ID,
			Logo:        *b.Band.Logo,
			Title:       b.Band.Title,
			Description: b.Band.Description,
		},
		Account: &accountoutputs.AccountOutput{
			ID:           b.Account.ID,
			Name:         b.Account.Name,
			Username:     b.Account.Username,
			Email:        b.Account.Email,
			Avatar:       *b.Account.Avatar,
			IsEmailValid: b.Account.IsEmailValid,
			Role:         b.Account.Role,
			IsActive:     b.Account.IsActive,
		},
		Role:     b.Role,
		JoinedAt: b.JoinedAt,
	}
}
