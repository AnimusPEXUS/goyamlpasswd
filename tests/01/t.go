package main

import (
	"log"

	"github.com/AnimusPEXUS/goyamlpasswd"
)

func main() {
	yf := goyamlpasswd.NewYAMLAuthFile("./auth.yaml")

	// yf.PutRecord(&goyamlpasswd.YAMLAuthFileSRecord{User: "u1", Key: &[]string{"superkey1"}[0]})
	// yf.PutRecord(&goyamlpasswd.YAMLAuthFileSRecord{User: "u2", Key: &[]string{"superkey2"}[0]})

	// err := yf.Save()
	// if err != nil {
	// 	log.Fatalln("error Save:", err)
	// }

	err := yf.Load()
	if err != nil {
		log.Fatalln("error Load:", err)
	}

	user, err := yf.UserByKey("key1")
	if err != nil {
		log.Fatalln("error UserByKey:", err)
	}

	log.Println("user:", user)
}
