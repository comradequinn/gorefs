### warning
unless you can't, use `go modules` instead

# gorefs
gorefs provides simple, yet powerful, dependency management for Go, capable of replacing Go Get (though in no way prevents its use should you wish to do so). gorefs allows more granular management of dependencies than `go get` while being more in the spirit of `go get` and simpler than some may consider `godeps` (archived) or `dep` to be. It doesn't use vendor directories and has a purposely  simple set-up and configuration mechanism. Fundamentally, it helps map import paths to a repo and a revision.

Dependencies can be specified centrally within a project using a gorefs.json file defined in its root. These can be tied to a specified version/release (*tag or branch*). godeps can auto-generate the references portion of a gorefs.json file by parsing your source code.

Once a valid gorefs.json file exists in an application, only this file need be stored in the application's repo with the application itself. The source code of any dependencies is not required as cloning the application's repo and running `gorefs` will rebuild the dependency graph, similar to how `go get` would do; unlike `go get` however, it uses the specific versions defined in the gorefs.json file and thereby also provides an explicit application dependency listing (*in the form of the gorefs.json file*) that can be used to describe and compare revisions and releases of the application.

# Walkthrough
Run "gorefs -i" in your application root to generate a sample gorefs.json file. This is shown below:

```
{
	"aliases": {
	  "repo.git/": "git@gitrepo.local:"
	},
	"references": {
	  "repo.git/example/toolA": "master",
	  "repo.git/example/toolB": "tag/v1.1"
	}
}
```
Next customise this file to specify your aliases in the Aliases section (*an alias is a mapping of a git address root to the first segment of an import statement*). These can be used in the standard `go get` manner; linking similarly named imports and urls; or be complelely unrelated. For example, they may re-map github.com import statements to your corporate git repo as opposed to github.com itself (*should you wish to host known, secure versions of code internally*). 

The References section maps an import statement value to a version (*the specific tag/branch in git*). 

An updated gorefs.json file is shown below, with examples of these use cases:

```
{
	"aliases": {
	  "github.com": "git@github.com:",
	  "mygit.com": "git@git.local.com:",
	},
	"references": {
	  "github.com/example/toolA": "master",
	  "mygit.com/cust-logger/logs": "v1.0.0"
	}
}
```
Now we have a valid gorefs.json file, we can run `gorefs` in the application root. This will cause the specified versions of the dependencies to be cloned from a repo address calculated from the aliases and placed in a directory named after the imports statement value. `Go build` will now recognise your dependencies' presence.

If you wish to add another reference, you can run `gorefs -a [reference] -v [version]`. The reference will be added to your gorefs.json file, and the source code fetched from the appropriate repository and placed in the appropriate directory. Running `gorefs -d [reference]` will delete an existing reference. Alternatively you can also edit the gorefs.json file by hand and run `gorefs` again.

If you wish to always take the latest version of a dependency, set the version to 'master', or another active branch, as opposed to a tag

# Migrating an Existing Project
If you have an existing project with 100s, or even 1000s, of existing import statements, gorefs can generate your references for you. Firstly define any aliases in your gorefs.json file, as shown below in the scenario of adding gorefs based versioning to all existing imports from github.com

```
{
	"aliases": {
	  "github.com": "git@github.com:"
	},
	"references": {
	  //It's fine to have pre-existing references here. 
	  //It's not fine to have comments though, this is json #smacksHand
	}
}
```

Now run `gorefs -g`. This will cause gorefs to recursivley scan all .go files in the application root, extracting any import paths that begin with any of the defined aliases as it does so. It will then add these to the references section of your gorefs.json file, if they do not already exist (*it will not overwite existing references, so it is safe to run the command repeatedly if required without losing existing version information*). 

When the above completes, review the command output and/or your amended gorefs.json file and, if you are happy, run `gorefs` to fetch the repos of all the newly added references

# Options
A number of options exist to customise the behaviour of gorefs. 

The -f (f)orce switch causes existing copies of dependency source code be to deleted and fetched again from the appropriate repo. Normally only missing dependencies are fetched. This is useful when you may have corrupted, or modified, your local copy, or to ensure a clean build during CI and other automated builds scenarios.

The -d (d)ev switch causes any dependency repos fetched to have their git repo bindings preserved. Normally .git* files that link the dependency source code to the repo are removed after a fetch. This is useful in development scenarios when an application dependency is itself being modified as part of the project.

The -r (r)ef switch specifies a specific reference to fetch. Normally all references are fetched. This can be useful when you want to force (-f) or dev (-d) fetch only a single reference but you are happy to leave the rest as they are; or you have a large number of references and require the optimisation of naming specific references to speed up the process.

The -p (p)ath switch specifies an application root path. This is useful if you are using gorefs in a script, or from a different location than the application you wish to target with gorefs.

The -h (h)help switch causes a list of supported options to be displayed.

# Install
An `Makefile` has been provided to build and un/install gorefs.

# Removing gorefs from an Application
Just delete the gorefs.json file from the application root. Gorefs does not modify any other source files in your application
