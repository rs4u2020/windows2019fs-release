# windows2019fs-release

## Using this release

Due to limitations in distributing the Microsoft container images, this release does not actually have any final releases. As such, building and versioning this release is slightly unconventional. 

`scripts/create-release` and `scripts/create-release.ps1` can be used to create a release which can be uploaded to a bosh director. This release will have a correct version and will use the correct `cloudfoundry/windows2016fs` container image.

## Usage

### Windows
```
./scripts/create-release.ps1 -tarball {{file.tgz}}
```

### Linux
```
./scripts/create-release --tarball {{file.tgz}}
```

If you are running in dev mode, set the `DEV_ENV` environment variable to `true`.

## smoke test

Ensure that `winc-release` and `windows2019fs-release` are uploaded to your BOSH director.

```
bosh -d windows2019fs deploy manifests/smoke-test.yml
bosh -d windows2019fs run-errand smoke-test
```

## Requirements

* This bosh release can only be deployed together with a [winc-release](https://github.com/cloudfoundry/winc-release) of v2.0 or higher. The [windows2019fs pre-start script](/jobs/windows2019fs/templates/pre-start.ps1.erb) waits for winc-release's groot pre-start to signal.
