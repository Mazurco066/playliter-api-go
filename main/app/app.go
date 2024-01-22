package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	accountrepo "github.com/mazurco066/playliter-api-go/data/repositories/account"
	bandrepo "github.com/mazurco066/playliter-api-go/data/repositories/band"
	concertrepo "github.com/mazurco066/playliter-api-go/data/repositories/concert"
	songrepo "github.com/mazurco066/playliter-api-go/data/repositories/song"
	accountusecase "github.com/mazurco066/playliter-api-go/data/usecases/account"
	authusecase "github.com/mazurco066/playliter-api-go/data/usecases/auth"
	bandusecase "github.com/mazurco066/playliter-api-go/data/usecases/band"
	concertusecase "github.com/mazurco066/playliter-api-go/data/usecases/concert"
	songusecase "github.com/mazurco066/playliter-api-go/data/usecases/song"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
	"github.com/mazurco066/playliter-api-go/domain/models/auth"
	"github.com/mazurco066/playliter-api-go/domain/models/band"
	"github.com/mazurco066/playliter-api-go/domain/models/concert"
	"github.com/mazurco066/playliter-api-go/domain/models/song"
	"github.com/mazurco066/playliter-api-go/infra/hmachash"
	"github.com/mazurco066/playliter-api-go/infra/middlewares"
	"github.com/mazurco066/playliter-api-go/main/config"
	accountcontroller "github.com/mazurco066/playliter-api-go/presentation/controllers/account"
	bandcontroller "github.com/mazurco066/playliter-api-go/presentation/controllers/band"
	concertcontroller "github.com/mazurco066/playliter-api-go/presentation/controllers/concert"
	songcontroller "github.com/mazurco066/playliter-api-go/presentation/controllers/song"
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
	bandRepo := bandrepo.NewBandRepo(db)
	bandRequestRepo := bandrepo.NewBandRequestRepo(db)
	memberRepo := bandrepo.NewMemberRepo(db)
	concertRepo := concertrepo.NewConcertRepo(db)
	songrepo := songrepo.NewSongRepo(db)

	/* ========= Setup usecases ========= */
	accountService := accountusecase.NewAccountUseCase(accountRepo, hm)
	authService := authusecase.NewAuthUseCase(configs.JWTSecret)
	bandService := bandusecase.NewBandUseCase(bandRepo)
	bandRequestService := bandusecase.NewBandRequestUseCase(bandRequestRepo)
	memberService := bandusecase.NewMemberUseCase(memberRepo)
	concertService := concertusecase.NewConcertUseCase(concertRepo)
	songService := songusecase.NewSongUseCase(songrepo)

	/* ========= Setup controllers ========= */
	accountController := accountcontroller.NewAccaccountController(accountService, authService)
	bandController := bandcontroller.NewBandController(accountService, bandService, bandRequestService, memberService)
	concertController := concertcontroller.NewConcertController(concertService)
	songController := songcontroller.NewSongController(songService)

	/* ========= Setup middlewares ========= */
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	/* ========= App default routes ========= */
	api := router.Group("/api")
	api.GET("/", HandleRoot)
	api.POST("/login", accountController.Login)
	api.POST("/register", accountController.Register)

	/* ========= App account routes ========= */
	accounts := api.Group("/accounts")
	accounts.Use(middlewares.RequiredLoggedIn(configs.JWTSecret))
	{
		accounts.GET("/me", accountController.CurrentAccount)
		accounts.PATCH("/me", accountController.Update)
		accounts.GET("/active_users", accountController.ListActiveAccounts)
	}

	/* ========= App band routes ========= */
	bands := api.Group("/bands")
	bands.Use(middlewares.RequiredLoggedIn(configs.JWTSecret))
	{
		bands.POST("/", bandController.Create)
		bands.GET("/", bandController.List)
		bands.GET("/:id", bandController.Get)
		bands.PATCH("/:id", bandController.Update)
		bands.DELETE("/:id", bandController.Remove)
		bands.PATCH("/:id/transfer/:member_id", bandController.Transfer)
		bands.POST("/:id/invite/:account_id", bandController.Invite)
		bands.PATCH("/:id/invite/:invite_id", bandController.RespondInvite)
		bands.PATCH("/:id/member/:member_id", bandController.UpdateMember)
		bands.DELETE(":id/member/:member_id", bandController.ExpelMember)
	}
	invites := api.Group("/invites")
	invites.Use(middlewares.RequiredLoggedIn(configs.JWTSecret))
	{
		invites.GET("/", bandController.PendingInvites)
	}

	/* ========= App concert routes ========= */
	concerts := api.Group("/concerts")
	concerts.Use(middlewares.RequiredLoggedIn(configs.JWTSecret))
	{
		concerts.POST("/", concertController.Create)
	}

	/* ========= App song routes ========= */
	songs := api.Group("/songs")
	songs.Use(middlewares.RequiredLoggedIn(configs.JWTSecret))
	{
		songs.POST("/", songController.Create)
	}

	/* ========= Server start ========= */
	host := fmt.Sprintf("%s:%s", configs.Host, configs.Port)
	router.Run(host)
}
