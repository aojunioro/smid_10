package main

import (
	"fmt"
	"github.com/aojunioro/smid_10/backend/internal/domain/admin"
)

func main() {
	hash, err := admin.HashPassword("Admin123!")
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
}
