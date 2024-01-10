package songcontroller

import (
	"github.com/gin-gonic/gin"
	songusecase "github.com/mazurco066/playliter-api-go/data/usecases/song"
)

type SongController interface {
	Create(*gin.Context)
}

type songController struct {
	SongUC songusecase.SongUseCase
}

func NewSongController(
	songUc songusecase.SongUseCase,
) SongController {
	return &songController{
		SongUC: songUc,
	}
}

func (ctl *songController) Create(c *gin.Context) {

}
