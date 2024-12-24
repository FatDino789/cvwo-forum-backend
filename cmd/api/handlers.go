package main

import (
	"backend/cmd/api/config"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type Post struct {
    ID        int    `json:"id"`
    UserID    int    `json:"user_id"`
    Title     string `json:"title"`
    Content   string `json:"content"`
    CreatedAt string `json:"created_at"`
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from %s", app.Domain)
}

func (app *application) GetPosts(w http.ResponseWriter, r *http.Request) {
    fmt.Println("GetPosts endpoint hit")
    
    connStr := config.GetDBConfig()
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        fmt.Println("Database connection error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    // Updated query to match your schema
    rows, err := db.Query("SELECT id, user_id, title, content, created_at FROM posts")
    if err != nil {
        fmt.Println("Query error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
        if err != nil {
            fmt.Println("Row scan error:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        posts = append(posts, post)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    if err := json.NewEncoder(w).Encode(posts); err != nil {
        fmt.Println("JSON encoding error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Printf("Returning %d posts\n", len(posts))
}