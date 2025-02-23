package auth

import "slices"

var admins = []int{
	6, // "whit@mail.com"
}

func IsAdmin(userID int) bool {
	return slices.Contains(admins, userID)
}
