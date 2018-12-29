package oauth

import (
	"github.com/devlover-id/api/pkg/database"
)

func checkUsernameRegistered(username string) (bool, error) {
	var query = `
		select exists(id) as exist
		from users
		where github_username=?
	`
	var temp struct{ Exist bool }
	if err := database.Query(&temp, query, username); err != nil {
		return false, err
	}
	return temp.Exist, nil
}
