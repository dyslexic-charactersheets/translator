package control

import (
	"github.com/dyslexic-charactersheets/translator/src/go/model"
	"github.com/dyslexic-charactersheets/translator/src/go/log"
	"fmt"
	// "code.google.com/p/go.crypto/bcrypt"
	// "crypto/md5"
	// "encoding/hex"
	// "html/template"
	// "math/rand"
	"net/http"
)

func APIEntriesHandler(w http.ResponseWriter, r *http.Request) {

}

func APITranslateHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)
	if user == nil {
		log.Warn("api", "Unknown user")
		return
	}

	entry := model.Entry{
		Original: r.FormValue("original"),
		PartOf:   r.FormValue("partOf"),
	}
	if entry.Original == "" {
		log.Warn("api", "Unknown string")
		return
	}
	language := r.FormValue("language")
	translation := r.FormValue("translation")

	if language == "" {
		log.Warn("api", "Unknown language:", language)
		return
	}
	// if translation == "" {
	// 	log.Warn("Blank translation:", entry.Original)
	// 	return
	// }
	log.Log("api", "Adding", language, "translation for:", entry.Original)

	t := &model.Translation{entry, language, translation, user.Email, false, false}
	t.Save(t.HasChanged())

	// recalculate conflicts
	stack := entry.GetStackedEntry()
	stack.MarkConflicts(language)

	fmt.Fprint(w, "OK")
}

const maxTranslations = 10

func APILookupHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)
	if user == nil {
		log.Warn("api", "Unknown user")
		return
	}

	lookup := r.FormValue("lookup")
	language := r.FormValue("language")

	// find stacked entries matching the search terms
	results := model.GetStackedEntries("", "0", "", "", lookup, true, "relevance", language, nil)

	// get those results with translations
	translationResults := make([]*model.StackedTranslation, 0, maxTranslations)
	for _, result := range results {
		translations := result.GetTranslations(language)
		if len(translations) > 0 {
			var preferredTranslation *model.StackedTranslation = nil
			for _, translation := range translations {
				if translation.IsPreferred {
					preferredTranslation = translation
				}
			}
			if preferredTranslation == nil {
				preferredTranslation = translations[0]
			}

			preferredTranslation.Entry = result
			log.Log("api", "Found result:", preferredTranslation)
			translationResults = append(translationResults, preferredTranslation)
		}
		if len(translationResults) >= maxTranslations {
			break;
		}
	}

	if len(translationResults) > 0 {
		fmt.Fprint(w, "<table>");
		for _, tr := range translationResults {
			log.Log("api", "Printing result:", tr)
			translated := tr.FullText
			original := tr.Entry.FullText
			fmt.Fprintf(w, "<tr><th>%s</th><td>%s</td></tr>", original, translated);
		}
		fmt.Fprint(w, "</table>");
	}
}

func APIVoteHandler(w http.ResponseWriter, r *http.Request) {
	if model.Debug >= 1 {
		log.Log("api", "Vote handler")
	}
	user := GetCurrentUser(r)
	if user == nil {
		log.Log("api", "Vote handler: Unknown user")
		return
	}

	entry := model.Entry{
		Original: r.FormValue("original"),
		PartOf:   r.FormValue("partOf"),
	}
	if entry.Original == "" {
		log.Warn("api", "Vote handler: Unknown string")
		return
	}
	language := r.FormValue("language")
	translation := r.FormValue("translation")

	// md5 := r.FormValue("voter")
	// voter := model.GetUserByMD5(md5)

	t := entry.GetMatchingTranslation(language, translation)

	up := r.FormValue("up") == "true"
	down := r.FormValue("down") == "true"
	if up && down {
		if model.Debug >= 1 {
			log.Warn("api", "Vote handler: Cannot vote both down and up")
		}
		return
	}

	if model.Debug >= 1 {
		log.Log("api", "Vote handler: Saving vote:", entry.Original, "=", translation, up, down)
	}

	model.ClearVotes(t)
	if up {
		model.ClearOtherVotes(t)
	}

	if up || down {
		v := &model.Vote{*t, user, up}
		v.Save()
	}

	// recalculate conflicts
	if model.Debug >= 1 {
		log.Log("api", "Checking for conflicts:")
	}
	stack := t.Entry.GetStackedEntry()
	stack.MarkConflicts(language)

	fmt.Fprintf(w, "OK")
}

func APISetLeadHandler(w http.ResponseWriter, r *http.Request) {
	me := GetCurrentUser(r)
	if !me.IsAdmin {
		log.Error("api", "Set lead: Not admin!")
		return
	}

	email := r.FormValue("user")
	user := model.GetUserByEmail(email)

	if user == nil {
		log.Error("api", "Set lead: Unknown user")
		return
	}
	log.Log("api", "Set lead: Setting", user.Name, "as language lead for", model.LanguageNames[user.Language])
	user.SetLanguageLead()
}

func APIClearLeadHandler(w http.ResponseWriter, r *http.Request) {
	me := GetCurrentUser(r)
	if !me.IsAdmin {
		log.Error("api", "Clear lead: Not admin!")
		return
	}

	email := r.FormValue("user")
	user := model.GetUserByEmail(email)

	if user == nil {
		log.Error("api", "Clear lead: Unknown user")
		return
	}
	log.Log("api", "Clear lead: Removing", user.Name, "as language lead for", model.LanguageNames[user.Language])
	user.ClearLanguageLead()
}
