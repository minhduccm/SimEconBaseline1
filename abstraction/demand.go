package abstraction

type Demand interface {
	GetDemand(float64, float64) float64
}
