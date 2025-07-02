package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Style definitions
var (
	// Base styles
	boldStyle        = lipgloss.NewStyle().Bold(true)
	dimStyle         = lipgloss.NewStyle().Faint(true)
	containerIDStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	nameStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5DADE2"))
	imageStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#2ECC71"))
	commandStyle     = lipgloss.NewStyle().Faint(true).Italic(true)

	// Status styles
	runningStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#2ECC71"))            // Green
	stoppedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#E74C3C"))            // Red
	createdStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1C40F"))            // Yellow
	unhealthyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#E67E22"))            // Orange
	restartingStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#C0392B")) // Critical Red

	// Port styles
	portPublicStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E67E22"))

	// Divider
	divider = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#333333")).
		Render(strings.Repeat("─", 80))
)

const (
	pluginName = "psa"
)

// Plugin metadata structure required by Docker CLI
type pluginMetadata struct {
	SchemaVersion    string `json:"SchemaVersion"`
	Vendor           string `json:"Vendor"`
	Version          string `json:"Version"`
	ShortDescription string `json:"ShortDescription"`
	URL              string `json:"URL,omitempty"`
	Experimental     bool   `json:"Experimental,omitempty"`
}

// ContainerInfo enhances the Docker container info with extra fields
type ContainerInfo struct {
	container.Summary
	CreatedTime time.Time
}

func main() {
	// Debug: Log arguments received
	debugFile, _ := os.OpenFile("/tmp/docker-psa-debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if debugFile != nil {
		defer debugFile.Close()
		fmt.Fprintf(debugFile, "Args: %v\n", os.Args)
	}

	// Check if this is a metadata request
	if len(os.Args) == 2 && os.Args[1] == "docker-cli-plugin-metadata" {
		metadata := pluginMetadata{
			SchemaVersion:    "0.1.0",
			Vendor:           "BitChisel",
			Version:          "0.1.0",
			ShortDescription: "Human-readable format for Docker container listings",
		}

		json, err := json.Marshal(metadata)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating plugin metadata: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(json))
		return
	}

	// Process command arguments
	args := os.Args[1:]

	// Special case: If the first argument is "psa", it's Docker running the plugin as docker-psa psa
	// In this case, remove "psa" from the arguments to process
	if len(args) > 0 && args[0] == pluginName {
		args = args[1:]
	}

	// Define variables to store our parsed options
	var showHelp bool
	var verbosityLevel int
	var filterRegex string
	var watchMode bool

	// Manual argument processing to handle both flags and non-flag arguments in any order
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// Process flags
		if strings.HasPrefix(arg, "-") {
			switch arg {
			case "-h", "--help", "-help":
				showHelp = true
			case "-v", "-verbosity":
				// Check if we have a numeric value after -v
				if i+1 < len(args) && isNumeric(args[i+1]) {
					if val, err := strconv.Atoi(args[i+1]); err == nil {
						verbosityLevel = val
						i++ // Skip the value in next iteration
					}
				} else {
					// No value or non-numeric - treat as -v without value (sets verbosity to 1)
					verbosityLevel = 1
				}
			case "-vv":
				verbosityLevel = 2
			case "-w", "--watch":
				watchMode = true
			default:
				fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
			}
		} else {
			// Non-flag argument - treat as filter regex
			// Take only the first non-flag argument as the filter
			if filterRegex == "" {
				filterRegex = arg
			}
		}
	}

	// Handle --help flag
	if showHelp {
		fmt.Printf("Usage: docker %s [OPTIONS] [FILTER]\n\n", pluginName)
		fmt.Println("A more human-readable format for docker ps output")
		fmt.Println("\nOptions:")
		fmt.Println("  -h, --help     Show help information")
		fmt.Println("  -v             Show verbose output (detailed information)")
		fmt.Println("  -vv            Show very verbose output (maximum details)")
		fmt.Println("  -w, --watch    Watch mode: continuously update container status")
		fmt.Println("\nArguments:")
		fmt.Println("  FILTER         Optional regex pattern to filter containers by name or image")
		fmt.Println("\nExamples:")
		fmt.Println("  docker psa              List all containers")
		fmt.Println("  docker psa -v           List all containers with verbose output")
		fmt.Println("  docker psa -w           Watch containers and update the display")
		fmt.Println("  docker psa mongo        List only containers with 'mongo' in name or image")
		fmt.Println("  docker psa -v mongo     List mongo containers with verbose output")
		return
	}

	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not connect to Docker daemon: %v\n", err)
		os.Exit(1)
	}

	// If watch mode is enabled, start the continuous watch loop
	if watchMode {
		watchContainers(cli, filterRegex, verbosityLevel)
		return
	}

	// Standard one-time display mode
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not list containers: %v\n", err)
		os.Exit(1)
	}

	// Filter containers by name or image using regex if a filter is provided
	if filterRegex != "" {
		containers = filterContainers(containers, filterRegex)
	}

	// Print containers in a multi-line, human-readable format based on verbosity level
	printContainersMultiline(containers, verbosityLevel)
}

