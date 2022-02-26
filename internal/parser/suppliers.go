package parser

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/parser/responses"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type SupplierParser struct {
	services    services.SupplierService
	suppliers   []*models.SupplierWeb
	supplierIdx map[uint64]uuid.UUID
	productIdx  map[uint64]uuid.UUID
}

func NewSupplierParser(services services.SupplierService) *SupplierParser {
	return &SupplierParser{
		services:    services,
		supplierIdx: make(map[uint64]uuid.UUID),
		productIdx:  make(map[uint64]uuid.UUID),
	}
}

func (parser *SupplierParser) ParseSuppliers() error {
	suppResp, err := http.Get("http://foodapi.true-tech.php.nixdev.co/suppliers")
	if err != nil {
		return err
	}
	defer suppResp.Body.Close()
	suppBody, err := ioutil.ReadAll(suppResp.Body)
	if err != nil {
		return err
	}
	var suppliers responses.SuppliersResponse
	err = json.Unmarshal(suppBody, &suppliers)
	if err != nil {
		return err
	}
	for i, supplier := range suppliers.Suppliers {
		menuResp, err := http.Get(fmt.Sprintf("http://foodapi.true-tech.php.nixdev.co/suppliers/%d/menu", supplier.ID))
		if err != nil {
			return err
		}
		defer menuResp.Body.Close()
		menuBody, err := ioutil.ReadAll(menuResp.Body)
		if err != nil {
			return err
		}
		var menu responses.MenuResponse
		err = json.Unmarshal(menuBody, &menu)
		if err != nil {
			return err
		}
		suppliers.Suppliers[i].Menu = menu.Products
	}
	for _, supplier := range suppliers.Suppliers {
		supplierId := uuid.New()
		parser.supplierIdx[supplier.ID] = supplierId
		if strings.HasPrefix(supplier.WorkingHours.Opening, "24:") {
			supplier.WorkingHours.Opening = strings.Replace(supplier.WorkingHours.Opening, "24:", "00:", 1)
		}
		openIn, err := time.Parse(time.RFC3339, fmt.Sprintf("2006-01-02T%s:05Z", supplier.WorkingHours.Opening))
		if err != nil {
			return err
		}
		if strings.HasPrefix(supplier.WorkingHours.Closing, "24:") {
			supplier.WorkingHours.Closing = strings.Replace(supplier.WorkingHours.Closing, "24:", "00:", 1)
		}
		closeIn, err := time.Parse(time.RFC3339, fmt.Sprintf("2006-01-02T%s:05Z", supplier.WorkingHours.Closing))
		if err != nil {
			return err
		}
		var menu []*models.MenuItem
		for _, product := range supplier.Menu {
			productId := uuid.New()
			parser.productIdx[supplier.ID] = productId
			menuItem := &models.MenuItem{
				ID:          productId,
				Name:        product.Name,
				Image:       product.Image,
				Description: strings.Join(product.Ingredients, ", "),
				Price:       float64(product.Price),
				Weight:      0,
				Type:        product.Type,
			}
			menu = append(menu, menuItem)
		}
		supplierWeb := &models.SupplierWeb{
			ID:          supplierId,
			Name:        supplier.Name,
			Image:       supplier.Image,
			Description: "",
			WorkingTime: models.WorkingSchedule{
				OpenIn:      openIn,
				CloseIn:     closeIn,
				WorkingDays: "",
			},
			Menu: menu,
		}
		parser.suppliers = append(parser.suppliers, supplierWeb)
	}
	return nil
}

func (parser *SupplierParser) GetSuppliers() []*models.SupplierWeb {
	return parser.suppliers
}

func (parser *SupplierParser) SaveToDB() error {
	for _, supplier := range parser.suppliers {
		_, err := parser.services.CreateSupplier(supplier)
		if err != nil {
			return err
		}
	}
	return nil
}
