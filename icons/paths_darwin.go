//+build darwin
package icons

import (
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// paths used only on darwin

// GetIconPathOnDisk ...
func GetIconPathOnDisk() string {
	relIcnsPath := path.Join(os.Args[0], "..", "..", "Resources", "octonotify.icns")
	icnsPath, err := filepath.Abs(relIcnsPath)
	if err != nil {
		logrus.WithError(err).Fatal("cannot build icon path")
	}
	return icnsPath
}
