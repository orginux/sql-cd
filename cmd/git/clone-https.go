package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Clones a git project into a directory via HTTPS
func cloneHTTPS(gitDest, gitUrl, branch string) error {

	_, err := git.PlainClone(gitDest, false, &git.CloneOptions{
		URL:           gitUrl,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		// Progress:      os.Stdout,
	})
	if err != nil {
		return err
	}
	return nil
}
