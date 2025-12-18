package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// Load database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Get command from arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Open database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migrate driver: %v", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Execute command
	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("✅ Migrations completed successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		log.Println("✅ Rollback completed successfully")

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Version number required for force command")
		}
		var version int
		fmt.Sscanf(os.Args[2], "%d", &version)
		if err := m.Force(version); err != nil {
			log.Fatalf("Failed to force version: %v", err)
		}
		log.Printf("✅ Forced version to %d\n", version)

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		if dirty {
			log.Printf("Current version: %d (dirty)\n", version)
		} else {
			log.Printf("Current version: %d\n", version)
		}

	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Migration name required for create command")
		}
		name := os.Args[2]
		log.Printf("Creating migration: %s\n", name)
		log.Println("Run: migrate create -ext sql -dir migrations -seq " + name)

	case "status":
		version, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				log.Println("No migrations have been applied yet")
			} else {
				log.Fatalf("Failed to get migration status: %v", err)
			}
			return
		}

		status := "clean"
		if dirty {
			status = "dirty (migration failed, needs manual intervention)"
		}

		fmt.Println("📊 Migration Status:")
		fmt.Printf("   Current Version: %d\n", version)
		fmt.Printf("   Status: %s\n", status)

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Database Migration Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run src/cmd/migrate/main.go <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up           Apply all pending migrations")
	fmt.Println("  down         Rollback the last migration")
	fmt.Println("  force <n>    Force set migration version (use with caution)")
	fmt.Println("  version      Show current migration version")
	fmt.Println("  status       Show detailed migration status")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  DATABASE_URL  PostgreSQL connection string (required)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run src/cmd/migrate/main.go up")
	fmt.Println("  go run src/cmd/migrate/main.go down")
	fmt.Println("  go run src/cmd/migrate/main.go version")
	fmt.Println("  go run src/cmd/migrate/main.go status")
	fmt.Println()
	fmt.Println("Creating new migrations:")
	fmt.Println("  migrate create -ext sql -dir migrations -seq <migration_name>")
}
