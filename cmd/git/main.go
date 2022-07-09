package git

import (
	"fmt"
	"os"
	"strings"

	"github.com/orginux/sql-cd/cmd/logging"
)

func Clone(gitDest, gitURL, gitBranch, gitPrivateKeyFile string, verbose bool) error {

	err := os.RemoveAll(gitDest)
	if err != nil {
		return fmt.Errorf("Cannot remove directory %s: %v", gitURL, err)
	}

	if verbose {
		logging.Info.Printf("git clone %s ", gitURL)
	}

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
