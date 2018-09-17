package abstraction

type Good interface {
	Increase(amt float64)
	Decrease(amt float64) float64
	GetQuantity() float64
}
