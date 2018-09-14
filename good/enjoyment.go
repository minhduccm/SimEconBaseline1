package good

type Enjoyment struct {
	Quantity float64
}

func (g *Enjoyment) Increase(amt float64) {
	g.Quantity += amt
}

func (g *Enjoyment) Decrease(amt float64) float64 {
	var ret float64 = 0
	if g.Quantity > amt {
		ret = amt
	} else {
		ret = g.Quantity
	}
	g.Quantity -= ret
	return ret
}

func (g *Enjoyment) GetQuantity() float64 {
	return g.Quantity
}
