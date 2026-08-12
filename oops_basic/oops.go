package oopsbasic

type Customer struct {
	name string
	Id   int
}

func (c *Customer) GetName() string {
	return c.name
}

func (c *Customer) GetId() int {
	return c.Id
}
