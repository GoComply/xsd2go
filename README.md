# XSD2Go - Automatically generate golang xml parser based on XSD

## Introduction

Did you ever got to implement XML parser? Perhaps for atom, or scap? You may have got XSD
(XML Schema Definition) files to verify adherence to given xml application? Wouldn't it be
cool to automatically generate XML parser based on XSD definition? That's exactly what we
are up to here.

### Related projects:
 - ![Metaschema](https://github.com/gocomply/metaschema) - generate golang code based on NIST metaschema input

## Usage

```
# Acquire latest some XSD file you want to convert
wget https://raw.githubusercontent.com/OpenSCAP/openscap/maint-1.3/schemas/xccdf/1.2/xccdf_1.2.xsd
# Parse XSD schema and generate golang structs
./gocomply_xsd2go convert xccdf_1.2.xsd github.com/complianceascode/librescap pkg/scap/models/xccdf/1.2
```

## Installation

```
go get -u -v github.com/gocomply/xsd2go/cli/gocomply_xsd2go
```
