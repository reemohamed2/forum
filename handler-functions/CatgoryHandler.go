package handlerfuncitons

import (
	"database/sql"
	"fmt"
	getFunctions "form/get-functions"
	"form/structs"
	"html/template"
	"net/http"
	"strconv"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	filename := "Templates/category.html"
	t, err := template.ParseFiles(filename) // Parse the category.html template file
	if err != nil {
		fmt.Println("error parsing:", err)
		return
	}

	var FilteredData structs.HomepageData

	// Open db to getPosts
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Get posts to display them
	FilteredData.Posts, err = getFunctions.GetPosts(db)
	if err != nil {
		fmt.Println("Error in GetPosts func: ", err)
		return
	}

	// Initially assume user is guest
	FilteredData.CurrentUser.Username = "guest"
	FilteredData.CurrentUser.Gender = "nil"

	if r.Method == "POST" {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			InternalServerError(w, r)
			return
		}
		categories := r.Form["category"]
		if !validateCategories(categories) {
			BadRequest(w,r)
			return
		}
		FilteredData.Posts = filterPosts(FilteredData.Posts, categories) // Corrected line
	}
	// Check if cookie already exists
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// No cookie, assume user is guest
		err = nil // Reset error to avoid interfering with template execution
		err = t.ExecuteTemplate(w, "category.html", FilteredData)
		if err != nil {
			fmt.Println("error executing:", err)
		}
		return
	}
	// If cookie exists, get the username from the token
	username, err := getFunctions.GetUsernameFromToken(db, cookie.Value)
	if err != nil || username == "" {
		// Invalid token or username, assume guest
		err = nil
		err = t.ExecuteTemplate(w, "category.html", FilteredData)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	// Valid token, get user details
	user, err := getFunctions.GetUser(username, db)
	if err != nil {
		fmt.Println("Error in getUser func: ", err)
		return
	}
	userids, err := getFunctions.GetUserID(db, username)
	if err != nil {
		fmt.Println("Error in GetUserID func: ", err)
		return
	}
	userid, _ := strconv.Atoi(userids)
	FilteredData.Posts, err = getFunctions.DisLikedpostsdis(db, userid, FilteredData.Posts)
	if err != nil {
		fmt.Println("Error in GetPosts func: ", err)
	}
	FilteredData.Posts, err = getFunctions.Likedpostsdis(db, userid, FilteredData.Posts)
	if err != nil {
		fmt.Println("Error in GetPosts func: ", err)
	}

	FilteredData.CurrentUser = user
	err = t.ExecuteTemplate(w, "category.html", FilteredData)
	if err != nil {
		fmt.Println("error executing:", err)
		return
	}
}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Helper function to check if all categories are present in post categories
func containsAllCategories(postCategories []string, categories []string) bool {
	for _, category := range categories {
		if !Contains(postCategories, category) { // Corrected function call
			return false
		}
	}
	return true
}

// Function to filter posts based on categories
func filterPosts(posts []structs.Post, categories []string) []structs.Post {
	var filteredPosts []structs.Post

	// Iterate through each post
	for _, post := range posts {
		// Check if the post categories contain all the required categories
		if containsAllCategories(post.Category, categories) {
			filteredPosts = append(filteredPosts, post)
		}
	}

	return filteredPosts
}

func Category(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/category" { // Check if the requested path is the root URL
		filename := "Templates/category.html"
		t, err := template.ParseFiles(filename) // Parse the index.html template file
		if err != nil {
			fmt.Println("error parsing:", err)
			return
		}
		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			fmt.Println("error opening database:", err)
		}
		defer db.Close()
		var Data structs.HomepageData

		Data.Posts, err = getFunctions.GetPosts(db)
		if err != nil {
			fmt.Println("Error in GetPosts func: ", err)
		}
		Data.CurrentUser.Username = "guest"
		Data.CurrentUser.Gender = "nil"
		cookie, err := r.Cookie("session_token")
		if err != nil {
			err = t.Execute(w, Data)
			if err != nil {
				fmt.Println("error executing:", err)
				return
			}
			return
		} else {
			username, err := getFunctions.GetUsernameFromToken(db, cookie.Value)
			// change username to current userName
			if err != nil || username == "" {
				err = nil
				err = t.ExecuteTemplate(w, "category.html", Data)
				if err != nil {
					fmt.Println(err)
				}
				return
			} else {
				user, err := getFunctions.GetUser(username, db)
				if err != nil {
					fmt.Println("Error in getUser func: ", err)
					return
				}
				userids, err := getFunctions.GetUserID(db, username)
				if err != nil {
					fmt.Println("Error in GetUserID func: ", err)
					return
				}
				userid, _ := strconv.Atoi(userids)
				Data.Posts, err = getFunctions.DisLikedpostsdis(db,userid,Data.Posts)
				if err != nil {
					fmt.Println("Error in GetPosts func: ", err)
				}
				Data.Posts, err = getFunctions.Likedpostsdis(db,userid,Data.Posts)
				if err != nil {
					fmt.Println("Error in GetPosts func: ", err)
				}
				Data.CurrentUser = user
				err = t.ExecuteTemplate(w, "category.html", Data)
				// Execute the template and write the output to the response writer
				if err != nil {
					fmt.Println("error executing:", err)
					return
				}
			}
		}
	} else { // Handle 404 error for other URLs
		NotFoundHandler(w, r)
	}
}

func isValidCategory(category string) bool {
	validCategories := map[string]struct{}{
		"Sport":     {},
		"Gaming":    {},
		"Art":       {},
		"Education": {},
		"Food":      {},
	}
	_, exists := validCategories[category]
	return exists
}

// validateCategories checks if all elements in the slice are valid categories
func validateCategories(categories []string) bool {
	for _, category := range categories {
		if !isValidCategory(category) {
			return false
		}
	}
	return true
}