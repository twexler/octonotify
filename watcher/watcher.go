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

// Watcher watches the ListNotifications endpoint for new notifications
type Watcher struct {
	lister           NotificationLister
	notificationChan chan<- Notification
}

// New returns a new Watcher instance
func New(lister NotificationLister, notificationChan chan<- Notification) *Watcher {
	return &Watcher{
		lister:           lister,
		notificationChan: notificationChan,
	}
}

// Run starts the Watcher
func (w *Watcher) Run(ctx context.Context, interval time.Duration) {
	logrus.Debug("Starting watcher")
	timer := time.NewTicker(interval)
	lastRun := time.Now()
	w.checkNotifications(ctx, lastRun)
	for t := range timer.C {
		w.checkNotifications(ctx, lastRun)
		lastRun = t
	}
}

func (w *Watcher) checkNotifications(ctx context.Context, lastRun time.Time) {
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
	for _, n := range notifs {
		w.notificationChan <- newNotification(n)
	}
}
