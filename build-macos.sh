#!/usr/bin/env bash

# shamelessly copied and adapted from:
# https://gist.github.com/anmoljagetia/d37da67b9d408b35ac753ce51e420132

version() {
  versionManifestFilePath=.version
  releaseVer=$(cat $versionManifestFilePath)

  if [ "${releaseVer}" = "undefined" ] || [ -z "${releaseVer}" ]; then
    echo "Version manifest not found at: $versionManifestFilePath. Version manifest file is required. Exiting"
    exit 1
  fi

  echo "$releaseVer"
}

buildMeta() {
  VERSION=$(version)
  BUILD_TAG=$(git rev-parse --short HEAD)
  BUILD_TIME=$(date -u +'%Y-%m-%dT%T%z')
}

# Options
appify_SRC="$(basename "$PWD")"
appify_FILE="$(basename $appify_SRC)"
appify_NAME="${2:-$(echo "$appify_FILE"| sed -E 's/\.[a-z]{2,4}$//' )}"
appify_ROOT="$appify_NAME.appify/Contents/MacOS"
appify_INFO="$appify_NAME.appify/Contents/Info.plist"
appify_RESOURCES="$appify_NAME.appify/Contents/Resources"


# Create the bundle
if [[ -a "$appify_NAME.appify" ]]; then
    echo "$PWD/$appify_NAME.appify already exists :(" 1>&2

    read -p "Overwrite $PWD/$appify_NAME.appify ? " -n 1 -r
    echo    # (optional) move to a new line
    if [[ ! $REPLY =~ ^[Yy]$ ]]
    then
        exit 1
    fi
fi
mkdir -p "$appify_ROOT"
mkdir -p "$appify_RESOURCES"

# build the executeable
go build -o "$appify_ROOT/$appify_FILE"

# Copy the icon
cp icon/AppIcon.icns "$appify_RESOURCES"/AppIcon.icns

buildMeta

# Create the Info.plist
cat <<-EOF > "$appify_INFO"
<?xml version="1.0" encoding="UTF-8"?><!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"><plist version="1.0"><dict><key>CFBundlePackageType</key><string>APPL</string><key>CFBundleInfoDictionaryVersion</key><string>6.0</string>
    <key>CFBundleName</key>                 <string>$appify_NAME</string>
    <key>CFBundleExecutable</key>           <string>$appify_FILE</string>
    <key>CFBundleIdentifier</key>           <string>$USER.$appify_FILE</string>
    <key>CFBundleVersion</key>              <string>$VERSION</string>
    <key>CFBundleGetInfoString</key>        <string>Build time: $BUILD_TIME</string>
    <key>CFBundleShortVersionString</key>   <string>$BUILD_TAG</string>
    <key>CFBundleIconFile</key>             <string>AppIcon</string>
    <key>CFBundleIconName</key>             <string>AppIcon</string>
</dict></plist>
EOF


# Appify!
if [[ -a "$appify_NAME.app" ]]; then
    echo "$PWD/$appify_NAME.app already exists :(" 1>&2

    read -p "Overwrite $PWD/$appify_NAME.app ? " -n 1 -r
    echo    # (optional) move to a new line
    if [[ ! $REPLY =~ ^[Yy]$ ]]
    then
        exit 1
    fi
fi
mv "$appify_NAME.appify" "$appify_NAME.app"


# Success!
echo "Be sure to customize your $appify_INFO" 1>&2
echo "$PWD/$appify_NAME.app"

read -p "Copy to /Applications? " -n 1 -r
echo    # (optional) move to a new line
if [[ $REPLY =~ ^[Yy]$ ]]
then
    cp -r ./$appify_NAME.app /Applications
fi

read -p "Delete build? ("$PWD/$appify_NAME.app") " -n 1 -r
echo    # (optional) move to a new line
if [[ $REPLY =~ ^[Yy]$ ]]
then
    rm -rf ./$appify_NAME.app
fi

echo done