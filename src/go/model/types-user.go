package model

import (
	// "crypto/md5"
	"database/sql"
	// "encoding/hex"
	// "encoding/binary"
	// "fmt"
	// "github.com/ziutek/mymysql/mysql"
	"strings"
)




// ** Users

type User struct {
	Email          string
	Password       string
	Secret         string
	Name           string
	IsAdmin        bool
	Language       string
	IsLanguageLead bool
}

func parseUser(rows *sql.Rows) (Result, error) {
	u := User{}
	err := rows.Scan(&u.Email, &u.Password, &u.Secret, &u.Name, &u.IsAdmin, &u.Language, &u.IsLanguageLead)
	u.Email = strings.ToLower(u.Email)
	return u, err
}

const userFields = "Email, Password, Secret, Name, IsAdmin, Language, IsLanguageLead"

func GetUsers() []*User {
	results := query("select " + userFields + " from Users order by IsAdmin desc, Language asc, Name asc").rows(parseUser)
	users := make([]*User, len(results))
	for i, result := range results {
		if user, ok := result.(User); ok {
			users[i] = &user
		}
	}
	return users
}

func GetUserByEmail(email string) *User {
	result := query("select "+userFields+" from Users where Email = ?", email).row(parseUser)
	if user, ok := result.(User); ok {
		return &user
	}
	return nil
}

func GetUsersByLanguage(language string) []*User {
	results := query("select "+userFields+" from Users where Language = ? order by IsLanguageLead desc, Name asc", language).rows(parseUser)
	users := make([]*User, len(results))
	for i, result := range results {
		if user, ok := result.(User); ok {
			users[i] = &user
		}
	}
	return users
}

func GetLanguageLead(language string) *User {
	result := query("select "+userFields+" from Users where Language = ? and IsLanguageLead = 1", language).row(parseUser)
	if result != nil {
		if user, ok := result.(User); ok {
			return &user
		}
	}

	// users := GetUsersByLanguage(language)
	// if len(users) > 0 {
	// 	return users[0]
	// }
	return nil
}

func (user *User) Save() bool {
	keyfields := map[string]interface{}{
		"Email": user.Email,
	}
	fields := map[string]interface{}{
		"Password":       user.Password,
		"Secret":         user.Secret,
		"Name":           user.Name,
		"IsAdmin":        user.IsAdmin,
		"Language":       user.Language,
		"IsLanguageLead": user.IsLanguageLead,
	}
	return saveRecord("Users", keyfields, fields)
}

func (user *User) Delete() {
	keyfields := map[string]interface{}{
		"Email": user.Email,
	}
	deleteRecord("Users", keyfields)
}

func (user *User) CountTranslations() map[string]int {
	counts := make(map[string]int, len(Languages))
	query("select Language, Count(*) from Translations where Translator = ? group by Language", user.Email).rows(func(rows *sql.Rows) (Result, error) {
		var language string
		var count int
		rows.Scan(&language, &count)
		counts[language] = count
		return nil, nil
	})
	return counts
}