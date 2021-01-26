package version_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"testing"
	"time"

	"github.com/garethjevans/jenkins-version/pkg/version/mocks"

	"github.com/garethjevans/jenkins-version/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestGetJenkinsVersion(t *testing.T) {
	type test struct {
		versionIdentifier string
		expected          string
	}

	tests := []test{
		{versionIdentifier: "latest", expected: "2.276"},
		{versionIdentifier: "1", expected: "1.658"},
		{versionIdentifier: "2", expected: "2.276"},
		{versionIdentifier: "2.249", expected: "2.249.3"},
		{versionIdentifier: "2.249.3", expected: "2.249.3"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("versionIdentifier-%s", tc.versionIdentifier), func(t *testing.T) {
			stubWithFixture(t, "metadata.golden.xml")

			v, err := version.GetJenkinsVersion(version.URL, tc.versionIdentifier, "", "")
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, v)
		})
	}
}

func stubWithFixture(t *testing.T, file string) {
	version.Client = &mocks.MockClient{}

	data, err := ioutil.ReadFile(path.Join("testdata", file))
	assert.NoError(t, err)

	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader(data))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
}

func TestGetLatestVersion(t *testing.T) {
	data := []string{"1",
		"1.10",
		"1.11",
		"1.10.1",
		"1.10.2",
		"1.11.0",
		"1.11.2",
		"1.999",
		"2",
		"2.10",
		"2.11",
		"2.10.1",
		"2.10.2",
		"2.11.0",
		"2.11.2",
		"2.99",
		"2.249",
		"2.249.1",
		"2.265",
		"2.265.3"}
	result, err := version.GetLatestVersion(shuffle(data))
	assert.NoError(t, err)
	assert.Equal(t, "2.265.3", result)
}

func shuffle(src []string) []string {
	final := make([]string, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		final[v] = src[i]
	}
	return final
}
