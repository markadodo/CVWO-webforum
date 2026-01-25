package database

import (
	"database/sql"
)

// database creation
func InitDB(db *sql.DB) error {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
 		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL,
		last_active TIMESTAMPTZ NOT NULL
		);
	`

	topicTable := `
	CREATE TABLE IF NOT EXISTS topics (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		created_by INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL,
		FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET DEFAULT
		);
	`
	topicFTSColumn := `
	ALTER TABLE topics 
	ADD COLUMN IF NOT EXISTS document tsvector;
	`
	topicFTSIdx := `
	CREATE INDEX IF NOT EXISTS topics_document_idx 
	ON topics USING GIN(document);
	`

	postTable := `
	CREATE TABLE IF NOT EXISTS posts(
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		topic_id INTEGER NOT NULL,
		likes INTEGER NOT NULL DEFAULT 0,
		dislikes INTEGER NOT NULL DEFAULT 0,
		is_edited INTEGER NOT NULL DEFAULT 0,
		views INTEGER NOT NULL DEFAULT 0,
		popularity INTEGER NOT NULL DEFAULT 0,
		created_by INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL,
		FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE,
		FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET DEFAULT
        );
    `
	postFTSColumn := `
	ALTER TABLE posts 
	ADD COLUMN IF NOT EXISTS document tsvector;
	`
	postFTSIdx := `
	CREATE INDEX IF NOT EXISTS posts_document_idx 
	ON posts USING GIN(document);
	`
	postReactionTable := `
	CREATE TABLE IF NOT EXISTS posts_reactions(
		id SERIAL PRIMARY KEY,
		post_id INTEGER NOT NULL,
		user_id INTEGER,
		reaction BOOLEAN NOT NULL,
		UNIQUE(post_id, user_id),
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		);
	`

	commentTable := `
	CREATE TABLE IF NOT EXISTS comments(
		id SERIAL PRIMARY KEY,
		description TEXT NOT NULL,
		likes INTEGER NOT NULL DEFAULT 0,
		dislikes INTEGER NOT NULL DEFAULT 0,
		is_edited INTEGER NOT NULL DEFAULT 0,
		post_id INTEGER NOT NULL,
		parent_comment_id INTEGER,
		created_by INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET DEFAULT,
		FOREIGN KEY (parent_comment_id) REFERENCES comments(id) ON DELETE CASCADE
		);
	`
	commentReactionTable := `
	CREATE TABLE IF NOT EXISTS comments_reactions(
		id SERIAL PRIMARY KEY,
		comment_id INTEGER NOT NULL,
		user_id INTEGER,
		reaction BOOLEAN NOT NULL,
		UNIQUE(comment_id, user_id),
		FOREIGN KEY (comment_id) REFERENCES comments(id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		);
	`
	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	FTSTriggerFunction := `
	CREATE OR REPLACE FUNCTION fts_trigger_handler() RETURNS trigger AS $$
	BEGIN 
		new.document := 
			to_tsvector('english', coalesce(new.title, '')) || 
			to_tsvector('english', coalesce(new.description, ''));

		return new;
	END;
	$$
	LANGUAGE plpgsql
	`

	insertPostReactionTriggerFunction := `
	CREATE OR REPLACE FUNCTION insert_post_reaction_handler() RETURNS trigger as $$
	BEGIN
		UPDATE posts SET
			likes = likes + CASE WHEN new.reaction = TRUE THEN 1 ELSE 0 END,
			dislikes = dislikes + CASE WHEN new.reaction = FALSE THEN 1 ELSE 0 END,
			popularity = (likes + CASE WHEN new.reaction = TRUE THEN 1 ELSE 0 END) * 10
                            - (dislikes + CASE WHEN new.reaction = FALSE THEN 1 ELSE 0 END) * 5
                            + views
		WHERE posts.id = new.post_id;

		return new;
	END;
	$$
	LANGUAGE plpgsql
	`
	deletePostReactionTriggerFunction := `
	CREATE OR REPLACE FUNCTION delete_post_reaction_handler() RETURNS trigger as $$
	BEGIN
		UPDATE posts SET
			likes = likes - CASE WHEN old.reaction = TRUE THEN 1 ELSE 0 END,
			dislikes = dislikes - CASE WHEN old.reaction = FALSE THEN 1 ELSE 0 END,
			popularity = (likes - CASE WHEN old.reaction = TRUE THEN 1 ELSE 0 END) * 10
                            - (dislikes - CASE WHEN old.reaction = FALSE THEN 1 ELSE 0 END) * 5
                            + views
		WHERE posts.id = old.post_id;

		return old;
	END;
	$$
	LANGUAGE plpgsql
	`

	insertCommentReactionTriggerFunction := `
	CREATE OR REPLACE FUNCTION insert_comment_reaction_handler() RETURNS trigger as $$
	BEGIN
		UPDATE comments SET
			likes = likes + CASE WHEN new.reaction = TRUE THEN 1 ELSE 0 END,
			dislikes = dislikes + CASE WHEN new.reaction = FALSE THEN 1 ELSE 0 END
		WHERE comments.id = new.comment_id;

		return new;
	END;
	$$
	LANGUAGE plpgsql
	`
	deleteCommentReactionTriggerFunction := `
	CREATE OR REPLACE FUNCTION delete_comment_reaction_handler() RETURNS trigger as $$
	BEGIN
		UPDATE comments SET
			likes = likes - CASE WHEN old.reaction = TRUE THEN 1 ELSE 0 END,
			dislikes = dislikes - CASE WHEN old.reaction = FALSE THEN 1 ELSE 0 END
		WHERE comments.id = old.comment_id;

		return old;
	END;
	$$
	LANGUAGE plpgsql
	`

	resetTopicFTSTrigger := `
	DROP TRIGGER IF EXISTS topics_ai_au ON topics;
	`
	topicFTSTrigger := `
	CREATE TRIGGER topics_ai_au
	BEFORE INSERT OR UPDATE ON topics
	FOR EACH ROW
	EXECUTE FUNCTION fts_trigger_handler();
	`

	resetPostFTSTrigger := `
	DROP TRIGGER IF EXISTS posts_ai_au ON posts;
	`
	postFTSTrigger := `
	CREATE TRIGGER posts_ai_au
	BEFORE INSERT OR UPDATE ON posts
	FOR EACH ROW
	EXECUTE FUNCTION fts_trigger_handler();
	`

	resetInsertPostReactionTrigger := `
	DROP TRIGGER IF EXISTS posts_reactions_ai ON posts_reactions;
	`
	resetDeletePostReactionTrigger := `
	DROP TRIGGER IF EXISTS posts_reactions_ad ON posts_reactions;
	`
	insertPostReactionTrigger := `
	CREATE TRIGGER posts_reactions_ai
	AFTER INSERT ON posts_reactions
	FOR EACH ROW
	EXECUTE FUNCTION insert_post_reaction_handler();
	`
	deletePostReactionTrigger := `
	CREATE TRIGGER posts_reactions_ad
	AFTER DELETE ON posts_reactions
	FOR EACH ROW
	EXECUTE FUNCTION delete_post_reaction_handler();
	`

	resetInsertCommentReactionTrigger := `
	DROP TRIGGER IF EXISTS comments_reactions_ai ON comments_reactions;
	`
	resetDeleteCommentReactionTrigger := `
	DROP TRIGGER IF EXISTS comments_reactions_ad ON comments_reactions;
	`
	insertCommentReactionTrigger := `
	CREATE TRIGGER comments_reactions_ai
	AFTER INSERT ON comments_reactions
	FOR EACH ROW
	EXECUTE FUNCTION insert_comment_reaction_handler();
	`
	deleteCommentReactionTrigger := `
	CREATE TRIGGER comments_reactions_ad
	AFTER DELETE ON comments_reactions
	FOR EACH ROW
	EXECUTE FUNCTION delete_comment_reaction_handler();
	`
	//////////////////////////////////////////////////////////////////////////////////////////////////////////

	insertsystemUser := `
	INSERT INTO users (id, username, password_hash, created_at, last_active)
	VALUES (0, 'deleted_users', '!', NOW(), NOW())
	ON CONFLICT(id) DO NOTHING;
	`

	tables := []string{
		userTable,
		topicTable,
		topicFTSColumn,
		topicFTSIdx,
		postTable,
		postFTSColumn,
		postFTSIdx,
		postReactionTable,
		commentTable,
		commentReactionTable,
	}

	triggers := []string{
		FTSTriggerFunction,
		insertPostReactionTriggerFunction,
		deletePostReactionTriggerFunction,
		insertCommentReactionTriggerFunction,
		deleteCommentReactionTriggerFunction,
		resetTopicFTSTrigger,
		topicFTSTrigger,
		resetPostFTSTrigger,
		postFTSTrigger,
		resetInsertPostReactionTrigger,
		resetDeletePostReactionTrigger,
		resetInsertCommentReactionTrigger,
		resetDeleteCommentReactionTrigger,
		insertPostReactionTrigger,
		deletePostReactionTrigger,
		insertCommentReactionTrigger,
		deleteCommentReactionTrigger,
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

	if _, err := db.Exec(insertsystemUser); err != nil {
		return err
	}

	return nil
}