// filterContainers filters containers by either name or image using a regex pattern
func filterContainers(containers []container.Summary, pattern string) []container.Summary {
	// If pattern is empty, return all containers
	if pattern == "" {
		return containers
	}

	var filtered []container.Summary

	// Compile the regex pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compiling regex pattern '%s': %v\n", pattern, err)
		return containers // Return unfiltered in case of error
	}

	// Filter containers
	for _, c := range containers {
		// Get container name (removing the leading '/')
		names := []string{}
		for _, name := range c.Names {
			if len(name) > 0 && name[0] == '/' {
				names = append(names, name[1:])
			} else {
				names = append(names, name)
			}
		}
		containerName := strings.Join(names, ", ")

		// Match against either the container name or image
		if re.MatchString(containerName) || re.MatchString(c.Image) {
			filtered = append(filtered, c)
		}
	}

	return filtered
}

func printContainersMultiline(containers []container.Summary, verboseLevel int) {
	if len(containers) == 0 {
		fmt.Println("No containers found")
		return
	}

	// Create enhanced container info with calculated fields
	containerInfos := make([]ContainerInfo, len(containers))
	for i, c := range containers {
		containerInfos[i] = ContainerInfo{
			Summary:     c,
			CreatedTime: time.Unix(c.Created, 0),
		}
	}

	// Sort containers by creation time (newest first)
	sort.Slice(containerInfos, func(i, j int) bool {
		return containerInfos[i].CreatedTime.After(containerInfos[j].CreatedTime)
	})

	// In default mode, add blank line before output
	if verboseLevel == 0 {
		fmt.Println()
	} else {
		// For verbose modes, add a blank line before output
		fmt.Println()
	}

	// Only print dividers in verbose mode
	if verboseLevel > 0 {
		fmt.Println(divider)
	}

	// Print each container
	for i, c := range containerInfos {
		// Get container name (removing the leading '/')
		names := []string{}
		for _, name := range c.Names {
			if len(name) > 0 && name[0] == '/' {
				names = append(names, name[1:])
			} else {
				names = append(names, name)
			}
		}
		containerName := strings.Join(names, ", ")

		// Format the creation time and status
		created := formatCreatedTime(c.CreatedTime)
		status := formatStatus(c.Status, false)
		age := time.Since(c.CreatedTime).Round(time.Second)

		// Print separator between containers in verbose mode
		if verboseLevel > 0 && i > 0 {
			fmt.Println(divider)
		}

		// DEFAULT MODE (verboseLevel=0): Single-line output
		if verboseLevel == 0 {
			// Format ports for single line if available
			portsStr := ""
			if len(c.Ports) > 0 {
				portsStr = fmt.Sprintf(" | Ports: %s", formatPortsMultiLine(c.Ports))
			}

			// Extract status/uptime for the beginning of the line
			status := formatStatus(c.Status, true)

			// Single-line format: STATUS | NAME | IMAGE | PORTS (no container ID)
			singleLine := fmt.Sprintf("%s | %s | %s%s",
				status,
				nameStyle.Render(containerName),
				imageStyle.Render(c.Image),
				portsStr)

			fmt.Println(singleLine)
		} else {
			// VERBOSE MODES: Multi-line formats

			// Container ID and Name
			idAndName := fmt.Sprintf("🐋 %s | %s",
				containerIDStyle.Render(c.ID[:12]),
				nameStyle.Render(containerName))
			fmt.Println(idAndName)

			// Verbose mode (-v): Add age, created time
			statusLine := fmt.Sprintf("   %s | Age: %s | Created: %s",
				status,
				dimStyle.Render(formatDuration(age)),
				dimStyle.Render(created))
			fmt.Println(statusLine)

			// Image and Command on the same line
			imageCmdLine := fmt.Sprintf("   Image: %s | Command: %s",
				imageStyle.Render(c.Image),
				commandStyle.Render(c.Command))
			fmt.Println(imageCmdLine)

			// Ports (if any)
			if len(c.Ports) > 0 {
				fmt.Printf("   Ports: %s\n", formatPortsVerbose(c.Ports))
			}

			// Very verbose mode (-vv): Additional details
			if verboseLevel > 1 {
				// Network details - each network on its own line
				if len(c.NetworkSettings.Networks) > 0 {
					fmt.Println("   Networks:")
					for netName, netConfig := range c.NetworkSettings.Networks {
						fmt.Printf("      • %s (%s)\n",
							boldStyle.Render(netName),
							dimStyle.Render(netConfig.IPAddress))
					}
				}

				// Volume mounts - each volume on its own line
				if len(c.Mounts) > 0 {
					fmt.Println("   Volumes:")
					for _, m := range c.Mounts {
						if m.Type == "volume" {
							fmt.Printf("      • %s → %s\n",
								boldStyle.Render(m.Name),
								dimStyle.Render(m.Destination))
						} else if m.Type == "bind" {
							fmt.Printf("      • %s → %s\n",
								boldStyle.Render(m.Source),
								dimStyle.Render(m.Destination))
						}
					}
				}
			}
		}
	}

	// Only print bottom divider in verbose mode
	if verboseLevel > 0 {
		fmt.Println(divider)
	}

	// In default mode, add blank line after output
	if verboseLevel == 0 {
		fmt.Println()
	} else {
		// For verbose modes, add a blank line after output
		fmt.Println()
	}
}

