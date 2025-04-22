package atom

type inData interface {
	Show()
	Description() string
	Add()
	Delite()
	Complite()
	Do()
}

type Atom struct {
	inter  []Atom
	inData inData
}

func New() Atom {
	return Atom{}
}
func (a Atom) Add() {

}
