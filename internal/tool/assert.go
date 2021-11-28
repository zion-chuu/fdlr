package tool

import (
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
)

func GetIPv4(ips []net.IP) []string {
	var ret = []string{}
	for _, ip := range ips {
		if ip.To4() != nil {
			ret = append(ret, ip.String())
		}
	}

	return ret
}

func Mkdir(folder string) error {
	if _, err := os.Stat(folder); err != nil {
		if err = os.MkdirAll(folder, 0700); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func IsFolderExisted(folder string) bool {
	_, err := os.Stat(folder)

	return err == nil
}

// Check if the current environment is terminal
func DisappearProgressBar() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

func GetFolderFrom(url string) (string, error) {
	path := filepath.Join(os.Getenv("HOME"), SaveFolder)

	absolutePath, err := filepath.Abs(filepath.Join(os.Getenv("HOME"), SaveFolder, filepath.Base(url)))
	if err != nil {
		return "", errors.WithStack(err)
	}

	// To prevent path traversal attack
	relative, err := filepath.Rel(path, absolutePath)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if strings.Contains(relative, "..") {
		return "", errors.WithStack(errors.New("Your download file may have a path traversal attack"))
	}

	return absolutePath, nil
}

func GetFilenameFrom(url string) string {
	filename := filepath.Base(url)

	return filename
}

func IsVaildURL(s string) bool {
	_, err := url.Parse(s)

	return err == nil
}
