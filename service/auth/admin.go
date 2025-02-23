package auth

import "slices"

var admins = []string{
	"whit@mail.com",
	"whitallee@gmail.com",
	"mariaelenamilan00@gmail.com",
}

func IsAdmin(user string) bool {
	return slices.Contains(admins, user)
}
