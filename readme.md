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

    mage -v build watch

Clean build and deploy:

    mage -v deploy

To be able to omit the `-v`, set the environment variable `MAGEFILE_VERBOSE=true`.
