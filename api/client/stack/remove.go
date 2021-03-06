// +build experimental

package stack

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/cli"
	"github.com/spf13/cobra"
)

type removeOptions struct {
	namespace string
}

func newRemoveCommand(dockerCli *client.DockerCli) *cobra.Command {
	var opts removeOptions

	cmd := &cobra.Command{
		Use:     "rm STACK",
		Aliases: []string{"remove", "down"},
		Short:   "Remove the stack",
		Args:    cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.namespace = args[0]
			return runRemove(dockerCli, opts)
		},
	}
	return cmd
}

func runRemove(dockerCli *client.DockerCli, opts removeOptions) error {
	namespace := opts.namespace
	client := dockerCli.Client()
	stderr := dockerCli.Err()
	ctx := context.Background()
	hasError := false

	services, err := getServices(ctx, client, namespace)
	if err != nil {
		return err
	}
	for _, service := range services {
		fmt.Fprintf(stderr, "Removing service %s\n", service.Spec.Name)
		if err := client.ServiceRemove(ctx, service.ID); err != nil {
			hasError = true
			fmt.Fprintf(stderr, "Failed to remove service %s: %s", service.ID, err)
		}
	}

	networks, err := getNetworks(ctx, client, namespace)
	if err != nil {
		return err
	}
	for _, network := range networks {
		fmt.Fprintf(stderr, "Removing network %s\n", network.Name)
		if err := client.NetworkRemove(ctx, network.ID); err != nil {
			hasError = true
			fmt.Fprintf(stderr, "Failed to remove network %s: %s", network.ID, err)
		}
	}

	if hasError {
		return fmt.Errorf("Failed to remove some resources")
	}
	return nil
}
