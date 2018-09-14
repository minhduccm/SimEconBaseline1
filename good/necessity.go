package good

type Necessity struct {
	Quantity float64
}

func (g *Necessity) Increase(amt float64) {
	g.Quantity += amt
}

func (g *Necessity) Decrease(amt float64) float64 {
	var ret float64 = 0
	if g.Quantity > amt {
		ret = amt
	} else {
		ret = g.Quantity
	}
	g.Quantity -= ret
	return ret
}

func (g *Necessity) GetQuantity() float64 {
	return g.Quantity
}
