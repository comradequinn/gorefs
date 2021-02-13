package cmds

import (
	"fmt"
	"os"
	"path"
)

//Returns the directory associated with a reference. If delete is set, it will delete the directory if it exists
func dir(reference string, delete bool) (string, error) {
	gp := os.Getenv("GOPATH")

	if gp == "" {
		return "", fmt.Errorf("error fetching reference %v. GOPATH unset", reference)
	}

	dest := path.Join(gp, "src", reference)

	if delete {
		os.RemoveAll(dest)
	}

	return dest, nil
}
