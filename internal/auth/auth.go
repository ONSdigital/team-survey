package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dghubble/gologin/google"
	"github.com/dghubble/sessions"
	"github.com/jinzhu/gorm"
	ggl "google.golang.org/api/oauth2/v2"
)

type key string

const (
	// AppUserKey is the session key for the app user
	AppUserKey key = "AppUser"
)

// Model is the basic data layout to be inherited
type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"updated"`
	DeletedAt *time.Time `json:"-"`
}

// User is a user object
type User struct {
	Model
	SSOID           string `gorm:"unique;not null;unique_index"`
	Email           string `gorm:"unique;not null;unique_index"`
	ProfileImageURL string
	Exists          bool
	APIToken        string `json:"-"`
}

// ProviderSession is a struct for a session
type ProviderSession struct {
	Store   *sessions.CookieStore
	Name    string
	Secret  string
	UserKey string
}

// Provider is a struct for auth
type Provider struct {
	DB      *gorm.DB
	Session ProviderSession
}

// AutoMigrateModels ensures all models are present in the database
func AutoMigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(&User{}).Error
}

// CreateProvider will initialize and return a new Provider
func CreateProvider(db *gorm.DB, sessionStore *sessions.CookieStore, sessionName, sessionSecret, sessionUserKey string) Provider {
	return Provider{
		DB: db,
		Session: ProviderSession{
			Store:   sessionStore,
			Name:    sessionName,
			Secret:  sessionSecret,
			UserKey: sessionUserKey,
		},
	}
}

// GetUserBySSOID returns a user by SSOID if it exists
func (a Provider) GetUserBySSOID(ssoID string) (User, error) {
	var user User
	return user, a.DB.Where("sso_id = ?", ssoID).First(&user).Error
}

// CreateUser creates a new user in the database
func (a Provider) CreateUser(user User) (User, error) {
	user.Exists = true
	return user, a.DB.Save(&user).Error
}

// WithUserFromSession loads a user from a session where it can
func (a Provider) WithUserFromSession(s interface{}) (User, error) {
	var user User
	return user, a.DB.First(&user, s).Error
}

// GetUserFromContext loads a user from context where it can
func (a Provider) GetUserFromContext(ctx context.Context) User {
	var user User

	if ctx.Value(AppUserKey) != nil {
		user = ctx.Value(AppUserKey).(User)
	}

	return user
}

// GetAllUsers returns all users
func (a Provider) GetAllUsers() ([]User, error) {
	var users []User
	return users, a.DB.Find(&users).Error
}

// AddAppUserContext adds user information context
func (a Provider) AddAppUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := a.Session.Store.Get(r, a.Session.Name)
		if err != nil {
			next.ServeHTTP(w, r)
		} else {
			if s := session.Values["app_user"]; s != nil {
				u, _ := a.WithUserFromSession(s)
				ctx := context.WithValue(r.Context(), AppUserKey, u)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				next.ServeHTTP(w, r)
			}
		}
	})
}

// IsAuthenticated returns true if the user has a signed session cookie.
func (a Provider) IsAuthenticated(req *http.Request) bool {
	if _, err := a.Session.Store.Get(req, a.Session.Name); err == nil {
		return true
	}
	return false
}

func getGoogleUserFromContext(req *http.Request) (*ggl.Userinfoplus, error) {
	return google.UserFromContext(req.Context())
}

// IssueSession issues a cookie session after successful Google login
func (a Provider) IssueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		googleUser, err := getGoogleUserFromContext(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		appUser, _ := a.GetUserBySSOID(googleUser.Id)

		if !appUser.Exists {
			appUser.Email = googleUser.Email
			appUser.SSOID = googleUser.Id
		}

		appUser.ProfileImageURL = googleUser.Picture

		appUser, err = a.CreateUser(appUser)
		if err != nil {
			log.Printf("%s", err)
		}

		if err := a.saveSession(w, appUser.ID, googleUser.Id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, req, "/admin/dashboard", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

// LogoutHandler destroys the session on POSTs and redirects to home.
func (a Provider) LogoutHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		a.Session.Store.Destroy(w, a.Session.Name)
	}
	http.Redirect(w, req, "/", http.StatusFound)
}

func (a Provider) saveSession(w http.ResponseWriter, appUserID uint, googleUserID string) error {
	session := a.Session.Store.New(a.Session.Name)
	session.Values["app_user"] = appUserID
	session.Values[a.Session.UserKey] = googleUserID
	return session.Save(w)
}
