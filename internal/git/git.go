package git

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type gitService struct {
	repoPath string
}

func New(repoPath string) (*gitService, error) {
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return nil, err
	}
	return &gitService{repoPath}, nil
}

func (s gitService) UploadPack(repo string, channel io.ReadWriter) error {
	repo = strings.Trim(repo, "'")
	repo = strings.TrimPrefix(repo, "/")

	gitPath := filepath.Join(s.repoPath, repo)
	cmd := exec.Command("git-upload-pack", gitPath)

	cmd.Stdin = channel
	cmd.Stdout = channel
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (s gitService) ReceivePack(repo string, channel io.ReadWriter) error {
	repo = strings.Trim(repo, "'")
	repo = strings.TrimPrefix(repo, "/")

	gitPath := filepath.Join(s.repoPath, repo)
	cmd := exec.Command("git-receive-pack", gitPath)

	cmd.Stdin = channel
	cmd.Stdout = channel
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (s gitService) UploadArchive(repo string, channel io.ReadWriter) error {
	repo = strings.Trim(repo, "'")
	repo = strings.TrimPrefix(repo, "/")

	gitPath := filepath.Join(s.repoPath, repo)
	cmd := exec.Command("git-upload-archive", gitPath)

	cmd.Stdin = channel
	cmd.Stdout = channel
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
