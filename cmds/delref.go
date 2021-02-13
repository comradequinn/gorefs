package cmds

import (
	"fmt"
	"gorefs/refs"
)

//DelRef deletes the specified reference from the configuration and then removes any local copy
func DelRef(grjc *refs.Content, ref string) error {
	delete(grjc.References, ref)

	refs.Write(grjc)

	if _, err := dir(ref, true); err != nil {
		return fmt.Errorf("unable to delete the directory associated with the reference")
	}

	return nil
}
