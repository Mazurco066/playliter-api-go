package accountcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	accountusecase "github.com/mazurco066/playliter-api-go/data/usecases/account"
	authusecase "github.com/mazurco066/playliter-api-go/data/usecases/auth"
	accountinputs "github.com/mazurco066/playliter-api-go/domain/inputs/account"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	accountoutputs "github.com/mazurco066/playliter-api-go/domain/outputs/account"
	"github.com/mazurco066/playliter-api-go/presentation/helpers"
)

type AccountController interface {
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

// @Summary Login into your application account
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /login [post]
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
// @Router /register [post]
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

// Map user struct aux function
func (ctl *accountController) mapToUserOutput(a *account.Account) *accountoutputs.AccountOutput {
	return &accountoutputs.AccountOutput{
		ID:       a.ID,
		Email:    a.Email,
		Name:     a.Name,
		Role:     a.Role,
		IsActive: a.IsActive,
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
