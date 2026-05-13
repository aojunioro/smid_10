package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash := "$2a$10$/5r1gHi4YWRRQEPtBvVSOei8HR0nD7QdsmNFrzkUtHvmImvFmFIca"
	password := "Admin123!"
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("Senha inválida: %v\n", err)
	} else {
		fmt.Printf("Senha válida\n")
	}
}
