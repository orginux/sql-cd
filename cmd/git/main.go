package git

import (
	"fmt"
	"os"
	"strings"
)

func Clone(gitDest, gitURL, gitBranch, gitPrivateKeyFile string) error {

	// Clean up before clone
	os.RemoveAll(gitDest)

	if strings.HasPrefix(gitURL, "https://") {
		err := cloneHTTPS(gitDest, gitURL, gitBranch)
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(gitURL, "git@") {
		err := cloneSSH(gitDest, gitURL, gitBranch, gitPrivateKeyFile)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unsupported protocol: %s", gitURL)
	}

	return nil
}
