package models

type User struct {
	id    int
	name  string
	email string
}

func GetNewUser(name, email string, id int) *User {
	return &User{
		id:    id,
		name:  name,
		email: email,
	}
}

func (u *User) GetId() int {
	return u.id
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetEmail() string {
	return u.email
}
