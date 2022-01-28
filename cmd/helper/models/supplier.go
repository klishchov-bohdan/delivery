package models

type Supplier struct {
	ID           uint64       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	Image        string       `json:"image"`
	WorkingHours WorkingHours `json:"workingHours"`
	Menu         []Product    `json:"menu"`
}

type SuppliersResponse struct {
	Suppliers []Supplier `json:"suppliers"`
}
