package models

type Product struct {
	BaseModelSoftDelete
	Name  string
	Price float64
}
