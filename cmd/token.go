package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os/user"

	"github.com/AlecAivazis/survey"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	keyring "github.com/zalando/go-keyring"
)

func findGithubToken() (githubToken string, err error) {
	if githubToken = viper.GetString("githubToken"); githubToken == "" {
		var currentUser *user.User
		if currentUser, err = user.Current(); err != nil {
			return
		}
		githubToken, err = loadFromKeychain(currentUser)
		if err != nil {
			logrus.WithError(err).Info("token not found in keychain")
			err = nil
		} else {
			return
		}
		prompt := survey.Input{
			Message: "Enter your github personal access token (requires the 'notifications' scope):",
		}
		survey.AskOne(&prompt, &githubToken, validateToken)
		if githubToken == "" {
			return githubToken, errors.New("github key not set")
		}
		var save bool
		survey.AskOne(&survey.Confirm{
			Message: "Save in keychain?",
		}, &save, nil)
		if save && currentUser != nil {
			if err = saveInKeychain(currentUser, githubToken); err != nil {
				logrus.WithError(err).Error("unable to save token in keychain")
			}
			logrus.Info("saved personal access token in keychain")
			return
		} else if save {
			return "", errors.New("unable to save, no current user")
		}
		return
	}
	return
}

func loadFromKeychain(currentUser *user.User) (githubToken string, err error) {
	key := fmt.Sprintf("octonotify-%s-pat", currentUser.Username)
	if githubToken, err = keyring.Get("octonotify", key); err == nil {
		logrus.Info("found personal access token in keychain")
		return
	}
	return
}

func saveInKeychain(currentUser *user.User, githubToken string) (err error) {
	key := fmt.Sprintf("octonotify-%s-pat", currentUser.Username)
	if err = keyring.Set("octonotify", key, githubToken); err != nil {
		return
	}
	return
}

func validateToken(a interface{}) error {
	var ans string
	var ok bool
	if ans, ok = a.(string); !ok {
		return errors.New("not a string")
	}
	if _, err := hex.DecodeString(ans); len(ans) == 40 && err == nil {
		return nil
	}
	return errors.New("not a valid personal access token")
}