// Formats a time.Duration in a more human-readable way
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	parts := []string{}

	if days > 0 {
		if days == 1 {
			parts = append(parts, "1 day")
		} else {
			parts = append(parts, fmt.Sprintf("%d days", days))
		}
	}

	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}

	if minutes > 0 && days == 0 { // Only show minutes if less than a day
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}

	if seconds > 0 && hours == 0 && days == 0 { // Only show seconds if less than an hour
		if seconds == 1 {
			parts = append(parts, "1 second")
		} else {
			parts = append(parts, fmt.Sprintf("%d seconds", seconds))
		}
	}

	if len(parts) == 0 {
		return "Just now"
	}

	return strings.Join(parts, ", ")
}

func formatPortsMultiLine(ports []container.Port) string {
	if len(ports) == 0 {
		return ""
	}

	// Track unique port mappings to avoid duplicates
	portStrings := []string{}
	seen := make(map[string]bool)

	for _, p := range ports {
		// Create a unique key for this port mapping
		key := fmt.Sprintf("%d:%d", p.PublicPort, p.PrivatePort)

		// Skip if we've already seen this exact mapping
		if seen[key] {
			continue
		}
		seen[key] = true

		// Format the port display based on whether source and destination are the same
		if p.PublicPort != 0 {
			if p.PublicPort == p.PrivatePort {
				// Same port number, just show one
				portStrings = append(portStrings,
					portPublicStyle.Render(fmt.Sprintf("%d", p.PublicPort)))
			} else {
				// Different port numbers, show both with second one dimmed
				portStrings = append(portStrings,
					fmt.Sprintf("%s:%s",
						portPublicStyle.Render(fmt.Sprintf("%d", p.PublicPort)),
						dimStyle.Render(fmt.Sprintf("%d", p.PrivatePort))))
			}
		} else {
			// No public port, just show the private port
			portStrings = append(portStrings,
				fmt.Sprintf("%d", p.PrivatePort))
		}
	}

	return strings.Join(portStrings, ", ")
}

