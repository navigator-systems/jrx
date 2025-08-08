package giteng

import (
	"fmt"
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
