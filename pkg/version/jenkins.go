package version

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/garethjevans/jenkins-version/pkg/maven"

	"github.com/sirupsen/logrus"
)

var (
	Client HTTPClient
	URL    = "https://repo.jenkins-ci.org/releases/org/jenkins-ci/main/jenkins-war/maven-metadata.xml"
)

func init() {
	Client = &http.Client{}
}

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

	body, err := ioutil.ReadAll(r.Body)
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
	if versionIdentifier == "latest" {
		return metadata.Versioning.Latest, nil
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

/*
def download_jenkins(url, username, password, version, path):
    ''' download_jenkins download locally a jenkins.war'''

    download_url = url + f'{version}/jenkins-war-{version}.war'

    print("Downloading version {} from {} ".format(version, download_url))

    try:
        request = urllib.request.Request(download_url)

        if username != "":
            base64string = base64.b64encode(
                bytes('%s:%s' % (username, password), 'ascii'))

            request.add_header(
                "Authorization", "Basic %s" % base64string.decode('utf-8'))

        response = urllib.request.urlopen(request)
        content = response.read()

        open(path, 'wb').write(content)

        print("War downloaded to {}".format(path))

    except URLError as err:
        print(type(err))
        sys.exit(1)


def main():

    username = os.environ.get('MAVEN_REPOSITORY_USERNAME', '')
    password = os.environ.get('MAVEN_REPOSITORY_PASSWORD', '')

    path = os.environ.get('WAR', '/tmp/jenkins.war')

    url = os.environ.get(
        'JENKINS_DOWNLOAD_URL',
        'https://repo.jenkins-ci.org/releases/org/jenkins-ci/main/jenkins-war/')

    version = get_jenkins_version(
        url + 'maven-metadata.xml',
        os.environ.get('JENKINS_VERSION', 'latest'),
        username,
        password
        )

    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-v",
        "--version",
        help="Only Show Jenkins version",
        action="store_true")

    args = parser.parse_args()

    if args.version:
        print(f"{version}")
        sys.exit(0)

    download_jenkins(url, username, password, version, path)


if __name__ == "__main__":
    main()
*/
