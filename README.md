# XSD2Go - Automatically generate golang xml parser based on XSD ![Build CI](https://github.com/GoComply/xsd2go/workflows/Build%20CI/badge.svg)

:warning: **You should run xsd2go, before ever importing `encoding/xml` to your project.** :warning:

You may want to start reading [blog introduction](http://isimluk.com/posts/2020/05/xsd2go-automatically-generate-golang-xml-parsers/) to this project.

## Motivation

Did you ever got to implement XML parser? Perhaps for atom, or scap? You may have got XSD
(XML Schema Definition) files to verify adherence to given xml application? Wouldn't it be
cool to automatically generate XML parser based on XSD definition? That's exactly what we
are up to here.

### Related projects:
 - ![Metaschema](https://github.com/gocomply/metaschema) - generate golang code based on NIST metaschema input
 - ![SCAP](https://github.com/gocomply/scap) - parsers of NIST SCAP family of standards

## Exemplary Usage

```
# Acquire latest some XSD file you want to convert - for instance XCCDF 1.2
git clone --depth 1 https://github.com/openscap/openscap
# Parse XSD schema and generate golang structs
./gocomply_xsd2go convert openscap/schemas/xccdf/1.2/xccdf_1.2.xsd github.com/complianceascode/librescap pkg/scap/models/xccdf/1.2
```

## Installation

```
go get -u -v github.com/gocomply/xsd2go/cli/gocomply_xsd2go
```
