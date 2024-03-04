package pattern

import "fmt"

// Позволяет объектам менять поведение в зависимости от состояния
// Плюсы: упрощает код
// Минусы: если состояний мало и они редко меняется, наоборот, усложнит код

func ExampleState() {
	car := NewCar()
	car.StepOnGas()

	car.TurnIgnitionKey()
	car.StepOnGas()
}

type State interface {
	TurnIgnitionKey()
	StepOnGas()
}

type StateStarted struct {
	car *Car
}

func NewStateStarted(car *Car) *StateStarted {
	return &StateStarted{
		car: car,
	}
}

func (c *StateStarted) TurnIgnitionKey() {
	c.car.ChangeState(NewStateNotStarted(c.car))
}

func (c *StateStarted) StepOnGas() {
	c.car.Go()
}

type StateNotStarted struct {
	car *Car
}

func NewStateNotStarted(car *Car) *StateNotStarted {
	return &StateNotStarted{
		car: car,
	}
}

func (c *StateNotStarted) TurnIgnitionKey() {
	c.car.ChangeState(NewStateStarted(c.car))
}

func (c *StateNotStarted) StepOnGas() {
}

type Car struct {
	state State
}

func NewCar() *Car {
	car := &Car{}
	state := NewStateNotStarted(car)
	car.ChangeState(state)
	return car
}

func (c *Car) ChangeState(state State) {
	c.state = state
}

func (c *Car) TurnIgnitionKey() {
	c.state.TurnIgnitionKey()
}

func (c *Car) StepOnGas() {
	c.state.StepOnGas()
}

func (c Car) Go() {
	fmt.Println("car rides")
}
