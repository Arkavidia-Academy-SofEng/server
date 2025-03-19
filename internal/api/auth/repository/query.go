package authRepository

const (
	queryCreateUser = `
    INSERT INTO users (
        id, email, password, name, role, created_at, updated_at, phone_number
    ) VALUES (
        :id, :email, :password, :name, :role, :created_at, :updated_at, :phone_number
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
    SET email = :email,
        password = :password,
        name = :name,
        role = :role,
        profile_picture = :profile_picture,
        banner_picture=:banner_picture,
        is_premium = :is_premium,
        premium_until = :premium_until,
        headline = :headline,
        location = :location,
        updated_at = :updated_at,
        phone_number = :phone_number,
    	deleted_at = :deleted_at
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

const (
	queryCreateCompany = `
    INSERT INTO companies (
        id, email, password, name, created_at, updated_at, phone_number
    ) VALUES (
        :id, :email, :password, :name, :created_at, :updated_at, :phone_number
    )`

	queryCheckCompanyEmailExists = `
    SELECT EXISTS (SELECT 1 FROM companies WHERE email = ? AND deleted_at IS NULL)
    `

	queryGetCompanyByID = `
   SELECT id, profile_picture, banner_picture, email, password, name, 
          about_us, industry_types, number_employees, established_date, 
          company_url, required_skill, location, phone_number,
          created_at, updated_at, deleted_at
   FROM companies
   WHERE id = ? AND deleted_at IS NULL
`

	queryUpdateCompanies = `
   UPDATE companies
   SET email = :email,
       password = :password,
       name = :name,
       profile_picture = :profile_picture,
       banner_picture = :banner_picture,
       about_us = :about_us,
       industry_types = :industry_types,
       number_employees = :number_employees,
       established_date = :established_date,
       company_url = :company_url,
       required_skill = :required_skill,
       location = :location,
       phone_number = :phone_number,
       updated_at = :updated_at,
       deleted_at = :deleted_at
   WHERE id = :id
   `

	queryGetCompanyByEmail = `
   SELECT id, email, password, name, profile_picture, banner_picture,
          about_us, industry_types, number_employees, established_date,
          company_url, required_skill, location, phone_number,
          created_at, updated_at, deleted_at
   FROM companies
   WHERE email = ? AND deleted_at IS NULL
   `

	querySoftDeleteCompany = `
   UPDATE companies
   SET deleted_at = ?
   WHERE id = ? AND deleted_at IS NULL
   `

	queryHardDeleteExpiredCompanies = `
   DELETE FROM companies
   WHERE deleted_at IS NOT NULL AND deleted_at <= ?
   `
)
