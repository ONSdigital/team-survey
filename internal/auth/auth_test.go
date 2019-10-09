package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/dghubble/sessions"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var authProvider Provider

func deleteTestDB() {
	if _, err := os.Stat("test.db"); err == nil {
		os.Remove("test.db")
	}
}

func mockAuthenticatedRequest() http.Request {
	var s = securecookie.New([]byte("very-secret-1234"), []byte("very-secret-1234"))
	recorder := httptest.NewRecorder()
	value := map[string]string{
		"user":     "joe-bloggs",
		"app_user": "1",
	}

	if encoded, err := s.Encode(authProvider.Session.UserKey, value); err == nil {
		cookie := &http.Cookie{
			Name:  authProvider.Session.Name,
			Value: encoded,
		}
		authProvider.Session.Store.Save(recorder, authProvider.Session.Store.New(authProvider.Session.Name))
		http.SetCookie(recorder, cookie)
	}

	return http.Request{
		Header: http.Header{
			"Cookie": recorder.HeaderMap["Set-Cookie"],
		},
	}
}

func TestSetup(t *testing.T) {
	deleteTestDB()
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		t.Errorf("Expected not nil, got %s", err)
	}

	authProvider = CreateProvider(db, sessions.NewCookieStore([]byte("team-surveytest"), nil), "team-surveytest", "team-surveytest", "team-surveytest")
	AutoMigrateModels(authProvider.DB)
}

func TestCreateProvider(t *testing.T) {
	p := CreateProvider(authProvider.DB, sessions.NewCookieStore([]byte("pvdr"), nil), "pvdr", "pvdr", "pvdr")

	if p.Session.Name != "pvdr" {
		t.Errorf("Expected 'pvdr', got %s", authProvider.Session.Name)
	}
	if p.Session.Secret != "pvdr" {
		t.Errorf("Expected 'pvdr', got %s", authProvider.Session.Secret)
	}
	if p.Session.UserKey != "pvdr" {
		t.Errorf("Expected 'pvdr', got %s", authProvider.Session.Secret)
	}
}

func TestCreateUser(t *testing.T) {
	if _, err := authProvider.CreateUser(User{Email: randomdata.Email()}); err != nil {
		t.Errorf("Expected not nil, got %s", err)
	}
}

func TestGetUserBySSOIDDoesNotExist(t *testing.T) {
	user, err := authProvider.GetUserBySSOID("123-123")
	if err.Error() != "record not found" {
		t.Errorf("Expected 'record nor found, got %s", err)
	}

	if user.Exists {
		t.Errorf("Expected false, got true")
	}
}

func TestGetUserBySSOIDDoesExist(t *testing.T) {
	email := randomdata.Email()
	authProvider.CreateUser(User{Email: email, SSOID: email})
	user, err := authProvider.GetUserBySSOID(email)
	if err != nil {
		t.Errorf("Expected not nil, got %s", err)
	}

	if !user.Exists {
		t.Errorf("Expected true, got false")
	}
}

func TestWithUserFromSessionInvalid(t *testing.T) {
	_, err := authProvider.WithUserFromSession(0)

	if err.Error() != "record not found" {
		t.Errorf("Expected 'record not found', got %s", err)
	}
}

func TestWithUserFromSessionValid(t *testing.T) {
	email := randomdata.Email()
	newUser, err := authProvider.CreateUser(
		User{
			Email: email,
			SSOID: email,
		},
	)

	if err != nil {
		t.Errorf("Expected not nil, got %s", err)
	}

	user, err := authProvider.WithUserFromSession(newUser.ID)

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if user.ID != newUser.ID {
		t.Errorf("Expected %d, got %d", user.ID, newUser.ID)
	}
}

func TestGetUserFromContext(t *testing.T) {
	email := randomdata.Email()
	newUser, err := authProvider.CreateUser(
		User{
			Email: email,
			SSOID: email,
		},
	)

	if err != nil {
		t.Errorf("Expected not nil, got %s", err)
	}
	r := http.Request{}
	ctx := context.WithValue(r.Context(), AppUserKey, newUser)

	user := authProvider.GetUserFromContext(ctx)

	if user.ID != newUser.ID {
		t.Errorf("Expected %d, got %d", user.ID, newUser.ID)
	}
}

func TestGetAllUsers(t *testing.T) {
	users, err := authProvider.GetAllUsers()

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	userCount := len(users)

	if userCount == 0 {
		t.Errorf("Expected > 0, got %d", userCount)
	}

	authProvider.CreateUser(
		User{
			Email: randomdata.Email(),
			SSOID: randomdata.Email(),
		},
	)

	users, err = authProvider.GetAllUsers()

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	expectedCount := (userCount + 1)

	if len(users) != expectedCount {
		t.Errorf("Expected %d, got %d", expectedCount, len(users))
	}
}

func TestIsAuthenticatedFalse(t *testing.T) {
	req := http.Request{}

	isAuthenticated := authProvider.IsAuthenticated(&req)

	if isAuthenticated {
		t.Errorf("Expected false, got true")
	}
}

func TestIsAuthenticatedTrue(t *testing.T) {
	request := mockAuthenticatedRequest()
	if _, err := authProvider.Session.Store.Get(&request, authProvider.Session.Name); err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if !authProvider.IsAuthenticated(&request) {
		t.Errorf("Expected true, got false")
	}
}

func TestSaveSession(t *testing.T) {
	r := httptest.NewRecorder()
	if err := authProvider.saveSession(r, uint(1), "12343876"); err != nil {
		t.Errorf("Expected nil, got %s", err)
	}
}

func TestAddAppUserContextUnauthenticated(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mock"))
	})
	handlerToTest := authProvider.AddAppUserContext(nextHandler)
	req := httptest.NewRequest("GET", "http://testing", nil)
	rr := httptest.NewRecorder()
	handlerToTest.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected 200, got %d", rr.Code)
	}
}

func TestAddAppUserContextAuthenticated(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Mock"))
	})
	handlerToTest := authProvider.AddAppUserContext(nextHandler)
	req := mockAuthenticatedRequest()
	rr := httptest.NewRecorder()
	handlerToTest.ServeHTTP(rr, &req)

	if rr.Code != 200 {
		t.Errorf("Expected 200, got %d", rr.Code)
	}
}

func TestLogoutHandler(t *testing.T) {
	req, _ := http.NewRequest("POST", "/logout", nil)
	rr := httptest.NewRecorder()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/logout", authProvider.LogoutHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != 302 {
		t.Errorf("Expected 302, got %d", rr.Code)
	}
}

func TestGetGoogleUserFromContext(t *testing.T) {
	request := mockAuthenticatedRequest()
	_, err := getGoogleUserFromContext(&request)

	if err.Error() != "google: Context missing Google User" {
		t.Errorf("Expected 'Context missing Google User', got %s", err)
	}
}

func TestTearDown(t *testing.T) {
	deleteTestDB()
}
