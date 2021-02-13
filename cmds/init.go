package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var template = []byte(`{
	"aliases": {
	  "repo.git/": "git@gitrepo.local:"
	},
	"references": {
	  "repo.git/example/toolA": "main",
	  "repo.git/example/toolB": "tags/v1.0.0"
	}
}`)

//Init writes a sample configuration file
func Init(root string) error {
	filepath := path.Join(root, "gorefs.json")

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return ioutil.WriteFile(filepath, template, 0666)
	}

	return fmt.Errorf("there is an existing %v. Delete or rename it before running init", filepath)
}
