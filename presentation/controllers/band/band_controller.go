package bandcontroller

import (
	"github.com/gin-gonic/gin"
	bandusecase "github.com/mazurco066/playliter-api-go/data/usecases/band"
)

type BandController interface {
	Create(*gin.Context)
}

type bandController struct {
	BandUC bandusecase.BandUseCase
}

func NewBandController(
	bandUc bandusecase.BandUseCase,
) BandController {
	return &bandController{
		BandUC: bandUc,
	}
}

func (ctl *bandController) Create(c *gin.Context) {

}
