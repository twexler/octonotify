package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/twexler/octonotify/notifier"

	"github.com/AlecAivazis/survey"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/twexler/octonotify/watcher"
	"golang.org/x/oauth2"
)

var (
	configFile string
	interval   time.Duration
	logLevel   string
	rootCmd    = cobra.Command{
		Version: version,
		Run:     cobraMain,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to a config file")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "warning", "set the log level")
	rootCmd.PersistentFlags().DurationVarP(&interval, "interval", "i", time.Minute, "interval for polling github")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("halp")
	}
}

func cobraMain(cmd *cobra.Command, _ []string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).Warn("invalid level")
	}
	logrus.SetLevel(level)
	var githubKey string
	if githubKey = viper.GetString("githubKey"); githubKey == "" {
		// TODO: use keychain/DPAPI+reg/etc
		prompt := survey.Input{
			Message: "Enter your github personal access token (requires the 'notifications' scope:",
		}
		survey.AskOne(&prompt, &githubKey, func(a interface{}) error {
			var ans string
			var ok bool
			if ans, ok = a.(string); !ok {
				return errors.New("not a string")
			}
			if _, err := hex.DecodeString(ans); len(ans) == 40 && err == nil {
				return nil
			}
			return errors.New("not a valid personal access token")
		})
		if githubKey == "" {
			logrus.Fatal("github key not set")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubKey},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	notificationChan := make(chan watcher.Notification)
	w := watcher.New(client.Activity, notificationChan)

	localNotifier := notifier.New(false)

	localNotifier.Notify("Octonotify", "Starting", "")
	go w.Run(ctx, interval)
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			localNotifier.Notify("Octonotify", "Exiting", "")
			cancel()
			os.Exit(0)
		}
	}()
	go func() {
		for n := range notificationChan {
			message := fmt.Sprintf("New notification on %s", *n.GetRespository())
			localNotifier.Notify("Octonotify", message, "")
		}
	}()
	select {}
}

func initConfig() {
	if configFile != "" {
		logrus.WithField("configFile", configFile).Info("loading configfile")
		viper.SetConfigFile(configFile)
	} else {
		// Find home directory.
		home := os.Getenv("HOME")
		if home == "" {
			logrus.Warn("$HOME not set")
			return
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".octonotify")
	}

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Error("Can't read config")
	}
}
