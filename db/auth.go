package db

import (
	"github.com/jonas747/fortia/errors"
)

// Represents a user
type AuthUserInfo struct {
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
	LoginUser(user string, duration int) (token string, err errors.FortiaError) // Logs the specified user in for "duration" seconds
	CheckUserPw(user, pw string) (correctpw bool, err errors.FortiaError)       // Returns true if the password provided is the right one
	ExtendSessionToken(token string, newDuration int) errors.FortiaError        // Sets the tll of specified token to newDuration
	CheckSessionToken(token string) (user string, err errors.FortiaError)       // Returns the user the specified token belongs to, "" if not found

	GetUserInfo(user string) (*AuthUserInfo, errors.FortiaError)                      // Returns info about the specified user
	SetUserInfo(info *AuthUserInfo) errors.FortiaError                                // Overwrite the users info with the info provided
	EditUserInfoFields(user string, fields map[string]interface{}) errors.FortiaError // Sets the users fields with the fields provided
	EditUserWorlds(user string, add []string, del []string) errors.FortiaError        // Add the worlds in add and removes the worlds in remove

	GetWorldListing() ([]*WorldInfo, errors.FortiaError)                          // Returns all worlds
	GetWorldInfo(world string) (*WorldInfo, errors.FortiaError)                   // Returns info about specified world
	SetWorldInfo(info *WorldInfo) errors.FortiaError                              // Saves the specified world info to the database
	EditWorldInfo(world string, fields map[string]interface{}) errors.FortiaError // (Redundant perhaps?) Overwrites the stored world info fields with the fields provided
}
