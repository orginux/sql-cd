package git

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
		return fmt.Errorf("%v, please check value for --private-key-file", err)
	}

	// Clone the given repository to the given directory
	logging.Info.Printf("git clone %s ", gitUrl)
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, password)
	if err != nil {
		return fmt.Errorf("generate publickeys failed: %v\n", err.Error())
	}

	gitHostname, err := getHostname(gitUrl)
	if err != nil {
		return fmt.Errorf("Error parse url: %v", err)
	}

	sshKeyscan(gitHostname, "/etc/ssh/ssh_known_hosts")

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

// ssh-keyscan -H github.com > /ssh_known_hosts
func sshKeyscan(host, knownHostsPath string) error {
	// /etc/ssh/ssh_known_hosts
	knownHostsDir := filepath.Dir(knownHostsPath)
	if _, err := os.Stat(knownHostsDir); os.IsNotExist(err) {
		err := os.MkdirAll(knownHostsDir, 0440)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("ssh-keyscan", "-H", host)

	stdout, err := cmd.Output()
	if err != nil {
		logging.Error.Printf("sshKeyscan: %v", err.Error())
		return err
	}

	knownHostsFile, err := os.Create(knownHostsPath)
	if err != nil {
		logging.Error.Printf("create file: %v", err)
		return err
	}
	defer knownHostsFile.Close()

	_, err = knownHostsFile.Write(stdout)
	if err != nil {
		logging.Error.Printf("WriteString: %v", err)
		return err
	}

	return nil
}

// git@github.com:orginux/clickhouse-test-env.git
func getHostname(sourceUrl string) (string, error) {
	withoutPointsURL := strings.ReplaceAll(sourceUrl, ":", "/")
	withoutGitURL := strings.ReplaceAll(withoutPointsURL, "git@", "https://")
	u, err := url.Parse(withoutGitURL)
	if err != nil {
		logging.Debug.Printf("Error URL: %v", err)
		return "", err
	}
	return u.Host, nil
}
