package scm

import (
	"fmt"
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func GitInit(repoPath string) {
	_, err := git.PlainInit(repoPath, false)
	if err != nil {
		fmt.Printf("Failed to initialize git repository at %s: %v\n", repoPath, err)
		return
	}

}

func GitBranchMain(repoPath string) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("failed to open repository: %v", err)
		return
	}

	err = repo.Storer.SetReference(
		plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName("main")),
	)
	if err != nil {
		log.Fatalf("failed to set HEAD to main: %v", err)
	}

	fmt.Println("Checked out main branch successfully.")
}

func GitAddCommmit(repoPath string) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("failed to open repository: %v", err)
		return
	}

	w, err := repo.Worktree()
	if err != nil {
		log.Fatalf("failed to get worktree: %v", err)
		return
	}

	_, err = w.Add(".")
	if err != nil {
		log.Fatalf("failed to add files: %v", err)
		return
	}
	message := "Initial commit by JRX cli"
	commit, err := w.Commit(message, &git.CommitOptions{
		All: true,
	})
	if err != nil {
		log.Fatalf("failed to commit changes: %v", err)
		return
	}

	fmt.Println("Committed changes successfully:", commit)
}

func GitAddRemote(repoPath, remoteName, remoteURL string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Check if remote already exists
	_, err = repo.Remote(remoteName)
	if err == nil {
		// Remote already exists, update it
		err = repo.DeleteRemote(remoteName)
		if err != nil {
			return fmt.Errorf("failed to delete existing remote: %w", err)
		}
	}

	// Create new remote
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{remoteURL},
	})
	if err != nil {
		return fmt.Errorf("failed to add remote: %w", err)
	}

	fmt.Printf("Added remote '%s' with URL: %s\n", remoteName, remoteURL)
	return nil
}

// GitPush pushes the local repository to the remote
func GitPush(repoPath string, remoteName string, branch string, sshKeyPath string, sshKeyPassphrase string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Setup SSH authentication
	publicKeys, err := ssh.NewPublicKeysFromFile("git", sshKeyPath, sshKeyPassphrase)
	if err != nil {
		return fmt.Errorf("failed to create SSH keys: %w", err)
	}

	// Push to remote
	refSpec := config.RefSpec(fmt.Sprintf("+refs/heads/%s:refs/heads/%s", branch, branch))
	err = repo.Push(&git.PushOptions{
		RemoteName: remoteName,
		RefSpecs:   []config.RefSpec{refSpec},
		Auth:       publicKeys,
	})
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	fmt.Printf("Pushed '%s' branch to '%s' successfully\n", branch, remoteName)
	return nil
}
