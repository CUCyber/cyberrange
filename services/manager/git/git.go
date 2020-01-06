package git

import (
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

var user string
var token string
var PlaybooksPath string

func ParseGitOptions() (map[string]string, error) {
	gitConfig := viper.New()
	gitConfig.SetConfigName("config")
	gitConfig.SetConfigType("yaml")
	gitConfig.AddConfigPath("./git/")

	if err := gitConfig.ReadInConfig(); err != nil {
		return nil, err
	}

	return gitConfig.GetStringMapString("git"), nil
}

func PullPlaybooks() error {
	r, err := git.PlainOpen(PlaybooksPath)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: user,
			Password: token,
		},
		Force: true,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}

func InitializeGit() {
	gitoptions, err := ParseGitOptions()
	if err != nil {
		panic(err.Error())
	}

	user = gitoptions["user"]
	token = gitoptions["token"]
	PlaybooksPath = gitoptions["path"]
}
