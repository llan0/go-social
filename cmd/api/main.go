package main

import (
	"time"

	"github.com/llan0/go-social/internal/db"
	"github.com/llan0/go-social/internal/env" // can use pkgs like godotenv (this is my own implimentation)
	"github.com/llan0/go-social/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			go-social
//	@description	WIP - Personal Blogging Platform
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONN", 30),
			maxIdelConns: env.GetInt("DB_MAX_IDEL_CONN", 30),
			maxIdelTime:  env.GetDuration("DB_MAX_IDEL_TIME", 15*time.Minute),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp: time.Hour * 24 * 3, //user invitation expires in 3 days
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdelConns, cfg.db.maxIdelTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Info("db connection pool established!")

	// Storage
	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))
}
