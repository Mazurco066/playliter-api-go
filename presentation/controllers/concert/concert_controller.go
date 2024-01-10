package concertcontroller

import (
	"github.com/gin-gonic/gin"
	concertusecase "github.com/mazurco066/playliter-api-go/data/usecases/concert"
)

type ConcertController interface {
	Create(*gin.Context)
}

type concertController struct {
	ConcertUC concertusecase.ConcertUseCase
}

func NewConcertController(
	concertUc concertusecase.ConcertUseCase,
) ConcertController {
	return &concertController{
		ConcertUC: concertUc,
	}
}

func (ctl *concertController) Create(c *gin.Context) {

}
