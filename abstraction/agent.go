package abstraction

type Agent interface {
	GetID() string
	Act()
	GetWalletAccountAddress() string
	GetGood(string) Good
	GetConsumption(string) float64
}
