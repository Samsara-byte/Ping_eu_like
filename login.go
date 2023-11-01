package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	store         *sessions.CookieStore
	sessionName   = "session-name"
	sessionSecret = "your-session-secret"
)

func init() {
	store = sessions.NewCookieStore([]byte(sessionSecret))
}

func validateUser(user User, isLogin bool) bool {
	fmt.Println("validate user")
	if isLogin {
		if strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
			return false
		}
	} else {
		if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
			return false
		}
	}

	return true
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Registering...")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := User{}
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	user.Username = r.FormValue("username")
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")

	if !validateUser(user, false) {
		http.Error(w, "Invalid user input.", http.StatusBadRequest)
		return
	}

	if err := saveUserToDatabase(user); err != nil {
		http.Error(w, fmt.Sprintf("Error saving user: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully registered.")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "website/login.html")
		return
	}

	fmt.Println("Login")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := User{}
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")

	if !validateUser(user, true) {
		http.Error(w, "Invalid user input.", http.StatusBadRequest)
		return
	}
	exists, err := userExistsInDatabase(user.Email, user.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking user existence: %v", err), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "User does not exist or invalid credentials.", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = true
	session.Save(r, w)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, "website/dashboard.html")
}

func saveUserToDatabase(user User) error {

	fmt.Println("Saving user to")
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func userExistsInDatabase(email, password string) (bool, error) {

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND password = ?", email, password).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
