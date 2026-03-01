package service

import (
	"testing"
)

func TestValidateAmount(t *testing.T) {
	tests := []struct {
		name    string
		amount  string
		wantErr bool
	}{
		{"valid amount", "10.15", false},
		{"valid integer", "100", false},
		{"one decimal", "5.5", false},
		{"zero", "0", true},
		{"negative", "-5.00", true},
		{"empty", "", true},
		{"too many decimals", "10.123", true},
		{"not a number", "abc", true},
		{"spaces", "  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAmount(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAmount(%q) error = %v, wantErr = %v", tt.amount, err, tt.wantErr)
			}
		})
	}
}

func TestValidateState(t *testing.T) {
	tests := []struct {
		name    string
		state   string
		wantErr bool
	}{
		{"win", "win", false},
		{"lose", "lose", false},
		{"invalid", "draw", true},
		{"empty", "", true},
		{"uppercase", "WIN", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateState(tt.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateState(%q) error = %v, wantErr = %v", tt.state, err, tt.wantErr)
			}
		})
	}
}

func TestValidateSourceType(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{"game", "game", false},
		{"server", "server", false},
		{"payment", "payment", false},
		{"invalid", "unknown", true},
		{"empty", "", true},
		{"uppercase", "GAME", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSourceType(tt.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateSourceType(%q) error = %v, wantErr = %v", tt.source, err, tt.wantErr)
			}
		})
	}
}