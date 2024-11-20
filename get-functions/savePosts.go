package getFunctions

import (
	"database/sql"
	"form/structs"
)

func SavePost(db *sql.DB, post structs.Post) error {
    query := `INSERT INTO Post (user_id, title, content) VALUES (?, ?, ?)`
    _, err := db.Exec(query, post.Username, post.Title, post.Content)
    return err
}

