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
emails := []string{"sem.responsavel@teste.com", "frontend.4201@teste.com"}
for _, email := range emails {
users, _ := client.GetUsers(ctx, token.AccessToken, "cepprs", gocloak.GetUsersParams{Email: &email})
if len(users) > 0 {
u := users[0]
fmt.Printf("User: %s | Required Actions: %v\n", email, *u.RequiredActions)
}
}
}
