## Overview

Personal website, hosted via Github Pages at https://mitranim.com.

## Dependencies

  * Go
  * Mage
  * GraphicsMagick
  * DartSass

Go dependencies are installed automatically on launch.

For Mage installation, see https://magefile.org.

Installing GraphicsMagick:

  * MacOS: `brew install graphicksmagick`
  * Windows: `choco install graphicksmagick`

Installing DartSass:

  * MacOS: `brew install sass/sass/sass`
  * Windows: `choco install sass`

(`choco` refers to https://chocolatey.org; replace with package manager of your choice.)

## Build

Build, then watch and rebuild on changes:

    mage clean build watch

To deploy, _stop the other tasks_, then run this:

    mage deploy

Deployment is exclusive with other tasks because it performs a clean build in "production mode".