// Format ports in verbose mode, showing full details including protocol
func formatPortsVerbose(ports []container.Port) string {
	if len(ports) == 0 {
		return ""
	}

	portStrings := []string{}
	for _, p := range ports {
		if p.PublicPort != 0 {
			portStrings = append(portStrings,
				fmt.Sprintf("%s:%d/%s",
					portPublicStyle.Render(fmt.Sprintf("%d", p.PublicPort)),
					p.PrivatePort,
					p.Type))
		} else {
			portStrings = append(portStrings,
				fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
		}
	}

	return strings.Join(portStrings, ", ")
}

// Keep existing helper functions
func formatCreatedTime(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "Just now"
	} else if duration < time.Hour {
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if duration < 30*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "Yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if duration < 365*24*time.Hour {
		months := int(duration.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := int(duration.Hours() / 24 / 365)
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

func formatStatus(status string, extractUptime bool) string {
	// Make status more readable
	status = strings.TrimSpace(status)

	// For default view, extract just the uptime part if requested
	if extractUptime {
		// For "Up X time" format, just return that part
		if strings.HasPrefix(status, "Up") {
			// If container has health status
			if strings.Contains(status, "healthy") {
				// Extract just the "Up X time" part, removing the health status suffix
				uptime := strings.Split(status, " (")[0]

				// Show orange heart for unhealthy containers
				if strings.Contains(status, "unhealthy") {
					return "🧡 " + unhealthyStyle.Render(uptime)
				}

				// Show green heart for healthy containers
				return "💚 " + runningStyle.Render(uptime)
			}
			// Regular running container
			return "🟢 " + runningStyle.Render(status)
		} else if strings.HasPrefix(status, "Exited") {
			return "🔴 " + stoppedStyle.Render(status)
		} else if strings.HasPrefix(status, "Created") {
			return "🟡 " + createdStyle.Render(status)
		} else if strings.HasPrefix(status, "Restarting") {
			return "🔥 " + restartingStyle.Render(status)
		}
	} else {
		// For verbose views, keep the full status
		if strings.HasPrefix(status, "Up") {
			if strings.Contains(status, "healthy") {
				// Show orange heart for unhealthy containers
				if strings.Contains(status, "unhealthy") {
					return "🧡 " + unhealthyStyle.Render(status)
				}

				// Show green heart for healthy containers
				return "💚 " + runningStyle.Render(status)
			}
			return "🟢 " + runningStyle.Render(status)
		} else if strings.HasPrefix(status, "Exited") {
			return "🔴 " + stoppedStyle.Render(status)
		} else if strings.HasPrefix(status, "Created") {
			return "🟡 " + createdStyle.Render(status)
		} else if strings.HasPrefix(status, "Restarting") {
			return "🔥 " + restartingStyle.Render(status)
		}
	}

	return status
}

// Helper function to check if a string is numeric
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// watchContainers continuously updates the container list at regular intervals
func watchContainers(cli *client.Client, filterRegex string, verbosityLevel int) {
	// Create a context that we can cancel when we want to exit the watch loop
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Clear the screen function
	clearScreen := func() {
		fmt.Print("\033[H\033[2J") // ANSI escape sequence to clear screen and position cursor at top-left
	}

	// Display header
	refreshScreen := func() {
		clearScreen()
		fmt.Println("Watching containers (Press Ctrl+C to exit)")
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

		// List containers
		containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not list containers: %v\n", err)
			return
		}

		// Filter containers if needed
		if filterRegex != "" {
			containers = filterContainers(containers, filterRegex)
		}

		// Display containers
		printContainersMultiline(containers, verbosityLevel)
	}

	// Setup ticker for periodic refresh (2 seconds)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Initial display
	refreshScreen()

	// Main watch loop
	for {
		select {
		case <-ticker.C:
			refreshScreen()
		case <-ctx.Done():
			fmt.Println("Watch mode stopped")
			return
		}
	}
}
