package main

import (
	"backend/cmd/api/config"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// data type of the post stored in the database
type Post struct {
    ID              string    `json:"id"`
    UserID          string    `json:"user_id"`
    Title           string    `json:"title"`
    Content         string    `json:"content"`
    PictureURL      string    `json:"picture_url,omitempty"`
    CreatedAt       time.Time `json:"created_at"`
    LikesCount      int       `json:"likes_count"`
    ViewsCount      int       `json:"views_count"`
    DiscussionThread string   `json:"discussion_thread,omitempty"`
    Comments        []Comment `json:"comments"`
    UpdatedAt       time.Time `json:"updated_at"`
    Tags            []string  `json:"tags"`  
}

// data type of the comment under each post
type Comment struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

// data type of users stored in the database
type User struct {
    ID           string    `json:"id"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}

// data type for tags
type Tag struct {
    ID       string `json:"id"`
    Text     string `json:"text"`
    Color    string `json:"color"`
    Searches int    `json:"searches"`
}

// data type received from the front end
type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// data type of JWT token
type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}
// default handler for testing
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from %s", app.Domain)
}

// API endpoint for getting posts
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
        SELECT id, user_id, title, content, picture_url, 
               created_at, likes_count, views_count, 
               discussion_thread, comments, updated_at, tags 
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
        var pictureURL, discussionThread sql.NullString
        var comments []byte // for JSONB data

        err := rows.Scan(
            &post.ID,
            &post.UserID,
            &post.Title,
            &post.Content,
            &pictureURL,
            &post.CreatedAt,
            &post.LikesCount,
            &post.ViewsCount,
            &discussionThread,
            &comments,
            &post.UpdatedAt,
            pq.Array(&post.Tags),  // Use pq.Array for scanning PostgreSQL array
        )
        if err != nil {
            fmt.Println("Row scan error:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Handle nullable fields
        if pictureURL.Valid {
            post.PictureURL = pictureURL.String
        }
        if discussionThread.Valid {
            post.DiscussionThread = discussionThread.String
        }

        // Parse JSONB comments
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

    if err := json.NewEncoder(w).Encode(posts); err != nil {
        fmt.Println("JSON encoding error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Printf("Returning %d posts\n", len(posts))
}
// API endpoint for logging in
func (app *application) Login(w http.ResponseWriter, r *http.Request) {

// parsing the incoming JSON request and assigning it to the creds variable
	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
				fmt.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}
		fmt.Printf("Received credentials - Email: %s, Password: %s\n", creds.Email, creds.Password)

	// connect to the database
	connStr := config.GetDBConfig()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
				fmt.Printf("Database connection error: %v\n", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
		fmt.Println("Database connection successful")

	// fetch user from the database and put then into our user variable
	var user User
	err = db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = $1", creds.Email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
				fmt.Printf("User lookup error: %v\n", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
		fmt.Printf("Found user with ID: %s, Email: %s\n", user.ID, user.Email)

	// Compare passwords by converting both to a byte slice
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
				fmt.Printf("Password comparison failed: %v\n", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return  // Added missing return statement
	}
		fmt.Println("Password verification successful")

	/* generate JWT token, claims are customised information while
	Registered claims are standardized information */
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims {
		UserID: user.ID, 
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		}, 
	}

	// generation of token with secret key to prevent tampering of token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(app.JwtSecret))
	if err != nil {
				fmt.Printf("Token generation error: %v\n", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
		fmt.Println("JWT token generated successfully")

	// send response 
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString, 
	})
}

// API endpoint for getting tags
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

    // Query all tags
    rows, err := db.Query(`
        SELECT id, text, color, searches 
        FROM tags`)  
    if err != nil {
        fmt.Println("Query error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var tags []Tag
    for rows.Next() {
        var tag Tag
        err := rows.Scan(
            &tag.ID,
            &tag.Text,
            &tag.Color,
            &tag.Searches,
        )
        if err != nil {
            fmt.Println("Row scan error:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        tags = append(tags, tag)
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    if err := json.NewEncoder(w).Encode(tags); err != nil {
        fmt.Println("JSON encoding error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Printf("Returning %d tags\n", len(tags))
}