package configuration

import (
	"errors"
	"github.com/eddieowens/axon"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const AuthMethodKey = "AuthMethod"

func authMethodFactory(inj axon.Injector, _ axon.Args) axon.Instance {
	config := inj.Get(ConfigKey).GetValue().(Config)
	username := config.Git.Username
	password := config.Git.Password
	if config.Git.SSHKey != "" {
		return axon.StructPtr(sshKeyFromFile(config.Git.SSHKey))
	} else if username != "" && password != "" {
		return axon.StructPtr(usernamePassword(username, password))
	} else {
		panic(errors.New("either an ssh key or a username/password are required for git access"))
	}
}

func usernamePassword(username, password string) transport.AuthMethod {
	return &http.BasicAuth{
		Username: username,
		Password: password,
	}
}

func sshKeyFromFile(fp string) transport.AuthMethod {
	if strings.HasPrefix(fp, "~/") {
		home, _ := os.UserHomeDir()
		fp = path.Join(home, fp[2:])
	}
	key, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}
	return sshKey(key)
}

func sshKey(key []byte) transport.AuthMethod {
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}

	pub := &gitssh.PublicKeys{
		User:   "git",
		Signer: signer,
	}

	return pub
}
