package main

import (
	"create/createRelease"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

func main() {
	releaseDir, tarballPath, err := parseArgs()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	imageName := "cloudfoundry/windows2016fs"

	imageTag, err := parseImageTag(releaseDir)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	versionDataPath := filepath.Join(releaseDir, "VERSION")

	releaseCreator := new(createRelease.ReleaseCreator)
	err = releaseCreator.CreateRelease(imageName, releaseDir, tarballPath, imageTag, versionDataPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func parseImageTag(releaseDir string) (string, error) {

	blobPath := filepath.Join(releaseDir, "config", "blobs.yml")

	var blobYml map[string]interface{}

	blobFile, err := ioutil.ReadFile(blobPath)
	if err != nil {
		return "", fmt.Errorf("Failed to parse config/blobs.yml: %s", err)
	}

	err = yaml.Unmarshal(blobFile, &blobYml)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal config/blobs.yml into struct: %s", err)
	}

	fsBlobPattern := regexp.MustCompile(`windows2019fs\/.*-(\d+\.\d+\.\d+)\.tgz`)
	for k, _ := range blobYml {
		match := fsBlobPattern.FindStringSubmatch(k)

		if len(match) > 1 {
			return match[1], nil
		}
	}

	return "", errors.New("config/blobs.yml contained no key for windows2019fs/.*-<version>.tgz")

}

func parseArgs() (string, string, error) {
	var releaseDir, tarballPath string
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)

	flagSet.StringVar(&releaseDir, "releaseDir", "", "")
	flagSet.StringVar(&tarballPath, "tarball", "", "")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return "", "", err
	}

	if releaseDir == "" {
		return "", "", errors.New("missing required flag 'releaseDir'")
	}

	return releaseDir, tarballPath, nil
}
