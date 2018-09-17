package abstraction

type LaborMarket interface {
	AddEmployee(string, string)
	AddEmployer(string, string, Good, float64)
	Perform()
}
