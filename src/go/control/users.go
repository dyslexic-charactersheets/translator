package control

import (
	"github.com/dyslexic-charactersheets/translator/src/go/config"
	"github.com/dyslexic-charactersheets/translator/src/go/model"
	"github.com/dyslexic-charactersheets/translator/src/go/log"
	"golang.org/x/crypto/bcrypt"
	"crypto/sha256"
	"fmt"
	// "github.com/bpowers/seshcookie"
	// "encoding/json"
	"encoding/hex"
	// "io/ioutil"
	"html/template"
	"net/http"
	// "strings"
	// "net/url"
	// "github.com/russross/blackfriday"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Log("users", "Dashboard")

	renderTemplate("home", w, r, func(data *TemplateData) {
		log.Log("users", "Dashboard params")
		data.LanguagesEnglish = model.LanguageNamesEnglish
		data.LanguageCompletion = model.GetLanguageCompletion()
		data.LiveLoginURL = GetLiveLoginURL(r)
		data.DevLoginURL = GetDevLoginURL(r)
		log.Log("users", "Dashboard params done")
	})
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		name := r.FormValue("name")
		language := r.FormValue("language")
		user := &model.User{
			Email:    email,
			Name:     name,
			Language: language,
			Password: "",
			Secret:   "",
		}
		user.Save()

		// send a welcome message
		if r.FormValue("welcome-email") == "on" {
			sendInvitationEmail(user)
		}

		http.Redirect(w, r, "/users", 303)
	} else {
		renderTemplate("users", w, r, func(data *TemplateData) {
			data.Users = model.GetUsers()
			data.UsersByLanguage = make(map[string][]*model.User, len(data.Languages))
			for _, user := range data.Users {
				data.UsersByLanguage[user.Language] = append(data.UsersByLanguage[user.Language], user)
			}
		})
	}
}

func UsersShowInviteHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := GetCurrentUser(r)
	if !currentUser.IsAdmin {
		http.Redirect(w, r, "/users", 303)
		return
	}

	email := r.FormValue("user")
	user := model.GetUserByEmail(email)
	if user == nil {
		log.Log("users", "User not found: "+email)
		return
	}
	
	protocol := "https"
	hostname := config.Config.Server.Hostname
	if hostname == "localhost" {
		protocol = "http"
		hostname = fmt.Sprintf("localhost:%d", config.Config.Server.Port)
	}
	secret := user.Secret

	renderTemplate("users_invite", w, r, func(data *TemplateData) {
		data.User = user
		
		if secret == "" {
			data.InviteURL = ""
		} else {
			data.InviteURL = fmt.Sprintf("%s://%s/account/reclaim?email=%s&secret=%s", protocol, hostname, email, secret);
		}
	})
}

func UsersRenewInviteHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := GetCurrentUser(r)
	if !currentUser.IsAdmin {
		http.Redirect(w, r, "/users", 303)
		return
	}

	email := r.FormValue("user")
	user := model.GetUserByEmail(email)
	if user == nil {
		log.Warn("users", "User not found:", email)
		return
	}

	user.GenerateSecret()
	http.Redirect(w, r, "/users/show-invite?user="+email, 303)
}

func UsersReinviteHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := GetCurrentUser(r)
	if !currentUser.IsAdmin {
		http.Redirect(w, r, "/users", 303)
		return
	}

	email := r.FormValue("user")
	user := model.GetUserByEmail(email)
	if user == nil {
		log.Warn("users", "User not found: "+email)
		return
	}
	log.Log("users", "Reinviting")
	sendInvitationEmail(user)

	http.Redirect(w, r, "/users/show-invite?user="+email, 303)
}

