package server

/*
import (
	"fmt"
	"net/http"
	"net/url"
)
*/

type AuthUser struct {
	username string
	password string
	valid    bool
}

func handleAuth(c *Connection, cs *CredentialSet, p AuthProcessor) string {
	message := c.GetMsg()
	cred := cs

	// handle the command
	cmd, args, err := parseCommand(message)
	if err != nil {
		return SyntaxErr
	}

	switch {
	case cmd == "USER" && args == "":
		return AnonUserDenied
	case cmd == "USER" && args != "":
		cred.Username = args
		return UsrNameOkNeedPass
	case cmd == "PASS" && args == "":
		return SyntaxErr
	case cmd == "PASS" && args != "" && cred.Username != "":
		cred.Password = args
	}
	if p != nil {
		cred.Valid = false
		// cred = p(cred)
	} else {
		cred.Valid = true
	}

	if cred.Valid == true {
		return UsrLoggedInProceed
	}
	cred.Username = ""
	cred.Password = ""
	return AuthFailureTryAgain
}

func handleLogin(message string, user *AuthUser) string {
	// Handle login operations

	cmd, args, err := parseCommand(message)
	if err != nil {
		return SyntaxErr
	}

	switch {
	case cmd == "USER" && args == "":
		return AnonUserDenied
	case cmd == "USER" && args != "":
		user.username = args
		return UsrNameOkNeedPass
	case cmd == "PASS" && args == "":
		return SyntaxErr
	case cmd == "PASS" && args != "" && user.username != "":
		user.password = args
	}

	user.Authenticate()

	if user.valid == true {
		return UsrLoggedInProceed
	} else {
		user.username = ""
		user.password = ""
		return AuthFailureTryAgain
	}
}

func (user *AuthUser) Authenticate() {
	// Authenticate user against data.ambition

	/*
		uri := fmt.Sprint(DataUrl, "/api/upload/ftp-auth/")
		resp, err := http.PostForm(uri,
			url.Values{
				"username": {user.username},
				"password": {user.password}})

		if resp.StatusCode == http.StatusOK && err == nil {
			user.valid = true
		} else {
			user.valid = false
		}
	*/
	user.valid = true
}
