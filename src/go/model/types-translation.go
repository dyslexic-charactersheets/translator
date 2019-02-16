package model

import (
	// "crypto/md5"
	"database/sql"
	// "encoding/hex"
	// "encoding/binary"
	// "fmt"
	// "github.com/ziutek/mymysql/mysql"
	// "strings"
)



// ** Translations

type Translation struct {
	Entry        Entry
	Language     string
	Translation  string
	Translator   string
	IsPreferred  bool
	IsConflicted bool
}

func (translation *Translation) ID() uint64 {
	if translation == nil {
		return 0
	}

	var str = translation.Language + "  ---  " + translation.Translator
	return translation.Entry.ID() + hash64(str)
}

func parseTranslation(rows *sql.Rows) (Result, error) {
	t := Translation{}
	var entryID string
	err := rows.Scan(&entryID, &t.Language, &t.Translation, &t.Translator, &t.IsPreferred, &t.IsConflicted)
	if entry := GetEntryByID(entryID); entry == nil {
		return nil, nil
	} else {
		t.Entry = *entry
	}
	return t, err
}

const translationFields = "EntryID, Language, Translation, Translator, IsPreferred, IsConflicted"

func GetTranslations() []*Translation {
	results := query("select " + translationFields + " from Translations").rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

func GetTranslationByID(id string) *Translation {
	result := query("select "+translationFields+" from Translations where TranslationID = ?", id).row(parseTranslation)
	if translation, ok := result.(Translation); ok {
		return &translation
	}
	return nil
}

func GetTranslationsForLanguage(language string) []*Translation {
	results := query("select "+translationFields+" from Translations where Language = ?", language).rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

func (entry *Entry) GetTranslations(language string) []*Translation {
	results := query("select "+translationFields+" from Translations where EntryID = ? and Language = ?", entry.ID(), language).rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

func (entry *Entry) GetTranslationBy(language, translator string) *Translation {
	result := query("select "+translationFields+" from Translations where EntryID = ? and Language = ? and Translator = ?", entry.ID(), language, translator).row(parseTranslation)
	if translation, ok := result.(Translation); ok {
		return &translation
	}
	return nil
}

func (entry *Entry) GetMatchingTranslation(language, translation string) *Translation {
	result := query("select "+translationFields+" from Translations where EntryID = ? and Language = ? and Translation = ?", entry.ID(), language, translation).row(parseTranslation)
	if translation, ok := result.(Translation); ok {
		return &translation
	}
	return nil
}

func (translation *Translation) HasChanged() bool {
	underlying := translation.Entry.GetTranslationBy(translation.Language, translation.Translator)
	return underlying != nil && underlying.Translation == translation.Translation
}

func (translation *Translation) Save(clearVotes bool) {
	keyfields := map[string]interface{}{
		"TranslationID": translation.ID(),
	}
	fields := map[string]interface{}{
		"EntryID":     translation.Entry.ID(),
		"Language":    translation.Language,
		"Translator":  translation.Translator,
		"Translation": translation.Translation,
		"IsPreferred": translation.IsPreferred,
		"IsConflicted": translation.IsConflicted,
	}
	saveRecord("Translations", keyfields, fields)
	if clearVotes {
		ClearVotes(translation)
	}
}

/*
func AddTranslation(entry *Entry, language, translation string, translator *User) {
	keyfields := map[string]interface{}{
		"EntryID":   entry.ID(),
		"Language":      language,
		"Translator":    translator.Email,
	}
	fields := map[string]interface{}{
		"Translation": translation,
	}
	saveRecord("Translations", keyfields, fields)
}*/
