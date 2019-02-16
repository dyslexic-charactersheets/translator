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



type EntrySource struct {
	Entry  Entry
	Source Source
	Count  int
}

func parseEntrySource(rows *sql.Rows) (Result, error) {
	es := EntrySource{}
	var entryID string
	var sourceID string
	err := rows.Scan(&entryID, &sourceID, &es.Count)
	if entry := GetEntryByID(entryID); entry == nil {
		return nil, nil
	} else {
		es.Entry = *entry
	}
	if source := GetSourceByID(sourceID); source == nil {
		return nil, nil
	} else {
		es.Source = *source
	}
	return es, err
}

const entrySourceFields = "EntryID, SourceID, Count"

func GetEntrySources() []*EntrySource {
	results := query("select EntryID, EntrySources.SourceID, Count" +
		" from EntrySources inner join Sources on EntrySources.SourceID = Sources.SourceID").rows(parseEntrySource)

	sources := make([]*EntrySource, len(results))
	for i, result := range results {
		if source, ok := result.(EntrySource); ok {
			sources[i] = &source
		}
	}
	return sources
}

func GetSourcesForEntry(entry *Entry) []*EntrySource {
	results := query("select EntryID, EntrySources.SourceID, Count"+
		" from EntrySources inner join Sources on EntrySources.SourceID = Sources.SourceID "+
		"where EntryID = ?", entry.ID()).rows(parseEntrySource)

	sources := make([]*EntrySource, len(results))
	for i, result := range results {
		if source, ok := result.(EntrySource); ok {
			sources[i] = &source
		}
	}
	return sources
}

func (es *EntrySource) Save() {
	keyfields := map[string]interface{}{
		"EntryID":    es.Entry.ID(),
		"SourceID": es.Source.ID(),
	}
	fields := map[string]interface{}{
		"Count": es.Count,
	}
	saveRecord("EntrySources", keyfields, fields)
}

type EntrySourcePlaceholder struct {
	SourceID uint64
	Count    int
}

func parseEntrySourcePlaceholder(rows *sql.Rows) (Result, error) {
	placeholder := EntrySourcePlaceholder{}
	err := rows.Scan(&placeholder.SourceID, &placeholder.Count)
	return placeholder, err
}

func GetSourceIDsForEntry(entry *Entry) []EntrySourcePlaceholder {
	results := query("select SourceID, Count from EntrySources where EntryID = ?", entry.ID()).rows(parseEntrySourcePlaceholder)

	sources := make([]EntrySourcePlaceholder, len(results))
	for i, result := range results {
		if id, ok := result.(EntrySourcePlaceholder); ok {
			sources[i] = id
		}
	}
	return sources
}