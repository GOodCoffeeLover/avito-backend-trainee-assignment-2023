package entity

type UserID uint32

type User struct {
	ID UserID `json:"id"`
}

func NewUser(userID UserID) (*User, error) {
	return &User{ID: userID}, nil
}
