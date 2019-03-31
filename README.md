# Forge Tools

Starting to (re-)write Autodesk forge tools in this go package for convinience.

# Forge cli
## Install

### Method 1

```bash
go get github.com/nicored/forge_tools/forge-cli
go install github.com/nicored/forge_tools/forge-cli
```

This will create a 'forge-cli' binary in your $GOBIN path that you should be able to use out of the box.

### Method 2
Alternatively you can clone this repository, ```cd``` into it and run ```go get && make```, which will create a 'bin' directory with the 'forge-cli' binary in there.

### Properties command
[x] Extracts and prints properties to stdout in Json format from a given directory
[-] Extracts from gzip and json (currently only supports json)
[-] Extracts from a given Autodesk urn instead of a directory
[-] Write more tests

#### Properties to Json
You must have a directory with the following files in it:
- objects_ids.json
- objects_offs.json
- objects_avs.json
- objects_vals.json
- objects_attrs.json

And then run:
```bash
# Print properties to stdout in json format.
# Argument is optional and defaults to current dir path '.'
forge-cli properties json /path/to/dir/with/objects/files/
```

### TODO
[-] More tests
[-] Dockerize it
[-] Implement server
[-] Heroku it
[-] Forge SDK
[-] Download derivatives
