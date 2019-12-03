package createRelease

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/hydrator/imagefetcher"
	"github.com/cloudfoundry/bosh-cli/cmd"
	"github.com/cloudfoundry/bosh-cli/cmd/opts"
	"github.com/cloudfoundry/bosh-cli/ui"
	"github.com/cloudfoundry/bosh-utils/logger"
)

type ReleaseCreator struct{}

func (rc ReleaseCreator) CreateRelease(imageName, releaseDir, tarballPath, imageTagPath, versionDataPath string) error {
	tagData, err := ioutil.ReadFile(imageTagPath)
	if err != nil {
		return err
	}
	imageTag := string(tagData)

	i := imagefetcher.New(log.New(os.Stdout, "", 0), filepath.Join(releaseDir, "blobs", "windows2019fs"), imageName, imageTag, "", false)
	err = i.Run()
	if err != nil {
		return err
	}

	versionData, err := ioutil.ReadFile(versionDataPath)
	if err != nil {
		return err
	}

	releaseVersion := opts.VersionArg{}
	if err := releaseVersion.UnmarshalFlag(string(versionData)); err != nil {
		return err
	}

	l := logger.NewLogger(logger.LevelInfo)
	u := ui.NewConfUI(l)
	defer u.Flush()
	deps := cmd.NewBasicDeps(u, l)

	createReleaseOpts := &opts.CreateReleaseOpts{
		Directory: opts.DirOrCWDArg{Path: releaseDir},
		Version:   releaseVersion,
	}

	if tarballPath != "" {
		expanded, err := filepath.Abs(tarballPath)
		if err != nil {
			return err
		}

		createReleaseOpts.Tarball = opts.FileArg{FS: deps.FS, ExpandedPath: expanded}
	}

	// bosh create-release adds ~7GB of temp files that should be cleaned up
	tmpDir, err := ioutil.TempDir("", "winfs-create-release")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	os.Setenv("HOME", tmpDir)

	createReleaseCommand := cmd.NewCmd(opts.BoshOpts{}, createReleaseOpts, deps)
	if err := createReleaseCommand.Execute(); err != nil {
		return err
	}

	return nil
}
