package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	env "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"modernc.org/sqlite"

	"github.com/taythebot/discord_chatbot"
	"github.com/taythebot/discord_chatbot/command"
	"github.com/taythebot/discord_chatbot/ent"
	"github.com/taythebot/discord_chatbot/handler"
	"github.com/taythebot/discord_chatbot/provider"
)

type Config struct {
	Debug        bool     `env:"DEBUG"`
	DiscordToken string   `env:"DISCORD_TOKEN,notEmpty"`
	DiscordAppID string   `env:"DISCORD_APP_ID,notEmpty"`
	AdminIDs     []string `env:"ADMIN_IDS" envSeparator:","`
	NanoGptToken string   `env:"NANOGPT_TOKEN,notEmpty"`
	ModelsFile   string   `env:"MODELS_FILE,notEmpty" envDefault:"models.yaml"`
}

func main() {
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	// Parse .env file if exists
	if _, err := os.Stat(".env"); !errors.Is(err, os.ErrNotExist) {
		if err := godotenv.Load(); err != nil {
			log.Fatal().Err(err).Msg("Failed to load .env file")
		}
	}

	// Parse env variables into struct
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse env variables")
	}

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Open sqlite3 database
	log.Debug().Msg("Opening sqlite3 database")

	sql.Register("sqlite3", &sqlite.Driver{})
	db, err := ent.Open("sqlite3", "file:database.sqlite3?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed opening connection to sqlite")
	}
	defer db.Close()

	// Run database auto migration
	log.Debug().Msg("Running database auto migration")

	if err := db.Schema.Create(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Create provider client
	providerClient := provider.NewClient(cfg.NanoGptToken)

	// Parse models
	models, err := discord_chatbot.ParseModels(context.TODO(), cfg.ModelsFile, providerClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse models file")
	}

	// Connect to Discord
	log.Debug().Msg("Connecting to Discord")

	s, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Discord")
	}
	defer s.Close()

	// Register all handlers
	handlerRegistry := handler.NewRegistry(db, providerClient)
	handlerRegistry.RegisterAll(s)

	// Register commands
	commandRegistry := command.NewRegistry(db, providerClient, cfg.AdminIDs, models)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create command registry")
	}

	// Register commands
	log.Debug().Msg("Registering commands")

	if err := commandRegistry.RegisterAll(s, cfg.DiscordAppID); err != nil {
		log.Fatal().Err(err).Msg("Failed to register commands")
	}

	fmt.Printf("\n\nBOT INSTALL LINK: https://discord.com/oauth2/authorize?client_id=%s&permissions=3088&integration_type=0&scope=bot\n\n\n", cfg.DiscordAppID)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
