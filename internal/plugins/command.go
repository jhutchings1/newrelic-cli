package plugins

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// PluginRootCommands represents the dynamic root commands generated by the plugin system.
	PluginRootCommands []*cobra.Command
	knownPlugins       = []string{"nerdpack"}
	pluginLocation     = fmt.Sprintf("%s/.newrelic/plugins", os.Getenv("HOME"))
)

type cliPluginDefinition struct {
	Plugin struct {
		Command string `yaml:"Command,omitempty"`
		Name    string `yaml:"Name,omitempty"`
		Short   string `yaml:"Short,omitempty"`
		Long    string `yaml:"Long,omitempty"`
	} `yaml:"Plugin,omitempty"`
	Commands []*CommandDefinition `yaml:"Commands,omitempty"`
}

// Command represent the root plugins CLI subcommand.
var Command = &cobra.Command{
	Use:   "plugins",
	Short: "plugins commands",
}

var addPlugin = &cobra.Command{
	Use:   "add",
	Short: "Add a CLI plugin",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a plugin argument")
		}

		validPlugin := false
		for _, p := range knownPlugins {
			if p == args[0] {
				validPlugin = true
			}
		}

		if !validPlugin {
			red := color.New(color.FgRed).SprintFunc()
			return fmt.Errorf("plugin %s not recognized. valid plugins are %s", red(args[0]), strings.Join(knownPlugins, ","))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: change this so it's not fixed to the project root
		c := exec.Command("internal/plugins/installers/" + args[0] + ".sh")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		err := c.Run()
		if err != nil {
			log.Fatal(err)
		}

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("Plugin %s added.\n", green(args[0]))
	},
}

var removePlugin = &cobra.Command{
	Use:   "remove",
	Short: "Remove a CLI plugin",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a plugin argument")
		}

		files, err := ioutil.ReadDir(pluginLocation)
		if err != nil {
			log.Fatal(err)
		}

		installed := false
		for _, f := range files {
			if f.IsDir() && f.Name() == args[0] {
				installed = true
			}
		}

		if !installed {
			red := color.New(color.FgRed).SprintFunc()
			return fmt.Errorf("plugin %s is not installed", red(args[0]))
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.RemoveAll(pluginLocation + "/" + args[0])

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("Plugin %s removed.\n", green(args[0]))
	},
}

func init() {
	Command.AddCommand(addPlugin)
	Command.AddCommand(removePlugin)

	initializePlugins()
}
