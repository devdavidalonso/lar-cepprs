package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Nerzal/gocloak/v13"
)

func main() {
	baseURL := "http://localhost:8081"
	client := gocloak.NewClient(baseURL)
	ctx := context.Background()

	realm := "cepprs"
	adminRealm := "master"
	username := "admin"
	password := "pigu@1025"

	token, err := client.LoginAdmin(ctx, username, password, adminRealm)
	if err != nil {
		log.Fatalf("failed to login: %v", err)
	}

	emails := []string{"sem.responsavel@teste.com", "frontend.4201@teste.com"}
	newPassword := "senha123"

	for _, email := range emails {
		users, err := client.GetUsers(ctx, token.AccessToken, realm, gocloak.GetUsersParams{
			Email: gocloak.StringP(email),
		})
		if err != nil || len(users) == 0 {
			fmt.Printf("User %s not found or error: %v\n", email, err)
			continue
		}

		userID := *users[0].ID
		err = client.SetPassword(ctx, token.AccessToken, userID, realm, newPassword, false)
		if err != nil {
			fmt.Printf("Failed to set password for %s: %v\n", email, err)
		} else {
			fmt.Printf("Successfully set password for %s to %s\n", email, newPassword)
		}
	}
}
