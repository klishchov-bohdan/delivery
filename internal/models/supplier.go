package models

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type Supplier struct {
	ID          uuid.UUID
	Name        string
	Image       string
	Description string
	WorkingTime *WorkingSchedule
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type WorkingSchedule struct {
	OpenIn      time.Time `json:"open_in"`
	CloseIn     time.Time `json:"close_in"`
	WorkingDays string    `json:"working_days"`
}

func NewWorkingSchedule(openIn time.Time, closeIn time.Time, workingDays []time.Weekday) *WorkingSchedule {
	var days string
	for _, day := range workingDays {
		days += day.String() + " "
	}
	days = strings.TrimSpace(days)
	return &WorkingSchedule{
		OpenIn:      openIn,
		CloseIn:     closeIn,
		WorkingDays: days,
	}
}

type SupplierWeb struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Image       string           `json:"image"`
	Description string           `json:"description"`
	WorkingTime *WorkingSchedule `json:"working_time"`
	Menu        []*MenuItem      `json:"menu"`
}

func NewSupplierWeb(supplier *Supplier, menu *[]Product) *SupplierWeb {
	supplierWeb := &SupplierWeb{
		ID:          supplier.ID,
		Name:        supplier.Name,
		Image:       supplier.Image,
		Description: supplier.Description,
		WorkingTime: supplier.WorkingTime,
	}
	for _, product := range *menu {
		supplierWeb.Menu = append(supplierWeb.Menu, product.ToMenuItem())
	}
	return supplierWeb
}
