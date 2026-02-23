package main
import (
"context"
"fmt"
"log"
"github.com/Nerzal/gocloak/v13"
)
func main() {
client := gocloak.NewClient("http://localhost:8081")
ctx := context.Background()
token, err := client.LoginAdmin(ctx, "admin", "pigu@1025", "master")
if err != nil { log.Fatal(err) }

realm := "cepprs"

// Test sem.responsavel
e1 := "sem.responsavel@teste.com"
u1, _ := client.GetUsers(ctx, token.AccessToken, realm, gocloak.GetUsersParams{Email: &e1})
if len(u1) > 0 {
user := u1[0]
user.RequiredActions = &[]string{}
err = client.UpdateUser(ctx, token.AccessToken, realm, *user)
fmt.Printf("User %s updated. Actions: %v | Error: %v\n", e1, *user.RequiredActions, err)
} else {
fmt.Printf("User %s not found\n", e1)
}

// Test frontend.4201
e2 := "frontend.4201@teste.com"
u2, _ := client.GetUsers(ctx, token.AccessToken, realm, gocloak.GetUsersParams{Email: &e2})
if len(u2) > 0 {
user := u2[0]
user.RequiredActions = &[]string{"UPDATE_PASSWORD"}
err = client.UpdateUser(ctx, token.AccessToken, realm, *user)
fmt.Printf("User %s updated. Actions: %v | Error: %v\n", e2, *user.RequiredActions, err)
} else {
fmt.Printf("User %s not found\n", e2)
}
}
