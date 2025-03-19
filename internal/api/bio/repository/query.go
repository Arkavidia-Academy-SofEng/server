package bioRepository

const (
	queryCreateExperience = `
    INSERT INTO experiences (
        id, user_id, image_url, job_title, job_location, skill_used, start_date, end_date, description, created_at, updated_at
    ) VALUES (
        :id, :user_id, :image_url, :job_title, :job_location, :skill_used, :start_date, :end_date, :description, :created_at, :updated_at
    )`

	queryGetExperienceByID = `
    SELECT id, user_id, image_url, job_title, skill_used, start_date, end_date, description, created_at, updated_at
    FROM experiences
    WHERE id = ?
    `

	queryGetExperiencesByUserID = `
    SELECT id, user_id, image_url, job_title, skill_used, start_date, end_date, description, created_at, updated_at
    FROM experiences
    WHERE user_id = ?
    ORDER BY start_date DESC
    `

	queryUpdateExperience = `
    UPDATE experiences
    SET image_url = :image_url,
        job_title = :job_title,
        skill_used = :skill_used,
        start_date = :start_date,
        end_date = :end_date,
        description = :description,
        updated_at = :updated_at
    WHERE id = :id
    `

	queryDeleteExperience = `
    DELETE FROM experiences
    WHERE id = ?
    `

	queryDeleteExperiencesByUserID = `
    DELETE FROM experiences
    WHERE user_id = ?
    `
)
