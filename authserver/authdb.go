package authserver

import (
	ferr "github.com/jonas747/fortia/error"
)

// Represents a user
type UserInfo struct {
	Name   string // username
	PwHash []byte // password hash
	Email  string // email

	Worlds   []string // Worlds this user is active in
	Role     int      // Users Role (0 - normal, 1 - server mod, 2 - global mod etc (These role numbers are just mockups))
	DonorLvl int      // Donor lvl
}

// Information about a world
type WorldInfo struct {
	Name    string // Name of the world
	Started int    // What time the world started in unix epoch time
	Players int    // Number of players on this world
	Size    int    // The size of this world in blocks
}

type AuthDB interface {
	LoginUser(user string, duration int) (token string, err ferr.FortiaError) // Logs the specified user in for "duration" seconds
	CheckUserPw(user, pw string) (correctpw bool, err ferr.FortiaError)       // Returns true if the password provided is the right one
	ExtendSessionToken(token string, newDuration int) ferr.FortiaError        // Sets the tll of specified token to newDuration
	CheckSessionToken(token string) (user string, err ferr.FortiaError)       // Returns the user the specified token belongs to, "" if not found

	GetUserInfo(user string) (*UserInfo, ferr.FortiaError) // Returns info about the specified user
	SetUserInfo(info *UserInfo) ferr.FortiaError           // Overwrite the users info with the info provided
	EditUserInfoFields(fields map[string]interface{})      // Sets the users fields with the fields provided

	GetWorldListing() ([]*WorldInfo, ferr.FortiaError)        // Returns all worlds
	GetWorldInfo(world string) (*WorldInfo, ferr.FortiaError) // Returns info about specified world
	SetWorldInfo(info *WorldInfo) ferr.FortiaError            // Saves the specified world info to the database
}
