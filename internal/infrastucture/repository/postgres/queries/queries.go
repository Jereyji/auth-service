package queries

const (
	QueryGetUserByID = `
		SELECT id, email, name, password 
		FROM users 
		WHERE id = $1;
	`
	QueryGetUserByEmail = `
		SELECT id, email, name, password 
		FROM users 
		WHERE email = $1;
	`
	QueryCreateUser = `
		INSERT INTO users (id, email, name, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	QueryUpdateUser = `
		UPDATE users
		SET email = $2, name = $3, password = $4
		WHERE id = $1;
	`
	QueryDeleteUser = `
		DELETE FROM users
		WHERE id = $1;
	`

	QueryGetRefreshToken = `
		SELECT refresh_token, user_id, created_at, expired_at 
		FROM refresh_sessions 
		WHERE refresh_token = $1;
	`
	QueryCreateRefreshToken = `
		INSERT INTO refresh_sessions (refresh_token, user_id, created_at, expired_at)
		VALUES ($1, $2, $3, $4)
		RETURNING refresh_token;
	`
	QueryUpdateRefreshToken = `
		UPDATE refresh_sessions
		SET refresh_token = $2, user_id = $3, expired_at = $4
		WHERE refresh_token = $1;
	`
	QueryDeleteRefreshToken = `
		DELETE FROM refresh_sessions
		WHERE refresh_token = $1;
	`
)
