# Contributing

## Coding Conventions

(tbd)

## Git Workflow

```
git clone git@github.com:resurfaceio/logger-go.git
cd logger-go
go mod download

```

Running unit tests:

From inside logger-go directory
```
go test
```

Committing changes:

```
git add -A
git commit -m "#123 Updated readme"       (123 is the GitHub issue number)
git pull --rebase                         (avoid merge bubbles)
git push origin master
```

## Release Process

Exercise using Golang test apps:
* test-mux
Push artifacts to ???:

```
(tbd)
```

Tag release version:

```
git tag v1.x.x
git push origin master --tags
```

Start the next version by incrementing the version number in BaseLogger.go.
