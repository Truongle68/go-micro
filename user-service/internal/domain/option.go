package domain

type UserOption func(*User)

func Role(role string) UserOption {
	return func(u *User) {
		u.Role = UserRole(role)
	}
}
