package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/ritchieridanko/stockflow/backend/configs"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/database"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/logger"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
)

func main() {
	// Config Initialization
	cfg, err := configs.Init("./configs")
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}

	// Logger Initialization
	l, err := logger.Init(cfg.App.Env)
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}

	// Database Initialization
	db, err := database.Init(&cfg.Database, l)
	if err != nil {
		log.Fatalln("[FATAL]:", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// SQL Files Reading
	seedsDir := "./seeds"
	files, err := os.ReadDir(seedsDir)
	if err != nil {
		log.Fatalln("[FATAL]: failed to read seeds directory:", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Database Seeding Execution
	txErr := database.NewTransactor(db).WithTx(ctx, func(ctx context.Context) *ce.Error {
		for _, sf := range sqlFiles {
			content, err := os.ReadFile(filepath.Join(seedsDir, sf))
			if err != nil {
				return ce.NewError(ce.CodeInternal, ce.MsgInternalServer, err)
			}
			if _, err := db.Exec(ctx, string(content)); err != nil {
				return ce.NewError(ce.CodeInternal, ce.MsgInternalServer, err)
			}
		}
		return nil
	})
	if txErr != nil {
		log.Fatalln("[FATAL]: database seeding failed:", txErr.Unwrap())
	}
}
