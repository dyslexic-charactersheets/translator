package control

import (
	"../config"
	"../model"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	// "github.com/bpowers/seshcookie"
	"encoding/json"
	"io/ioutil"
	"html/template"
	"net/http"
	"net/smtp"
	"crypto/tls"
	"strings"
	// "net/url"
	"github.com/russross/blackfriday"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Dashboard")

	renderTemplate("home", w, r, func(data *TemplateData) {
		data.LanguageCompletion = model.GetLanguageCompletion()
		data.Issues, data.NumIssues = GetGithubIssues()
		data.WebsiteIssues, data.NumWebsiteIssues = GetWebsiteIssues()
		data.TranslatorIssues, data.NumTranslatorIssues = GetTranslatorIssues()
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

func UsersReinviteHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := GetCurrentUser(r)
	if !currentUser.IsAdmin {
		http.Redirect(w, r, "/users", 303)
		return
	}

	email := r.FormValue("user")
	user := model.GetUserByEmail(email)
	if user == nil {
		fmt.Println("User not found: "+email)
		return
	}
	fmt.Println("Reinviting")
	sendInvitationEmail(user)

	http.Redirect(w, r, "/users", 303)
}

func sendInvitationEmail(user *model.User) {
	mailConfig := config.Config.Mail

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

	// to := []string{user.Email}
	from := mailConfig.From

	fmt.Println("Sending message to", user.Email, "\n", msg)
	// auth := smtp.CRAMMD5Auth(mailConfig.Username, mailConfig.Password)
	// err := smtp.SendMail(mailConfig.Hostname, auth, from, to, []byte(msg))

	// SEND EMAIL THE HARD WAY
	// connect
	client, err := smtp.Dial(mailConfig.Hostname)
	defer client.Quit()
	if err != nil {
		fmt.Println("Error sending mail:", err)
		return
	}
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: mailConfig.Hostname, InsecureSkipVerify: true} 
		client.StartTLS(tlsConfig)
	}
	// err = client.Auth(auth)
	// if err != nil {
	// 	fmt.Println("Error sending mail:", err)
	// 	return
	// }

	// set recipient
	client.Mail(from)
	client.Rcpt(user.Email)

	// write the body
	writer, err := client.Data()
	if err != nil {
		fmt.Println("Error sending mail:", err)
		return
	}
	_, err = fmt.Fprintf(writer, msg)
	if err != nil {
		fmt.Println("Error sending mail:", err)
		return
	}
	err = writer.Close()
	if err != nil {
		fmt.Println("Error sending mail:", err)
		return
	}

	// Send the QUIT command and close the connection.
	err = client.Quit()
	if err != nil {
		fmt.Println("Error sending mail:", err)
		return
	}

	fmt.Println("Invitation email sent to "+user.Email)
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

func GetGithubIssues() ([]Issue, int) {
	issues, num := getGithubAPIIssues("marcusatbang/charactersheets")
	// issues, num := getGithubAPIIssues("dyslexic-charactersheets/pages")
	return issues, num
}

func GetWebsiteIssues() ([]Issue, int) {
	issues, num := getGithubAPIIssues("marcusatbang/charactersheets-website")
	// issues, num := getGithubAPIIssues("dyslexic-charactersheets/website")
	return issues, num
}

func GetTranslatorIssues() ([]Issue, int) {
	issues, num := getGithubAPIIssues("marcusatbang/charactersheets-translator")
	// issues, num := getGithubAPIIssues("dyslexic-charactersheets/translator")
	return issues, num
}

func getGithubAPIIssues(repo string) ([]Issue, int) {
	resp, err := http.Get("https://api.github.com/repos/"+repo+"/issues?state=open&sort=updated&access_token="+config.Config.Github.AccessToken)
	if err != nil {
		fmt.Println("Error fetching issues from GitHub:", err)
		return []Issue{}, 0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading issues from GitHub:", err)
		return []Issue{}, 0
	}

	// fmt.Println(string(body))
	issues := make([]Issue, 0)
	err = json.Unmarshal(body, &issues)
	if err != nil {
		fmt.Println("Error decoding issues from GitHub:", err)
		return []Issue{}, 0
	}
	numIssues := len(issues)
	if len(issues) > 30 {
		issues = issues[0:30]
	}

	for i, issue := range issues {
		issues[i].URL = strings.Replace(issue.URL, "https://api.github.com/repos/", "https://www.github.com/", 1)

		if issue.SummaryMarkdown != "" {
			html := blackfriday.MarkdownCommon([]byte(issue.SummaryMarkdown))
			if html != nil {
				issues[i].SummaryHTML = template.HTML(html)
			}
		}
		// if issue.SummaryMarkdown != "" {
		// 	fmt.Println("Parsing Markdown:", issue.SummaryMarkdown)
		// 	resp, err = http.PostForm("https://api.github.com/markdown", url.Values{"text": {issue.SummaryMarkdown}})
		// 	if err != nil {
		// 		fmt.Println("Error parsing Markdown:", err)
		// 	} else {
		// 		html, err := ioutil.ReadAll(resp.Body)
		// 		resp.Body.Close()
		// 		if err != nil {
		// 			fmt.Println("Error parsing Markdown:", err)
		// 		} else {
		// 			issues[i].SummaryHTML = string(html)
		// 			fmt.Println("Parsed into HTML:", issues[i].SummaryHTML)
		// 		}
		// 	}
		// }

		for _, label := range issue.Labels {
			fmt.Println("Located label:", label.Name)
			if label.Name == "bug" {
				issues[i].CssClass = "danger"
			} else if label.Name == "enhancement" {
				issues[i].CssClass = "success"
			}
		}
	}

	fmt.Println("Loaded", len(issues), "issues from GitHub")
	return issues, numIssues
}
