package atom

type inData interface {
	Show()
	Description()
	Add()
	Delite()
	Complite()
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
