package structs

type Post struct{
	ID int
	Username string
	Title string
	Content string
	CreatedAt string
	Category []string 
	Like int
	Dislike int
	Comment int
	Gender string
	Likedbyuser bool
	Dislikedbyuser bool
}

type Comment struct{
	CommentID int
	PostID int
	UserID int 
	Username string
	Content string
	CreatedAt string
	Like int
	Dislike int
	Comment int
	Gender string
	Likedbyuser bool
	Dislikedbyuser bool
}

type CurrentUser struct {
	Username string
	Gender   string
}

type HomepageData struct{
	CurrentUser CurrentUser
	Posts []Post
}

type Profiledata struct{
	Email string
	Username string
	Gender string
	RequestedPosts []Post
}

type CommentPost struct {
	CurrentUser CurrentUser
	Post Post
	Comments []Comment
}

type PostData struct {
	MainPost Post
	Comments []Comment
}