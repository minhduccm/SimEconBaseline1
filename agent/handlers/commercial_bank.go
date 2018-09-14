package handlers

// import "math"

// const (
// 	// balance type - checking account
// 	CHECKING = 1

// 	// account type - saving account
// 	SAVING = 2

// 	// primary income: wage for laborers, sales revenue for firms
// 	PRIIC = 1

// 	// secondary income: e.g. dividend
// 	SECIC = 2

// 	// other payment
// 	OTHER = 3
// )

// type CommercialBank struct {
// 	// acounts
// 	Accounts map[int]*Account

// 	// total amount of loans
// 	TotalLoan float64

// 	// total amount of deposits
// 	TotalDeposit float64

// 	// loan interest rate
// 	LoanIR float64

// 	// deposit interest rate
// 	DepositIR float64

// 	// long-term loan interest rate
// 	LtLoanIR float64

// 	// long-term deposit interest rate
// 	LtDepositIR float64

// 	// classes used to compute the average interest rate within LT_IR_WIN
// 	// private static Averager depositIRAvger = new Averager(LT_IR_WIN),
// 	// 		loanIRAvger = new Averager(LT_IR_WIN);

// 	TargetIR float64
// }

// type Account struct {
// 	Balances map[int]float64 // checking/saving acc balance
// 	PriIC    float64         // primary income in the last step
// 	SecIC    float64         // secondary income in the last step
// 	Interest float64         // interest in the last step
// }

// func NewAccount(checkingBal float64, savingBal float64) *Account {
// 	return &Account{
// 		Balances: map[int]float64{
// 			CHECKING: checkingBal,
// 			SAVING:   savingBal,
// 		},
// 		PriIC:    0,
// 		SecIC:    0,
// 		Interest: 0,
// 	}
// }

// func (acc *Account) GetBalance(balanceType int) float64 {
// 	return acc.Balances[balanceType]
// }

// // Bank section

// func (comBank *CommercialBank) Act() {

// }

// func (comBank *CommercialBank) OpenAccount(
// 	agentID int,
// 	initCheckingBal float64,
// 	initSavingsBal float64,
// ) {
// 	acc := NewAccount(initCheckingBal, initSavingsBal)
// 	comBank.Accounts[agentID] = acc
// }

// func (comBank *CommercialBank) CloseAccount(
// 	agentID int,
// ) {
// 	if _, ok := comBank.Accounts[agentID]; ok {
// 		delete(comBank.Accounts, agentID)
// 	}
// }

// func (comBank *CommercialBank) GetBalance(
// 	agentId int,
// 	balanceType int,
// ) float64 {
// 	acc, ok := comBank.Accounts[agentId]
// 	if !ok {
// 		return 0
// 	}
// 	return acc.GetBalance(balanceType)
// }

// func (comBank *CommercialBank) PayFrom(
// 	payerID int,
// 	amt float64,
// ) {
// 	fromAcc := comBank.Accounts[payerID]
// 	if fromAcc.Balances[CHECKING] < amt {
// 		diff := amt - fromAcc.Balances[CHECKING]
// 		fromAcc.Balances[CHECKING] += diff
// 		fromAcc.Balances[SAVING] -= diff
// 	}
// 	fromAcc.Balances[CHECKING] -= amt
// }

// func (comBank *CommercialBank) PayTo(
// 	payeeID int,
// 	amt float64,
// 	purpose int, // either PRIIC or SECIC
// ) {
// 	toAcc := comBank.Accounts[payeeID]
// 	toAcc.Balances[CHECKING] += amt
// 	if purpose == PRIIC {
// 		toAcc.PriIC += amt
// 		return
// 	}
// 	toAcc.SecIC += amt
// }

// func (comBank *CommercialBank) Pay(
// 	payerID int,
// 	payeeID int,
// 	amt float64,
// 	purpose int, // either PRIIC or SECIC
// ) {
// 	fromAcc := comBank.Accounts[payerID]
// 	toAcc := comBank.Accounts[payeeID]
// 	if fromAcc.Balances[CHECKING] < amt {
// 		diff := amt - fromAcc.Balances[CHECKING]
// 		fromAcc.Balances[CHECKING] += diff
// 		fromAcc.Balances[SAVING] -= diff
// 	}
// 	fromAcc.Balances[CHECKING] -= amt
// 	toAcc.Balances[CHECKING] += amt
// 	if purpose == PRIIC {
// 		toAcc.PriIC += amt
// 		return
// 	}
// 	toAcc.SecIC += amt
// }

// func (comBank *CommercialBank) Deposit(
// 	agentId int,
// 	amt float64,
// ) float64 {
// 	acc := comBank.Accounts[agentId]
// 	realAmt := math.Min(amt, acc.GetBalance(CHECKING))
// 	acc.Balances[CHECKING] -= realAmt
// 	acc.Balances[SAVING] += realAmt
// 	return realAmt
// }
