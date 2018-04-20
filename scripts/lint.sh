#!/bin/bash

set -eou pipefail

# relative path calculation shim, I hope python is installed!
function relative_path {
    if hash realpath 2> /dev/null; then
        realpath --relative-to="$2" "$1"
    else
        if hash python 2> /dev/null; then  
            python -c "import os.path; print os.path.relpath('$1', '$2')"
        else
            echo "unable to calculate a relative path"
            exit 1
        fi
    fi    
}

# This is all in a single file for simplicity of future updates

gometalinter_file="$(mktemp)"
errcheck_excludes_file="$(mktemp)"

# TODO: download the metalinter config and base errcheck file from a shared repo?
cat <<JSON > "$gometalinter_file"
{
  "Deadline": "160s",
  "Linters": {
    "errcheck": "errcheck -exclude $errcheck_excludes_file -abspath:PATH:LINE:COL:MESSAGE"
  },
  "Exclude": [
    ".*declaration of \\"err\\" shadows declaration at.*\\\\(vetshadow\\\\)",
    ".*error return value not checked \\\\(defer .+\\\\(errcheck\\\\)"
  ],
  "Enable": [
    "deadcode",
    "errcheck",
    "goconst",
    "gofmt",
    "goimports",
    "golint",
    "gotype",
    "gotypex",
    "ineffassign",
    "interfacer",
    "megacheck",
    "misspell",
    "nakedret",
    "structcheck",
    "unconvert",
    "varcheck",
    "vet",
    "vetshadow"
  ],
  "Vendor": true,
  "WarnUnmatchedDirective": true
}
JSON

# This is calculated as vendor packages include their full path in their type names
pkg="$(relative_path "$(pwd)" "$GOPATH/src")"
cat <<ERRCHECK > "$errcheck_excludes_file"
(*$pkg/vendor/github.com/hashicorp/terraform/helper/schema.ResourceData).Set
ERRCHECK

# cleanup the file but still return the proper exit code
lintexit="0"
echo "Running gometalinter"
gometalinter --config "$gometalinter_file" ./... || lintexit="$?" && :

# cleanup the tmp files
rm -f "$gometalinter_file" "$errcheck_excludes_file"

# use the lint exit code
exit "$lintexit"
