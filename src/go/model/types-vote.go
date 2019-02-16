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



// ** Votes

type Vote struct {
	Translation Translation
	Voter       *User
	Vote        bool
}

const voteFields = "TranslationID, Voter, Vote"

func parseVote(rows *sql.Rows) (Result, error) {
	v := Vote{}
	var translationID, voter string
	err := rows.Scan(&translationID, &voter, &v.Vote)
	if err != nil {
		return nil, err
	}

	if translation := GetTranslationByID(translationID); translation == nil {
		return nil, nil
	} else {
		v.Translation = *translation
	}

	v.Voter = GetUserByEmail(voter)
	return v, err
}

func (translation *Translation) GetVote(voter *User) *Vote {
	result := query("select " + voteFields + " from Votes").row(parseVote)
	if vote, ok := result.(Vote); ok {
		vote.Translation = *translation
		vote.Voter = voter
		return &vote
	}
	return nil
}

func (entry *Entry) GetTranslationVotes(language string) []*Vote {
	results := query("select "+voteFields+" from Votes where EntryID = ? and Language = ?", entry.ID(), language).rows(parseVote)
	votes := make([]*Vote, len(results))
	for i, result := range results {
		if vote, ok := result.(Vote); ok {
			votes[i] = &vote
		}
	}
	return votes
}

func (vote *Vote) Save() {
	keyfields := map[string]interface{}{
		"TranslationID": vote.Translation.ID(),
		"Voter":      vote.Voter.Email,
	}
	fields := map[string]interface{}{
		"Vote": vote.Vote,
	}
	saveRecord("Votes", keyfields, fields)
}

func DeleteVote(vote *Vote) {
	keyfields := map[string]interface{}{
		"TranslationID": vote.Translation.ID(),
		"Voter":      vote.Voter.Email,
	}
	deleteRecord("Votes", keyfields)
}

func ClearVotes(translation *Translation) {
	keyfields := map[string]interface{}{
		"TranslationID": translation.ID(),
	}
	deleteRecord("Votes", keyfields)
}

func ClearOtherVotes(translation *Translation) {
	keyfields := map[string]interface{}{
		"TranslationID": translation.ID(),
		"Vote":     true,
	}
	deleteRecord("Votes", keyfields)
}