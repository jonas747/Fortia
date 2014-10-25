package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	//"fmt"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/resterrors"
	"net/http"
)

type BodyRegister struct {
	Username string
	Pw       string
	Email    string
}

type BodyLogin struct {
	Username string
	Pw       string
}

// Returns true if logged on, false if not
// Also writes the status code and error
func checkSession(w http.ResponseWriter, r *http.Request) (bool, string) {
	// Get the session cookie
	cookie, nerr := r.Cookie("fortia-session")
	if nerr != nil {
		// No session cookie
		server.HandleUnauthorized(w, r, "")
		return false, ""
	}

	session := cookie.Value
	user, err := authDb.CheckSessionToken(session)
	if server.HandleFortiaError(w, r, err) {
		return false, ""
	}

	if user == "" {
		// Expired
		server.HandleUnauthorized(w, r, "")
		return false, ""
	}

	return true, user
}

// /login
func handleLogin(w http.ResponseWriter, r *http.Request, body interface{}) {
	bl := body.(*BodyLogin)
	correctPw, err := authDb.CheckUserPw(bl.Username, bl.Pw)
	if server.HandleFortiaError(w, r, err) {
		return
	}
	if !correctPw {
		logger.Warna("User tried logging in with invalid password", map[string]interface{}{"remoteaddr": r.RemoteAddr, "user": bl.Username})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resterrors.ApiError(resterrors.ErrWrongLoginDetails, "Username or password incorrect"))
		return
	}
	session, err := authDb.LoginUser(bl.Username, 3600*24) // 24 hours, make this expand as the users does stuff
	if server.HandleFortiaError(w, r, err) {
		return
	}
	// Assemble the cookie
	cookie := &http.Cookie{
		Name:  "fortia-session",
		Value: session,
		Path:  "/",
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("{\"ok\": true}"))
	logger.Info("User Logged in: " + bl.Username)
}

// TODO: Actual validation
// TODO: Return a string with description of what was wrong
// Returns true if the email is okay
func validateEmail(email string) bool {
	// TODO
	if len(email) < 4 {
		return false
	}
	return true
}

// Returns true if the username is okay
func validateUsername(user string) bool {
	// TODO
	if len(user) < 1 {
		return false
	}
	return true
}

func validatePassword(pw string) bool {
	// TODO
	if len(pw) < 1 {
		return false
	}
	return true
}

// /register
func handleRegister(w http.ResponseWriter, r *http.Request, body interface{}) {
	rBody := body.(*BodyRegister)

	// Make sure all details are valid
	if !validateUsername(rBody.Username) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resterrors.ApiError(resterrors.ErrInvalidUsername, "Username is not valid"))
		return
	}
	if !validateEmail(rBody.Email) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resterrors.ApiError(resterrors.ErrInvalidEmail, "Email is not valid"))
		return
	}
	if !validatePassword(rBody.Pw) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resterrors.ApiError(resterrors.ErrInavlidPassword, "Password is not valid"))
		return
	}

	existingInfo, err := authDb.GetUserInfo(rBody.Username)
	if server.HandleFortiaError(w, r, err) {
		return
	}
	_, ok := existingInfo["name"]
	if ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resterrors.ApiError(resterrors.ErrUsernameTaken, "Username is taken"))
		return
	}

	pwHash, nerr := bcrypt.GenerateFromPassword([]byte(rBody.Pw), bcrypt.DefaultCost)
	if nerr != nil {
		err := ferr.Wrapa(err, "", map[string]interface{}{"user": rBody.Username})
		server.HandleFortiaError(w, r, err)
		return
	}

	infoMap := make(map[string]interface{})
	infoMap["name"] = rBody.Username
	infoMap["pw"] = string(pwHash)
	infoMap["mail"] = rBody.Email

	err = authDb.SetUserInfo(rBody.Username, infoMap)
	if !server.HandleFortiaError(w, r, err) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"ok\": true}"))
	}
	logger.Info("User " + rBody.Username + " Sucessfully registered!")
}

// /me
func handleMe(w http.ResponseWriter, r *http.Request, body interface{}) {
	loggedIn, user := checkSession(w, r)
	if !loggedIn {
		return
	}

	info, err := authDb.GetUserInfo(user)
	if server.HandleFortiaError(w, r, err) {
		return
	}

	// Add ok true, because everythign is ok
	info["ok"] = "true"
	// We dont want to send the password hash...
	info["pw"] = ""

	out, nerr := json.Marshal(info)
	if nerr != nil {
		err = ferr.Wrapa(nerr, "Marshal error", map[string]interface{}{"user": user})
		server.HandleFortiaError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// /worlds
func handleWorlds(w http.ResponseWriter, r *http.Request, body interface{}) {
	params := r.URL.Query()
	wname := params.Get("world")

	info, err := authDb.GetWorldInfo(wname)
	if server.HandleFortiaError(w, r, err) {
		return
	}
	out, nerr := json.Marshal(info)
	if nerr != nil {
		server.HandleFortiaError(w, r, ferr.Wrap(nerr, ""))
		return
	}

	w.Write(out)
}
