package helper_test

import (
	"errors"
	"github.com/mohsenabedy91/polyglot-sentences/mocks"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetEnvFilePath(t *testing.T) {
	tests := []struct {
		name           string
		currentDir     string
		envFileName    string
		setupMocks     func(mockOS *mocks.MockOS, mockStat *mocks.MockStat)
		expectedResult string
		expectedError  bool
	}{
		{
			name:        "Env file found in current directory",
			currentDir:  "/home/user/project",
			envFileName: ".env",
			setupMocks: func(mockOS *mocks.MockOS, mockStat *mocks.MockStat) {
				mockOS.On("Getwd").Return("/home/user/project", nil)
				mockStat.On("Stat", "/home/user/project/.env").Return(nil, nil)
			},
			expectedResult: "/home/user/project/.env",
			expectedError:  false,
		},
		{
			name:        "Env file found in parent directory",
			currentDir:  "/home/user/project/sub_dir",
			envFileName: ".env",
			setupMocks: func(mockOS *mocks.MockOS, mockStat *mocks.MockStat) {
				mockOS.On("Getwd").Return("/home/user/project/sub_dir", nil)
				mockStat.On("Stat", "/home/user/project/sub_dir/.env").Return(nil, os.ErrNotExist)
				mockStat.On("Stat", "/home/user/project/.env").Return(nil, nil)
			},
			expectedResult: "/home/user/project/.env",
			expectedError:  false,
		},
		{
			name:        "Env file not found",
			currentDir:  "/home/user/project/sub_dir",
			envFileName: ".env",
			setupMocks: func(mockOS *mocks.MockOS, mockStat *mocks.MockStat) {
				mockOS.On("Getwd").Return("/home/user/project/sub_dir", nil)
				mockStat.On("Stat", "/home/user/project/sub_dir/.env").Return(nil, os.ErrNotExist)
				mockStat.On("Stat", "/home/user/project/.env").Return(nil, os.ErrNotExist)
				mockStat.On("Stat", "/home/user/.env").Return(nil, os.ErrNotExist)
				mockStat.On("Stat", "/home/.env").Return(nil, os.ErrNotExist)
				mockStat.On("Stat", "/.env").Return(nil, os.ErrNotExist)
			},
			expectedResult: "",
			expectedError:  true,
		},
		{
			name:        "Error getting current directory",
			currentDir:  "",
			envFileName: ".env",
			setupMocks: func(mockOS *mocks.MockOS, mockStat *mocks.MockStat) {
				mockOS.On("Getwd").Return("", errors.New("cannot get current directory"))
			},
			expectedResult: "",
			expectedError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockOS := new(mocks.MockOS)
			mockStat := new(mocks.MockStat)

			test.setupMocks(mockOS, mockStat)

			helper.OSGetwd = mockOS.Getwd
			helper.OSStat = mockStat.Stat

			result, err := helper.GetEnvFilePath(test.envFileName)

			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectedResult, result)
			}

			mockOS.AssertExpectations(t)
			mockStat.AssertExpectations(t)
		})
	}
}
