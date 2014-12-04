package authserver

import (
	"code.google.com/p/go.crypto/bcrypt"
	"code.google.com/p/goprotobuf/proto"
	ferr "github.com/jonas747/fortia/error"
	"github.com/jonas747/fortia/messages"
	"github.com/jonas747/fortia/rest"
	"net/http"
	"time"
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

func AuthRequiredMiddleWare(r *rest.Request) bool {
	loggedIn, user, err := CheckSession(r)
	if server.HandleFortiaError(r, err) {
		return false
	}
	if !loggedIn {
		return false
	}

	// Set the user additional info
	r.AdditionalData["user"] = user
	return true
}

// Returns true if logged on, false if not, also returns the user if logged on
func CheckSession(r *rest.Request) (bool, string, ferr.FortiaError) {
	// Get the session cookie
	cookie, nErr := r.Request.Cookie("fortia-session")
	if nErr != nil {
		return false, "", ferr.Wrap(nErr, "")
	}

	session := cookie.Value
	user, err := authDb.CheckSessionToken(session)
	if err != nil {
		return false, "", err
	}
	if user == "" {
		// Expired
		return false, "", nil
	}
	return true, user, nil
}

// /login
func handleLogin(r *rest.Request, body interface{}) {
	bl := body.(*BodyLogin)
	correctPw, err := authDb.CheckUserPw(bl.Username, bl.Pw)
	if server.HandleFortiaError(r, err) {
		return
	}
	if !correctPw {
		logger.Warna("User tried logging in with invalid password", map[string]interface{}{"remoteaddr": r.Request.RemoteAddr, "user": bl.Username})
		// w.WriteHeader(http.StatusBadRequest)
		// w.Write(rest.ApiError(rest.ErrWrongLoginDetails, "Username or password incorrect"))
		r.WriteResponse(rest.NewPlainResponse(rest.ErrWrongLoginDetails, "Login details incorrect"), http.StatusBadRequest)
		return
	}
	session, err := authDb.LoginUser(bl.Username, 3600*24) // 24 hours
	if server.HandleFortiaError(r, err) {
		return
	}

	expires := time.Now().Add(time.Duration(24) * time.Hour)
	// Assemble the cookie
	cookie := &http.Cookie{
		Name:    "fortia-session",
		Value:   session,
		Path:    "/",
		Expires: expires,
	}
	http.SetCookie(r.RW, cookie)

	r.WriteResponse(rest.NewPlainResponse(0, ""), http.StatusOK)
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
func handleRegister(r *rest.Request, body interface{}) {
	rBody := body.(*BodyRegister)

	// Make sure all details are valid
	if !validateUsername(rBody.Username) {
		r.WriteResponse(rest.NewPlainResponse(rest.ErrInvalidUsername, "Username is not valid"), http.StatusBadRequest)
		return
	}
	if !validateEmail(rBody.Email) {
		r.WriteResponse(rest.NewPlainResponse(rest.ErrInvalidEmail, "Email is not valid"), http.StatusBadRequest)
		return
	}
	if !validatePassword(rBody.Pw) {
		r.WriteResponse(rest.NewPlainResponse(rest.ErrInavlidPassword, "Password is not valid"), http.StatusBadRequest)
		return
	}

	existingInfo, err := authDb.GetUserInfo(rBody.Username)
	if server.HandleFortiaError(r, err) {
		return
	}
	if existingInfo.Email != "" {
		r.WriteResponse(rest.NewPlainResponse(rest.ErrUsernameTaken, "Username taken"), http.StatusBadRequest)
		return
	}

	pwHash, nErr := bcrypt.GenerateFromPassword([]byte(rBody.Pw), bcrypt.DefaultCost)
	if nErr != nil {
		err := ferr.Wrapa(nErr, "", map[string]interface{}{"user": rBody.Username})
		server.HandleFortiaError(r, err)
		return
	}

	user := &UserInfo{
		Name:   rBody.Username,
		PwHash: pwHash,
		Email:  rBody.Email,
	}

	err = authDb.SetUserInfo(user)
	if server.HandleFortiaError(r, err) {
		return
	}
	r.WriteResponse(rest.NewPlainResponse(0, ""), http.StatusOK)
	logger.Info("User " + rBody.Username + " Sucessfully registered!")
}

// /me
func handleMe(r *rest.Request, body interface{}) {
	user, _ := r.AdditionalData.GetString("user") // Set by authentication middleware

	info, err := authDb.GetUserInfo(user)
	if server.HandleFortiaError(r, err) {
		return
	}

	wireInfo := messages.UserInfo{
		Name:  proto.String(info.Name),
		Email: proto.String(info.Email),

		Worlds:   info.Worlds,
		Role:     proto.Int(info.Role),
		DonorLvl: proto.Int(info.DonorLvl),
	}

	wireResp := messages.MeResponse{
		Info: &wireInfo,
	}

	r.WriteResponse(&wireResp, http.StatusOK)
}

// /worlds
func handleWorlds(r *rest.Request, body interface{}) {
	params := r.Request.URL.Query()
	wname := params.Get("world")
	var wireWorlds []*messages.WorldInfo
	if wname == "" {
		// Return info on all worlds
		worlds, err := authDb.GetWorldListing()
		if server.HandleFortiaError(r, err) {
			return
		}

		wireWorlds := make([]*messages.WorldInfo, len(worlds))
		for k, v := range worlds {
			wireWorld := &messages.WorldInfo{
				Name:    proto.String(v.Name),
				Started: proto.Int(v.Started),
				Players: proto.Int(v.Players),
				Size:    proto.Int(v.Size),
			}
			wireWorlds[k] = wireWorld
		}
	} else {
		info, err := authDb.GetWorldInfo(wname)
		if server.HandleFortiaError(r, err) {
			return
		}
		wireInfo := &messages.WorldInfo{
			Name:    proto.String(info.Name),
			Started: proto.Int(info.Started),
			Players: proto.Int(info.Players),
			Size:    proto.Int(info.Size),
		}
		wireWorlds = []*messages.WorldInfo{wireInfo}
	}

	wireResp := &messages.WorldsResponse{
		Info: wireWorlds,
	}
	r.WriteResponse(wireResp, http.StatusOK)
}
