//+build !windows !darwin
package icons

import "path"

const (
	IconPath  = "octonotify-small.png"
	sharePath = "/usr/share"
)

func GetIconPathOnDisk() string {
	return path.Join(sharePath, IconPath)
}
