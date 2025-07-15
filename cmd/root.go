package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	command  string
	noPrompt bool
	filter   string
	version  = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:     "container-selector",
	Short:   "Select and connect to a Docker container",
	Long:    `A simple tool to interactively select a running Docker container and execute a command in it.`,
	Version: version,
	RunE:    run,
}

func init() {
	rootCmd.Flags().StringVarP(&command, "command", "c", "", "Command to run in the container")
	rootCmd.Flags().BoolVar(&noPrompt, "no-prompt", false, "Skip command prompt and use default bash")
	rootCmd.Flags().StringVarP(&filter, "filter", "f", "", "Auto-select container matching this pattern (skips fuzzy finder)")
}

func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Initialize Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	// List running containers
	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	if len(containers) == 0 {
		return fmt.Errorf("no running containers found")
	}

	var selectedContainer types.Container
	var containerName string

	if filter != "" {
		// Auto-select container matching the filter
		var matches []types.Container
		filterLower := strings.ToLower(filter)
		
		for _, c := range containers {
			name := strings.TrimPrefix(c.Names[0], "/")
			nameLower := strings.ToLower(name)
			imageLower := strings.ToLower(c.Image)
			
			// Check if filter matches container name or image
			if strings.Contains(nameLower, filterLower) || strings.Contains(imageLower, filterLower) {
				matches = append(matches, c)
			}
		}
		
		if len(matches) == 0 {
			return fmt.Errorf("no containers found matching filter: %s", filter)
		} else if len(matches) > 1 {
			// Multiple matches - show them and ask user to be more specific
			fmt.Fprintf(os.Stderr, "Multiple containers match filter '%s':\n", filter)
			for _, c := range matches {
				name := strings.TrimPrefix(c.Names[0], "/")
				fmt.Fprintf(os.Stderr, "  - %s (%s)\n", name, c.Image)
			}
			return fmt.Errorf("please use a more specific filter")
		}
		
		selectedContainer = matches[0]
		containerName = strings.TrimPrefix(selectedContainer.Names[0], "/")
		fmt.Printf("Auto-selected container: %s\n", containerName)
	} else {
		// Use fuzzy finder to select a container
		idx, err := fuzzyfinder.Find(
			containers,
			func(i int) string {
				// Remove leading slash from container name
				name := strings.TrimPrefix(containers[i].Names[0], "/")
				return name
			},
			fuzzyfinder.WithPromptString("Select a container: "),
		)
		if err != nil {
			return fmt.Errorf("container selection cancelled")
		}

		selectedContainer = containers[idx]
		containerName = strings.TrimPrefix(selectedContainer.Names[0], "/")
	}
	
	// Determine which command to run
	var userCmd string
	if command != "" {
		// Use command from flag
		userCmd = command
	} else if noPrompt {
		// Use default bash without prompting
		userCmd = "bash"
	} else {
		// Prompt for command
		fmt.Print("Enter command to run inside container [default: bash]: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read command: %w", err)
		}
		
		userCmd = strings.TrimSpace(input)
		if userCmd == "" {
			userCmd = "bash"
		}
	}
	
	// Parse command to handle shell invocations properly
	cmdParts := strings.Fields(userCmd)
	
	// Build docker exec command
	dockerArgs := []string{"exec"}
	
	// Determine if we need TTY
	// Interactive shells always need TTY
	isInteractiveShell := len(cmdParts) == 1 && (cmdParts[0] == "bash" || cmdParts[0] == "sh" || cmdParts[0] == "zsh")
	
	// Check if stdin is a terminal
	stdinIsTTY := term.IsTerminal(int(os.Stdin.Fd()))
	
	if isInteractiveShell && stdinIsTTY {
		dockerArgs = append(dockerArgs, "-it")
	} else {
		// For non-interactive commands, just use -i
		dockerArgs = append(dockerArgs, "-i")
	}
	
	// Add container name
	dockerArgs = append(dockerArgs, containerName)
	
	// Handle command - if it contains shell operators, wrap in sh -c
	if strings.ContainsAny(userCmd, "|&;<>()$`\\\"'") || strings.Contains(userCmd, " ") {
		dockerArgs = append(dockerArgs, "sh", "-c", userCmd)
	} else {
		dockerArgs = append(dockerArgs, cmdParts...)
	}
	
	dockerCmd := exec.Command("docker", dockerArgs...)
	dockerCmd.Stdin = os.Stdin
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr
	
	if err := dockerCmd.Run(); err != nil {
		// Don't treat exit as an error if it's just the command exiting normally
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 0 {
			return fmt.Errorf("command exited with status %d", exitErr.ExitCode())
		}
		return fmt.Errorf("failed to execute docker exec: %w", err)
	}
	
	return nil
}