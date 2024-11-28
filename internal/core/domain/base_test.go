package domain_test

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChallengeType_String(t *testing.T) {
	tests := []struct {
		name           string
		challengeType  domain.ChallengeType
		expectedResult string
	}{
		{
			name:           "TOTP challenge type",
			challengeType:  domain.ChallengeTOTP,
			expectedResult: domain.ChallengeTOTPStr,
		},
		{
			name:           "Recovery Codes challenge type",
			challengeType:  domain.ChallengeRecoveryCodes,
			expectedResult: domain.ChallengeRecoveryCodesStr,
		},
		{
			name:           "Unknown challenge type",
			challengeType:  domain.ChallengeType("UNKNOWN_TYPE"),
			expectedResult: domain.UnknownStr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedResult, test.challengeType.String())
		})
	}
}

func TestChallengeStatusType_String(t *testing.T) {
	tests := []struct {
		name           string
		status         domain.ChallengeStatusType
		expectedResult string
	}{
		{
			name:           "Enable challenge status",
			status:         domain.ChallengeEnable,
			expectedResult: domain.ChallengeEnableStr,
		},
		{
			name:           "Disabled challenge status",
			status:         domain.ChallengeDisable,
			expectedResult: domain.ChallengeDisableStr,
		},
		{
			name:           "Unknown challenge status",
			status:         domain.ChallengeStatusType("UNKNOWN_STATUS"),
			expectedResult: domain.UnknownStr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedResult, test.status.String())
		})
	}
}
