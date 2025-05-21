package entity

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUserInput(t *testing.T) {
	tests := []struct {
		name        string
		input       WalletData
		expected    WalletData
		expectedErr error
	}{
		{
			name: "valid deposit operation",
			input: WalletData{
				Uuid:           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Balance:        100.123,
				Operation_type: "DEPOSIT",
			},
			expected: WalletData{
				Uuid:           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Balance:        100.12,
				Operation_type: "deposit",
			},
			expectedErr: nil,
		},
		{
			name: "valid withdraw operation",
			input: WalletData{
				Uuid:           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Balance:        50.567,
				Operation_type: "withdraw",
			},
			expected: WalletData{
				Uuid:           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Balance:        50.57, // Rounded up
				Operation_type: "withdraw",
			},
			expectedErr: nil,
		},
		{
			name: "zero balance",
			input: WalletData{
				Uuid:           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Balance:        0,
				Operation_type: "deposit",
			},
			expected: WalletData{
				Uuid:           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Balance:        0,
				Operation_type: "deposit",
			},
			expectedErr: nil,
		},

		{
			name: "invalid operation type",
			input: WalletData{
				Uuid:           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Balance:        100,
				Operation_type: "transfer",
			},
			expected:    WalletData{},
			expectedErr: errors.New("invalid operation"),
		},
		{
			name: "empty operation type",
			input: WalletData{
				Uuid:           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Balance:        100,
				Operation_type: "",
			},
			expected:    WalletData{},
			expectedErr: errors.New("invalid operation"),
		},

		{
			name: "negative balance",
			input: WalletData{
				Uuid:           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Balance:        -10.555,
				Operation_type: "deposit",
			},
			expected: WalletData{
				Uuid:           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Balance:        -10.56,
				Operation_type: "deposit",
			},
			expectedErr: nil,
		},
		{
			name: "very large balance",
			input: WalletData{
				Uuid:           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Balance:        999999999.999,
				Operation_type: "withdraw",
			},
			expected: WalletData{
				Uuid:           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				Balance:        1000000000.00,
				Operation_type: "withdraw",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Validate(tt.input)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
