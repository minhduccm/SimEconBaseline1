package abstraction

type Agent interface {
	GetID() string
	Act()
	PrintLastState()
	GetWalletAccountAddress() string
	GetGood(string) Good
	GetConsumption(string) float64
}
