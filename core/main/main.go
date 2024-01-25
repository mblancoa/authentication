package main

import (
	"fmt"
	"github.com/x/authentication/core"
	"github.com/x/authentication/core/tools"
)

func main() {
	credentials := core.UserCredentials{
		ID:       "maria.blanco.arranz@gmail.com",
		Password: "estaEsC0nEñe",
		State:    "todoCorrecto",
	}

	hCredentials := core.UserCredentials{}
	_ = tools.MarshalHash(credentials, &hCredentials)
	fmt.Println(credentials)
	fmt.Println(hCredentials)

	hCredentials2 := core.UserCredentials{}
	_ = tools.MarshalHash(credentials, &hCredentials2)
	fmt.Println(hCredentials2)
}
