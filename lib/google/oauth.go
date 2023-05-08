package google

import (
	"os"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
)

func CreateGoogleOAuthURL() string {
	tokFile := "credential/google/credentials.json"
	path := getAbsPath(tokFile)

	data, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	scopes := []string{
		oauth2.UserinfoEmailScope,
		oauth2.UserinfoProfileScope,
	}

	config, err := google.ConfigFromJSON(data, scopes...)

	if err != nil {
		panic(err)
	}

	return config.AuthCodeURL("state")
}
