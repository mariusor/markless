package markless

func NewLine(bytes []byte) *Line {
	return &Line{Data: bytes, Size: len(bytes)}
}

type Line struct {
	Data     []byte
	Size     int
	previous *Line
	next     *Line
}
