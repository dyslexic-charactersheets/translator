package control

import (
	"github.com/dyslexic-charactersheets/translator/src/go/model"
	"io/ioutil"
	"fmt"
	"github.com/dyslexic-charactersheets/translator/src/go/log"
	"net/http"
	"strconv"
	"strings"
	"regexp"
	"github.com/robfig/gettext-go/gettext/po"
)

func ImportPotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("POT import")
		
		file, _, err := r.FormFile("import-file")
		if err != nil {
			fmt.Println("Error reading file:", err)
			http.Redirect(w, r, "/import", 303)
			return
		}
		if file == nil {
			fmt.Println("Missing file")
			http.Redirect(w, r, "/import", 303)
			return
		}

		// read the file
		// potData := readPo(file)

		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("Error loading file:", err)
			http.Redirect(w, r, "/import", 303)
			return
		}

		progress := new(TaskProgress)
		progress.ID = <- nextProgressID
		CurrentProgress[progress.ID] = progress

		go importPot(data, progress)
		
		http.Redirect(w, r, "/import/progress?id="+strconv.Itoa(progress.ID), 303)
	}
}

func importPot(data []byte, progress *TaskProgress) {
	progress.Progress = 0

	poFile, err := po.LoadData(data)
	if err != nil {
		fmt.Println("Error:", err)
		progress.Abort = true
		return
	}
	
	progress.Scale = len(poFile.Messages) + 2
	progress.Progress = 1

	// header metadata
	log.Log("POT", "Translator Comment:", poFile.MimeHeader.TranslatorComment)

	globalMeta := readPoMetaLines(poFile.MimeHeader.ExtractedComment)
	log.Log("POT", "Global meta: %v", globalMeta)

	// find the game
	gameCode := "pathfinder2"
	filepathBase := "Pathfinder 2e/"
	if games, ok := globalMeta["Game"]; ok {
		for _, game := range games {
			log.Log("POT", "Found game:", game)
			switch game {
			case "Pathfinder 2nd Edition":
				gameCode = "pathfinder2"
				filepathBase = "Pathfinder 2e/"

			case "Pathfinder 1st Edition":
			case "Pathfinder":
				gameCode = "pathfinder"
				filepathBase = "Pathfinder/"

			case "Starfinder":
				gameCode = "starfinder"
				filepathBase = "Starfinder/"

			case "Dungeone & Dragons 3.5":
			case "Dungeons &amp; Dragons 3.5":
				gameCode = "3.5"
				filepathBase = "3.5/"
			}
		}
	}

	// messages
	for i, message := range poFile.Messages {
		progress.Progress = 2 + i
		log.Log("POT", "MSG:", message.MsgId)

		messageMeta := readPoMetaLines(message.Comment.ExtractedComment)
		messageMeta = mergeMeta(messageMeta, globalMeta)
		messageRefs := readPoReferences(message.Comment)
		log.Log("POT", " - meta: %v", messageMeta)
		log.Log("POT", " - refs: %v", messageRefs)

		// check if the message is part of a whole
		context := message.MsgContext
		partOf := ""
		if ix := strings.Index(context, message.MsgId); ix != -1 {
			partOf = context
		}

		// create the entry
		entry := &model.Entry{
			Original: message.MsgId,
			PartOf:   partOf,
		}
		entry.Save()

		for _, ref := range messageRefs {
			level := 1
			volume := ""
			if sources, ok := messageMeta["Source"]; ok {
				level1 := false
				level2 := false
				level3 := false
				for _, source := range sources {
					for _, src := range strings.Split(source, ",") {
						src = strings.TrimSpace(src)
						level = 4  // if there's at least one source, don't default to core any more
						volume = src
						switch src {
						case "Core Rulebook":
							level1 = true
						case "Advanced Player's Guide", "Secrets of Magic":
							level2 = true
						case "Lost Omens World Guide", "Lost Omens Chatracter Guide", "Lost Omens Gods and Magic", "Lost Omens Legends", "Lost Omens Pathfinder Society Guide":
							level3 = true
						}
					}
				}
				if level1 {
					level = 1
				} else if level2 {
					level = 2
				} else if level3 {
					level = 3
				} else {
					level = 4
				}
			}
			log.Log("POT", "Unit:", messageMeta["Unit"]);
			fileParts := strings.Split(ref.File, "/")
			page := fileParts[len(fileParts)-1]
			if unit, ok := messageMeta["Unit"]; ok && unit != nil {
				page = unit[0]+" ("+page+")"
			}
			source := &model.Source{
				Filepath: filepathBase+ref.File,
				Page:     page,
				Volume:   volume,
				Level:    level,
				Game:     gameCode,
			}
			source.Save()
			
			count := len(ref.Lines)
			entrySource := &model.EntrySource{
				Entry:  *entry,
				Source: *source,
				Count:  count,
			}
			entrySource.Save()
		}
	}

	progress.Finished = true
}

