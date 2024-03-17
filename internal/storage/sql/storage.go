package sqlstorage

//nolint:depguard
import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Baraulia/anti_bruteforce_service/internal/app"
	// Empty import to ensure execution of code in the package's init function.
	_ "github.com/lib/pq"
)

const (
	MaxConnections = 10
	WhiteListTable = "white_list"
	BlackListTable = "black_list"
)

type PostgresStorage struct {
	dataSource string
	db         *sql.DB
	logger     app.Logger
}

type PgConfig struct {
	Host     string
	Username string
	Password string
	Port     string
	Database string
}

func NewPostgresStorage(conf PgConfig, logger app.Logger, migrate bool) *PostgresStorage {
	dataSource := fmt.Sprintf("host= %s port= %s user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.Username, conf.Password, conf.Database)
	storage := &PostgresStorage{dataSource: dataSource, logger: logger}

	db, err := sql.Open("postgres", storage.dataSource)
	if err != nil {
		storage.logger.Fatal("Unable to connect to database", map[string]interface{}{"error": err, "dataSource": storage.dataSource})
	}

	err = db.Ping()
	if err != nil {
		storage.logger.Fatal("error while pinging", map[string]interface{}{"error": err})
	}

	db.SetMaxOpenConns(MaxConnections)

	if migrate {
		whiteListQuery := fmt.Sprintf(`
		CREATE TABLE %s IF NOT EXISTS (
			id serial PRIMARY KEY,
			ip varchar(255) NOT NULL UNIQUE,
			CREATED_AT timestamp DEFAULT now()
		);`, WhiteListTable)

		blackListQuery := fmt.Sprintf(`
		CREATE TABLE %s IF NOT EXISTS (
			id serial PRIMARY KEY,
			ip varchar(255) NOT NULL UNIQUE,
			CREATED_AT timestamp DEFAULT now()
		);`, BlackListTable)

		_, err = db.Exec(whiteListQuery)
		if err != nil {
			storage.logger.Error("error while creating white list table", map[string]interface{}{"error": err})
		}

		_, err = db.Exec(blackListQuery)
		if err != nil {
			storage.logger.Error("error while creating black list table", map[string]interface{}{"error": err})
		}
	}

	storage.db = db

	return storage
}

func (s *PostgresStorage) Close() {
	s.db.Close()
}

func (s *PostgresStorage) AddToWhiteList(ctx context.Context, ip string) error {
	sqlString := fmt.Sprintf(
		"INSERT INTO %s (ip) VALUES($1) ON CONFLICT (ip) DO NOTHING", WhiteListTable)
	_, err := s.db.ExecContext(ctx, sqlString, ip)
	if err != nil {
		s.logger.Error("error while adding to white list", map[string]interface{}{"error": err})
		return fmt.Errorf("error while adding to white list: %w", err)
	}

	return nil
}

func (s *PostgresStorage) AddToBlackList(ctx context.Context, ip string) error {
	sqlString := fmt.Sprintf(
		"INSERT INTO %s (ip) VALUES($1) ON CONFLICT (ip) DO NOTHING", BlackListTable)
	_, err := s.db.ExecContext(ctx, sqlString, ip)
	if err != nil {
		s.logger.Error("error while adding to black list", map[string]interface{}{"error": err})
		return fmt.Errorf("error while adding to black list: %w", err)
	}

	return nil
}

func (s *PostgresStorage) RemoveFromWhiteList(ctx context.Context, ip string) error {
	sqlString := fmt.Sprintf(
		"DELETE FROM %s WHERE ip = $1", WhiteListTable)
	_, err := s.db.ExecContext(ctx, sqlString, ip)
	if err != nil {
		s.logger.Error("error while deleting from white list", map[string]interface{}{"error": err, "ip": ip})
		return fmt.Errorf("error while deleting from white list: %w", err)
	}

	return nil
}

func (s *PostgresStorage) RemoveFromBlackList(ctx context.Context, ip string) error {
	sqlString := fmt.Sprintf(
		"DELETE FROM %s WHERE ip = $1", BlackListTable)
	_, err := s.db.ExecContext(ctx, sqlString, ip)
	if err != nil {
		s.logger.Error("error while deleting from black list", map[string]interface{}{"error": err, "ip": ip})
		return fmt.Errorf("error while deleting from black list: %w", err)
	}

	return nil
}

func (s *PostgresStorage) CheckIPInWhiteList(ctx context.Context, ip string) (bool, error) {
	var exists bool
	sqlString := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE ip = $1)", WhiteListTable)

	err := s.db.QueryRowContext(ctx, sqlString, ip).Scan(&exists)
	if err != nil {
		s.logger.Error("error while checking in white list", map[string]interface{}{"error": err, "ip": ip})
		return false, fmt.Errorf("error while checking in white list: %w", err)
	}

	return exists, nil
}

func (s *PostgresStorage) CheckIPInBlackList(ctx context.Context, ip string) (bool, error) {
	var exists bool
	sqlString := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE ip = $1)", BlackListTable)

	err := s.db.QueryRowContext(ctx, sqlString, ip).Scan(&exists)
	if err != nil {
		s.logger.Error("error while checking in black list", map[string]interface{}{"error": err, "ip": ip})
		return false, fmt.Errorf("error while checking in black list: %w", err)
	}

	return exists, nil
}

func (s *PostgresStorage) CheckReadness() (bool, error) {
	err := s.db.Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}
