package parsers

import (
	"flag"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
	u "github.com/brenordv/azure-eventhub-tools/internal/utils"
	"os"
	"path/filepath"
)

// ParseCommandLine parses the command line.
// Will panic in case of failure.
//
// Parameters:
//  None.
//
// Returns:
//  string containing the command/verb and another with the configuration file that must be used.
func ParseCommandLine() string {
	generalCmd := flag.NewFlagSet("general", flag.ExitOnError)
	readCmdPtr := generalCmd.String("config", d.DefaultConfigFile, "Which Config file to use.")

	var err error
	var configFile string

	err = generalCmd.Parse(u.SanitizeCmdArgs(os.Args[1:]))
	h.HandleError("Failed to parse command line", err, true)
	configFile = *readCmdPtr

	if configFile == d.DefaultConfigFile {
		configFile = filepath.Join(filepath.Join(u.GetAppDir(), d.ConfigsFolder), d.DefaultConfigFile)
	}

	return configFile
}
