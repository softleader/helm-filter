#!/usr/bin/env bash

BINARY="filter"
PROJECT_NAME="helm-"$BINARY
PROJECT_GH="softleader/$PROJECT_NAME"

initOS() {
  OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Msys support
    msys*) OS='windows';;
    # Minimalist GNU for Windows
    mingw*) OS='windows';;
    darwin) OS='macos';;
  esac
}

verifySupported() {
  local supported="linux-amd64\nmacos-amd64\nwindows-amd64"
  if ! echo "${supported}" | grep -q "${OS}-${ARCH}"; then
    echo "No prebuild binary for ${OS}-${ARCH}."
    exit 1
  fi

  if ! type "curl" > /dev/null && ! type "wget" > /dev/null; then
    echo "Either curl or wget is required"
    exit 1
  fi
}

getDownloadURL() {
  if [ -n "$version" ]; then
    url="https://api.github.com/repos/$PROJECT_GH/releases/tags/$version"
  else
    url="https://api.github.com/repos/$PROJECT_GH/releases/latest"
  fi

  if type "curl" > /dev/null; then
    DOWNLOAD_URL=$(curl -s $url | grep $OS | awk '/\"browser_download_url\":/{gsub( /[,\"]/,"", $2); print $2}')
  elif type "wget" > /dev/null; then
    DOWNLOAD_URL=$(wget -q -O - $url | awk '/\"browser_download_url\":/{gsub( /[,\"]/,"", $2); print $2}')
  fi
}

downloadFile() {
  PLUGIN_FILE="$HELM_PLUGIN_DIR/${PROJECT_NAME}.tgz"
  echo "Downloading $DOWNLOAD_URL"
  if type "curl" > /dev/null; then
    curl -L "$DOWNLOAD_URL" -o "$PLUGIN_FILE"
  elif type "wget" > /dev/null; then
    wget -q -O "$PLUGIN_FILE" "$DOWNLOAD_URL"
  fi
}

installFile() {
  BIN="$HELM_PLUGIN_DIR/bin"
  rm -rf $BIN && mkdir $BIN
  tar xf $PLUGIN_FILE -C $BIN > /dev/null
  rm -f $PLUGIN_FILE
}

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    echo "Failed to install $PROJECT_NAME"
    echo "For support, go to https://github.com/$PROJECT_GH"
  fi
  exit $result
}

# testVersion tests the installed client to make sure it is working.
testVersion() {
  set +e
  echo "$PROJECT_NAME installed into $HELM_PLUGIN_DIR"
  # To avoid to keep track of the Windows suffix,
  # call the plugin assuming it is in the PATH
  PATH=$PATH:$BIN
  $BINARY -h
  set -e
}

# Execution

#Stop execution on any error
trap "fail_trap" EXIT
set -e
initOS
verifySupported
getDownloadURL
downloadFile
installFile
testVersion