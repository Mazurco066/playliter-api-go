package accountcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	accountusecase "github.com/mazurco066/playliter-api-go/data/usecases/account"
	accountinputs "github.com/mazurco066/playliter-api-go/domain/inputs/account"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"github.com/mazurco066/playliter-api-go/presentation/helpers"
)

type AccountController interface {
	Login(*gin.Context)
	Register(*gin.Context)
}

type accountController struct {
	AccountUC accountusecase.AccountUseCase
}

func NewAccaccountController(
	accountUC accountusecase.AccountUseCase,
) AccountController {
	return &accountController{
		AccountUC: accountUC,
	}
}

func (ctl *accountController) Login(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "login com sucesso"})
}

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

	helpers.HTTPRes(c, http.StatusCreated, "Account successfully created!", account)
	return
}