package entity

import (
	"errors"
	"math"
	"strings"
	"time"
)

type WalletData struct {
	Uuid           string    `json:"uuid"`
	Balance        float64   `json:"balance"`
	Operation_type string    `json:"type"`
	Created_At     time.Time `json:"time"`
}

func roundTo(f float64, dec int) float64 {
	presicion := math.Pow(10, float64(dec))
	return math.Round(f*presicion) / presicion
}

func Validate(w WalletData) (WalletData, error) {
	w.Balance = roundTo(w.Balance, 2)
	w.Operation_type = strings.ToLower(w.Operation_type)

	switch w.Operation_type {
	case "deposit":
	case "withdraw":
	default:
		return WalletData{}, errors.New("invalid operation")
	}

	return w, nil
}
