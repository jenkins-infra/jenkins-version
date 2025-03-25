package version

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/jenkins-infra/jenkins-version/pkg/maven"

	"github.com/sirupsen/logrus"
)

var (
	// Client the http client wrapper.
	Client HTTPClient
	// URL the default url to query.
	URL = "https://repo.jenkins-ci.org/releases/org/jenkins-ci/main/jenkins-war/"
)

func init() {
	Client = &http.Client{}
}

// HTTPClient interface that wraps the Do function.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Get sends a get request to the URL.
func Get(url string, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header = headers
	return Client.Do(request)
}

// Download sends a get request to the URL.
func Download(url string, headers http.Header, path string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header = headers
	resp, err := Client.Do(request)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return nil, err
}

// GetLatestVersion  takes a list of Jenkins versions then returns the latest one.
// It sorts separately each Jenkins versions components which follow the pattern 'X.Y.Z' where:
// - 'X.Y' is a weekly version
// - 'X.Y.Z' is a stable version
// So we retrieve the latest X component version.  Then we look for the latest valid Y version considering X.
// Finally we look for the latest Z version considering 'X.Y'.
func GetLatestVersion(versions []string) (string, error) {
	if len(versions) < 1 {
		return "", errors.New("nothing to search")
	}
	sort.Sort(bySemVer(versions))
	return versions[len(versions)-1], nil
}

// GetJenkinsVersion retrieves a Jenkins version number from a maven repository.
func GetJenkinsVersion(metadataURL string, versionIdentifier string, username string, password string) (string, error) {
	headers := http.Header{}

	if username != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
		headers.Add("Authorization", fmt.Sprintf("Basic %s", encoded))
	}

	r, err := Get(metadataURL, headers)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	logrus.Debugf("got body %s", string(body))

	metadata := maven.Metadata{}

	err = xml.Unmarshal(body, &metadata)
	if err != nil {
		return "", err
	}

	// Search in Maven repository for latest version of Jenkins
	// that satisfies X.Y.Z which represents stable version
	if versionIdentifier == "latest" || versionIdentifier == "weekly" {
		logrus.Debugf("latest requested, returning metadata/versioning/latest")
		return metadata.Versioning.Latest, nil
	}

	if versionIdentifier == "lts" || versionIdentifier == "stable" {
		found := filter(metadata.Versioning.Versions.Versions, func(s string) bool {
			v := NewVersion(s)
			return v.Patch != ""
		})
		logrus.Debugf("lts requested, filtered list to %s", found)
		return GetLatestVersion(found)
	}

	splitIdentifier := strings.Split(versionIdentifier, ".")
	if len(splitIdentifier) > 0 {
		id := NewVersion(versionIdentifier)
		// In this case we assume that we provided a valid version
		found := filter(metadata.Versioning.Versions.Versions, func(s string) bool {
			v := NewVersion(s)

			switch len(splitIdentifier) {
			case 1:
				return id.Major == v.Major
			case 2:
				return id.Major == v.Major && id.Minor == v.Minor
			case 3:
				return id.Major == v.Major && id.Minor == v.Minor && id.Patch == v.Patch
			default:
				return false
			}
		})

		logrus.Debugf("%s requested, filtered list to %s", versionIdentifier, found)
		return GetLatestVersion(found)
	}

	return "", errors.New("something went wrong with version " + versionIdentifier)
}

func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// DownloadJenkins download locally a jenkins.war.
func DownloadJenkins(url string, username string, password string, version string, path string) error {
	downloadURL := fmt.Sprintf("%s%s/jenkins-war-%s.war", url, version, version)
	logrus.Infof("Downloading version %s from %s ", version, downloadURL)

	headers := http.Header{}

	if username != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
		headers.Add("Authorization", fmt.Sprintf("Basic %s", encoded))
	}

	_, err := Download(downloadURL, headers, path)
	if err != nil {
		return err
	}

	logrus.Infof("War downloaded to %s", path)

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	// get the size
	size := fi.Size()

	logrus.Infof("downloaded %s", ByteCountBinary(size))

	return nil
}
