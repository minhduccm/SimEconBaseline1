package transaction_manager

import (
	"github.com/ninjadotorg/SimEconBaseline1/common"
	"github.com/ninjadotorg/SimEconBaseline1/util"
)

type TransactionManager struct {
	WalletAccounts map[string]*WalletAccount
}

func (transManager *TransactionManager) OpenWalletAccount(
	agentID string,
	balance float64,
) {
	newAddress := util.NewUUID()
	acc := NewWalletAccount(newAddress, balance)
	transManager.WalletAccounts[agentID] = acc
}

func (transManager *TransactionManager) CloseWalletAccount(
	agentID string,
) {
	if _, ok := transManager.WalletAccounts[agentID]; ok {
		delete(transManager.WalletAccounts, agentID)
	}
}

func (transManager *TransactionManager) GetBalance(
	agentId string,
) float64 {
	acc, ok := transManager.WalletAccounts[agentId]
	if !ok {
		return 0
	}
	return acc.Balance
}

func (transManager *TransactionManager) PayFrom(
	payerID string,
	amt float64,
) {
	fromAcc := transManager.WalletAccounts[payerID]
	// if fromAcc.Balance < amt { // TODO: will handle this case later

	// }
	fromAcc.Balance -= amt
}

func (transManager *TransactionManager) PayTo(
	payeeID string,
	amt float64,
	purpose int, // either PRIIC or SECIC
) {
	toAcc := transManager.WalletAccounts[payeeID]
	toAcc.Balance += amt
	if purpose == common.PRIIC {
		toAcc.PriIC += amt
		return
	}
	toAcc.SecIC += amt
}

func (transManager *TransactionManager) Pay(
	payerID string,
	payeeID string,
	amt float64,
	purpose int, // either PRIIC or SECIC
) {
	fromAcc := transManager.WalletAccounts[payerID]
	toAcc := transManager.WalletAccounts[payeeID]

	fromAcc.Balance -= amt
	toAcc.Balance += amt
	if purpose == common.PRIIC {
		toAcc.PriIC += amt
		return
	}
	toAcc.SecIC += amt
}