func mergeMeta(left, right map[string][]string) map[string][]string {
	merged := make(map[string][]string, len(left) + len(right))
	for key, values := range left {
		merged[key] = values
	}
	for key, values := range right {
		if _, ok := merged[key]; ok {
			merged[key] = append(merged[key], values...)
		} else {
			merged[key] = values
		}
	}
	return merged
}

func readPoMetaLines(str string) map[string][]string {
	lines := strings.Split(str, "\n")
	meta := make(map[string][]string, len(lines))

	metaRx := regexp.MustCompile("(.*?): (.*)");

	for _, line := range lines {
		if metaRx.MatchString(line) {
			submatch := metaRx.FindStringSubmatch(line)
			key := strings.Trim(submatch[1], " ")
			value := strings.Trim(submatch[2], " ")
			if _, ok := meta[key]; !ok {
				meta[key] = make([]string, 0, 8)
			}
			meta[key] = append(meta[key], value)
		}
	}

	return meta
}

type poReference struct {
	File  string
	Lines []int
}

func readPoReferences(comment po.Comment) []poReference {
	refs := make(map[string]poReference, len(comment.ReferenceFile))

	for i, file := range comment.ReferenceFile {
		if _, ok := refs[file]; !ok {
			refs[file] = poReference{file, make([]int, 0, 8)}
		}
		ref := refs[file]
		ref.Lines = append(ref.Lines, comment.ReferenceLine[i])
		refs[file] = ref
	}
	poRefs := make([]poReference, 0, len(refs))
    for _, ref := range refs {
        poRefs = append(poRefs, ref)
    }
	return poRefs
}


func ExportPoHandler(w http.ResponseWriter, r *http.Request) {
	language := r.FormValue("language")
	if language != "" {
		fmt.Println("Exporting in", language)
		translations := model.GetPreferredTranslations(language, true)

		messages := make([]po.Message, 0, len(translations))

		for _, translation := range translations {
			messages = append(messages, po.Message{
				po.Comment{
					0,                                   // start line
					"",                                  // translator comments
					"",                                  // extracted comments
					[]string{},                          // references
					[]int{},                             //     ''
					[]string{},                          // flags
					"",                                  // previous context
					"",                                  // previous untranslated string
				},
				"",                                      // msgcontext
				translation.Entry.FullText,              // msgid
				"",                                      // plural
				translation.FullText,                    // msgstr
				[]string{},                              // plural
			})
		}

		poFile := po.File{
			po.Header{
				po.Comment{
					0,                                   // start line
					"",                                  // translator comments
					"",                                  // extracted comments
					[]string{},                          // references
					[]int{},                             //     ''
					[]string{},                          // flags
					"",                                  // previous context
					"",                                  // previous untranslated string
				},
				"dyslexic-charactersheets 0.12.0",       // project version
				"Marcus Downing <marcus@bang-on.net>",   // report bugs to
				"",                                      // pot creation date
				"",                                      // po revision date
				"",                                      // last translator
				"",                                      // language team
				language,                                // language
				"1.0",                                   // mime version
				"text/plain; charset=UTF-8",             // content-type
				"8bit",                                  // content transfer encoding
				"",                                      // plural-forms
				"Dyslexic Character Sheets Translator",  // x-generator
				map[string]string{},                     // others
			},
			messages,
		}

		data := poFile.String()
		
		w.Header().Set("Content-Encoding", "UTF-8")
		w.Header().Set("Content-Type", "text/x-gettext-translation; charset=UTF-8")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+model.LanguageNamesEnglish[language]+".po\"")

		w.Write([]byte(data))
		
	} else {
		renderTemplate("export", w, r, nil)
	}
}