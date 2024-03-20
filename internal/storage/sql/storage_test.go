package sqlstorage

import (
	"context"
	"fmt"
	"testing"

	"github.com/Baraulia/anti_bruteforce_service/pkg/logger"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestAddToWhiteList(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logg, err := logger.GetLogger("debug", false)
	require.NoError(t, err)

	storage := &PostgresStorage{
		db:     db,
		logger: logg,
	}

	mock.ExpectExec(fmt.Sprintf(
		"INSERT INTO %s", WhiteListTable)).WithArgs("192.1.1.0/25").WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.AddToWhiteList(context.Background(), "192.1.1.0/25")
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestAddToBlackList(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logg, err := logger.GetLogger("debug", false)
	require.NoError(t, err)

	storage := &PostgresStorage{
		db:     db,
		logger: logg,
	}

	mock.ExpectExec(fmt.Sprintf(
		"INSERT INTO %s", BlackListTable)).WithArgs("192.1.1.0/25").WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.AddToBlackList(context.Background(), "192.1.1.0/25")
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestRemoveFromWhiteList(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logg, err := logger.GetLogger("debug", false)
	require.NoError(t, err)

	storage := &PostgresStorage{
		db:     db,
		logger: logg,
	}

	mock.ExpectExec(fmt.Sprintf(""+
		"DELETE FROM %s", WhiteListTable)).WithArgs("192.1.1.0/25").WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.RemoveFromWhiteList(context.Background(), "192.1.1.0/25")
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestRemoveFromBlackList(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logg, err := logger.GetLogger("debug", false)
	require.NoError(t, err)

	storage := &PostgresStorage{
		db:     db,
		logger: logg,
	}

	mock.ExpectExec(fmt.Sprintf(
		"DELETE FROM %s", BlackListTable)).WithArgs("192.1.1.0/25").WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.RemoveFromBlackList(context.Background(), "192.1.1.0/25")
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCheckIPInWhiteList(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logg, err := logger.GetLogger("debug", false)
	require.NoError(t, err)

	storage := &PostgresStorage{
		db:     db,
		logger: logg,
	}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("SELECT EXISTS").WithArgs("192.1.1.0/25").WillReturnRows(rows)

	exists, err := storage.CheckIPInWhiteList(context.Background(), "192.1.1.0/25")
	require.NoError(t, err)
	require.Equal(t, true, exists)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCheckIPInBlackList(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logg, err := logger.GetLogger("debug", false)
	require.NoError(t, err)

	storage := &PostgresStorage{
		db:     db,
		logger: logg,
	}

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	mock.ExpectQuery("SELECT EXISTS").WithArgs("192.1.1.0/25").WillReturnRows(rows)

	exists, err := storage.CheckIPInBlackList(context.Background(), "192.1.1.0/25")
	require.NoError(t, err)
	require.Equal(t, false, exists)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
