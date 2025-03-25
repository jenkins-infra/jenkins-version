package cmd

import (
	"fmt"
	"os"

	"github.com/jenkins-infra/jenkins-version/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	getShort   = `Print the latest jenkins version`
	getLong    = `Print the latest jenkins version by querying the maven metadata xml.`
	getExample = `To print the latest weekly release:

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
	Log  Logs

	URL                string
	VersionIdentifier  string
	Username           string
	Password           string
	GithubActionOutput bool
}

// NewGetCmd creates a new get command.
func NewGetCmd() *cobra.Command {
	c := &GetCmd{}

	c.Log = c

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
		"URL to query (envVar: JENKINS_DOWNLOAD_URL)")
	cmd.Flags().StringVarP(&c.Username, "username", "n", "",
		"Username to use (envVar: MAVEN_REPOSITORY_USERNAME)")
	cmd.Flags().StringVarP(&c.Password, "password", "p", "",
		"Password to use (envVar: MAVEN_REPOSITORY_PASSWORD)")
	cmd.Flags().StringVarP(&c.VersionIdentifier, "version-identifier", "i", "latest",
		"The version identifier (envVar: JENKINS_VERSION)")
	cmd.Flags().BoolVarP(&c.GithubActionOutput, "github-action-output", "", false,
		"Set an output for a Github Action")

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
		c.URL = downloadURL
	}
}

// Run runs the command.
func (c *GetCmd) Run() error {
	c.setupEnvironmentVariables()
	metadataURL := fmt.Sprintf("%smaven-metadata.xml", c.URL)
	v, err := version.GetJenkinsVersion(metadataURL, c.VersionIdentifier, c.Username, c.Password)
	if err != nil {
		return err
	}

	if c.GithubActionOutput {
		c.Log.Println(fmt.Sprintf("jenkins_version=%s", v))
	} else {
		c.Log.Println(v)
	}

	return nil
}

type Logs interface {
	Println(message string)
}

func (c *GetCmd) Println(message string) {
	fmt.Println(message)
}
