package main

import (
	"./config"
	"./model"
	// "database/sql"
	"fmt"
	"strings"
	"strconv"
	_ "github.com/ziutek/mymysql/godrv"
)


func main() {
	db, err := config.Config.Database.Open()
	if err != nil {
		fmt.Println("Error opening database:", err)
	}

	// convert users to Brazilian
	brazilianUsers := []string{
		"erunnot@gmail.com",
		"bricio@live.com",
		"jmariossilva@gmail.com",
		"tsanzfre@gmail.com",
		"franciscocabralfilho@gmail.com",
	}

	fmt.Print("Converting... ")
	_, _ = db.Exec("update Users set Language = 'br' where Language = 'pt' and Email in ('"+strings.Join(brazilianUsers, "', '")+"')")
	fmt.Println("done.")

	// get all current Portuguese translations
	portugueseTranslations := model.GetTranslationsForLanguage("pt")
	fmt.Print("Converting "+strconv.Itoa(len(portugueseTranslations))+" translations... ")

	// insert Brazilian version
	for _, translation := range portugueseTranslations {
		// save a Brazilian Portuguese version of it
		translation.Language = "br"
		translation.Save(true);
	}
	fmt.Println("done.")
}