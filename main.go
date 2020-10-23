package main

import (
	"fmt"
	"io/ioutil"
	"ops-backend/pkg"
	"ops-backend/pkg/util"
	"ops-backend/router"
)

var (
	mysupersecretpassword = "unicornsAreAwesome"
)

func init() {
	pkg.Setup()

}

func main() {
	content, _ := ioutil.ReadFile("conf/start.txt")
	fmt.Println(util.Red(string(content)))
	router.Setup()
}
