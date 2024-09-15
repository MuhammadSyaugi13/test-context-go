package main

import (
	"coba-konteks/controller"
	"coba-konteks/helper"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {

	router := httprouter.New()
	// router.GET("/api/cancel-proccess", controller.CancelProccess)
	// router.POST("/api/categories", controller.InsertKelas)

	router.GET("/api/jumlah-goroutine", controller.NumGoroutine)

	router.GET("/api/bulk-insert", controller.BulkInsertHandler)
	router.GET("/api/cancel-bulk-insert", controller.CancelBulkInsertHandler)

	// router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}

	fmt.Println("Menjalankan server")
	err := server.ListenAndServe()
	helper.PanicIfError(err)

}
