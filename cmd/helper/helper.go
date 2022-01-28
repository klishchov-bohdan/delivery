package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/klishchov-bohdan/delivery/cmd/helper/models"
	"github.com/klishchov-bohdan/delivery/internal/store/db"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()
	ctxWithCancel, cancelFunction := context.WithCancel(ctx)

	defer func() {
		fmt.Println("Main Defer: canceling context")
		cancelFunction()
	}()

	suppliers, err := parseSuppliers(ctxWithCancel)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := db.Dial()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	err = createSuppliers(conn, suppliers, ctxWithCancel)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				suppliers, err := parseSuppliers(ctxWithCancel)
				if err != nil {
					log.Fatal(err)
				}
				err = updateSuppliersPrice(conn, suppliers, ctxWithCancel)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()
	time.Sleep(10 * time.Second)
}

func parseSuppliers(ctx context.Context) (*[]models.Supplier, error) {
	_, cancelFunc := context.WithTimeout(ctx, time.Second)
	defer func() {
		log.Println("parseSuppliers: context cancel")
		cancelFunc()
	}()
	suppResp, err := http.Get("http://foodapi.true-tech.php.nixdev.co/suppliers")
	if err != nil {
		return nil, err
	}
	defer suppResp.Body.Close()
	suppBody, err := ioutil.ReadAll(suppResp.Body)
	if err != nil {
		return nil, err
	}
	var suppliers models.SuppliersResponse
	err = json.Unmarshal(suppBody, &suppliers)
	if err != nil {
		return nil, err
	}
	for i, supplier := range suppliers.Suppliers {
		menuResp, err := http.Get(fmt.Sprintf("http://foodapi.true-tech.php.nixdev.co/suppliers/%d/menu", supplier.ID))
		if err != nil {
			return nil, err
		}
		defer menuResp.Body.Close()
		menuBody, err := ioutil.ReadAll(menuResp.Body)
		if err != nil {
			return nil, err
		}
		var menu models.MenuResponse
		err = json.Unmarshal(menuBody, &menu)
		if err != nil {
			return nil, err
		}
		suppliers.Suppliers[i].Menu = menu.Products
	}
	return &suppliers.Suppliers, nil
}

func createSuppliers(db *sql.DB, suppliers *[]models.Supplier, ctx context.Context) error {
	_, cancelFunc := context.WithTimeout(ctx, time.Second)
	defer func() {
		log.Println("createSuppliers: context cancel")
		cancelFunc()
	}()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for _, supplier := range *suppliers {
		stmt1, err := tx.Prepare("INSERT INTO suppliers(id, name, type, image, opening, closing) VALUES(?, ?, ?, ?, ?, ?)")
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}
		_, err = stmt1.Exec(supplier.ID, supplier.Name, supplier.Type, supplier.Image, supplier.WorkingHours.Opening, supplier.WorkingHours.Closing)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}
		for _, product := range supplier.Menu {
			stmt2, err := tx.Prepare("INSERT INTO products(id, name, price, image, type, ingredients, supplier_id) VALUES(?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				err := tx.Rollback()
				if err != nil {
					return err
				}
				return err
			}
			ingredients := strings.Join(product.Ingredients, ", ")
			_, err = stmt2.Exec(product.ID, product.Name, product.Price, product.Image, product.Type, ingredients, supplier.ID)
			if err != nil {
				err := tx.Rollback()
				if err != nil {
					return err
				}
				return err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func updateSuppliersPrice(db *sql.DB, suppliers *[]models.Supplier, ctx context.Context) error {
	_, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
	defer func() {
		log.Println("updateSuppliersMenu: context cancel")
		cancelFunc()
	}()
	for _, supplier := range *suppliers {
		for _, product := range supplier.Menu {
			stmt, err := db.Prepare("UPDATE products SET price = ? WHERE id = ?")
			if err != nil {
				return err
			}
			_, err = stmt.Exec(product.Price, product.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
