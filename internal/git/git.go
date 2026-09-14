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

func New() gitService {
	return gitService{
		repoPath: "/var/lib/gitria/repo",
	}
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
