package domain

type UserOption func(*User)

func Role(role string) UserOption {
	return func(u *User) {
		u.Role = UserRole(role)
	}
}

type AddressOption func(*Address)

func Coordinates(lat, lng float64) AddressOption {
	return func(a *Address) {
		a.Lat = lat
		a.Lng = lng
	}
}
