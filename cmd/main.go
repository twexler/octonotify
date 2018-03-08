package main

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"github.com/twexler/octonotify/watcher"
	"golang.org/x/oauth2"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ""},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	w := watcher.New(client.Activity)

	w.Run(ctx, time.Minute)
}
