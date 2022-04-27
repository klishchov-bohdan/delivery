package controller

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/klishchov-bohdan/delivery/config"
	"github.com/klishchov-bohdan/delivery/internal/models"
	"github.com/klishchov-bohdan/delivery/internal/requests"
	"github.com/klishchov-bohdan/delivery/internal/services"
	"io/ioutil"
	"net/http"
)

type SupplierController struct {
	services *services.Manager
	cfg      *config.Config
}

func NewSupplierController(services *services.Manager, cfg *config.Config) *SupplierController {
	return &SupplierController{
		services: services,
		cfg:      cfg,
	}
}

// GetAllSuppliers godoc
// @Summary Get Suppliers
// @Description Get Suppliers
// @Tags suppliers
// @ID get-suppliers
// @Accept  json
// @Produce  json
// @Success 200 {array} models.SupplierWeb
// @Router /suppliers [get]
func (ctr *SupplierController) GetAllSuppliers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		suppliers, err := ctr.services.Supplier.GetAllSuppliers()
		if err != nil {
			http.Error(w, "Cant get all suppliers", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(suppliers)
	default:
		http.Error(w, "Only GET method", http.StatusMethodNotAllowed)
	}
}

// GetSupplier godoc
// @Summary Get Supplier
// @Description Get Supplier
// @Tags suppliers
// @ID get-supplier
// @Accept  json
// @Produce  json
// @Param supplierId path string true "Supplier ID"
// @Success 200 {object} models.SupplierWeb
// @Router /suppliers/{supplierId} [get]
func (ctr *SupplierController) GetSupplier(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := chi.URLParam(r, "supplierId")
		uid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		supplier, err := ctr.services.Supplier.GetSupplierByID(uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(supplier)
	default:
		http.Error(w, "Only GET method", http.StatusMethodNotAllowed)
	}
}

// GetSupplierMenu godoc
// @Summary Get Supplier Menu
// @Description Get Supplier Menu
// @Tags suppliers
// @ID get-supplier-menu
// @Accept  json
// @Produce  json
// @Param supplierId path string true "Supplier ID"
// @Success 200 {array} models.MenuItem
// @Router /suppliers/{supplierId}/menu [get]
func (ctr *SupplierController) GetSupplierMenu(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := chi.URLParam(r, "supplierId")
		uid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		menu, err := ctr.services.Product.GetMenuBySupplierID(uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(menu)
	default:
		http.Error(w, "Only GET method", http.StatusMethodNotAllowed)
	}
}

// GetSupplierMenuItem godoc
// @Summary Get Supplier Menu Item
// @Description Get Supplier Menu Item
// @Tags suppliers
// @ID get-supplier-menu-item
// @Accept  json
// @Produce  json
// @Param productId path string true "Product ID"
// @Success 200 {object} models.MenuItem
// @Router /suppliers/menu/{productId} [get]
func (ctr *SupplierController) GetSupplierMenuItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := chi.URLParam(r, "productId")
		uid, err := uuid.Parse(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		supplier, err := ctr.services.Product.GetMenuItemByID(uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(supplier)
	default:
		http.Error(w, "Only GET method", http.StatusMethodNotAllowed)
	}
}

// CreateSupplier godoc
// @Summary Create Supplier
// @Description Create Supplier
// @Tags suppliers
// @ID create-supplier
// @Accept  json
// @Produce  json
// @Param user body requests.SupplierWebRequest true "Create Supplier Input"
// @Success 200 {object} models.SupplierWeb
// @Router /suppliers [post]
func (ctr *SupplierController) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		body, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		var supplierReq requests.SupplierWebRequest
		err := json.Unmarshal(body, &supplierReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createdSupplierID, err := ctr.services.Supplier.CreateSupplier(
			models.NewSupplierWeb(supplierReq.GetSupplier(), supplierReq.GetProducts()))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		createdSupplier, _ := ctr.services.Supplier.GetSupplierByID(createdSupplierID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdSupplier)
	default:
		http.Error(w, "Only POST method", http.StatusMethodNotAllowed)
	}
}
