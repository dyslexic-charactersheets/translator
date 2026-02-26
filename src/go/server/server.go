package server

import (
	"github.com/dyslexic-charactersheets/translator/src/go/config"
	"github.com/dyslexic-charactersheets/translator/src/go/model"
	"fmt"
	"github.com/dyslexic-charactersheets/translator/src/go/log"
	// "io/ioutil"
	"github.com/dyslexic-charactersheets/translator/src/go/control"
	"golang.org/x/crypto/bcrypt"
	"github.com/bpowers/seshcookie"
	"html/template"
	"net/http"
	"strings"
	"strconv"
)

const SESSIONKEY = "93Yb8c59aASAf3kfT5xU8wz2GmfP4CbSNdhuvLxAdUqZnThbxuAAZu5AVWUrpsmXz47SYnvDcqr7TfNgLP8CpEpAmzGXNvMu72Scd4EAZGuepTQ7kWENemqr"

func RunTranslator(host string, debug int) {
	model.Debug = debug

	log.Space()
	log.Log("server", "Starting web server:", host)

	handler := http.NewServeMux()
	handler.HandleFunc("/home", control.DashboardHandler)
	handler.HandleFunc("/sources", control.SourcesHandler)
	handler.HandleFunc("/entries", control.EntriesHandler)
	handler.HandleFunc("/translate", control.TranslationHandler)
	handler.HandleFunc("/import", control.ImportHandler)
	handler.HandleFunc("/import/pot", control.ImportPotHandler)
	handler.HandleFunc("/import/progress", control.ImportProgressHandler)
	handler.HandleFunc("/import/abort", control.ImportAbortHandler)
	handler.HandleFunc("/export", control.ExportHandler)
	handler.HandleFunc("/export/po", control.ExportPoHandler)
	handler.HandleFunc("/live-export", control.LiveExportHandler)
	handler.HandleFunc("/users", control.UsersHandler)
	handler.HandleFunc("/users/add", control.UsersAddHandler)
	handler.HandleFunc("/users/del", control.UsersDelHandler)
	handler.HandleFunc("/users/masq", control.UsersMasqueradeHandler)
	handler.HandleFunc("/users/reinvite", control.UsersReinviteHandler)
	handler.HandleFunc("/users/show-invite", control.UsersShowInviteHandler)
	handler.HandleFunc("/users/renew-invite", control.UsersRenewInviteHandler)
	handler.HandleFunc("/account", control.AccountHandler)
	handler.HandleFunc("/account/password", control.SetPasswordHandler)
	handler.HandleFunc("/account/reclaim", control.AccountReclaimHandler)
	handler.HandleFunc("/authorise", control.AuthRedirectHandler)

	handler.HandleFunc("/api/setlead", control.APISetLeadHandler)
	handler.HandleFunc("/api/clearlead", control.APIClearLeadHandler)
	handler.HandleFunc("/api/entries", control.APIEntriesHandler)
	handler.HandleFunc("/api/translate", control.APITranslateHandler)
	handler.HandleFunc("/api/vote", control.APIVoteHandler)
	handler.HandleFunc("/api/lookup", control.APILookupHandler)

	// static files
	handler.Handle("/css/", http.FileServer(http.Dir("dist/htdocs")))
	handler.Handle("/bootstrap/", http.FileServer(http.Dir("dist/htdocs")))
	handler.Handle("/images/", http.FileServer(http.Dir("dist/htdocs")))
	handler.Handle("/js/", http.FileServer(http.Dir("dist/htdocs")))

	handler.Handle("/pdf/", http.StripPrefix("/pdf/", http.FileServer(http.Dir(config.Config.PDF.Path))))

	handler.HandleFunc("/", defaultHandler)

	authHandler := AuthHandler{handler}
	sessionHandler := seshcookie.NewHandler(&authHandler, SESSIONKEY, nil)

	listenPort := ":"+strconv.Itoa(config.Config.Server.Port)
	log.Log("server", "Listening on port:", listenPort)
	if err := http.ListenAndServe(listenPort, sessionHandler); err != nil {
		log.Log("server", "Error in ListenAndServe:", err)
	}

	log.Log("server", "Done.")
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Log("server", "Default handler")
	user := control.GetCurrentUser(r)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		http.Redirect(w, r, "/home", http.StatusFound)
	}
	return
}

//  AuthHandler: handle login/out, then pass other requests onto the basic handler
type AuthHandler struct {
	Handler http.Handler
}

