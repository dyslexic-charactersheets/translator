package model

import (
	// "crypto/md5"
	"database/sql"
	// "encoding/hex"
	// "encoding/binary"
	"github.com/dyslexic-charactersheets/translator/src/go/log"
	// "github.com/ziutek/mymysql/mysql"
	// "strings"
)



// ** Sources

type Source struct {
	Filepath string
	Page     string
	Volume   string
	Level    int
	Game     string
}

func GetSourceByID(id string) *Source {
	result := query("select "+sourceFields+" from Sources where SourceID = ?", id).row(parseSource)
	if source, ok := result.(Source); ok {
		return &source
	}
	return nil
}

func GetSourceByPath(path string) *Source {
	result := query("select "+sourceFields+" from Sources where Filepath = ?", path).row(parseSource)
	if source, ok := result.(Source); ok {
		return &source
	}
	return nil
}

func (source *Source) ID() uint64 {
	if source == nil {
		return 0
	}

	return hash64(source.Filepath)
	// hasher := md5.New()
	// hasher.Write([]byte(source.Filepath))
	// return hex.EncodeToString(hasher.Sum(nil))
}

func parseSource(rows *sql.Rows) (Result, error) {
	s := Source{}
	err := rows.Scan(&s.Filepath, &s.Page, &s.Volume, &s.Level, &s.Game)
	return s, err
}

const sourceFields = "Filepath, Page, Volume, Level, Game"

func GetSources() []*Source {
	results := query("select " + sourceFields + " from Sources").rows(parseSource)

	sources := make([]*Source, len(results))
	for i, result := range results {
		if source, ok := result.(Source); ok {
			sources[i] = &source
		}
	}
	return sources
}
func GetSourcesAt(game string, level int, show string) []*Source {
	if game == "" && level == 0 && show == "" {
		return GetSources()
	}
	args := make([]interface{}, 0, 2)
	sql := "select " + sourceFields + " from Sources "

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

	// sql = sql+" group by Original, PartOf"
	if show == "translated" || show == "untranslated" {
		sql = sql + " and Sources.SourceID"
		if show == "untranslated" {
			sql = sql + " not"
		}
		sql = sql + " in (select EntrySources.SourceID from EntrySources" +
			" inner join Translations on EntrySources.EntryID = Translations.EntryID)"
	}

	sql = sql + " order by Page"

	log.Log("types:source", "Get entries:", sql)
	results := query(sql, args...).rows(parseSource)

	sources := make([]*Source, 0, len(results))
	for _, result := range results {
		if source, ok := result.(Source); ok && source.Page != "" {
			sources = append(sources, &source)
		}
	}
	return sources
}

func (source *Source) Save() {
	keyfields := map[string]interface{}{
		"SourceID": source.ID(),
	}
	fields := map[string]interface{}{
		"Filepath": source.Filepath,
		"Page":   source.Page,
		"Volume": source.Volume,
		"Level":  source.Level,
		"Game":   source.Game,
	}
	saveRecord("Sources", keyfields, fields)
}

func (source *Source) GetLanguageCompletion() map[string]int {
	var completion = make(map[string]int, len(Languages))

	total := query("select count(distinct Entries.EntryID) from Entries "+
		"inner join EntrySources on Entries.EntryID = EntrySources.EntryID "+
		"where EntrySources.SourceID = ?", source.ID()).count()
	if total > 0 {
		for _, lang := range Languages {
			count := query("select count(distinct Translations.EntryID) from Translations "+
				"inner join EntrySources on Translations.EntryID = EntrySources.EntryID "+
				"where EntrySources.SourceID = ? and Language = ?", source.ID(), lang).count()
			completion[lang] = 100 * count / total
		}
	}
	return completion
}