package main

import (
	"context"
	"log"

	"github.com/findardi/Wadi/server/internal/app"
	"github.com/findardi/Wadi/server/internal/platform/config"
	"github.com/findardi/Wadi/server/internal/platform/database"
	"github.com/findardi/Wadi/server/internal/platform/storage"
)

func main() {
	if err := config.LoadEnvFile("configs/.env"); err != nil {
		log.Fatal(err)
	}

	dbCfg, err := config.LoadDatabaseConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.New(context.Background(), dbCfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	minioCfg := config.LoadMinioConfig()
	store, err := storage.NewMinio(minioCfg)
	if err != nil {
		log.Fatal(err)
	}

	otpSecret := config.GetEnv("OTP_SECRET", "")
	jwtSecret := config.GetEnv("JWT_SECRET", "")
	addr := config.GetEnv("ADDR", ":8181")

	if otpSecret == "" || jwtSecret == "" {
		log.Fatal("OTP_SECRET and JWT_SECRET must be set")
	}

	if err := app.New(db, otpSecret, addr, jwtSecret, store).Run(); err != nil {
		log.Fatal(err)
	}
}
