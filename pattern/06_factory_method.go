package pattern

import (
	"fmt"
	"math/rand"
)

// Определяет общий интерфейс для создания объектов в суперклассе
// Плюсы: избавляет класс от привязки к конкретным классам продуктов, упрощает поддержку кода
// Минусы: может привести к созданию больших параллельных иерархий классов

func ExampleFactoryMethod() {
	var factory Factory
	if rand.Intn(2) == 0 {
		factory = FactoryCar{}
	} else {
		factory = FactoryToyCar{}
	}
	product := factory.CreateProduct()
	product.DoSomething()
}

type Factory interface {
	CreateProduct() Product
}

type FactoryCar struct {
}

func (p FactoryCar) CreateProduct() Product {
	return ProductCar{}
}

type FactoryToyCar struct {
}

func (p FactoryToyCar) CreateProduct() Product {
	return ProductToyCar{}
}

type Product interface {
	DoSomething()
}

type ProductCar struct {
}

func (p ProductCar) DoSomething() {
	fmt.Println("Car is controlled by the driver")
}

type ProductToyCar struct {
}

func (p ProductToyCar) DoSomething() {
	fmt.Println("Toy car is driven by hand")
}
