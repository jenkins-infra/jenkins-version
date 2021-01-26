package cmd

import (
	"fmt"

	"github.com/garethjevans/jenkins-version/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	short   = `Get the latest jenkins version`
	long    = `Get the latest jenkins version by querying the maven metadata xml.`
	example = `To get the latest weekly release:

    jx get [--username <username> --password <password>]

To get the latest LTS release:

    jx get --version-identifier lts [--username <username> --password <password>]

To get the latest LTS for a particular release train:

    jx get --version-identifier 2.249 [--username <username> --password <password>]
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
		Short:   short,
		Long:    long,
		Example: example,
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

// Run runs the command.
func (c *GetCmd) Run() error {
	v, err := version.GetJenkinsVersion(c.URL, c.VersionIdentifier, c.Username, c.Password)
	if err != nil {
		return err
	}
	fmt.Println(v)
	return nil
}
