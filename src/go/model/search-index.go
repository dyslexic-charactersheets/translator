package model
/*
var stopwords = []string{
	
}

const searchIndexFields = "Entry, Term"

type SearchIndexEntry struct {
	Entry Entry
	Term  string
}

func parseSearchIndex(rows *sql.Rows) (Result, error) {
	e := SearchIndexEntry{}
	err := rows.Scan(&e.EntryID, &e.Term)
	return e, err
}

func Search(term string) {
	results := query("select "+searchIndexFields+" from SearchIndex where Term like ?", term).rows(parseSearchIndex)
	q.rows()
}

func CheckSearchIndexEntry(entry *Entry, term string) {

}

func AddSearchIndexEntry(entry *Entry, term string) {

}

func RemoveSearchIndexEntry(entry *Entry, term string) {

}

func (entry *Entry) IndexTerms() []string {

}
*/