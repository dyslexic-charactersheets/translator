package model

import (
	"crypto/md5"
	"database/sql"
	// "encoding/hex"
	// "encoding/binary"
	// "fmt"
	"github.com/ziutek/mymysql/mysql"
	// "strings"
)

func hash64(data string) uint64 {
	hasher := md5.New()
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)

	hash64 :=
		uint64(hash[0])<<24 +
			uint64(hash[1])<<16 +
			uint64(hash[2])<<8 +
			uint64(hash[3])
	return hash64
}

func parseID(rows *sql.Rows) (Result, error) {
	var id uint64
	err := rows.Scan(&id)
	return id, err
}






// ** Comments

type Comment struct {
	Entry       Entry
	Language    string
	Commenter   string
	Comment     string
	CommentDate mysql.Timestamp
}
