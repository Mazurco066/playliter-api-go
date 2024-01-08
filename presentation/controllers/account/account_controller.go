package accountcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	accountusecase "github.com/mazurco066/playliter-api-go/data/usecases/account"
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
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Registrado com sucesso"})
}
