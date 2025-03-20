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

	queryCreateEducation = `
    INSERT INTO educations (
        id, image, user_id, title_degree, institutional_name, start_date, end_date, description, created_at, updated_at
    ) VALUES (
        :id, :image, :user_id, :title_degree, :institutional_name, :start_date, :end_date, :description, :created_at, :updated_at
    )`

	queryGetEducationByID = `
    SELECT id, user_id, image, title_degree, institutional_name, start_date, end_date, description, created_at, updated_at
    FROM educations
    WHERE id = ?
    `

	queryGetEducationsByUserID = `
    SELECT id, user_id, image, title_degree, institutional_name, start_date, end_date, description, created_at, updated_at
    FROM educations
    WHERE user_id = ?
    ORDER BY start_date DESC
    `

	queryUpdateEducation = `
    UPDATE educations
    SET image = :image,
        user_id = :user_id,
        title_degree = :title_degree,
        institutional_name = :institutional_name,
        start_date = :start_date,
        end_date = :end_date,
        description = :description,
        updated_at = :updated_at
    WHERE id = :id
    `

	queryDeleteEducation = `
    DELETE FROM educations
    WHERE id = ?
    `

	queryDeleteEducationsByUserID = `
    DELETE FROM educations
    WHERE user_id = ?
    `

	queryCreatePortfolio = `
   INSERT INTO portfolios (
       id, user_id, image, project_name, project_location, description_image, project_link, start_date, end_date, description, created_at, updated_at
   ) VALUES (
       :id, :user_id, :image, :project_name, :project_location, :description_image, :project_link, :start_date, :end_date, :description, :created_at, :updated_at
   )`

	queryGetPortfolioByID = `
   SELECT id, user_id, image, project_name, project_location, description_image, project_link, start_date, end_date, description, created_at, updated_at
   FROM portfolios
   WHERE id = ?
   `

	queryGetPortfoliosByUserID = `
   SELECT id, user_id, image, project_name, project_location, description_image, project_link, start_date, end_date, description, created_at, updated_at
   FROM portfolios
   WHERE user_id = ?
   ORDER BY start_date DESC
   `

	queryUpdatePortfolio = `
   UPDATE portfolios
   SET image = :image,
       user_id = :user_id,
       project_name = :project_name,
       project_location = :project_location,
       description_image = :description_image,
       project_link = :project_link,
       start_date = :start_date,
       end_date = :end_date,
       description = :description,
       updated_at = :updated_at
   WHERE id = :id
   `

	queryDeletePortfolio = `
   DELETE FROM portfolios
   WHERE id = ?
   `

	queryDeletePortfoliosByUserID = `
   DELETE FROM portfolios
   WHERE user_id = ?
   `
)
