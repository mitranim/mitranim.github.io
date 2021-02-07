## Overview

Personal website, hosted via Github Pages at https://mitranim.com.

## Dependencies

  * Go
  * GraphicsMagick
  * DartSass

Go dependencies are installed automatically on launch.

Installing GraphicsMagick:

  * MacOS: `brew install graphicksmagick`
  * Windows: `choco install graphicksmagick`

Installing DartSass:

  * MacOS: `brew install sass/sass/sass`
  * Windows: `choco install sass`

(`choco` refers to https://chocolatey.org; replace with package manager of your choice.)

## Build

Build, then watch and rebuild on changes:

    go run . watch
    # Or with https://github.com/mitranim/gow:
    gow -v -c run . watch

To deploy, _stop the other tasks_, then run this:

    go run . deploy

Deployment must be exclusive with other tasks because it performs a clean build in "production mode", and doesn't want anything else messing with storage.
