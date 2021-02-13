package cmds

import (
	"gorefs/refs"
	"sync"
)

//FetchAll fetches all references
func FetchAll(grjc *refs.Content, devmode bool, force bool) []error {
	errs := []error{}
	wg := sync.WaitGroup{}

	for ref := range grjc.References {
		wg.Add(1)
		r := ref
		go func() {
			defer wg.Done()
			if err := Fetch(grjc, devmode, force, r); err != nil {
				errs = append(errs, err)
			}
		}()
	}

	wg.Wait()

	return errs
}
