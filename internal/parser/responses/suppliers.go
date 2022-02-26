package responses

type SupplierParsed struct {
	ID           uint64          `json:"id"`
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	Image        string          `json:"image"`
	WorkingHours WorkingHours    `json:"workingHours"`
	Menu         []ProductParsed `json:"menu"`
}

type ProductParsed struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`
	Price       float32  `json:"price"`
	Image       string   `json:"image"`
	Type        string   `json:"type"`
	Ingredients []string `json:"ingredients"`
}

type MenuResponse struct {
	Products []ProductParsed `json:"menu"`
}

type SuppliersResponse struct {
	Suppliers []SupplierParsed `json:"suppliers"`
}

type WorkingHours struct {
	Opening string `json:"opening"`
	Closing string `json:"closing"`
}
