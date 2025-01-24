package helper

import "fmt"

func PanicIfError(err error) {
	if err != nil {
		fmt.Println("Error:", err.Error())
		panic(err)
	}
}
