package pattern

import (
	"fmt"
	"math/rand"
)

// Определяет семейство схожих алгоритмов и помещает каждый из них в собственный класс
// Плюсы: упрощенная замена алгоритмов, изоляция кода и данных алгоритмов, уход от наследования к делегированию
// Минусы: усложняет программу за счет новых классов, клиент должен знать, в чем разница между стратегиями

func ExampleStrategy() {
	context := Context{}
	if rand.Intn(2) == 0 {
		context.SetStrategy(Add{})
	} else {
		context.SetStrategy(Sub{})
	}
	fmt.Println(context.ExecuteStrategy(1, 2))
}

type Strategy interface {
	Execute(a, b int) int
}

type Add struct {
}

func (s Add) Execute(a, b int) int {
	return a + b
}

type Sub struct {
}

func (s Sub) Execute(a, b int) int {
	return a - b
}

type Context struct {
	strategy Strategy
}

func (c *Context) SetStrategy(s Strategy) {
	c.strategy = s
}

func (c *Context) ExecuteStrategy(a, b int) int {
	return c.strategy.Execute(a, b)
}
