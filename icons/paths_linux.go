//+build !windows !darwin

package icons

import "path"

const (
	sharePath = "/usr/share"
)

func GetIconPathOnDisk() string {
	return path.Join(sharePath, IconPath)
}
