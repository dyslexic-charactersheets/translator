package model

import (
	// "crypto/md5"
	"database/sql"
	// "encoding/hex"
	// "encoding/binary"
	"../log"
	// "github.com/ziutek/mymysql/mysql"
	"strings"
)



// ** Entries

type Entry struct {
	Original string
	PartOf   string
}

func (entry *Entry) ID() uint64 {
	if entry == nil {
		return 0
	}

	var str = entry.Original
	if entry.PartOf != "" && entry.PartOf != entry.Original {
		str = entry.Original + "  ----  " + entry.PartOf
	}
	return hash64(str)
}

func DeleteAllEntrySources() {
	if Debug >= 2 {
		log.Warn("types:entry", " ***** Deleting ALL entry sources")
	}
	// if ok := query("delete from Entries").exec(); !ok {
	// 	log.Warn("types:entry", " ***** Error deleting entries")
	// }
	if ok := query("delete from Sources").exec(); !ok {
		log.Warn("types:entry", " ***** Error deleting sources")
	}
	if ok := query("delete from EntrySources").exec(); !ok {
		log.Warn("types:entry", " ***** Error deleting entrysources")
	}
	if Debug >= 2 {
		log.Warn("types:entry", " ***** Deleted ALL entry sources")
	}
}

const entryFields = "Original, PartOf"

func parseEntry(rows *sql.Rows) (Result, error) {
	e := Entry{}
	err := rows.Scan(&e.Original, &e.PartOf)
	// log.Log("types:entry", "Entry ID: " + e.ID() + " (" + string(len(e.ID())) + ")")
	return e, err
}

func makeEntries(results []Result) []*Entry {
	entries := make([]*Entry, len(results))
	for i, result := range results {
		if entry, ok := result.(Entry); ok {
			entries[i] = &entry
		}
	}
	return entries
}

func CountEntries() int {
	return query("select count(*) from Entries").count()
}

func GetEntryByID(id string) *Entry {
	result := query("select "+entryFields+" from Entries where EntryID = ?", id).row(parseEntry)
	if entry, ok := result.(Entry); ok {
		return &entry
	}
	return nil
}

func GetEntries() []*Entry {
	results := query("select " + entryFields + " from Entries").rows(parseEntry)
	return makeEntries(results)
}

func GetEntriesPartOf(partOf string) []*Entry {
	results := query("select " + entryFields + " from Entries where Partof = ?", partOf).rows(parseEntry)
	return makeEntries(results)
}

func GetEntriesAt(game string, level int, file, show, search string, fuzzySearch bool, language string, translator *User) []*Entry {
	if game == "" && level == 0 && show == "" && search == "" {
		return GetEntries()
	}
	args := make([]interface{}, 0, 2)
	sql := "select Original, PartOf from Entries " +
		"inner join EntrySources on Entries.EntryID = EntrySources.EntryID " +
		"inner join Sources on EntrySources.SourceID = Sources.SourceID"
	if show == "conflicts" {
		sql = sql + " inner join Translations on Entries.EntryID = Translations.EntryID and Translations.Language = ?"
		args = append(args, language)
		// sql = sql + " inner join Translations Mine on EntryID = Mine.EntryID and Mine.Language = ? and Mine.Translator = ?" +
		// 	"inner join Translations Others on Entries.EntryID = Others.EntryID and Others.Language = ? and Others.Translator != ?"
		// args = append(args, language)
		// args = append(args, translator.Email)
		// args = append(args, language)
		// args = append(args, translator.Email)
	} else if show == "mine" {
		sql = sql + " inner join Translations Mine on Entries.EntryID = Mine.EntryID and Mine.Language = ? and Mine.Translator = ?"
		args = append(args, language)
		args = append(args, translator.Email)
	} else if show == "others" {
		sql = sql + " inner join Translations Others on Entries.EntryID = Others.EntryID and Others.Language = ? and Others.Translator = ?"
		args = append(args, language)
		args = append(args, translator.Email)
	} else if show != "" {
		sql = sql + " left join Translations on Entries.EntryID = Translations.EntryID and Translations.Language = ?"
		args = append(args, language)
	}
	sql = sql + " where 1 = 1"

	if game != "" {
		if game == "dnd35" {
			game = "3.5"
		}
		sql = sql + " and Game = ?"
		args = append(args, game)
	}
	if level != 0 {
		sql = sql + " and Level = ?"
		args = append(args, level)
	}
	if file != "" {
		sql = sql + " and Filepath = ?"
		args = append(args, file)
	}
	if show == "conflicts" {
		sql = sql + " and Translations.IsConflicted = 1"
	}
	// if show != "" {
	// 	sql = sql+" and Translations.Language = ?"
	// 	args = append(args, language)
	// }
	if search != "" {
		searchTerms := strings.Split(search, " ")
		log.Log("type:entry", "Searching for:", search)

		if fuzzySearch {
			// todo make it more fuzzy?
			sql = sql + " and ("
			first := true
			for _, term := range searchTerms {
				if first {
					first = false
				} else {
					sql = sql + " or "
				}
				term = strings.ToLower(term)
				sql = sql + "lower(Original) like ?"
				args = append(args, "%"+term+"%")
			}
			sql = sql + ")"
		} else {
			for _, term := range searchTerms {
				term = strings.ToLower(term)
				sql = sql + " and lower(Original) like ?"
				args = append(args, "%"+term+"%")
			}
		}
	}

	sql = sql + " group by Entries.EntryID"
	if show == "translated" {
		sql = sql + " having count(Translations.Translation) > 0"
	} else if show == "untranslated" {
		sql = sql + " having count(Translations.Translation) = 0"
	}
	log.Log("types:entry", "Get entries:", sql)
	results := query(sql, args...).rows(parseEntry)
	return makeEntries(results)
}

func (entry *Entry) Save() {
	keyfields := map[string]interface{}{
		"EntryID": entry.ID(),
	}
	fields := map[string]interface{}{
		"Original": entry.Original,
		"PartOf":   entry.PartOf,
	}
	saveRecord("Entries", keyfields, fields)
}

func (entry *Entry) CountTranslations() map[string]int {
	counts := make(map[string]int, len(Languages))
	query("select Language, Count(*) from Translations where EntryID = ? group by Language", entry.ID()).rows(func(rows *sql.Rows) (Result, error) {
		var language string
		var count int
		rows.Scan(&language, &count)
		counts[language] = count
		return nil, nil
	})
	return counts
}

func (entry *Entry) GetParts() []*Entry {
	if entry.PartOf == "" {
		entries := make([]*Entry, 1)
		entries[0] = entry
		return entries
	}
	results := query("select "+entryFields+" from Entries where PartOf = ?", entry.PartOf).rows(parseEntry)
	return makeEntries(results)
}