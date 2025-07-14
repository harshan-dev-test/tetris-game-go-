package grid

type Grid struct {
	Width  int
	Height int
	Data   [][]int
}

func NewGrid(width, height int) *Grid {
	data := make([][]int, height)
	for i := range data {
		data[i] = make([]int, width)
	}

	return &Grid{
		Width:  width,
		Height: height,
		Data:   data,
	}
}