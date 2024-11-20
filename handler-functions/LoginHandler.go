package handlerfuncitons

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	getFunctions "form/get-functions"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// open database
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	cookie, err := r.Cookie("session_token")
	// if there is no cookie then assume user is user guest
	if err == nil { // if cookie exists then do the following
		_, err := getFunctions.GetUsernameFromToken(db, cookie.Value)
		// change username to current userName
		if err == nil {
			http.Redirect(w, r, "/welcome", http.StatusFound)
			return
		}
	}
	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			InternalServerError(w, r)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		username = strings.TrimSpace(username)
		password = strings.TrimSpace(password)

		// This is to handle the reddirection from the like and also to validate
		if username == "" || password == "" {
			tmpl, err := template.ParseFiles("Templates/login.html")
			if err != nil {
				fmt.Println(err)
				InternalServerError(w, r)
				return
			}
			tmpl.ExecuteTemplate(w, "login.html", nil)
			return
		}

		var storedHashedPassword string
		err = db.QueryRow("SELECT password FROM User WHERE username = ?", username).Scan(&storedHashedPassword) //query pass and store it
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				tmpl, err := template.ParseFiles("Templates/login.html")
				if err != nil {
					fmt.Println(err)
					InternalServerError(w, r)
					return
				}
				tmpl.ExecuteTemplate(w, "login.html", "invalidUser")
				return
			} else {
				log.Printf("Database error: %v", err) // Log the detailed database error
				if err != nil {
					InternalServerError(w, r)
					return
				}
			}
			return
		}
		// Compare the password
		err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
		if err != nil {
			// if password isn't correct then user unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			tmpl, err := template.ParseFiles("Templates/login.html")
			if err != nil {
				fmt.Println(err)
				InternalServerError(w, r)
				return
			}
			tmpl.ExecuteTemplate(w, "login.html", "invalidUser")
			return
		} else {
			// if username and password correct then add cookie and go to home page
			getFunctions.GetUser(username, db)
			sessionToken, err := generateSessionToken() // function to generate token
			if err != nil {
				InternalServerError(w, r)
				return
			}
			if err := setSessionToken(db, username, sessionToken); err != nil {
				InternalServerError(w, r)
				return
			}
			// Assuming tmpl is properly initialized
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Path:    "/",
				Expires: time.Now().Add(24 * time.Hour),
			})
			// if successful then redirect user to main page
			http.Redirect(w, r, "/welcome", http.StatusFound)
		}
	} else { // Display the form
		tmpl, err := template.ParseFiles("Templates/login.html")
		if err != nil {
			InternalServerError(w, r)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			InternalServerError(w, r)
			return
		}
	}
}

func generateSessionToken() (string, error) {
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func setSessionToken(db *sql.DB, username, token string) error {
	// Invalidate old session
	_, err := db.Exec("DELETE FROM sessions WHERE username = ?", username)
	if err != nil {
		return err
	}

	// Set new session token
	_, err = db.Exec("INSERT INTO sessions (username, token) VALUES (?, ?)", username, token)
	return err
}
