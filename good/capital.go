package good

import (
	"math/rand"

	"github.com/ninjadotorg/SimEconBaseline1/common"
)

type Capital struct {
	Quantity         float64
	OwnerID          string
	Machines         []*Machine
	ScrappedMachines []*Machine // machines scrapped in the current step
}

type Machine struct {
	// (per step) price of the machine at purchase time
	Price float64

	// number of steps the machine could be used for
	Life int

	// remaining number of steps the machine could be used for
	RemainingLife int

	// producer of this machine
	ProducerID string
}

func NewCapital(
	quantity int,
	ownerID string,
	producerIDs []string,
) *Capital {
	machines := []*Machine{}
	for i := 0; i < quantity; i++ {
		machine := &Machine{}
		machine.Price = common.INIT_CAPITAL_PRICE
		machine.Life = common.CAPITAL_LIFE
		machine.RemainingLife = rand.Intn(machine.Life/2) + machine.Life/2
		machine.ProducerID = producerIDs[rand.Intn(len(producerIDs))]
		machines = append(machines, machine)
	}
	return &Capital{
		Quantity:         float64(quantity),
		OwnerID:          ownerID,
		Machines:         machines,
		ScrappedMachines: []*Machine{},
	}
}

func (c *Capital) Increase(amt float64) {
	return
}

func (c *Capital) Decrease(amt float64) float64 {
	return 0.0
}

func (c *Capital) GetQuantity() float64 {
	return c.Quantity
}

func (c *Capital) Add(
	qty int,
	price float64,
	life int,
	producerID string,
) {
	c.Quantity += float64(qty)
	for i := 0; i < qty; i++ {
		machine := &Machine{}
		machine.Price = price
		machine.Life = life
		machine.RemainingLife = life
		machine.ProducerID = producerID
		c.Machines = append(c.Machines, machine)
	}
}

func (c *Capital) GetPresentValue() float64 {
	var val float64 = 0
	for _, m := range c.Machines {
		val += m.Price * float64(m.RemainingLife)
	}
	return val
}
