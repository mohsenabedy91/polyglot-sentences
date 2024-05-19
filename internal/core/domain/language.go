package domain

type Language struct {
	Base

	Name   string
	Code   string
	Status StatusType
}

func (r Language) IsActive() bool {
	return r.Status == StatusActive
}
