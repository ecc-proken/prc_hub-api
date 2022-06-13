package migration

import (
	"prc_hub-api/users"

	"golang.org/x/crypto/bcrypt"
)

func MigrateAdminUser(email string, passwd string) (adminFound bool, updated bool, invalidEmail bool, usedEmail bool, err error) {
	// Admin権限のuserを取得
	u, notFound, err := users.GetMigratedAdmin()
	if err != nil {
		return
	}

	if !notFound {
		// 登録済み
		adminFound = true

		var verify bool
		verify, err = u.Verify(passwd)
		if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
			return
		}
		if email != u.Email || !verify {
			updated = true
			// 更新
			u, invalidEmail, usedEmail, _, err = users.Patch(u.Id, users.PatchBody{Email: &email, Password: &passwd})
		}
		return
	}

	_, invalidEmail, usedEmail, err = users.PostAdmin(users.PostBody{Name: "admin", Email: email, Password: passwd})
	if err != nil {
		return
	}
	return
}
