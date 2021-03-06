[![Build Status](https://travis-ci.org/softleader/helm-filter.svg?branch=master)](https://travis-ci.org/softleader/helm-filter)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/softleader/helm-filter/blob/master/LICENSE)
[![Build Status](https://github-basic-badges.herokuapp.com/release/softleader/helm-filter.svg)](https://github.com/softleader/helm-filter/releases)

# Helm Filter Plugin

Helm plugin to filter out template files

[![asciicast](https://asciinema.org/a/195346.png)](https://asciinema.org/a/195346)

## Install

Fetch the latest binary release of helm-filter and install it:
 
```sh
$ helm plugin install https://github.com/softleader/helm-filter
```

Or have fun with Docker!

```sh
$ docker pull softleader/helm
```

## Usage
 
```sh
$ helm filter [flags] CHART
```

### Flags

```sh
Flags:
  -h, --help                help for helm
  -o, --output-dir string   copy all files to output-dir and filter there instead filter in chart path
      --overwrite-values    overwrite values after filtered out
  -f, --values string       specify values in a YAML file to filter (default "values.yaml")
```

## Example

The structure is like:

```js
.
├── mychart
│   ├── Chart.yaml
│   ├── charts
│   ├── templates
│   └── values.yaml
└── myenv
    ├── client-a.yaml
    ├── sit.yaml
    └── uat.yaml
```

*client-a.yaml* is like:

```yaml
ingress:
  __filter_out: ingress*
```

The script for package different environments chart archive:

```sh
# Merge sit and client-a to values.yaml
$ helm values mychart -f myenv/sit.yaml -f myenv/client-a.yaml -o mychart

# Filter out files in templates-dir and copy to tmp folder
$ helm filter mychart -o tmp

# Pack chart archive
$ helm package tmp/mychart

# Remove tmp folder 
$ rm -rf tmp 

# Restore values.yaml from backup file
$ mv mychart/values.yaml.bak mychart/values.yaml
```
