package entity

type User struct {
	ID             string
	Email          string
	Password       string
	Name           string
	Role           userRole
	ProfilePicture string
	headline       string
	createdAt      string
	updatedAt      string
}

type userRole string

const (
	Admin     userRole = "admin"
	Recruiter userRole = "recruiter"
	Candidate userRole = "candidate"
)

type UserLoginData struct {
	ID       string
	Username string
	Email    string
}
