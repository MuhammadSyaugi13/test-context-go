package controller

import (
	"coba-konteks/app"
	"coba-konteks/helper"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/xuri/excelize/v2"
)

var (
	bulkInsertCtx    context.Context
	bulkInsertCancel context.CancelFunc
	mu               sync.Mutex
)

// func readExcel(filename string) [][]string {
func readExcel(filename string) [][]string {
	// f, err := excelize.OpenFile("asset/20.xlsx")
	// f, err := excelize.OpenFile("asset/500.xlsx")
	// f, err := excelize.OpenFile("asset/1000.xlsx")

	f, err := excelize.OpenFile("asset/" + filename + ".xlsx")

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	// Get all sheet names
	sheetNames := f.GetSheetList()
	fmt.Println("Sheets in the file:")
	for _, name := range sheetNames {
		fmt.Println(name)
	}

	// Read data from a specific sheet
	sheetName := sheetNames[0] // Assuming you want the first sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Error reading rows: %v", err)
		helper.PanicIfError(err)
	}

	//tampilkan excel
	// for _, row := range rows {
	// 	for _, cell := range row {
	// 		fmt.Printf("%s\t", cell)
	// 	}
	// 	fmt.Println()
	// }

	return rows

}

var numGoroutine int

func BulkInsertHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()

	if bulkInsertCancel != nil {
		http.Error(w, "Bulk insert is already running", http.StatusConflict)
		return
	}

	db := app.NewDB()

	var cancelCtx context.Context
	cancelCtx, bulkInsertCancel = context.WithCancel(context.Background())

	// Membaca body dari request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() // Pastikan untuk menutup body setelah selesai

	fmt.Println(string(body))

	data := readExcel(string(body))

	// chanNumGoroutine := make(channel)

	go func(ctx context.Context) {

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}

		stmt, err := tx.PrepareContext(ctx, "INSERT INTO category (name) VALUES (?)")
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to prepare statement", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		for _, row := range data {
			select {
			case <-ctx.Done():
				tx.Rollback()
				http.Error(w, "Bulk insert cancelled", http.StatusRequestTimeout)
				return
			default:
				_, err := stmt.ExecContext(ctx, row[0])
				if err != nil {
					tx.Rollback()
					http.Error(w, "Failed to execute statement", http.StatusInternalServerError)
					return
				}
				fmt.Println("insert data", row[0], "Go Num :", strconv.Itoa(runtime.NumGoroutine()))
			}

			// time.Sleep(500 * time.Millisecond)
		}

		err = tx.Commit()
		if err != nil {
			fmt.Println("masuk error")
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		// numGoroutine = runtime.NumGoroutine()
		// fmt.Println("num goroutine :", strconv.Itoa(runtime.NumGoroutine()))

		bulkInsertCancel()
		bulkInsertCancel = nil

		// http.Error(w, "Bulk insert successful", http.StatusOK)

	}(cancelCtx)

	fmt.Println("selesai")

}

func CancelBulkInsertHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mu.Lock()
	defer mu.Unlock()

	if bulkInsertCancel != nil {
		bulkInsertCancel()
		bulkInsertCancel = nil
		http.Error(w, "Bulk insert cancelled", http.StatusOK)
	} else {
		http.Error(w, "No bulk insert process to cancel", http.StatusNotFound)
	}
}
