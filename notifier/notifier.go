package notifier

import "github.com/gen2brain/beeep"

type notifyFunc func(string, string, string) error

// Notifier ...
type Notifier interface {
	Notify(string, string, string) error
}

type notifier struct {
	notifyFunc notifyFunc
}

// New returns a new Notifier, using the default implementation
func New(sound bool) Notifier {
	n := notifier{}
	if sound {
		n.notifyFunc = beeep.Alert
	} else {
		n.notifyFunc = beeep.Notify
	}
	return &n
}

func (n *notifier) Notify(title, message, iconPath string) error {
	return n.notifyFunc(title, message, iconPath)
}
