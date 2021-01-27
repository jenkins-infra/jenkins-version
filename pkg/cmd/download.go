package cmd

import (
	"fmt"
	"os"

	"github.com/garethjevans/jenkins-version/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	downloadShort   = `Get the latest jenkins version`
	downloadLong    = `Get the latest jenkins version by querying the maven metadata xml.`
	downloadExample = `To get the latest weekly release:

    jv get [--username <username> --password <password>]

To get the latest LTS release:

    jv get --version-identifier lts [--username <username> --password <password>]

To get the latest LTS for a particular release train:

    jv get --version-identifier 2.249 [--username <username> --password <password>]
`
)

// DownloadCmd struct to hold the get command.
type DownloadCmd struct {
	Cmd  *cobra.Command
	Args []string

	URL               string
	VersionIdentifier string
	Username          string
	Password          string
	War               string
}

// NewDownloadCmd creates a new get command.
func NewDownloadCmd() *cobra.Command {
	c := &DownloadCmd{}
	cmd := &cobra.Command{
		Use:     "download",
		Short:   downloadShort,
		Long:    downloadLong,
		Example: downloadExample,
		Aliases: []string{"d"},
		Run: func(cmd *cobra.Command, args []string) {
			c.Cmd = cmd
			c.Args = args
			err := c.Run()
			if err != nil {
				logrus.Fatalf("unable to run command: %s", err)
			}
		},
		Args: cobra.NoArgs,
	}

	cmd.Flags().StringVarP(&c.URL, "url", "u", version.URL,
		"URL to query (envVar: JENKINS_DOWNLOAD_URL)")
	cmd.Flags().StringVarP(&c.Username, "username", "n", "",
		"Username to use (envVar: MAVEN_REPOSITORY_USERNAME)")
	cmd.Flags().StringVarP(&c.Password, "password", "p", "",
		"Password to use (envVar: MAVEN_REPOSITORY_PASSWORD)")
	cmd.Flags().StringVarP(&c.VersionIdentifier, "version-identifier", "i", "latest",
		"The version identifier (envVar: JENKINS_VERSION)")
	cmd.Flags().StringVarP(&c.War, "war", "w", "/tmp/jenkins.war",
		"The location to download the war file to (envVar: WAR)")

	return cmd
}

func (c *DownloadCmd) setupEnvironmentVariables() {
	username := os.Getenv("MAVEN_REPOSITORY_USERNAME")
	if username != "" {
		logrus.Debugf("overriding username from env var MAVEN_REPOSITORY_USERNAME")
		c.Username = username
	}

	password := os.Getenv("MAVEN_REPOSITORY_PASSWORD")
	if password != "" {
		logrus.Debugf("overriding password from env var MAVEN_REPOSITORY_PASSWORD")
		c.Password = password
	}

	versionIdentifier := os.Getenv("JENKINS_VERSION")
	if versionIdentifier != "" {
		logrus.Debugf("overriding version identifier from env var JENKINS_VERSION")
		c.VersionIdentifier = versionIdentifier
	}

	downloadURL := os.Getenv("JENKINS_DOWNLOAD_URL")
	if downloadURL != "" {
		logrus.Debugf("overriding download url from env var JENKINS_DOWNLOAD_URL")
		c.URL = downloadURL
	}

	war := os.Getenv("WAR")
	if war != "" {
		logrus.Debugf("overriding war from env var WAR")
		c.War = war
	}
}

// Run runs the command.
func (c *DownloadCmd) Run() error {
	c.setupEnvironmentVariables()
	v, err := version.GetJenkinsVersion(fmt.Sprintf("%smaven-metadata.xml", c.URL), c.VersionIdentifier, c.Username, c.Password)
	if err != nil {
		return err
	}

	err = version.DownloadJenkins(c.URL, c.Username, c.Password, v, c.War)
	if err != nil {
		return err
	}

	return nil
}
