package refs

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"strings"
)

//Content defines the options specified for sourcing packages
type Content struct {
	Aliases    map[string]string `json:"aliases"`
	References map[string]string `json:"references"`
	source     string
}

//Write the passed content to the gorefs.json file. Failure to write the file is fatal
func Write(c *Content) {

	if err := os.RemoveAll(c.source); err != nil {
		log.Fatalf("unable to overwrite remove existing configuration file. %v", err)
	}

	file, err := os.Create(c.source)

	if err != nil {
		log.Fatalf("unable to create configuration file at '%v'. Error: %v", c.source, err)
	}

	defer file.Close()

	ec := json.NewEncoder(file)
	ec.SetIndent("", "\t")

	err = ec.Encode(c)

	if err != nil {
		log.Fatalf("unable to serialise configuration. %v", err)
	}
}

//Read parses the content from the gorefs.json file. Failure to read the (valid) file is fatal
func Read(approot string) *Content {
	filepath := path.Join(approot, "gorefs.json")

	file, err := os.Open(filepath)

	if err != nil {
		log.Fatalf("unable to open configuration file at '%v'. Error: %v", filepath, err)
	}

	defer file.Close()

	dc := json.NewDecoder(file)
	c := &Content{source: filepath}

	err = dc.Decode(c)

	if err != nil {
		log.Fatalf("unable to parse configuration file at '%v'. Error: %v", filepath, err)
	}

	for alias, baseurl := range c.Aliases {
		if !strings.HasSuffix(alias, "/") || (!strings.HasSuffix(baseurl, ":") && !strings.HasSuffix(baseurl, "/")) {
			log.Fatalf("invalid alias, include the trailing dividers and ensure the base url is ssh. Example: \"members.git/\": \"git@mem-git-prod.lb.local:\"")
		}
	}

	return c
}
