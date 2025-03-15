package authRepository

const (
	queryCreateUser = `
INSERT INTO Users (id, email, username, password)
VALUES (:id, :email, :username, :password)`

	queryGetById = `
SELECT id, username, email, password
FROM Users 
WHERE id = :id`

	queryGetByEmail = `
SELECT id, email, username, password
FROM Users
WHERE email = :email
LIMIT 1`

	queryUpdateUser = `
UPDATE Users 
SET username = :username, password = :password 
WHERE id = :id`

	queryDeleteUser = `
DELETE FROM Users
WHERE id = :id`
)
