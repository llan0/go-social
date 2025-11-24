package main

import (
	"time"

	"github.com/llan0/go-social/internal/auth"
	"github.com/llan0/go-social/internal/db"
	"github.com/llan0/go-social/internal/env"
	"github.com/llan0/go-social/internal/mailer"
	"github.com/llan0/go-social/internal/store"
	"github.com/llan0/go-social/internal/store/cache"
	"github.com/redis/go-redis/v9"
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
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONN", 30),
			maxIdelConns: env.GetInt("DB_MAX_IDEL_CONN", 30),
			maxIdelTime:  env.GetDuration("DB_MAX_IDEL_TIME", 15*time.Minute),
		},
		redisCfg: redisConfig{
			addr:   env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:     env.GetString("REDIS_PW", ""),
			db:     env.GetInt("REDIS_DB", 0),
			enable: env.GetBool("REDIS_ENABLED", true),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			resend: resendConfig{
				apiKey: env.GetString("RESEND_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				// TODO: dont add default creds in PRODUCTION!
				username: env.GetString("AUTH_BASIC_USERNAME", "admin"),
				password: env.GetString("AUTH_BASIC_PASSWORD", "admin"),
			},
			token: tokenConfig{
				// TODO: dont add default creds in PRODUCTION!
				secret: env.GetString("AUTH_TOKEN_SECRET", "dev"),
				exp:    time.Hour * 24 * 3, // 3 days
				iss:    "gosocial",
			},
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

	// Mailer
	resendClient := mailer.NewResendClient(cfg.mail.resend.apiKey, cfg.mail.fromEmail)

	// Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	// Redis cache
	var rdb *redis.Client
	if cfg.redisCfg.enable {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("redis cache connection established")
		defer rdb.Close()
	}

	app := &application{
		config:        cfg,
		store:         store.NewStorage(db),
		cacheStore:    cache.NewRedisStorage(rdb),
		logger:        logger,
		mailer:        resendClient,
		authenticator: jwtAuthenticator,
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))
}
