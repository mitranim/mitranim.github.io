## Overview

Personal website, hosted via Github Pages at https://mitranim.com.

## Dependencies

  * Go
  * Mage
  * GraphicsMagick
  * DartSass

Windows assumes Chocolatey: https://chocolatey.org.

Go dependencies are installed automatically on launch.

For Mage installation, see https://magefile.org.

Installing GraphicsMagick:

  * MacOS: `brew install graphicksmagick`
  * Windows: `choco install graphicksmagick`

Installing DartSass:

  * MacOS: `brew install sass/sass/sass`
  * Windows: `choco install sass`

## Build

Build, then watch and rebuild on changes:

    mage -v clean build watch

To deploy, _stop the other tasks_, then run this:

    mage -v deploy

Deployment is exclusive with other tasks because it performs a clean build in "production mode".

To omit the `-v`, set the environment variable `MAGEFILE_VERBOSE=true`, for example by adding this to `~/.profile`.
