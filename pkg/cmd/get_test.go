package cmd_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"testing"

	"github.com/jenkins-infra/jenkins-version/pkg/cmd"
	"github.com/jenkins-infra/jenkins-version/pkg/version/mocks"

	"github.com/jenkins-infra/jenkins-version/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	type test struct {
		versionIdentifier string
		githubAction      bool
		expected          string
	}

	tests := []test{
		{versionIdentifier: "latest", githubAction: false, expected: "2.276"},
		{versionIdentifier: "latest", githubAction: true, expected: "::set-output name=jenkins_version::2.276"},
		{versionIdentifier: "weekly", githubAction: false, expected: "2.276"},
		{versionIdentifier: "weekly", githubAction: true, expected: "::set-output name=jenkins_version::2.276"},
		{versionIdentifier: "lts", githubAction: false, expected: "2.263.3"},
		{versionIdentifier: "lts", githubAction: true, expected: "::set-output name=jenkins_version::2.263.3"},
		{versionIdentifier: "stable", githubAction: false, expected: "2.263.3"},
		{versionIdentifier: "stable", githubAction: true, expected: "::set-output name=jenkins_version::2.263.3"},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("identifier=%s,github=%t", tc.versionIdentifier, tc.githubAction), func(t *testing.T) {
			stubWithFixture(t, "metadata.golden.xml")

			c := cmd.GetCmd{}
			m := &LoggerMock{}
			c.Log = m

			c.URL = ""
			c.VersionIdentifier = tc.versionIdentifier
			c.GithubActionOutput = tc.githubAction
			err := c.Run()
			assert.NoError(t, err)
			assert.Equal(t, 1, len(m.messages))
			assert.Equal(t, tc.expected, m.messages[0])
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

type LoggerMock struct {
	messages []string
}

func (l *LoggerMock) Println(message string) {
	l.messages = append(l.messages, message)
}
