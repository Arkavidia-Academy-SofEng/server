package recruitmentRepository

const (
	queryCreateJobVacancy = `
INSERT INTO job_vacancies (id, recruiter_id, title, description, requirements, location, job_type, deadline, is_active, created_at, updated_at)
VALUES (:id, :recruiter_id, :title, :description, :requirements, :location, :job_type, :deadline, :is_active, :created_at, :updated_at)`

	queryGetJobVacancies = `
    SELECT id, recruiter_id, title, description, requirements, location, job_type, 
           deadline, is_active, created_at, updated_at
    FROM job_vacancies
    ORDER BY created_at DESC
    LIMIT ? OFFSET ?
    `

	queryCheckJobVacancyExists = `
    SELECT EXISTS (SELECT 1 FROM job_vacancies WHERE id = ?)
    `

	queryUpdateJobVacancy = `
    UPDATE job_vacancies
    SET title = :title,
        description = :description,
        requirements = :requirements,
        location = :location,
        job_type = :job_type,
        deadline = :deadline,
        is_active = :is_active,
        updated_at = :updated_at
    WHERE id = :id
    `

	queryCountJobVacancies = `
    SELECT COUNT(*) FROM job_vacancies
    `
	queryGetAllJobVacancies = `
SELECT *
FROM job_vacancies`

	queryDeleteJobVacancy = `
    DELETE FROM job_vacancies 
    WHERE id = ?
    `
)
