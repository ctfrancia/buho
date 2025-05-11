package domain

type User struct {
	UUID               string
	Email              string
	FirstName          string
	LastName           string
	Website            string
	Password           string
	RefreshToken       string
	RefreshTokenExpiry string
}
