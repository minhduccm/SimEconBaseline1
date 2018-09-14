package models

import (
	"github.com/ninjadotorg/SimEconBaseline1/good"
)

type Agent interface {
	Act()
	GetWalletAccountAddress() string
	GetGood(string) good.Good
	GetConsumption(string) float64
}
