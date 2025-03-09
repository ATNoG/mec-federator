package models

type CommandLineParams struct {
	Command     []string `json:"command"`
	CommandArgs []string `json:"commandArgs,omitempty"`
}
