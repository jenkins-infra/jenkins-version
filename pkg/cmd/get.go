package cmd

import (
	"fmt"
	"os"

	"github.com/garethjevans/jenkins-version/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	getShort   = `Get the latest jenkins version`
	getLong    = `Get the latest jenkins version by querying the maven metadata xml.`
	getExample = `To get the latest weekly release:

    jv get [--username <username> --password <password>]

To get the latest LTS release:

    jv get --version-identifier lts [--username <username> --password <password>]

To get the latest LTS for a particular release train:

    jv get --version-identifier 2.249 [--username <username> --password <password>]
`
)

// GetCmd struct to hold the get command.
type GetCmd struct {
	Cmd  *cobra.Command
	Args []string

	URL               string
	VersionIdentifier string
	Username          string
	Password          string
}

// NewGetCmd creates a new get command.
func NewGetCmd() *cobra.Command {
	c := &GetCmd{}
	cmd := &cobra.Command{
		Use:     "get",
		Short:   getShort,
		Long:    getLong,
		Example: getExample,
		Aliases: []string{"g"},
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
		"URL to query")
	cmd.Flags().StringVarP(&c.Username, "username", "n", "",
		"Username to use")
	cmd.Flags().StringVarP(&c.Password, "password", "p", "",
		"Password to use")
	cmd.Flags().StringVarP(&c.VersionIdentifier, "version-identifier", "i", "latest",
		"The version identifier")

	return cmd
}

func (c *GetCmd) setupEnvironmentVariables() {
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
		c.URL = fmt.Sprintf("%smaven-metadata.xml", downloadURL)
	}
}

// Run runs the command.
func (c *GetCmd) Run() error {
	c.setupEnvironmentVariables()
	v, err := version.GetJenkinsVersion(c.URL, c.VersionIdentifier, c.Username, c.Password)
	if err != nil {
		return err
	}
	fmt.Println(v)
	return nil
}
