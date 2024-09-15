package controller

import (
	"coba-konteks/app"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

var (
	bulkInsertCtx    context.Context
	bulkInsertCancel context.CancelFunc
	mu               sync.Mutex
)

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

	go func(ctx context.Context) {
		data := []struct {
			Field1 string
			Field2 string
		}{
			{"Value1", "Value2"},
			{"Value2", "Value4"},
			{"Value3", "Value4"},
			{"Value3", "Value4"},
			{"Value4", "Value4"},
			{"Value6", "Value4"},
			{"Value7", "Value4"},
			{"Value8", "Value4"},
			{"Value9", "Value4"},
			{"Value10", "Value4"},
			{"Value11", "Value4"},
			{"Value12", "Value4"},
			{"Value13", "Value4"},
			{"Value14", "Value4"},
			{"Value15", "Value4"},
			// Add more data if needed
		}

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
				_, err := stmt.ExecContext(ctx, row.Field1)
				if err != nil {
					tx.Rollback()
					http.Error(w, "Failed to execute statement", http.StatusInternalServerError)
					return
				}
				fmt.Println("insert data", row.Field1)
			}

			time.Sleep(500 * time.Millisecond)
		}

		err = tx.Commit()
		if err != nil {
			fmt.Println("masuk error")
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		bulkInsertCancel()
		bulkInsertCancel = nil

		http.Error(w, "Bulk insert successful", http.StatusOK)

	}(cancelCtx)
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
