package helper

import (
	"fmt"
	"runtime"
)

func PanicIfError(err error) {

	if err != nil {

		// menampilkan pesan file dan baris yang terjadi error pada terminal
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("Error pada file %s, line: %d \n", file, line)

		// menampilkan pesan error pada terminal
		fmt.Println(err.Error())

	}

}
