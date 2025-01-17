package main

import (
	"fmt"
	"vr-shope/internal/app"
)

func main() {
	err := app.Run("internal/config/config.yaml")
	if err != nil {
		fmt.Println(err)
	}
}
