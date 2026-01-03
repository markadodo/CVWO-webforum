package database

import (
	"database/sql"
)

// database creation
func InitDB(db *sql.DB) error {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
 		password_hash TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		last_active DATETIME NOT NULL
		);
	`

	topicTable := `
	CREATE TABLE IF NOT EXISTS topics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		created_by INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (created_by) REFERENCES users(id)
		);
	`
	topicFTSTable := `
	CREATE VIRTUAL TABLE IF NOT EXISTS topics_fts USING fts5(
		title,
		description,
		content = "topics",
		content_rowid = "id"
		);
	`

	postTable := `
	CREATE TABLE IF NOT EXISTS posts(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		topic_id INTEGER NOT NULL,
		likes INTEGER NOT NULL DEFAULT 0,
		dislikes INTEGER NOT NULL DEFAULT 0,
		is_edited INTEGER NOT NULL DEFAULT 0,
		views INTEGER NOT NULL DEFAULT 0,
		popularity INTEGER NOT NULL DEFAULT 0,
		created_by INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (topic_id) REFERENCES topics(id),
		FOREIGN KEY (created_by) REFERENCES users(id)
        );
    `
	postFTSTable := `
	CREATE VIRTUAL TABLE IF NOT EXISTS posts_fts USING fts5(
		title,
		description,
		content = "posts",
		content_rowid = "id"
		);
	`
	postReactionTable := `
	CREATE TABLE IF NOT EXISTS posts_reactions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		reaction BOOLEAN NOT NULL,
		UNIQUE(post_id, user_id),
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`

	//attachment

	commentTable := `
	CREATE TABLE IF NOT EXISTS comments(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		likes INTEGER NOT NULL DEFAULT 0,
		dislikes INTEGER NOT NULL DEFAULT 0,
		is_edited INTEGER NOT NULL DEFAULT 0,
		post_id INTEGER NOT NULL,
		parent_comment_id INTEGER,
		created_by INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		FOREIGN KEY (created_by) REFERENCES users(id),
		FOREIGN KEY (parent_comment_id) REFERENCES comments(id)
		);
	`
	commentReactionTable := `
	CREATE TABLE IF NOT EXISTS comments_reactions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		comment_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		reaction BOOLEAN NOT NULL,
		UNIQUE(comment_id, user_id),
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`

	createTopicTrigger := `
	CREATE TRIGGER IF NOT EXISTS topics_ai 
	AFTER INSERT ON topics 
	BEGIN
		INSERT INTO topics_fts(rowid, title, description)
		VALUES(new.id, new.title, new.description);
	END;
	`
	deleteTopicTrigger := `
	CREATE TRIGGER IF NOT EXISTS topics_ad 
	AFTER DELETE ON topics 
	BEGIN
		INSERT INTO topics_fts(topics_fts, rowid, title, description)
		VALUES("delete", old.id, old.title, old.description);
	END;
	`
	updateTopicTrigger := `
	CREATE TRIGGER IF NOT EXISTS topics_au
	AFTER UPDATE ON topics
	BEGIN
		INSERT INTO topics_fts(topics_fts, rowid, title, description)
		VALUES("delete", old.id, old.title, old.description);

		INSERT INTO topics_fts(rowid, title, description)
		VALUES(new.id, new.title, new.description);
	END;
	`

	createPostTrigger := `
	CREATE TRIGGER IF NOT EXISTS posts_ai 
	AFTER INSERT ON posts 
	BEGIN
		INSERT INTO posts_fts(rowid, title, description)
		VALUES(new.id, new.title, new.description);
	END;
	`
	deletePostTrigger := `
	CREATE TRIGGER IF NOT EXISTS posts_ad
	AFTER DELETE ON posts
	BEGIN
		INSERT INTO posts_fts(posts_fts, rowid, title, description)
		VALUES("delete", old.id, old.title, old.description);
	END;
	`
	updatePostTrigger := `
	CREATE TRIGGER IF NOT EXISTS posts_au
	AFTER UPDATE ON posts
	BEGIN
		INSERT INTO posts_fts(posts_fts, rowid, title, description)
		VALUES("delete", old.id, old.title, old.description);

		INSERT INTO posts_fts(rowid, title, description)
		VALUES(new.id, new.title, new.description);
	END;
	`

	insertPostReactionTrigger := `
	CREATE TRIGGER IF NOT EXISTS posts_reactions_ai
	AFTER INSERT ON posts_reactions
	BEGIN
		UPDATE posts SET
			likes = likes + CASE WHEN new.reaction = 1 THEN 1 ELSE 0 END,
			dislikes = dislikes + CASE WHEN new.reaction = 0 THEN 1 ELSE 0 END
		WHERE posts.id = new.post_id;
	END;
	`
	deletePostReactionTrigger := `
	CREATE TRIGGER IF NOT EXISTS posts_reactions_ad
	AFTER DELETE ON posts_reactions
	BEGIN
		UPDATE posts SET
			likes = likes - CASE WHEN old.reaction = 1 THEN 1 ELSE 0 END,
			dislikes = dislikes - CASE WHEN old.reaction = 0 THEN 1 ELSE 0 END
		WHERE posts.id = old.post_id;
	END;
	`

	insertCommentReactionTrigger := `
	CREATE TRIGGER IF NOT EXISTS comments_reactions_ai
	AFTER INSERT ON comments_reactions
	BEGIN
		UPDATE comments SET
			likes = likes + CASE WHEN new.reaction = 1 THEN 1 ELSE 0 END,
			dislikes = dislikes + CASE WHEN new.reaction = 0 THEN 1 ELSE 0 END
		WHERE comments.id = new.comment_id;
	END;
	`
	deleteCommentReactionTrigger := `
	CREATE TRIGGER IF NOT EXISTS comments_reactions_ad
	AFTER DELETE ON comments_reactions
	BEGIN
		UPDATE comments SET
			likes = likes - CASE WHEN old.reaction = 1 THEN 1 ELSE 0 END,
			dislikes = dislikes - CASE WHEN old.reaction = 0 THEN 1 ELSE 0 END
		WHERE comments.id = old.comment_id;
	END;
	`

	tables := []string{
		userTable,
		topicTable,
		topicFTSTable,
		postTable,
		postFTSTable,
		postReactionTable,
		commentTable,
		commentReactionTable,
	}

	triggers := []string{
		createTopicTrigger,
		deleteTopicTrigger,
		updateTopicTrigger,
		createPostTrigger,
		deletePostTrigger,
		updatePostTrigger,
		insertPostReactionTrigger,
		deletePostReactionTrigger,
		insertCommentReactionTrigger,
		deleteCommentReactionTrigger,
	}

	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return err
		}
	}

	for _, trigger := range triggers {
		_, err := db.Exec(trigger)
		if err != nil {
			return err
		}
	}
	return nil
}
