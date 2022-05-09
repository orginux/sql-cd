package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Clone clones a git project into a directory
func Clone(dir, gitUrl, branch string) error {

	// Clean up
	os.RemoveAll(dir)

	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           gitUrl,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})
	if err != nil {
		return err
	}
	return nil
}