func sendInvitationEmail(user *model.User) {
	msg := `Subject: Welcome to the Dyslexic Character Sheets Translator
Content-Type: text/plain; charset="UTF-8"
Reply-to: marcus@dyslexic-charactersheets.com

Welcome to the Dyslexic Character Sheets Translator.

Your account has been created in the %s group. Click this link to set a password:

https://%s/account/reclaim?email=%s&secret=%s

Click the "Translate" link at the top to start translating or to correct existing translations. 
You can use the options at the top to limit it to [Pathfinder, Core, Untranslated]; or you can 
search for specific words. If the original text includes a vertical bar "|", that indicates a line 
break; try to include a "|" at a similar point in the translation. If the original is broken into 
two or more boxes, that means the font changes mid-line (eg, "Greater RAGE!") and you should try to 
make sure the right bit of the translation lines up with the right font... as best you can. 
It won't always work out neatly.

If other people have already translated a line, you can use the tick and cross boxes to vote for 
the version you like or dislike. When I export the translations, it will take these votes into 
account when choosing which one to give me. Writing the exact same translation as somebody else is 
equivalent to voting for their translation.

If you need to see writing on the page, the button on the left of each line brings up a list of the 
pages it appears on, with links to the PDFs. The "Sources" link at the top gives you a more 
complete list.

Once you're logged in, click the "Sources" link at the top to see the most recent preview of the 
translations in PDF. The previews *don't* update automatically as you translate, only when I spend 
an evening or two making all the PDFs.


https://charactersheets.slack.com/

You should receive an invitation to this site as well. It's a place for translators to chat and 
compare notes. If you need to ask me a question, this the place to do it.


Marcus Downing
https://www.dyslexic-charactersheets.com/
`

	language := user.Language
	email := user.Email
	hostname := config.Config.Server.Hostname
	secret := user.GenerateSecret()
	msg = fmt.Sprintf(msg, model.LanguageNamesEnglish[language], hostname, email, secret)

	log.Log("users", "Sending message to", user.Email, "\n", msg)

	if ok := config.SendMail(email, msg); ok {
		log.Log("users", "Invitation email sent to "+user.Email)
	}
}

func UsersAddHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("users_add", w, r, nil)
}

func UsersMasqueradeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("users_masq", w, r, nil)
}

func UsersDelHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := GetCurrentUser(r)
	if !currentUser.IsAdmin {
		http.Redirect(w, r, "/users", 303)
		return
	}

	email := r.FormValue("user")
	user := model.GetUserByEmail(email)
	if user == nil {
		http.Redirect(w, r, "/users", 303)
		return
	}

	gonow := r.FormValue("go")
	if r.Method == "POST" && gonow == "yes" {
		user.Delete()
		http.Redirect(w, r, "/users", 303)
		return
	} else {
		renderTemplate("users_del", w, r, func(data *TemplateData) {
			data.User = user
		})
	}
}

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)

	if r.Method == "POST" {
		user.Name = r.FormValue("name")
		language := r.FormValue("language")
		if language != "" {
			user.Language = language
		}
		user.Save()

		http.Redirect(w, r, "/home", 303)
	} else {
		renderTemplate("account", w, r, nil)
	}
}

func AccountReclaimHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("account_reclaim", w, r, nil)
}

func SetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)

	if r.Method == "POST" {
		password := r.FormValue("password")
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err == nil {
			user.Password = string(hash)
			user.Save()
		}
		http.Redirect(w, r, "/account", 303)
	} else {
		renderTemplate("account_set_password", w, r, nil)
	}
}

type Issue struct {
	Number          int    `json:"number"`
	Name            string `json:"title"`
	SummaryMarkdown string `json:"body"`
	SummaryHTML     template.HTML
	URL             string `json:"url"`
	CssClass        string
	Avatar          string
	User            struct{
		Avatar string `json:"avatar_url"`
	} `json:"user"`
	Labels          []struct {
		URL   string `json:"url"`
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"labels"`
}

// type issueLabel struct {
// 	URL   string `json:"url"`
// 	Name  string `json:"name"`
// 	Color string `json:"color"`
// }


func AuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	log.Log("users", "Auth redirect")

	liveLoginURL := GetLiveLoginURL(r)
	http.Redirect(w, r, liveLoginURL, 303)
}


func GetLiveLoginURL(r *http.Request) string {
	base := config.Config.Live.LiveLoginURL
	sharedSecret := config.Config.Live.SharedSecret

	log.Log("users", "Live login: shared secret:", sharedSecret)

	currentUser := GetCurrentUser(r)
	h := sha256.New()
	h.Write([]byte(currentUser.Email))
	hash := h.Sum(nil)
	token := hex.EncodeToString(hash)

	log.Log("users", "Live login: token:", token)

	h = sha256.New()
	h.Write([]byte(token))
	h.Write([]byte(sharedSecret))
	hash = h.Sum(nil)
	signature := hex.EncodeToString(hash)
	
	log.Log("users", "Live login: signature:", signature)

	return base+"?login="+token+":"+signature
}


func GetDevLoginURL(r *http.Request) string {
	base := config.Config.Dev.DevLoginURL
	sharedSecret := config.Config.Dev.SharedSecret

	currentUser := GetCurrentUser(r)
	h := sha256.New()
	h.Write([]byte(currentUser.Email))
	hash := h.Sum(nil)
	token := hex.EncodeToString(hash)

	h = sha256.New()
	h.Write([]byte(token))
	h.Write([]byte(sharedSecret))
	hash = h.Sum(nil)
	signature := hex.EncodeToString(hash)

	return base+"?login="+token+":"+signature
}
