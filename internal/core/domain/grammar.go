package domain

type Grammar struct {
	Base
	Modifier

	Title  string
	Status StatusType
}

func (r Grammar) IsActive() bool {
	return r.Status == StatusActive
}
