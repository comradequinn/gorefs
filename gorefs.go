package main

import (
	"flag"
	"gorefs/cmds"
	"gorefs/refs"
	"log"
)

func main() {

	log.Println("gorefs @ v1.0.0")

	flag.Bool("_", true, "DEFAULT: fetch all references in target application's gorefs.json")
	root := flag.String("p", "./", "the (p)ath to the root of the target application")
	init := flag.Bool("i", false, "(i)nit creates a sample gorefs.json file in the specified (p)ath, if one does not already exist. Ignores any other switches specified aside from (p)ath")
	gen := flag.Bool("g", false, "(g)enerate populates the References section of the target application's gorefs.json with any imports statements found in the target application's source code that match one of the Aliases defined in it's gorefs.json. Safe to run repeatedly as it does not overwrite existing references. Ignores any other switches specified aside from (p)ath")
	force := flag.Bool("f", false, "(f)orce a re-fetch. Previously fetched references will be overwritten and any local changes will be lost")
	dev := flag.Bool("d", false, "(d)ev mode maintains references links with their backing repo. Useful when a dependency is being developed alongside the application")
	ref := flag.String("r", "", "fetch a single (r)eference to fetch, overrides (a)ll if specified")
	del := flag.String("x", "", "deletes an existing reference, removing the local copy if one exists")
	add := flag.String("a", "", "(a)dds and fetches a new reference. Used in conjuction with (v)er if a specific version is to be added")
	ver := flag.String("v", "master", "specifies the version of a reference to (a)dd. Used in conjuction with (a)dd")

	flag.Parse()

	switch {
	case *init:
		if err := cmds.Init(*root); err != nil {
			log.Fatal(err)
		}
		break
	case *gen:
		if err := cmds.GenRefs(refs.Read(*root), *root); err != nil {
			log.Fatal(err)
		}
		break
	case *add != "":
		if errs := cmds.AddRef(refs.Read(*root), *dev, *force, *add, *ver); len(errs) > 0 {
			log.Fatal(errs)
		}
		break
	case *del != "":
		if err := cmds.DelRef(refs.Read(*root), *del); err != nil {
			log.Fatal(err)
		}
		break
	case *ref != "":
		if err := cmds.Fetch(refs.Read(*root), *dev, *force, *ref); err != nil {
			log.Fatal(err)
		}
		break
	default:
		if errs := cmds.FetchAll(refs.Read(*root), *dev, *force); len(errs) > 0 {
			log.Fatal(errs)
		}
		break
	}
}
