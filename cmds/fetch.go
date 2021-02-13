package cmds

import (
	"fmt"
	"gorefs/refs"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

//Fetch retrieves the passed reference from the source repo and places it in the application root in the correct location for the reference
func Fetch(grjc *refs.Content, devmode bool, force bool, ref string) error {
	dest, err := dir(ref, force)

	if err != nil {
		return fmt.Errorf("error preparing reference path for %v. %v", ref, err)
	}

	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		log.Printf("skipping %v as it has already been fetched. Re-run with (f)orce switch to overwrite the existing version", ref)
	} else {

		repo, err := repo(ref, grjc.Aliases)

		if err != nil {
			return fmt.Errorf("error parsing reference source for %v. %v", ref, err)
		}

		version, err := version(ref, grjc.References)

		if err != nil {
			return fmt.Errorf("error parsing reference version for %v. %v", ref, err)
		}

		_, stderr := exec.Command("git", "clone", "-b", version, "--single-branch", repo, dest).Output()

		if stderr != nil {
			return fmt.Errorf("error fetching reference %v. %v", ref, stderr)
		}

		if devmode {
			return rmSCC(dest, ref)
		}

		log.Printf("fetched %v at version '%v' from %v into %v", ref, version, repo, dest)
	}

	return nil
}

func rmSCC(dir string, ref string) error {
	if err := os.RemoveAll(path.Join(dir, ".git")); err != nil {
		return fmt.Errorf("error breaking links between %v and its backing repo. %v", ref, err)
	}

	if err := os.RemoveAll(path.Join(dir, ".gitignore")); err != nil {
		return fmt.Errorf("error breaking links between %v and its backing repo. %v", ref, err)
	}

	return nil
}

func version(reference string, refs map[string]string) (string, error) {
	for ref, ver := range refs {
		if strings.ToLower(ref) == strings.ToLower(reference) {
			return ver, nil
		}
	}

	return "", fmt.Errorf("unable to find version for reference '%v'", reference)
}

func repo(reference string, aliases map[string]string) (string, error) {
	for alias, baseurl := range aliases {
		if strings.HasPrefix(reference, alias) {
			return strings.Replace(reference, alias, baseurl, 1), nil
		}
	}

	return "", fmt.Errorf("unable to find source for reference '%v'", reference)
}
