package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := fs.WalkDir(os.DirFS(os.Getenv("HOME")+"/.vim"), ".", getStartDir); err != nil {
		log.Fatal(err)
	}
}

func getStartDir(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if filepath.Base(path) == "start" {
		return gitPull(os.Getenv("HOME") + "/.vim/" + path)
	}
	return nil
}

func gitPull(path string) error {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		if err := os.Chdir(path + "/" + dir.Name()); err != nil {
			return (err)
		}
		repo, err := git.PlainOpen(path + "/" + dir.Name())
		if err != nil {
			return (err)
		}
		workTree, err := repo.Worktree()
		if err != nil {
			return (err)
		}
		pubKey, err := ssh.NewPublicKeysFromFile("git", os.Getenv("HOME")+"/.ssh/id_ed25519", "")
		if err != nil {
			return (err)
		}
		if err := workTree.Pull(&git.PullOptions{RemoteName: "origin", Auth: pubKey}); err != nil {
			if errors.Is(err, git.NoErrAlreadyUpToDate) {
				fmt.Println(path+"/"+dir.Name(), "\t", err)
				continue
			}
			return err
		}
		fmt.Println(path+"/"+dir.Name(), "\t", "updated")
	}
	return nil
}
