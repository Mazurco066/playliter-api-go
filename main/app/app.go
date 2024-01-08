package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mazurco066/playliter-api-go/domain/account"
	"github.com/mazurco066/playliter-api-go/domain/auth"
	"github.com/mazurco066/playliter-api-go/domain/band"
	"github.com/mazurco066/playliter-api-go/domain/concert"
	"github.com/mazurco066/playliter-api-go/domain/song"
	"github.com/mazurco066/playliter-api-go/main/config"
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
	// Loads .env file
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	// App enviroment values
	configs := config.GetConfig()

	// Postgres connection
	dci := config.GetPostgresConfig().GetPostgresConnectionInfo()
	db, dbErr := gorm.Open(
		postgres.Open(dci),
		&gorm.Config{},
	)
	if dbErr != nil {
		panic(dbErr)
	}

	// Db auto migrate according to schemas
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

	// App routes
	router.GET("/", HandleRoot)

	// Starting Http server
	host := fmt.Sprintf("%s:%s", configs.Host, configs.Port)
	router.Run(host)
}
