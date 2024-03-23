package app

//nolint:depguard
import (
	"context"
	"errors"
	"testing"

	mocklimiter "github.com/Baraulia/anti_bruteforce_service/internal/app/mocks/limiter"
	mockstorage "github.com/Baraulia/anti_bruteforce_service/internal/app/mocks/storage"
	"github.com/Baraulia/anti_bruteforce_service/internal/models"
	"github.com/Baraulia/anti_bruteforce_service/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	loginLimit    = 10
	passwordLimit = 100
	ipLimit       = 1000
)
var ctx = context.Background()

func TestCheck(t *testing.T) {
	type mockBehavior func(storage *mockstorage.MockStorage, limiter *mocklimiter.MockLimiter)
	testTable := []struct {
		name           string
		mockBehavior   mockBehavior
		inputData      models.Data
		expectedResult bool
		expectedError  bool
	}{
		{
			name:      "IP in white list",
			inputData: models.Data{Login: "Test", Password: "Test", IP: "192.1.1.0/25"},
			mockBehavior: func(storage *mockstorage.MockStorage, _ *mocklimiter.MockLimiter) {
				storage.EXPECT().CheckIPInWhiteList(gomock.Any(), "192.1.1.0/25").Return(true, nil)
			},
			expectedResult: true,
			expectedError:  false,
		},
		{
			name: "IP in black list and not in white list",
			inputData: models.Data{
				Login:    "Test",
				Password: "Test",
				IP:       "192.1.1.0/25",
			},
			mockBehavior: func(storage *mockstorage.MockStorage, _ *mocklimiter.MockLimiter) {
				storage.EXPECT().CheckIPInWhiteList(gomock.Any(), "192.1.1.0/25").Return(false, nil)
				storage.EXPECT().CheckIPInBlackList(gomock.Any(), "192.1.1.0/25").Return(true, nil)
			},
			expectedResult: false,
			expectedError:  false,
		},
		{
			name:      "all limits in range",
			inputData: models.Data{Login: "Test", Password: "Test", IP: "192.1.1.0/25"},
			mockBehavior: func(storage *mockstorage.MockStorage, limiter *mocklimiter.MockLimiter) {
				storage.EXPECT().CheckIPInWhiteList(gomock.Any(), "192.1.1.0/25").Return(false, nil)
				storage.EXPECT().CheckIPInBlackList(gomock.Any(), "192.1.1.0/25").Return(false, nil)
				limiter.EXPECT().CheckLimit(gomock.Any(), "192.1.1.0/25").Return(999, nil)
				limiter.EXPECT().CheckLimit(gomock.Any(), "Test").Return(9, nil)
				limiter.EXPECT().CheckLimit(gomock.Any(), "Test").Return(99, nil)
			},
			expectedResult: true,
			expectedError:  false,
		},
		{
			name:      "ip limit exceed",
			inputData: models.Data{Login: "Test", Password: "Test", IP: "192.1.1.0/25"},
			mockBehavior: func(storage *mockstorage.MockStorage, limiter *mocklimiter.MockLimiter) {
				storage.EXPECT().CheckIPInWhiteList(gomock.Any(), "192.1.1.0/25").Return(false, nil)
				storage.EXPECT().CheckIPInBlackList(gomock.Any(), "192.1.1.0/25").Return(false, nil)
				limiter.EXPECT().CheckLimit(gomock.Any(), "192.1.1.0/25").Return(1001, nil)
			},
			expectedResult: false,
			expectedError:  false,
		},
		{
			name:      "login limit exceed",
			inputData: models.Data{Login: "Test", Password: "Test", IP: "192.1.1.0/25"},
			mockBehavior: func(storage *mockstorage.MockStorage, limiter *mocklimiter.MockLimiter) {
				storage.EXPECT().CheckIPInWhiteList(gomock.Any(), "192.1.1.0/25").Return(false, nil)
				storage.EXPECT().CheckIPInBlackList(gomock.Any(), "192.1.1.0/25").Return(false, nil)
				limiter.EXPECT().CheckLimit(gomock.Any(), "192.1.1.0/25").Return(999, nil)
				limiter.EXPECT().CheckLimit(gomock.Any(), "Test").Return(11, nil)
			},
			expectedResult: false,
			expectedError:  false,
		},
		{
			name:      "password limit exceed",
			inputData: models.Data{Login: "Test", Password: "Test", IP: "192.1.1.0/25"},
			mockBehavior: func(storage *mockstorage.MockStorage, limiter *mocklimiter.MockLimiter) {
				storage.EXPECT().CheckIPInWhiteList(gomock.Any(), "192.1.1.0/25").Return(false, nil)
				storage.EXPECT().CheckIPInBlackList(gomock.Any(), "192.1.1.0/25").Return(false, nil)
				limiter.EXPECT().CheckLimit(gomock.Any(), "192.1.1.0/25").Return(999, nil)
				limiter.EXPECT().CheckLimit(gomock.Any(), "Test").Return(9, nil)
				limiter.EXPECT().CheckLimit(gomock.Any(), "Test").Return(101, nil)
			},
			expectedResult: false,
			expectedError:  false,
		},
		{
			name:      "internal error from storage",
			inputData: models.Data{Login: "Test", Password: "Test", IP: "192.1.1.0/25"},
			mockBehavior: func(storage *mockstorage.MockStorage, _ *mocklimiter.MockLimiter) {
				storage.EXPECT().CheckIPInWhiteList(gomock.Any(), "192.1.1.0/25").Return(
					false, errors.New("internal error"),
				)
			},
			expectedResult: false,
			expectedError:  true,
		},
		{
			name:      "internal error from limiter",
			inputData: models.Data{Login: "Test", Password: "Test", IP: "192.1.1.0/25"},
			mockBehavior: func(storage *mockstorage.MockStorage, limiter *mocklimiter.MockLimiter) {
				storage.EXPECT().CheckIPInWhiteList(gomock.Any(), "192.1.1.0/25").Return(false, nil)
				storage.EXPECT().CheckIPInBlackList(gomock.Any(), "192.1.1.0/25").Return(false, nil)
				limiter.EXPECT().CheckLimit(gomock.Any(), "192.1.1.0/25").Return(999, errors.New("internal error"))
			},
			expectedResult: false,
			expectedError:  true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			storageMock := mockstorage.NewMockStorage(c)
			limiterMock := mocklimiter.NewMockLimiter(c)

			logg, err := logger.GetLogger("INFO", false)
			require.NoError(t, err)

			app := New(logg, storageMock, limiterMock, loginLimit, passwordLimit, ipLimit)
			testCase.mockBehavior(storageMock, limiterMock)

			result, err := app.Check(ctx, testCase.inputData)

			if testCase.expectedError {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}
