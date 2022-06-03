package migration

import (
	"prc_hub-api/users"
)

func MigrateAdminUser(email string, passwd string) (adminFound bool, invalidEmail bool, usedEmail bool, err error) {
	// Admin権限のuserを取得
	queryAdmin := true
	adminUsers, err := users.Get(users.GetQuery{Name: &email, Admin: &queryAdmin})
	if err != nil {
		return
	}

	// admin用メールアドレスが使用済みでないか確認
	adminFound = false
	if len(adminUsers) != 0 && adminUsers[0].Email == email {
		usedEmail = true
		return
	}

	_, invalidEmail, usedEmail, err = users.PostAdmin(users.PostBody{Name: "admin", Email: email, Password: passwd})
	if err != nil {
		return
	}
	return
}
