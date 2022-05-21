package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/orginux/sql-cd/cmd/logging"
)

// Clone clones a git project into a directory
func Clone(dir, gitUrl, branch string) error {

	// Clean up
	os.RemoveAll(dir)

	privateKeyFile := "/tmp/key"
	password := ""

	_, err := os.Stat(privateKeyFile)
	if err != nil {
		logging.Warning.Printf("read file %s failed %s\n", privateKeyFile, err)
	}

	// Clone the given repository to the given directory
	logging.Info.Printf("git clone %s ", gitUrl)
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, password)
	if err != nil {
		logging.Warning.Printf("generate publickeys failed: %s\n", err.Error())
	}

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:           gitUrl,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		Progress:      os.Stdout,
		Auth:          publicKeys,
	})
	if err != nil {
		return err
	}
	return nil
}
