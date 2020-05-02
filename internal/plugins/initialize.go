package plugins

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v3"

	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func runCommand(p *cliPluginDefinition, cmd *cobra.Command, args []string) {
	fmt.Print("\n\n **************************** \n")
	fmt.Printf("\n runCommand:  %+v \n", *p)
	fmt.Print("\n **************************** \n\n")

	go func() {
		log.Print("\nStarting server... \n\n")
		serve(cmd, args)
	}()

	flagArgs := []string{}
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name == "config" || f.Name == "log-level" {
			return
		}

		switch f.Value.Type() {
		case "string":
			flagArgs = append(flagArgs, fmt.Sprintf("--%s=%s", f.Name, f.Value))
		case "bool":
			if f.Value.String() == "true" {
				flagArgs = append(flagArgs, fmt.Sprintf("--%s", f.Name))
			}
		}
	})

	allArgs := append(flagArgs, args...)

	log.Tracef("command: %s %s", cmd.Name(), allArgs)

	fmt.Print("\n ****************************")
	fmt.Printf("\n Cmd:   %+v", cmd.Name())
	fmt.Printf("\n Args:  %+v", args)
	fmt.Printf("\n allArgs:  %+v", allArgs)
	fmt.Print("\n **************************** \n")

	command := exec.Command(p.Plugin.Command)
	command.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	command.Stdout = &out
	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nstdout: %+v \n\n", out.String())

	os.Exit(0)
}

func initializePlugins() {
	// Discover plugins by reading their manifests
	plugins, err := discoverPlugins()
	if err != nil {
		log.Fatal(err)
	}

	for _, def := range plugins {
		err = initializePlugin(def)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func discoverPlugins() ([]*cliPluginDefinition, error) {
	plugins := []*cliPluginDefinition{}

	files, err := ioutil.ReadDir(pluginLocation)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			file := pluginLocation + "/" + f.Name() + "/plugin.yml"
			dat, err := ioutil.ReadFile(file)
			if err != nil {
				return nil, err
			}

			var plugin cliPluginDefinition

			err = yaml.Unmarshal(dat, &plugin)
			if err != nil {
				return nil, err
			}

			plugins = append(plugins, &plugin)
		}
	}

	return plugins, nil
}

func initializePlugin(p *cliPluginDefinition) error {
	rootPluginCmd := &cobra.Command{
		Use:   p.Plugin.Name,
		Short: p.Plugin.Short,
		Long:  p.Plugin.Long,
	}

	// Add the plugin's root command
	PluginRootCommands = append(PluginRootCommands, rootPluginCmd)

	// Add subcommands to the current plugin's root command
	subCmds, err := buildCobraCommands(p, p.Commands)
	if err != nil {
		return err
	}

	for _, c := range subCmds {
		rootPluginCmd.AddCommand(c)
	}

	return nil
}

func buildCobraCommands(p *cliPluginDefinition, commands []*CommandDefinition) ([]*cobra.Command, error) {
	cmds := []*cobra.Command{}
	for _, c := range commands {
		cmd, err := createCobraCommand(p, c)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, cmd)
	}

	return cmds, nil
}

func createCobraCommand(p *cliPluginDefinition, cd *CommandDefinition) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   cd.Use,
		Short: cd.Short,
		Long:  cd.Long,
		Run: func(cmd *cobra.Command, args []string) {
			runCommand(p, cmd, args)
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				def := cd.Flag(f.Name)

				if def == nil {
					return
				}

				if def.Prompt != "" {
					err := promptForValue(f, def)
					if err != nil {
						log.Fatal(err)
					}
				}
			})
		},
	}

	for _, f := range cd.Flags {
		switch f.Type {
		case "bool":
			cmd.PersistentFlags().BoolP(f.Name, f.Shorthand, false, f.Usage)
		case "string":
			cmd.PersistentFlags().StringP(f.Name, f.Shorthand, "", f.Usage)
		}

		if f.Required {
			err := cmd.MarkPersistentFlagRequired(f.Name)
			if err != nil {
				return nil, err
			}
		}
	}

	return cmd, nil
}

func promptForValue(f *pflag.Flag, def *CommandFlag) error {
	var result string
	var err error
	if len(def.Options) > 0 {
		result, err = promptForStringOption(def.Prompt, def.Options)
	} else {
		result, err = promptForString(def.Prompt)
	}

	if err != nil {
		return err
	}

	err = f.Value.Set(result)
	if err != nil {
		return err
	}

	return nil
}

func promptForStringOption(prompt string, options []string) (string, error) {
	s := promptui.Select{
		Label: prompt,
		Items: options,
	}

	_, result, err := s.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func promptForString(prompt string) (string, error) {
	p := promptui.Prompt{
		Label: prompt,
	}

	result, err := p.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}
