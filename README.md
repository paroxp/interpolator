# Interpolator

It is a replacement for the popular tool `envsubst` which allows you to fill in
the environment variables in a given template in a bash-like style. This tool
however, introduces failover while the variable filling when one is missing.
Also it allows users to provide custom syntax for the variable format and
defaults to bash-like style.

## Installation

```sh
go get -u github.com/paroxp/interpolator
```

## Usage

Out of the box, the tool supports the `${ENV_VAR}` syntax:

```sh
export GREET="World"; echo "Hello \${GREET}" | interpolator -f
# Hello World
```

You can amend the pattern of the variables, by adding the `-m` | `--pattern`
flag:

```sh
export GREET="World"; echo "Hello ((GREET))" | interpolator -m "\\(\\(([A-Z0-9]+)\\)\\)"
# Hello World
```

You can also dry run the thing, to establish which variables are used in the
template.

```sh
export GREET="World"; echo "Hello \${GREET}" | interpolator --dry
# GREET:		'World'
```

## Testing

```sh
go test ./...
```
