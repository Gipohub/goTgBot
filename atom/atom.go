package atom

type Atom struct {
	inter []Atom
	data  any
}

func New() Atom {
	return Atom{}
}
func (a Atom) Add() {

}
