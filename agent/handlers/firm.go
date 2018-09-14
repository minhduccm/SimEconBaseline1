package handlers

import (
	agentModels "github.com/ninjadotorg/SimEconBaseline1/agent/models"
	"github.com/ninjadotorg/SimEconBaseline1/economy"
	"github.com/ninjadotorg/SimEconBaseline1/good"
	"github.com/ninjadotorg/SimEconBaseline1/util"
)

func NewFirm(initWalletBal float64) *agentModels.Firm {
	firmID := util.NewUUID()
	econ := economy.GetEconInstance()
	econ.TransactionManager.OpenWalletAccount(
		firmID,
		initWalletBal,
	)

	return &agentModels.Firm{
		ID: firmID,
		Labor: &good.Labor{
			Quantity: 0,
		},
	}
}
