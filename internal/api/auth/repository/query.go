package authRepository

const (
	queryCreateUser = `
    INSERT INTO users (
        id, email, password, name, role, profile_picture, 
        is_premium, premium_until, headline, created_at, updated_at
    ) VALUES (
        :id, :email, :password, :name, :role, :profile_picture, 
        :is_premium, :premium_until, :headline, :created_at, :updated_at
    )`

	queryCheckEmailExists = `
    SELECT EXISTS (SELECT 1 FROM users WHERE email = ? AND deleted_at IS NULL)
    `

	queryGetUserByID = `
    SELECT id, email, password, name, role, profile_picture, 
           is_premium, premium_until, headline, created_at, updated_at, deleted_at
    FROM users
    WHERE id = ? AND deleted_at IS NULL
    `

	queryUpdateUser = `
    UPDATE users
    SET name = :name,
        role = :role,
        profile_picture = :profile_picture,
        is_premium = :is_premium,
        premium_until = :premium_until,
        headline = :headline,
        updated_at = :updated_at
    WHERE id = :id
    `

	queryGetUserByEmail = `
    SELECT id, email, password, name, role, profile_picture, 
           is_premium, premium_until, headline, created_at, updated_at, deleted_at
    FROM users
    WHERE email = ? AND deleted_at IS NULL
    `

	querySoftDeleteUser = `
    UPDATE users
    SET deleted_at = ?
    WHERE id = ? AND deleted_at IS NULL
    `

	queryHardDeleteExpiredUsers = `
    DELETE FROM users
    WHERE deleted_at IS NOT NULL AND deleted_at <= ?
    `
)
