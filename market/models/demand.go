package models

type Demand interface {
	GetDemand(price float64, consumption float64) float64
}
