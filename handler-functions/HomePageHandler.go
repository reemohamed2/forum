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

func Homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" { // Check if the requested path is the root URL
		filename := "Templates/home.html"
		t, err := template.ParseFiles(filename) // Parse the index.html template file
		if err != nil {
			fmt.Println("error parsing:", err)
			return
		}
		// open db to getPosts
		db, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			fmt.Println("Error opening database:", err)
		}
		defer db.Close()

		var Data structs.HomepageData

		// getPosts to display them
		Data.Posts, err = getFunctions.GetPosts(db)
		if err != nil {
			fmt.Println("Error in GetPosts func: ", err)
		}

		// initially assume user is guest
		Data.CurrentUser.Username = "guest"
		Data.CurrentUser.Gender = "nil"

		// here we check if cookie already exists
		cookie, err := r.Cookie("session_token")
		// if there is no cookie then assume user is user guest
		if err != nil {
			err = nil // to make cookie error nil and notice the execute error
			err = t.ExecuteTemplate(w, "home.html", Data)
			if err != nil {
				fmt.Println("error executing:", err)
			}
		} else {
			// if cookie exists then do the following
			username, err := getFunctions.GetUsernameFromToken(db, cookie.Value)
			// change username to current userName
			if err != nil || username == "" {
				err = nil
				err = t.ExecuteTemplate(w, "home.html", Data)
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
				err = t.ExecuteTemplate(w, "home.html", Data)
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
