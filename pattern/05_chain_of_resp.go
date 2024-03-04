package pattern

import (
	"fmt"
	"log"
)

// Позволяет передавать запросы по цепочке обработчиков
// Плюсы: уменьшает зависимость между клиентом и обработчиками, реализует принципы единственной обязаноости и открытости/закрытости
// Минусы: запрос может остаться никем не обработанным

func ExampleChainOfResponsibility() {
	helperAnsweringMachine := &HelperAnsweringMachine{}
	helperTechnicalSupport := &HelperTechnicalSupport{}
	helperTechnicalSpecialist := &HelperTechnicalSpecialist{}

	helperAnsweringMachine.SetNext(helperTechnicalSupport)
	helperTechnicalSupport.SetNext(helperTechnicalSpecialist)

	helperAnsweringMachine.Help()
}

type Helper interface {
	Help()
	SetNext(h Helper)
}

type HelperBasic struct {
	Next Helper
}

func (helper *HelperBasic) SetNext(h Helper) {
	helper.Next = h
}

type HelperAnsweringMachine struct {
	HelperBasic
}

func (helper HelperAnsweringMachine) Help() {
	fmt.Println("Hello, I'm answering machine")
	fmt.Println("Have you tried turning it off and on?")
	var str string
	_, err := fmt.Scan(&str)
	if err != nil {
		log.Fatal(err)
	}
	if str == "no" {
		fmt.Println("Do it")
		return
	}
	helper.Next.Help()
}

type HelperTechnicalSupport struct {
	HelperBasic
}

func (helper HelperTechnicalSupport) Help() {
	fmt.Println("Hello, I'm technical support")
	fmt.Println("It's windows?")
	var str string
	_, err := fmt.Scan(&str)
	if err != nil {
		log.Fatal(err)
	}
	if str == "yes" {
		fmt.Println("Press F")
		return
	}
	helper.Next.Help()
}

type HelperTechnicalSpecialist struct {
	HelperBasic
}

func (helper HelperTechnicalSpecialist) Help() {
	fmt.Println("Hello, I'm technical specialist")
	fmt.Println("Now I can help you")
}
