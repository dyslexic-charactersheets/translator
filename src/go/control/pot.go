package control

import (
	"../model"
	"io/ioutil"
	"fmt"
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
	// fmt.Println("POT:", poFile.MimeHeader.TranslatorComment)

	globalMeta := readPoMetaLines(poFile.MimeHeader.ExtractedComment)
	// fmt.Printf(" - %v\n", globalMeta)

	// messages
	for i, message := range poFile.Messages {
		progress.Progress = 2 + i
		// fmt.Println("MSG:", message.MsgId)

		messageMeta := readPoMetaLines(message.Comment.ExtractedComment)
		messageMeta = mergeMeta(messageMeta, globalMeta)
		messageRefs := readPoReferences(message.Comment)
		// fmt.Printf(" - meta: %v\n", messageMeta)
		// fmt.Printf(" - refs: %v\n", messageRefs)

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
			source := &model.Source{
				Filepath: "Pathfinder 2e/"+ref.File,
				Page:     "",
				Volume:   "",
				Level:    1,
				Game:     "pathfinder2",
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
	lines := strings.Split(str, `\n`)
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





/*
type PoEntry struct {
	msgctxt string
	msgid string
	msgstr string
}

type lastLineEnum int;

const (
	noLine lastLineEnum = iota + 1
	msgctxtLine
	msgidLine
	msgstrLine
)

// read a PO translation file or template

func readPo(reader io.Reader) []PoEntry {
	// regexes
	tagRx := regexp.MustCompile("#. (.*)\\: (.*)")
	tcommentRx := regexp.MustCompile("#. (.*)")
	refRx := regexp.MustCompile("#: (.*)")
	
	msgctxtRx := regexp.MustCompile("msgctxt \"(.*)\"")
	msgidRx := regexp.MustCompile("msgid \"(.*)\"")
	msgstrRx := regexp.MustCompile("msgstr \"(.*)\"")
	quoteRx := regexp.MustCompile("\"(.*)\"")


	entries := make([]PoEntry, 128)

	basemeta := make([string][]string, 32)
	meta := make([string][]string, 32)
	msgctxt := ""
	msgid := ""
	msgstr := ""

	scanner := bufio.NewScanner(reader)

	lastLine := noLine
	primed := false
	var accum strings.Builder

	pushEntry := func () {
		
	}

	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading pot:", err)
			break
		}

		// actually read the line!
		if quoteRx.MatchString(line) {
			content := quoteRx.FindStringSubmatch()[1]
			accum.writeString(content)
		} else {
			str := accum.String()
			accum.Reset()
		}
		
		if msgctxtRx.MatchString(line) {
			content := msgctxtRx.FindStringSubmatch()[1]
			accum.Reset()
			accum.writeString(content)

			lastLine = msgctxtLine
			primed = true
		} else if msgidRx.MatchString(line) {
			content := msgidRx.FindStringSubmatch()[1]
			accum.Reset()
			accum.writeString(content)

			lastLine = msgidLine
			primed = true
		} else if msgstrRx.MatchString(line) {
			content := msgstrRx.FindStringSubmatch()[1]
			accum.Reset()
			accum.writeString(content)

			lastLine = msgstrLine
			primed = true
		} else if tagRx.Match(line) {

		}

		// ... 


		// we found one!


		entries = append(entries, PoEntry{
			msgctxt: msgctxt,
			msgid: msgid,
			msgstr: msgstr,
		})
	}

	if primed {

	}
	
	file.Close()

	return entries
}
*/
