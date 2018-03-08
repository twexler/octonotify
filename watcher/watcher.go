package watcher

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// NotificationLister represents the github.ActivityService.ListNotifications method
// primarily for easy testability
type NotificationLister interface {
	ListNotifications(context.Context, *github.NotificationListOptions) ([]*github.Notification, *github.Response, error)
}

type watcher struct {
	lister NotificationLister
}

// New returns a new Watcher instance
func New(lister NotificationLister) *watcher {
	return &watcher{
		lister: lister,
	}
}

func (w *watcher) Run(ctx context.Context, interval time.Duration) {
	timer := time.NewTimer(interval)
	lastRun := time.Now()
	w.checkNotifications(ctx, lastRun)
	for t := range timer.C {
		w.checkNotifications(ctx, lastRun)
		lastRun = t
	}
}

func (w *watcher) checkNotifications(ctx context.Context, lastRun time.Time) {
	logrus.Info("checking for notifications")
	subCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	opts := &github.NotificationListOptions{
		Since: lastRun,
	}
	notifs, _, err := w.lister.ListNotifications(subCtx, opts)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("error calling github")
	}
	logrus.WithField("notifs", notifs).Debug("got notifications")
}
