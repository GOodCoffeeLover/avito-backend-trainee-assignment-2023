package entity

type UserID uint

type User struct {
	ID UserID
}

func NewUser(id UserID) (*User, error) {
	return &User{
		ID: id,
	}, nil
}
