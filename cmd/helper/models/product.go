package models

type Product struct {
	ID          uint64
	Name        string
	Price       float32
	Image       string
	Type        string
	Ingredients []string
}
