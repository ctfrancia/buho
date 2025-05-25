package repository

import (
	"log"
	"os"

	"github.com/ctfrancia/buho/internal/core/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database holds the database connection
type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase() (*Database, error) {
	// Database connection string
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=buho_admin password=pa55word dbname=buho port=5432 sslmode=disable"
	}

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
			Colorful: true,
		},
	)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

// Migrate runs database migrations
func (d *Database) Migrate() error {
	return d.DB.AutoMigrate(
		&models.Player{},
		&models.Club{},
		&models.Tournament{},
		&models.Match{},
		&models.ClubMembership{},
		&models.TournamentRegistration{},
		&models.MatchParticipant{},
		&models.APIClient{},
	)
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
