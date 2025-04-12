package main

import "github.com/yehiamoh/go-fem-workshop/pkg/server"

func main() {

	if err := server.Run(); err != nil {
		panic(err.Error())
	}

}
