package auth

import "slices"

var admins = []int{
	2, // "admin@mail.com"
}

func IsAdmin(userID int) bool {
	return slices.Contains(admins, userID)
}
