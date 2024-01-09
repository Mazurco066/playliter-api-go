package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	accountrepo "github.com/mazurco066/playliter-api-go/data/repositories/account"
	accountusecase "github.com/mazurco066/playliter-api-go/data/usecases/account"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"github.com/mazurco066/playliter-api-go/domain/models/auth"
	"github.com/mazurco066/playliter-api-go/domain/models/band"
	"github.com/mazurco066/playliter-api-go/domain/models/concert"
	"github.com/mazurco066/playliter-api-go/domain/models/song"
	"github.com/mazurco066/playliter-api-go/infra/hmachash"
	"github.com/mazurco066/playliter-api-go/main/config"
	accountcontroller "github.com/mazurco066/playliter-api-go/presentation/controllers/account"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	router = gin.Default()
)

func HandleRoot(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Playliter api go version 1.0.0"})
}

func Run() {
	/* ========= Loads .env file ========= */
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	/* ========= Enviroment values ========= */
	configs := config.GetConfig()

	/* ========= Postgres connection ========= */
	dci := config.GetPostgresConfig().GetPostgresConnectionInfo()
	db, dbErr := gorm.Open(
		postgres.Open(dci),
		&gorm.Config{},
	)
	if dbErr != nil {
		panic(dbErr)
	}

	/* ========= Database auto migration (schemas) ========= */
	db.AutoMigrate(
		&account.Account{},
		&account.EmailVerification{},
		&auth.Auth{},
		&band.Band{},
		&band.BandRequest{},
		&band.Member{},
		&concert.Concert{},
		&concert.ConcertSong{},
		&song.Song{},
	)

	/* ========= Setup common ========= */
	hm := hmachash.NewHMAC(configs.HMACKey)

	/* ========= Setup infra ========= */

	/* ========= Setup repositories ========= */
	accountRepo := accountrepo.NewAccountRepo(db)

	/* ========= Setup usecases ========= */
	accountService := accountusecase.NewAccountUseCase(accountRepo, hm)

	/* ========= Setup controllers ========= */
	accountController := accountcontroller.NewAccaccountController(accountService)

	/* ========= Setup middlewares ========= */
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	/* ========= App routes ========= */
	router.GET("/", HandleRoot)
	router.POST("/login", accountController.Login)
	router.POST("/register", accountController.Register)

	/* ========= Server start ========= */
	host := fmt.Sprintf("%s:%s", configs.Host, configs.Port)
	router.Run(host)
}
