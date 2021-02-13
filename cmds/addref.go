package cmds

import (
	"fmt"
	"gorefs/refs"
)

//AddRef adds the specified reference to the configuration and then fetches it
func AddRef(grjc *refs.Content, devmode bool, force bool, ref string, ver string) []error {
	grjc.References[ref] = ver

	errs := []error{}

	refs.Write(grjc)

	if err := Fetch(grjc, devmode, force, ref); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		errs = append(errs, fmt.Errorf("attempting to clean up after failed add reference attempt"))

		if err := DelRef(grjc, ref); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
