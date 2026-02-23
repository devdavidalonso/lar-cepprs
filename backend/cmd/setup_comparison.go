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

// Ensure frontend.4201 has UPDATE_PASSWORD
email := "frontend.4201@teste.com"
users, _ := client.GetUsers(ctx, token.AccessToken, "cepprs", gocloak.GetUsersParams{Email: &email})
if len(users) > 0 {
u := users[0]
u.RequiredActions = &[]string{"UPDATE_PASSWORD"}
err = client.UpdateUser(ctx, token.AccessToken, "cepprs", *u)
if err != nil { fmt.Printf("Error updating: %v\n", err) }
fmt.Printf("Updated %s with UPDATE_PASSWORD\n", email)
}

// Ensure sem.responsavel has NO required actions
email2 := "sem.responsavel@teste.com"
users2, _ := client.GetUsers(ctx, token.AccessToken, "cepprs", gocloak.GetUsersParams{Email: &email2})
if len(users2) > 0 {
u := users2[0]
u.RequiredActions = &[]string{}
err = client.UpdateUser(ctx, token.AccessToken, "cepprs", *u)
fmt.Printf("Ensured %s has NO actions\n", email2)
}
}
