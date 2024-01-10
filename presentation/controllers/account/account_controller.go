package accountcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	accountusecase "github.com/mazurco066/playliter-api-go/data/usecases/account"
	authusecase "github.com/mazurco066/playliter-api-go/data/usecases/auth"
	accountinputs "github.com/mazurco066/playliter-api-go/domain/inputs/account"
	commoninputs "github.com/mazurco066/playliter-api-go/domain/inputs/common"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	accountoutputs "github.com/mazurco066/playliter-api-go/domain/outputs/account"
	"github.com/mazurco066/playliter-api-go/presentation/helpers"
)

type AccountController interface {
	CurrentAccount(*gin.Context)
	ListActiveAccounts(*gin.Context)
	Login(*gin.Context)
	Register(*gin.Context)
}

type accountController struct {
	AccountUC accountusecase.AccountUseCase
	AuthUc    authusecase.AuthUseCase
}

func NewAccaccountController(
	accountUC accountusecase.AccountUseCase,
	authUc authusecase.AuthUseCase,
) AccountController {
	return &accountController{
		AccountUC: accountUC,
		AuthUc:    authUc,
	}
}

// @Summary Returns current account
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/accounts/me [get]
func (ctl *accountController) CurrentAccount(c *gin.Context) {
	account := ctl.validateTokenData(c)
	if account == nil {
		helpers.HTTPRes(c, http.StatusForbidden, "Forbidden", nil)
		return
	}

	userOutput := ctl.mapToUserOutput(account)
	helpers.HTTPRes(c, http.StatusOK, "Authenticated account", userOutput)
}

// @Summary List active user accounts
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/accounts/active_users [get]
func (ctl *accountController) ListActiveAccounts(c *gin.Context) {
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
	if paging.Limit == 0 {
		paging.Limit = 100
	}

	results, err := ctl.AccountUC.ListActiveAccounts(account, &paging)
	if err != nil {
		helpers.HTTPRes(c, http.StatusInternalServerError, "Internal server error", nil)
		return
	}

	var resultOutput []*accountoutputs.AccountPublicOutput
	for _, a := range results {
		output := ctl.mapToUserPublicOutput(a)
		resultOutput = append(resultOutput, output)
	}

	// Empty array if no results
	if resultOutput == nil {
		helpers.HTTPRes(c, http.StatusOK, "Active user accounts", []string{})
		return
	}

	// Formatted array
	helpers.HTTPRes(c, http.StatusOK, "Active user accounts", resultOutput)
}

// @Summary Login into your application account
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/login [post]
func (ctl *accountController) Login(c *gin.Context) {
	var loginInput accountinputs.LoginInput
	if err := c.BindJSON(&loginInput); err != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", nil)
		return
	}

	validate := validator.New()
	if validationErr := validate.Struct(loginInput); validationErr != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", validationErr.Error())
		return
	}

	account, err := ctl.AccountUC.GetAccountByUsernameOrEmail(loginInput.UsernameOrEmail)
	if err != nil {
		helpers.HTTPRes(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	err = ctl.AccountUC.ComparePassword(loginInput.Password, account.Password)
	if err != nil {
		helpers.HTTPRes(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	err = ctl.login(c, account)
	if err != nil {
		helpers.HTTPRes(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
}

// @Summary Register a new user account
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/register [post]
func (ctl *accountController) Register(c *gin.Context) {
	var newAccount accountinputs.RegisterInput
	if err := c.BindJSON(&newAccount); err != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", nil)
		return
	}

	validate := validator.New()
	if validationErr := validate.Struct(newAccount); validationErr != nil {
		helpers.HTTPRes(c, http.StatusBadRequest, "Invalid Payload", validationErr.Error())
		return
	}

	account := account.Account{
		Email:    newAccount.Email,
		Username: newAccount.Username,
		Name:     newAccount.Name,
		Password: newAccount.Password,
	}

	if persistErr := ctl.AccountUC.Create(&account); persistErr != nil {
		helpers.HTTPRes(c, http.StatusInternalServerError, "Error persisting account!", persistErr.Error())
		return
	}

	err := ctl.login(c, &account)
	if err != nil {
		helpers.HTTPRes(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
}

/* =========== PRIVATE METHODS =========== */

func (ctl *accountController) validateTokenData(c *gin.Context) *account.Account {
	id, exists := c.Get("user_email")
	if exists == false {
		return nil
	}

	account, err := ctl.AccountUC.GetAccountByEmail(id.(string))
	if err != nil {
		return nil
	}

	return account
}

// Map user struct aux function
func (ctl *accountController) mapToUserOutput(a *account.Account) *accountoutputs.AccountOutput {
	return &accountoutputs.AccountOutput{
		ID:           a.ID,
		Email:        a.Email,
		Name:         a.Name,
		Username:     a.Username,
		Avatar:       *a.Avatar,
		Role:         a.Role,
		IsEmailValid: a.IsEmailValid,
		IsActive:     a.IsActive,
	}
}

func (ctl *accountController) mapToUserPublicOutput(a *account.Account) *accountoutputs.AccountPublicOutput {
	return &accountoutputs.AccountPublicOutput{
		ID:     a.ID,
		Name:   a.Name,
		Avatar: *a.Avatar,
	}
}

// Login aux function
func (ctl *accountController) login(c *gin.Context, a *account.Account) error {
	token, err := ctl.AuthUc.IssueToken(*a)
	if err != nil {
		return err
	}
	userOutput := ctl.mapToUserOutput(a)
	out := gin.H{"token": token, "user": userOutput}
	helpers.HTTPRes(c, http.StatusOK, "Successfully logged in!", out)
	return nil
}
