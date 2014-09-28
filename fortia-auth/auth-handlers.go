package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	//"fmt"
	ferr "github.com/jonas747/fortia/error"
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

// Returns true if there was an error, false is error == nil
func handleFortiaError(w http.ResponseWriter, r *http.Request, err ferr.FortiaError) bool {
	// Add some additional details
	if err == nil {
		return false
	}
	err.SetData("remoteaddr", r.RemoteAddr)
	err.SetData("path", r.URL.Path)
	logger.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(ErrServerError)
	return true
}

func handleUnauthorized(w http.ResponseWriter, r *http.Request, user string) {
	logger.Warna("Unauthorized", map[string]interface{}{"remoteaddr": r.RemoteAddr, "user": user, "path": r.URL.Path})
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(ApiError(ErrCodeInvalidSessionCookie, "Session expired"))
}

// Returns true if logged on, false if not
// Also writes the status code and error
func checkSession(w http.ResponseWriter, r *http.Request, user string) bool {
	// Get the session cookie
	cookie, nerr := r.Cookie("fortia-session")
	if nerr != nil {
		// No session cookie
		handleUnauthorized(w, r, user)
		return false
	}

	session := cookie.Value
	ttl, err := authDb.CheckSessionToken(user, session)
	if handleFortiaError(w, r, err) {
		return false
	}

	if ttl < 0 {
		// Expired
		handleUnauthorized(w, r, user)
		return false
	}

	return true
}

// /login
func handleLogin(w http.ResponseWriter, r *http.Request, body interface{}) {
	bl := body.(*BodyLogin)
	correctPw, err := authDb.CheckUserPw(bl.Username, bl.Pw)
	if handleFortiaError(w, r, err) {
		return
	}
	if !correctPw {
		logger.Warna("User tried logging in with invalid password", map[string]interface{}{"remoteaddr": r.RemoteAddr, "user": bl.Username})
		return
	}
	session, err := authDb.LoginUser(bl.Username)
	if handleFortiaError(w, r, err) {
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

// Returns true if the email is okay
func validateEmail(email string) bool {
	// TODO
	return true
}

// Returns true if the username is okay
func validateUsername(user string) bool {
	// TODO
	return true
}

func validatePassword(pw string) bool {
	// TODO
	return true
}

// /register
func handleRegister(w http.ResponseWriter, r *http.Request, body interface{}) {
	rBody := body.(*BodyRegister)

	// Make sure all details are valid
	ok := validateUsername(rBody.Username)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ApiError(ErrCodeInvalidUsername, "Username is not valid"))
		return
	}
	ok = validateUsername(rBody.Email)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ApiError(ErrCodeInvalidEmail, "Email is not valid"))
		return
	}
	ok = validateUsername(rBody.Pw)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ApiError(ErrCodeInavlidPassword, "Password is not valid"))
		return
	}

	existingInfo, err := authDb.GetUserInfo(rBody.Username)
	if handleFortiaError(w, r, err) {
		return
	}
	_, ok = existingInfo["name"]
	if ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ApiError(ErrCodeUsernameTaken, "Username is taken"))
		return
	}

	pwHash, nerr := bcrypt.GenerateFromPassword([]byte(rBody.Pw), bcrypt.DefaultCost)
	if nerr != nil {
		err := ferr.Wrapa(err, "", map[string]interface{}{"user": rBody.Username})
		handleFortiaError(w, r, err)
		return
	}

	infoMap := make(map[string]interface{})
	infoMap["name"] = rBody.Username
	infoMap["pw"] = string(pwHash)
	infoMap["mail"] = rBody.Email

	err = authDb.SetUserInfo(rBody.Username, infoMap)
	if !handleFortiaError(w, r, err) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"ok\": true}"))
	}
	logger.Info("User " + rBody.Username + " Sucessfully registered!")
}

// /getinfo
func handleGetInfo(w http.ResponseWriter, r *http.Request, body interface{}) {
	params := r.URL.Query()
	user := params.Get("user")
	loggedIn := checkSession(w, r, user)
	if !loggedIn {
		return
	}

	info, err := authDb.GetUserInfo(user)
	if handleFortiaError(w, r, err) {
		return
	}

	// Add ok true, because everythign is ok
	info["ok"] = "true"
	// We dont want to send the password hash...
	info["pw"] = ""

	out, nerr := json.Marshal(info)
	if nerr != nil {
		err = ferr.Wrapa(nerr, "Marshal error", map[string]interface{}{"user": user})
		handleFortiaError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