type ReclaimFormData struct {
	Email  string
	Secret string
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if model.Debug >= 1 {
		log.Log("server", " -- %d database queries so far", model.QueryCount)
	}
	session := seshcookie.GetSession(r.Context())

	log.Space()
	log.Log("server", "Processing", r.Method, r.URL.Path)
	log.Log("server", "Using session: %#v\n", session)

	// bypass auth for static files
	segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	log.Log("server", "URL segments:", segments)
	first := segments[0]
	log.Log("server", "Checking URL segment:", first)
	switch first {
	case "css", "bootstrap", "images", "js", "pdf":
		log.Log("server", "Bypassing auth for", first)
		h.Handler.ServeHTTP(w, r)
		return
	}

	switch r.URL.Path {
	case "/users/masq":
		currentUser := control.GetCurrentUser(r)
		if currentUser != nil && currentUser.IsAdmin {
			email := r.FormValue("user")
			user := model.GetUserByEmail(email)
			if user != nil {
				// actually become that user
				log.Log("server", "Masquerading as user:", user.Email)
				session["user"] = user.Email
				session["masquerade"] = currentUser.Email
				log.Log("server", "altered session: %#v\n", session)
				http.Redirect(w, r, "/home", http.StatusFound)
				return
			}
		}
		log.Error("server", "Masquerade: Not admin!")
		http.Redirect(w, r, "/home", http.StatusFound)
		return

	case "/login":
		if r.Method != "POST" {
			http.ServeFile(w, r, "view/login.html")
			return
		}
		err := r.ParseForm()
		if err != nil {
			log.Error("server", "Error '%s' parsing form for %#v\n", err, r)
		}
		email := r.Form.Get("email")
		log.Log("server", "Login attempt", email);
		user := model.GetUserByEmail(email)
		password := r.Form.Get("password")

		if user == nil {
			log.Error("server", "Unknown user, redirecting")
			http.Redirect(w, r, "/login", 303)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			log.Error("server", "Password incorrect, redirecting", err)
			http.Redirect(w, r, "/login", 303)
			return
		}

		log.Log("server", "Login successful!", user.Name, "/", email)
		control.PingUser(user.Email)
		session["user"] = user.Email
		http.Redirect(w, r, "/home", http.StatusFound)
		return

	case "/logout":
		if email, ok := session["user"].(string); ok {
			log.Warn("server", "Logging out", email)
		}
		delete(session, "user")
		http.Redirect(w, r, "/login", http.StatusFound)
		return

	case "/account/reclaim/sent":
		http.ServeFile(w, r, "view/account_reclaim_sent.html")
		return

	case "/account/reclaim/done":
		http.ServeFile(w, r, "view/account_reclaim_done.html")
		return

	case "/account/reclaim/incorrect":
		http.ServeFile(w, r, "view/account_reclaim_incorrect.html")
		return

	case "/account/reclaim/nouser":
		http.ServeFile(w, r, "view/account_reclaim_nouser.html")
		return

	case "/account/reclaim":
		err := r.ParseForm()
		if err != nil {
			log.Error("server", "Error '%s' parsing form for %#v\n", err, r)
		}
		email := r.Form.Get("email")
		secret := r.Form.Get("secret")
		user := model.GetUserByEmail(email)
		log.Log("server", "Account reclaim: User at", email)

		if r.Method == "POST" {
			if user == nil {
				log.Log("server", "Account reclaim: Unknown user:", email)
				http.Redirect(w, r, "/account/reclaim/nouser", http.StatusFound)
				return
			}
			if secret == "" {
				secret := user.GenerateSecret()
				log.Log("server", "Account reclaim: Generating new secret:", secret)
				sendSecretEmail(user, secret)
				http.Redirect(w, r, "/account/reclaim/sent", http.StatusFound)
				return
			}

			log.Log("server", "Account reclaim: Comparing secret", secret)
			log.Log("server", "Account reclaim: Against hash", user.Secret)

			if user.VerifySecret(secret) {
				password := r.Form.Get("password")
				password2 := r.Form.Get("password2")
				if password != "" && password == password2 {
					log.Log("server", "Account reclaim: Setting password")
					hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
					if err == nil {
						user.Password = string(hash)
						user.Secret = ""
						user.Save()
						log.Log("server", "Account reclaim: Saved", string(hash))
					} else {
						log.Error("server", "Account reclaim: Error", err)
					}
					http.Redirect(w, r, "/account/reclaim/done", http.StatusFound)
					return
				} else {
					log.Warn("server", "Account reclaim: Redirecting to password form")
					http.Redirect(w, r, "/account/reclaim?email="+email+"&secret="+secret, http.StatusFound)
					return
				}
				return
			} else {
				log.Error("server", "Account reclaim: Cannot verify secret")
				http.Redirect(w, r, "/account/reclaim/incorrect", http.StatusFound)
				return
			}
		} else if r.Method == "GET" {
			if secret == "" {
				http.ServeFile(w, r, "view/account_reclaim.html")
				return
			}

			if user == nil {
				log.Error("server", "Account reclaim: Unknown user:", email)
				http.Redirect(w, r, "/account/reclaim/nouser", http.StatusFound)
				return
			}

			log.Log("server", "Account reclaim: Comparing secret", secret)
			log.Log("server", "Account reclaim: Against hash", user.Secret)
			if user.VerifySecret(secret) {
				log.Log("server", "Account reclaim: Showing password form")
				data := ReclaimFormData{
					Email:  email,
					Secret: secret,
				}
				t, _ := template.ParseFiles("view/account_reclaim_set_password.html")
				t.Execute(w, data)
				return
			} else {
				http.Redirect(w, r, "/account/reclaim/incorrect", http.StatusFound)
				return
			}
		}
	}

	if _, ok := session["user"]; !ok {
		log.Warn("server", "Not logged in, redirecting to login")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	log.Log("server", "Delivering")
	control.PingCurrentUser(r)
	h.Handler.ServeHTTP(w, r)
}

func sendSecretEmail(user *model.User, secret string) {
	scheme := "http"
	if config.Config.Server.SSL {
		scheme = "https"
	}
	host := config.Config.Server.Hostname
	url := "%s://%s/account/reclaim?email=%s&secret=%s"
	url = fmt.Sprintf(url, scheme, host, user.Email, secret)
	log.Log("server", "Recovery URL", url)

	msg := `Subject: Your account at the Character Sheets Translator
Content-Type: text/plain; charset="UTF-8"

This is your password reclaim email for the Dyslexic Character Sheets Translator

To set your password, click here:

`
	msg = msg + url

	log.Log("server", "Sending message to", user.Email)

	if ok := config.SendMail(user.Email, msg); ok {
		log.Log("server", "Sent password reclaim email to:", user.Email)
	}
}
