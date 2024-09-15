package controller

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

// func CancelProccess(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

// 	if cancelFunc != nil {
// 		// cancelFunc()
// 		fmt.Fprintln(w, "Slow process cancellation requested")
// 	} else {
// 		http.Error(w, "No slow process to cancel", http.StatusBadRequest)
// 	}

// }

// func InsertKelas(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

// 	ctx := r.Context()
// 	var cancel context.CancelFunc

// 	ctx, cancel = context.WithCancel(ctx)
// 	cancelFunc = cancel

// 	for i := 0; i < 100; i++ {

// 		select {
// 		// case <-time.After(2 * time.Second):
// 		//     fmt.Fprintln(w, "Slow process completed")
// 		case <-ctx.Done():
// 			// http.Error(w, "Slow process canceled", http.StatusRequestTimeout)
// 			w.Header().Add("Content-Type", "application/json")
// 			dataJson, err := json.Marshal(`{"status": "success cancel"}`)
// 			if err != nil {
// 				helper.PanicIfError(err)
// 			}

// 			w.Write(dataJson)
// 			return
// 		default:

// 			var errMasukanData error

// 			_, errMasukanData = masukanData(w, r)

// 			helper.PanicIfError(errMasukanData)

// 		}

// 		// time.Sleep(100 * time.Millisecond)
// 	}

// 	fmt.Println("jumlah goroutine :", strconv.Itoa(runtime.NumGoroutine()))

// 	w.Header().Add("Content-Type", "application/json")
// 	dataJson, err := json.Marshal(`{"status": "success cancel"}`)
// 	if err != nil {
// 		helper.PanicIfError(err)
// 	}

// 	w.Write(dataJson)

// }

// func masukanData(w http.ResponseWriter, r *http.Request) (sql.Result, error) {

// 	db := app.NewDB()

// 	tx, err := db.Begin()
// 	helper.PanicIfError(err)
// 	defer helper.CommitOrRollback(tx)

// 	SQL := "insert into category(name) values (?)"

// 	return tx.ExecContext(r.Context(), SQL, "celana")
// }

// func InsertBaru(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

// }

var mtx sync.Mutex

func NumGoroutine(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	startGoroutine := runtime.NumGoroutine()
	fmt.Println("start num goroutine :", strconv.Itoa(startGoroutine))

	for i := 0; i < 10000; i++ {
		go func() {
			defer mtx.Unlock()
			mtx.Lock()
			_ = 2 + 9
		}()
	}

	fmt.Println("num goroutine :", strconv.Itoa(runtime.NumGoroutine()))
}
