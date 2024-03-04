package pattern

import "fmt"

// Создает сложные объекты пошагово
// Плюсы: использование одного и того же кода для создания разных объектов, изоляция кода от бизнес-логики
// Минусы: усложняет код программы из-за введения новых классов, клиент будет привязан к конкретным классам строителей

func ExampleBuilder() {
	d := Director{}
	b := &BuilderHouse{}
	d.CreateBox(b)
	b.PrintRes()
	d.CreateRichHouse(b)
	b.PrintRes()
}

type Builder interface {
	SetWalls(int)
	SetWindows(int)
	SetRooms(int)
	SetGarden(bool)
}

type house struct {
	ExternalWalls int
	Windows       int
	Rooms         int
	Garden        bool
}

func (h house) String() string {
	strGarden := ""
	if h.Garden {
		strGarden = " and a beautiful garden"
	}
	return fmt.Sprintf("House has %d external walls, %d windows, %d rooms%s.", h.ExternalWalls, h.Windows, h.Rooms, strGarden)
}

type BuilderHouse struct {
	houseRes house
}

func (b *BuilderHouse) PrintRes() {
	fmt.Println(b.houseRes)
}

func (b *BuilderHouse) SetWalls(count int) {
	if count < 3 {
		panic("house can't have less than 3 external walls")
	}
	b.houseRes.ExternalWalls = count
}

func (b *BuilderHouse) SetWindows(count int) {
	b.houseRes.Windows = count
}

func (b *BuilderHouse) SetRooms(count int) {
	if count < 1 {
		panic("house can't have less than 1 room")
	}
	b.houseRes.Rooms = count
}

func (b *BuilderHouse) SetGarden(has bool) {
	b.houseRes.Garden = has
}

type Director struct {
}

func (d Director) CreateBox(b Builder) {
	b.SetGarden(false)
	b.SetRooms(1)
	b.SetWalls(4)
	b.SetWindows(0)
}

func (d Director) CreateRichHouse(b Builder) {
	b.SetGarden(true)
	b.SetRooms(100)
	b.SetWalls(20)
	b.SetWindows(40)
}
