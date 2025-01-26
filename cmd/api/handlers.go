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
    Username    string `json:"username"`
    IconIndex   int    `json:"icon_index"`
    ColorIndex  int    `json:"color_index"`
}

type Comment struct {
    ID          string    `json:"id"`
    UserID      string    `json:"user_id"`
    Content     string    `json:"content"`
    CreatedAt   time.Time `json:"created_at"`
    Username    string    `json:"username"`
    IconIndex   int       `json:"icon_index"`
    ColorIndex  int       `json:"color_index"`
 }

type User struct {
    ID            string    `json:"id"`
    Username      string    `json:"username"` 
    Email         string    `json:"email"`
    PasswordHash  string    `json:"-"`
    PlainPassword string    `json:"password"`
    CreatedAt     time.Time `json:"created_at"`
    IconIndex     int       `json:"icon_index"`   
    ColorIndex    int       `json:"color_index"`   
    Likes         []string  `json:"likes"`
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
    db, err := sql.Open("postgres", config.GetDBConfig())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer db.Close()
 
    rows, err := db.Query(`
        SELECT p.id, p.user_id, p.title, p.content,
               p.created_at, p.likes_count, p.views_count, 
               p.comments, p.updated_at, p.tags,
               u.username, u.icon_index, u.color_index
        FROM posts p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC`)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()
 
    var posts []Post
    for rows.Next() {
        var post Post
        var commentsStr []byte
 
        err := rows.Scan(
            &post.ID,
            &post.UserID,
            &post.Title,
            &post.Content,
            &post.CreatedAt,
            &post.LikesCount,
            &post.ViewsCount,
            &commentsStr,
            &post.UpdatedAt,
            pq.Array(&post.Tags),
            &post.Username,
            &post.IconIndex,
            &post.ColorIndex,
        )
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
 
        // Parse comments array
        if len(commentsStr) > 0 {
            if err := json.Unmarshal(commentsStr, &post.Comments); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            
            // For each comment, fetch the user details
            for i := range post.Comments {
                var username string
                var iconIndex, colorIndex int
                err := db.QueryRow(`
                    SELECT username, icon_index, color_index 
                    FROM users WHERE id = $1`, 
                    post.Comments[i].UserID).Scan(&username, &iconIndex, &colorIndex)
                if err != nil {
                    continue // Skip if user not found
                }
                post.Comments[i].Username = username
                post.Comments[i].IconIndex = iconIndex
                post.Comments[i].ColorIndex = colorIndex
            }
        }
 
        posts = append(posts, post)
    }
 
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(posts)
 }

 func (app *application) StreamPosts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")
 
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "SSE not supported", http.StatusInternalServerError)
        return
    }
 
    // Initial posts
    posts, err := app.fetchPosts()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    data, _ := json.Marshal(posts)
    fmt.Fprintf(w, "data: %s\n\n", data)
    flusher.Flush()
 
    // Poll for updates
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
 
    for {
        select {
        case <-ticker.C:
            posts, err := app.fetchPosts()
            if err != nil {
                continue
            }
            data, _ := json.Marshal(posts)
            fmt.Fprintf(w, "data: %s\n\n", data)
            flusher.Flush()
        case <-r.Context().Done():
            return
        }
    }
 }

 // UpdatePost handler
 func (app *application) UpdatePost(w http.ResponseWriter, r *http.Request) {
    postID := chi.URLParam(r, "id") // Extract the post ID from the URL path
    if postID == "" {
        http.Error(w, "Post ID is required", http.StatusBadRequest)
        return
    }

    var requestBody struct {
        Field string      `json:"field"`
        Value interface{} `json:"value"`
    }

    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        fmt.Printf("Error decoding request body: %v\n", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    allowedFields := map[string]bool{
        "likes_count": true,
        "views_count": true,
        "title":       true,
        "content":     true,
        "tags":        true,
    }

    if !allowedFields[requestBody.Field] {
        http.Error(w, "Invalid field for update", http.StatusBadRequest)
        return
    }

    updatedAt := time.Now().UTC()

    query := fmt.Sprintf("UPDATE posts SET %s = $1, updated_at = $2 WHERE id = $3 RETURNING id, %s, updated_at", requestBody.Field, requestBody.Field)
    var updatedValue interface{}

    db, err := sql.Open("postgres", config.GetDBConfig())
    if err != nil {
        fmt.Printf("Database connection error: %v\n", err)
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    err = db.QueryRow(query, requestBody.Value, updatedAt, postID).Scan(&postID, &updatedValue, &updatedAt)
    if err != nil {
        fmt.Printf("Error updating post: %v\n", err)
        http.Error(w, "Error updating post", http.StatusInternalServerError)
        return
    }

    response := struct {
        PostID      string      `json:"post_id"`
        UpdatedField string      `json:"field"`
        UpdatedValue interface{} `json:"value"`
        UpdatedAt   time.Time   `json:"updated_at"`
    }{
        PostID:      postID,
        UpdatedField: requestBody.Field,
        UpdatedValue: updatedValue,
        UpdatedAt:   updatedAt,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


 
 func (app *application) fetchPosts() ([]Post, error) {
    db, err := sql.Open("postgres", config.GetDBConfig())
    if err != nil {
        return nil, fmt.Errorf("database connection error: %v", err)
    }
    defer db.Close()

    rows, err := db.Query(`
        SELECT p.id, p.user_id, p.title, p.content,
               p.created_at, p.likes_count, p.views_count, 
               p.comments, p.updated_at, p.tags,
               u.username, u.icon_index, u.color_index
        FROM posts p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC`)
    if err != nil {
        return nil, fmt.Errorf("query error: %v", err)
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        var commentsStr []byte

        err := rows.Scan(
            &post.ID,
            &post.UserID,
            &post.Title,
            &post.Content,
            &post.CreatedAt,
            &post.LikesCount,
            &post.ViewsCount,
            &commentsStr,
            &post.UpdatedAt,
            pq.Array(&post.Tags),
            &post.Username,
            &post.IconIndex,
            &post.ColorIndex,
        )
        if err != nil {
            return nil, fmt.Errorf("row scan error: %v", err)
        }

        // Parse comments array
        if len(commentsStr) > 0 {
            if err := json.Unmarshal(commentsStr, &post.Comments); err != nil {
                return nil, fmt.Errorf("comments parsing error: %v", err)
            }

            // Fetch user details for each comment
            for i := range post.Comments {
                var username string
                var iconIndex, colorIndex int
                err := db.QueryRow(`
                    SELECT username, icon_index, color_index 
                    FROM users WHERE id = $1`, 
                    post.Comments[i].UserID).Scan(&username, &iconIndex, &colorIndex)
                if err != nil {
                    continue // Skip if user details are not found
                }
                post.Comments[i].Username = username
                post.Comments[i].IconIndex = iconIndex
                post.Comments[i].ColorIndex = colorIndex
            }
        }

        posts = append(posts, post)
    }

    return posts, nil
}


// Login handler
// Login handler update
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid request Body", http.StatusBadRequest)
        return
    }

    db, err := sql.Open("postgres", config.GetDBConfig())
    if err != nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    var user User
    var likes []string
    err = db.QueryRow("SELECT id, username, email, password_hash, icon_index, color_index, likes FROM users WHERE email = $1", creds.Email).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.IconIndex,
        &user.ColorIndex,
        pq.Array(&likes),
    )
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
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
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    response := struct {
        Token string `json:"token"`
        User  struct {
            ID         string   `json:"id"`
            Username   string   `json:"username"`
            Email      string   `json:"email"`
            IconIndex  int      `json:"icon_index"`
            ColorIndex int      `json:"color_index"`
            Likes      []string `json:"likes"`
        } `json:"user"`
    }{
        Token: tokenString,
        User: struct {
            ID         string   `json:"id"`
            Username   string   `json:"username"`
            Email      string   `json:"email"`
            IconIndex  int      `json:"icon_index"`
            ColorIndex int      `json:"color_index"`
            Likes      []string `json:"likes"`
        }{
            ID:         user.ID,
            Username:   user.Username,
            Email:      user.Email,
            IconIndex:  user.IconIndex,
            ColorIndex: user.ColorIndex,
            Likes:      likes,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(response)
}

// Register handler update
func (app *application) Register(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PlainPassword), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error processing registration", http.StatusInternalServerError)
        return
    }
    user.PasswordHash = string(hashedPassword)

    db, err := sql.Open("postgres", config.GetDBConfig())
    if err != nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    likes := []string{}
    err = db.QueryRow(`
        INSERT INTO users (id, username, email, password_hash, icon_index, color_index, likes)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, username, email, icon_index, color_index, likes`,
        user.ID,
        user.Username,
        user.Email,
        user.PasswordHash,
        user.IconIndex,
        user.ColorIndex,
        pq.Array(likes),
    ).Scan(&user.ID, &user.Username, &user.Email, &user.IconIndex, &user.ColorIndex, pq.Array(&likes))

    if err != nil {
        http.Error(w, "Error creating user", http.StatusInternalServerError)
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
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    response := struct {
        Token string `json:"token"`
        User  struct {
            ID         string   `json:"id"`
            Username   string   `json:"username"`
            Email      string   `json:"email"`
            IconIndex  int      `json:"icon_index"`
            ColorIndex int      `json:"color_index"`
            Likes      []string `json:"likes"`
        } `json:"user"`
    }{
        Token: tokenString,
        User: struct {
            ID         string   `json:"id"`
            Username   string   `json:"username"`
            Email      string   `json:"email"`
            IconIndex  int      `json:"icon_index"`
            ColorIndex int      `json:"color_index"`
            Likes      []string `json:"likes"`
        }{
            ID:         user.ID,
            Username:   user.Username,
            Email:      user.Email,
            IconIndex:  user.IconIndex,
            ColorIndex: user.ColorIndex,
            Likes:      likes,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (app *application) UpdateUserLikes(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")
    
    var requestBody struct {
        Field string `json:"field"`
        Value string `json:"value"`
    }
 
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
 
    db, err := sql.Open("postgres", config.GetDBConfig())
    if err != nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()
 
    var user User
    var likes []string
    err = db.QueryRow(`
        UPDATE users 
        SET likes = array_append(likes, $1)
        WHERE id = $2
        RETURNING id, username, email, icon_index, color_index, likes`,
        requestBody.Value, userID,
    ).Scan(&user.ID, &user.Username, &user.Email, &user.IconIndex, &user.ColorIndex, pq.Array(&likes))
 
    if err != nil {
        http.Error(w, "Error updating user", http.StatusInternalServerError)
        return
    }
 
    response := struct {
        User struct {
            ID         string   `json:"id"`
            Username   string   `json:"username"`
            Email      string   `json:"email"`
            IconIndex  int      `json:"icon_index"`
            ColorIndex int      `json:"color_index"`
            Likes      []string `json:"likes"`
        } `json:"user"`
    }{
        User: struct {
            ID         string   `json:"id"`
            Username   string   `json:"username"`
            Email      string   `json:"email"`
            IconIndex  int      `json:"icon_index"`
            ColorIndex int      `json:"color_index"`
            Likes      []string `json:"likes"`
        }{
            ID:         user.ID,
            Username:   user.Username,
            Email:      user.Email,
            IconIndex:  user.IconIndex,
            ColorIndex: user.ColorIndex,
            Likes:      likes,
        },
    }
 
    w.Header().Set("Content-Type", "application/json")
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

    currentTime := time.Now()
    post.CreatedAt = currentTime
    post.UpdatedAt = currentTime
    post.LikesCount = 0
    post.ViewsCount = 0


    if post.Comments == nil {
        post.Comments = []Comment{}
    }
    if post.Tags == nil {
        post.Tags = []string{}
    }

    commentsJSON, err := json.Marshal(post.Comments)
    if err != nil {
        fmt.Printf("Error marshaling comments: %v\n", err)
        http.Error(w, "Error processing comments", http.StatusInternalServerError)
        return
    }

    fmt.Printf("Inserting post with values: %+v\n", post)

    err = db.QueryRow(`
        INSERT INTO posts (
            id,           
            user_id,      
            title,        
            content,      
            created_at,   
            updated_at,   
            likes_count,  
            views_count,  
            comments,     
            tags          
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

func (app *application) StreamTags(w http.ResponseWriter, r *http.Request) {
    
    // Set necessary headers for SSE
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "SSE not supported", http.StatusInternalServerError)
        return
    }

    // Fetch initial tags
    tags, err := app.fetchTags()
    if err != nil {
        http.Error(w, "Error fetching initial tags", http.StatusInternalServerError)
        return
    }

    // Send initial tags to the client
    data, _ := json.Marshal(tags)
    fmt.Fprintf(w, "data: %s\n\n", data)
    flusher.Flush()

    // Create a ticker to periodically fetch updates
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            tags, err := app.fetchTags()
            if err != nil {
                // Log and continue if there's an error fetching tags
                fmt.Printf("Error fetching tags: %v\n", err)
                continue
            }

            data, _ := json.Marshal(tags)
            fmt.Fprintf(w, "data: %s\n\n", data)
            flusher.Flush()

        case <-r.Context().Done():
            // Stop streaming if the client disconnects
            return
        }
    }
}

func (app *application) fetchTags() ([]Tag, error) {
    db, err := sql.Open("postgres", config.GetDBConfig())
    if err != nil {
        return nil, fmt.Errorf("database connection error: %v", err)
    }
    defer db.Close()

    rows, err := db.Query(`SELECT id, text, color, searches FROM tags`)
    if err != nil {
        return nil, fmt.Errorf("query error: %v", err)
    }
    defer rows.Close()

    var tags []Tag
    for rows.Next() {
        var tag Tag
        err := rows.Scan(&tag.ID, &tag.Text, &tag.Color, &tag.Searches)
        if err != nil {
            return nil, fmt.Errorf("row scan error: %v", err)
        }
        tags = append(tags, tag)
    }

    return tags, nil
}


func (app *application) AddComment(w http.ResponseWriter, r *http.Request) {
    postID := chi.URLParam(r, "id")
    fmt.Printf("AddComment endpoint hit for post ID: %s\n", postID)
 
    // Parse request body
    var requestBody struct {
        Field string      `json:"field"`
        Value Comment     `json:"value"`
        PostID string    `json:"postId"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
        fmt.Printf("Error decoding request: %v\n", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
 
    db, err := sql.Open("postgres", config.GetDBConfig())
    if err != nil {
        fmt.Printf("Database error: %v\n", err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer db.Close()
 
    // Get existing post and comments
    var post Post
    var commentsJSON []byte

    err = db.QueryRow(`
        SELECT p.id, p.user_id, p.title, p.content, p.created_at, 
            p.likes_count, p.views_count, p.comments, p.updated_at, 
            p.tags, u.username, u.icon_index, u.color_index
        FROM posts p
        JOIN users u ON p.user_id = u.id  
        WHERE p.id = $1`, postID).Scan(
            &post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt,
            &post.LikesCount, &post.ViewsCount, &commentsJSON, &post.UpdatedAt,
            pq.Array(&post.Tags), &post.Username, &post.IconIndex, &post.ColorIndex)
 
    if err != nil {
        fmt.Printf("Error fetching post: %v\n", err)
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }
 
    // Parse existing comments
    var comments []Comment
    if len(commentsJSON) > 0 {
        if err := json.Unmarshal(commentsJSON, &comments); err != nil {
            fmt.Printf("Error parsing comments: %v\n", err)
            http.Error(w, "Error parsing comments", http.StatusInternalServerError)
            return
        }
    }
 
    // Add new comment
    comments = append(comments, requestBody.Value)
    updatedCommentsJSON, err := json.Marshal(comments)
    if err != nil {
        fmt.Printf("Error marshaling comments: %v\n", err)
        http.Error(w, "Error processing comments", http.StatusInternalServerError)
        return
    }
 
    // Update post
    _, err = db.Exec(`
        UPDATE posts 
        SET comments = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2`,
        updatedCommentsJSON, postID)
 
    if err != nil {
        fmt.Printf("Error updating post: %v\n", err)
        http.Error(w, "Error updating post", http.StatusInternalServerError)
        return
    }
 
    post.Comments = comments
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
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
        INSERT INTO tags (id, text, color, searches)
        VALUES ($1, $2, $3, $4)
        RETURNING id`,
        tag.ID,      
        tag.Text,
        tag.Color,
        tag.Searches,
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

func (app *application) GetUserLikes(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")
    fmt.Printf("GetUserLikes called for userID: %s\n", userID)

    db, err := sql.Open("postgres", config.GetDBConfig())
    if err != nil {
        http.Error(w, "Database connection error", http.StatusInternalServerError)
        return
    }
    defer db.Close()
 
    var likes []string
    err = db.QueryRow("SELECT likes FROM users WHERE id = $1", userID).Scan(pq.Array(&likes))
    if err != nil {
        http.Error(w, "Error fetching likes", http.StatusInternalServerError)
        return
    }
 
    response := struct {
        Likes []string `json:"likes"`
    }{
        Likes: likes,
    }
 
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
 }
 