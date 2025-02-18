package queries

const (
	GetUserByIDQuery = `
		SELECT id, username, password, access_level 
		FROM users 
		WHERE id = $1;
	`
	GetUserByUsernameQuery = `
		SELECT id, username, password, access_level 
		FROM users 
		WHERE username = $1;
	`
	CreateUserQuery = `
		INSERT INTO users (id, username, password, access_level)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	UpdateUserQuery = `
		UPDATE users
		SET username = $2, password = $3, access_level = $4
		WHERE id = $1;
	`
	DeleteUserQuery = `
		DELETE FROM users
		WHERE id = $1;
	`

	GetRefreshTokenQuery = `
		SELECT refresh_token, user_id, created_at, expired_at 
		FROM refresh_sessions 
		WHERE refresh_token = $1;
	`
	CreateRefreshTokenQuery = `
		INSERT INTO refresh_sessions (refresh_token, user_id, created_at, expired_at)
		VALUES ($1, $2, $3, $4)
		RETURNING refresh_token;
	`
	UpdateRefreshTokenQuery = `
		UPDATE refresh_sessions
		SET refresh_token = $2, user_id = $3, expired_at = $4
		WHERE refresh_token = $1;
	`
	DeleteRefreshTokenQuery = `
		DELETE FROM refresh_sessions
		WHERE refresh_token = $1;
	`
)
