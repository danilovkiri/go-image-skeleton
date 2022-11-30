package neighbors

type Point struct {
	X     int
	Y     int
	Color int
}

type Neighbours struct {
	P1 Point
	P2 Point
	P3 Point
	P4 Point
	P5 Point
	P6 Point
	P7 Point
	P8 Point
	P9 Point
}

func (ng *Neighbours) GetTransitions() int {
	cycle := []Point{ng.P2, ng.P3, ng.P4, ng.P5, ng.P6, ng.P7, ng.P8, ng.P9, ng.P2}
	var count int
	for i := 0; i < len(cycle)-1; i++ {
		if cycle[i].Color == 0 && cycle[i+1].Color == 1 {
			count++
		}
	}
	return count
}

func (ng *Neighbours) GetNonZeros() int {
	cycle := []Point{ng.P2, ng.P3, ng.P4, ng.P5, ng.P6, ng.P7, ng.P8, ng.P9}
	var count int
	for i := 0; i < len(cycle); i++ {
		if cycle[i].Color == 1 {
			count++
		}
	}
	return count
}
