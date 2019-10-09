package main

import "net/http"

// DashboardMiddleware protects dashboard links behind google auth
func DashboardMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if authProvider.IsAuthenticated(r) {
		session, _ := authProvider.Session.Store.Get(r, sessionName)

		if s := session.Values["app_user"]; s != nil {
			u, err := authProvider.WithUserFromSession(s)

			if err != nil {
				http.Redirect(w, r, "/google/login", http.StatusFound)
			}
			if u.ID != s {
				http.Redirect(w, r, "/google/login", http.StatusFound)
			}

			next(w, r)
		} else {
			http.Redirect(w, r, "/google/login", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/google/login", http.StatusFound)
	}
}

// ShareCodeOrSSOMiddleware protects dashboard links behind google auth
func ShareCodeOrSSOMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == "POST" {
		next(w, r)
		return
	}

	if r.URL.Query().Get("share_code") != "" {
		next(w, r)
		return
	}

	if authProvider.IsAuthenticated(r) {
		session, _ := authProvider.Session.Store.Get(r, sessionName)

		if s := session.Values["app_user"]; s != nil {
			u, err := authProvider.WithUserFromSession(s)

			if err != nil {
				http.Redirect(w, r, "/google/login", http.StatusFound)
			}
			if u.ID != s {
				http.Redirect(w, r, "/google/login", http.StatusFound)
			}

			next(w, r)
		} else {
			http.Redirect(w, r, "/google/login", http.StatusFound)
		}
	} else {
		http.Redirect(w, r, "/google/login", http.StatusFound)
	}
}

// ShareCodeMiddleware protects dashboard links behind google auth
func ShareCodeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Query().Get("share_code") != "" {
		next(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
