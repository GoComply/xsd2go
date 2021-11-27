#!/bin/bash

set -e -o pipefail
set -x

xsd2go=$(pwd)/gocomply_xsd2go

workdir=./scap/
mkdir -p $workdir
cat <<__END__ > $workdir/go.mod
module github.com/gocomply/scap
go 1.17
__END__

pushd $workdir

    # Acquire XSDs
    [ -d .scap_schemas ] || git clone --depth 1 https://github.com/openscap/openscap .scap_schemas

    # Clean-up the workspace
    [ -d pkg/scap/models ] && find pkg/scap/models -name models.go | xargs rm --

    # Generage go code based on XSDs
    $xsd2go convert .scap_schemas/schemas/cpe/2.3/cpe-dictionary_2.3.xsd github.com/gocomply/scap pkg/scap/models

    # Ensure the code can be compiled
    go vet ./...

popd

