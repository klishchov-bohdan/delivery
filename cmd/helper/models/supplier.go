package models

type Supplier struct {
	ID           uint64
	Name         string
	Type         string
	Image        string
	WorkingHours WorkingHours
	Menu         *[]Product
}
