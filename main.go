package main

import (
	"fmt"

	"github.com/faxter/bloggator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	err = cfg.SetUser("kat")
	if err != nil {
		fmt.Println(err)
	}
	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}
