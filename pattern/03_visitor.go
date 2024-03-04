package pattern

import (
	"fmt"
	"math"
)

// Позволяет добавлять в программу новые операции, не изменяя классы объектов
// Плюсы: упрощает добавление операций, работающих со сложными структурами, объединяет родственные операции в одном классе
// Минусы: может привести к нарушению инкапсуляции элементов, не оправдан, если иерархия элементов часто меняется

func ExampleVisitor() {
	sPoint := Point{0, 0}

	point1 := Point{1, 0}

	point2 := Point{2, 4}

	point3 := Point{2, -3}

	line := Line{point2, point3}

	v := &VisitorDistance{
		SPoint: sPoint,
	}
	point1.Accept(v)
	fmt.Printf("point 1: x: %.2f, y: %.2f;\npoint 2: x: %.2f, y: %.2f;\n", sPoint.X, sPoint.Y, point1.X, point1.Y)
	fmt.Printf("dist: %.2f;\n", v.Res)
	line.Accept(v)
	fmt.Printf("point 1: x: %.2f, y: %.2f;\nline: x1: %.2f, y1: %.2f; x2: %.2f, y2: %.2f;\n", sPoint.X, sPoint.Y, line.X1.X, line.X1.Y, line.X2.X, line.X2.Y)
	fmt.Printf("dist: %.2f;\n", v.Res)
}

type Point struct {
	X float64
	Y float64
}

type Line struct {
	X1 Point
	X2 Point
}

type Visitor interface {
	DoPoint(x Point)
	DoLine(x Line)
}

func (p *Point) Accept(v Visitor) {
	v.DoPoint(*p)
}

func (l *Line) Accept(v Visitor) {
	v.DoLine(*l)
}

type VisitorDistance struct {
	SPoint Point
	Res    float64
}

func (vd *VisitorDistance) DoPoint(x Point) {
	vd.Res = math.Sqrt(math.Pow(float64(x.X)-float64(vd.SPoint.X), 2) + math.Pow(float64(x.Y)-float64(vd.SPoint.Y), 2))
}

func (vd *VisitorDistance) DoLine(x Line) {
	vd.Res = math.Abs((x.X2.Y-x.X1.Y)*vd.SPoint.X-(x.X2.X-x.X1.X)*vd.SPoint.Y+
		x.X2.X*x.X1.Y-x.X2.Y*x.X1.X) / math.Sqrt(math.Pow(x.X2.Y-x.X1.Y, 2)+
		math.Pow(x.X2.X-x.X1.X, 2))
}
