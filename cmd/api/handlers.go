package main

import (
	"backend/cmd/api/config"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Data types
type Post struct {
    ID          string    `json:"id"`
    UserID      string    `json:"user_id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    CreatedAt   time.Time `json:"created_at"`
    LikesCount  int       `json:"likes_count"`
    ViewsCount  int       `json:"views_count"`
    Comments    []Comment `json:"comments"`
    UpdatedAt   time.Time `json:"updated_at"`
    Tags        []string  `json:"tags"`
}

type Comment struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

type User struct {
    ID           string    `json:"id"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}

type Tag struct {
    ID       string `json:"id"`
    Text     string `json:"text"`
    Color    string `json:"color"`
    Searches int    `json:"searches"`
}

type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

// Home handler
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from %s", app.Domain)
}

// GetPosts handler
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

    rows, err := db.Query(`
        SELECT id, user_id, title, content,
               created_at, likes_count, views_count, 
               comments, updated_at, tags 
        FROM posts
        ORDER BY created_at DESC`)
    if err != nil {
        fmt.Println("Query error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        var comments []byte

        err := rows.Scan(
            &post.ID,
            &post.UserID,
            &post.Title,
            &post.Content,
            &post.CreatedAt,
            &post.LikesCount,
            &post.ViewsCount,
            &comments,
            &post.UpdatedAt,
            pq.Array(&post.Tags),
        )
        if err != nil {
            fmt.Println("Row scan error:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        if len(comments) > 0 {
            err = json.Unmarshal(comments, &post.Comments)
            if err != nil {
                fmt.Println("Comments parsing error:", err)
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        }

        posts = append(posts, post)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(posts)
}

// Login handler
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        fmt.Printf("Error decoding request body: %v\n", err)
        http.Error(w, "Invalid request Body", http.StatusBadRequest)
        return
    }

    connStr := config.GetDBConfig()
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        fmt.Printf("Database connection error: %v\n", err)
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var user User
    err = db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = $1", creds.Email).Scan(&user.ID, &user.Email, &user.PasswordHash)
    if err != nil {
        fmt.Printf("User lookup error: %v\n", err)
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
        fmt.Printf("Password comparison failed: %v\n", err)
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(app.JwtSecret))
    if err != nil {
        fmt.Printf("Token generation error: %v\n", err)
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    response := struct {
        Token string `json:"token"`
        User  struct {
            ID    string `json:"id"`
            Email string `json:"email"`
        } `json:"user"`
    }{
        Token: tokenString,
        User: struct {
            ID    string `json:"id"`
            Email string `json:"email"`
        }{
            ID:    user.ID,
            Email: user.Email,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(response)
}

// CreatePost handler
func (app *application) CreatePost(w http.ResponseWriter, r *http.Request) {
    fmt.Println("CreatePost endpoint hit")

    var post Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        fmt.Printf("Error decoding request body: %v\n", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    connStr := config.GetDBConfig()
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        fmt.Printf("Database connection error: %v\n", err)
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    // Set creation time and updated time to current time
    currentTime := time.Now()
    post.CreatedAt = currentTime
    post.UpdatedAt = currentTime
    post.LikesCount = 0
    post.ViewsCount = 0

    // Initialize empty arrays if they're nil
    if post.Comments == nil {
        post.Comments = []Comment{}
    }
    if post.Tags == nil {
        post.Tags = []string{}
    }

    // Convert comments to JSONB format
    commentsJSON, err := json.Marshal(post.Comments)
    if err != nil {
        fmt.Printf("Error marshaling comments: %v\n", err)
        http.Error(w, "Error processing comments", http.StatusInternalServerError)
        return
    }

    // Debug print
    fmt.Printf("Inserting post with values: %+v\n", post)

    err = db.QueryRow(`
    INSERT INTO posts (
        id,           -- 1
        user_id,      -- 2
        title,        -- 3
        content,      -- 4
        created_at,   -- 5
        updated_at,   -- 6
        likes_count,  -- 7
        views_count,  -- 8
        comments,     -- 9
        tags          -- 10
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING id`,
    post.ID,
    post.UserID,
    post.Title,
    post.Content,
    post.CreatedAt,
    post.UpdatedAt, 
    post.LikesCount,
    post.ViewsCount,
    commentsJSON,
    pq.Array(post.Tags),
).Scan(&post.ID)

    if err != nil {
        fmt.Printf("Error inserting post: %v\n", err)
        http.Error(w, "Error creating post", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(post)
}

// GetTags handler
func (app *application) GetTags(w http.ResponseWriter, r *http.Request) {
    fmt.Println("GetTags endpoint hit")
    connStr := config.GetDBConfig()
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        fmt.Println("Database connection error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()

    rows, err := db.Query(`SELECT id, text, color, searches FROM tags`)
    if err != nil {
        fmt.Println("Query error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var tags []Tag
    for rows.Next() {
        var tag Tag
        err := rows.Scan(&tag.ID, &tag.Text, &tag.Color, &tag.Searches)
        if err != nil {
            fmt.Println("Row scan error:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        tags = append(tags, tag)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(tags)
}

// CreateTag handler
func (app *application) CreateTag(w http.ResponseWriter, r *http.Request) {
    fmt.Println("CreateTag endpoint hit")

    var tag Tag
    if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
        fmt.Printf("Error decoding request body: %v\n", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    connStr := config.GetDBConfig()
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        fmt.Printf("Database connection error: %v\n", err)
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    err = db.QueryRow(`
        INSERT INTO tags (text, color, searches)
        VALUES ($1, $2, $3)
        RETURNING id`,
        tag.Text,
        tag.Color,
        0, // Initial searches count
    ).Scan(&tag.ID)

    if err != nil {
        fmt.Printf("Error inserting tag: %v\n", err)
        http.Error(w, "Error creating tag", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(tag)
}

// UpdateTagSearchCount handler
func (app *application) UpdateTagSearchCount(w http.ResponseWriter, r *http.Request) {
    fmt.Println("UpdateTagSearchCount endpoint hit")
    
    // Get tag ID from URL using Chi instead of mux.Vars
    tagID := chi.URLParam(r, "id")
    if tagID == "" {
        http.Error(w, "Tag ID is required", http.StatusBadRequest)
        return
    }

    connStr := config.GetDBConfig()
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        fmt.Printf("Database connection error: %v\n", err)
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var tag Tag
    err = db.QueryRow(`
        UPDATE tags 
        SET searches = searches + 1
        WHERE id = $1
        RETURNING id, text, color, searches`,
        tagID,
    ).Scan(&tag.ID, &tag.Text, &tag.Color, &tag.Searches)

    if err != nil {
        fmt.Printf("Error updating tag search count: %v\n", err)
        http.Error(w, "Error updating tag", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(tag)
}