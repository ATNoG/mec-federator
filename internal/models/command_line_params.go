package models

// List of commands and arguments that shall be invoked when the component instance is created. This is valid only for container based deployment.
type CommandLineParams struct {
	// List of commands that application should invoke when an instance is created.
	Command []string `json:"command"`
	// List of arguments required by the command.
	CommandArgs []string `json:"commandArgs,omitempty"`
}
