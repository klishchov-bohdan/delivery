package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/klishchov-bohdan/delivery/cmd/helper/models"
	"github.com/klishchov-bohdan/delivery/internal/worker_pool"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	wc1 := func() worker_pool.Worker {
		return ParserJSON{}
	}
	wpool1 := worker_pool.NewPool(4, wc1)
	wpool1.DataSource = make(chan interface{})

	go wpool1.Run()
	for i := 1; i <= 7; i++ {
		wpool1.DataSource <- i
	}
	wpool1.Stop()
	//////////////
	wc2 := func() worker_pool.Worker {
		return WriterDB{}
	}
	wpool2 := worker_pool.NewPool(4, wc2)
	wpool2.DataSource = make(chan interface{})

	go wpool2.Run()
	for i := 1; i <= 7; i++ {
		wpool2.DataSource <- i
	}
	wpool2.Stop()
}

type ParserJSON struct {
}

func (w ParserJSON) Do(data interface{}, i int) {
	get, err := http.Get(fmt.Sprintf("http://foodapi.true-tech.php.nixdev.co/public/supplier_%d.json", data))
	if err != nil {
		log.Fatal(err)
	}
	defer get.Body.Close()
	created, err := os.Create(fmt.Sprintf("./cmd/helper/suppliers/supplier_%d", data))
	if err != nil {
		log.Fatal(err)
	}
	supplier, err := ioutil.ReadAll(get.Body)
	if err != nil {
		log.Fatal(err)
	}
	_, err = created.WriteString(string(supplier))
	if err != nil {
		log.Fatal(err)
	}
}
func (w ParserJSON) Stop() {

}

////////////////////////////////////////
type WriterDB struct {
}

func (w WriterDB) Do(data interface{}, i int) {
	fmt.Printf("Routine %d: Writing supplier %d to db\n", i, data)
	supplier := &models.Supplier{}
	supplierRepo, err := ioutil.ReadDir("./cmd/helper/suppliers")
	if err != nil {
		log.Fatal(err)
	}
	for _, fileInfo := range supplierRepo {
		file, err := os.Open("./cmd/helper/suppliers/" + fileInfo.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		var fileData []byte
		for {
			chunk := make([]byte, 64)
			n, err := file.Read(chunk)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fileData = append(fileData, chunk[:n]...)
		}
		err = json.Unmarshal(fileData, &supplier)
		if err != nil {
			log.Fatal(err)
		}
		supplierID, err := strconv.ParseUint(fmt.Sprint(data), 10, 64)
		if err != nil {
			return
		}
		if supplier.ID == supplierID {
			db, err := sql.Open(
				"mysql",
				"bohdan:Bohdan.2525@tcp(127.0.0.1:3306)/shop?parseTime=true",
			)
			if err != nil {
				log.Fatal(err)
			}
			err = db.Ping()
			if err != nil {
				log.Fatal(err)
			}
			stmt1, err := db.Prepare("INSERT INTO suppliers(id, name, type, image, opening, closing) VALUES(?, ?, ?, ?, ?, ?)")
			if err != nil {
				log.Fatal(err)
			}
			_, err = stmt1.Exec(supplier.ID, supplier.Name, supplier.Type, supplier.Image, supplier.WorkingHours.Opening, supplier.WorkingHours.Closing)
			if err != nil {
				log.Fatal(err)
			}
			for _, product := range *supplier.Menu {
				stmt2, err := db.Prepare("INSERT INTO products(id, name, price, image, type, ingredients, supplier_id) VALUES(?, ?, ?, ?, ?, ?, ?)")
				if err != nil {
					log.Fatal(err)
				}
				ingredients := strings.Join(product.Ingredients, ", ")
				_, err = stmt2.Exec(product.ID, product.Name, product.Price, product.Image, product.Type, ingredients, supplier.ID)
				if err != nil {
					log.Fatal(err)
				}
			}
			fmt.Printf("Routine %d: Done\n", i)
			break
		}
	}
}
func (w WriterDB) Stop() {

}
