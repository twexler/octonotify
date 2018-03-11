package watcher

import "github.com/google/go-github/github"

// Notification ...
type Notification interface {
	GetRespository() *string
}

type githubNotification struct {
	github.Notification
}

func newNotification(gn *github.Notification) *githubNotification {
	return &githubNotification{*gn}
}

func (gn *githubNotification) GetRespository() *string {
	return gn.Repository.Name
}
