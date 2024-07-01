package domain

type Language struct {
	Base

	Name   string
	Code   string
	Status StatusType
}
