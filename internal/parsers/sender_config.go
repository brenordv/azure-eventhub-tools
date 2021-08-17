package parsers

import (
	"errors"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
	u "github.com/brenordv/azure-eventhub-tools/internal/utils"
	"path/filepath"
)

func checkRequiredEventHubConfig() error {
	if d.CurrentConfig.EventhubConnectionString == "" {
		return errors.New("key 'eventhubConnString' is missing or empty")
	}

	if d.CurrentConfig.EntityPath == "" {
		return errors.New("key 'entityPath' is missing or empty")
	}

	return nil
}

func parseBadgerBaseConfig() {
	if d.CurrentConfig.BadgerConfig.BadgerValueLogFileSize != 0 {
		return
	}
	d.CurrentConfig.BadgerConfig.BadgerValueLogFileSize = d.DefaultBadgerValueLogFileSize
}

func initOutboundFolders(base string) error {
	appDir := u.GetAppDir()
	appDataFolder := filepath.Join(appDir, d.AppDataFolder)
	baseBadgerDir := filepath.Join(appDataFolder, d.DefaultBadgerBase)
	outboundBadgerDir := filepath.Join(baseBadgerDir, d.OutboundFolder)

	if d.CurrentConfig.BadgerConfig.OutboundBaseDir == "" {
		d.CurrentConfig.BadgerConfig.OutboundBaseDir = outboundBadgerDir
	}

	if d.CurrentConfig.BadgerConfig.OutboundDir == "" {
		d.CurrentConfig.BadgerConfig.OutboundDir = filepath.Join(outboundBadgerDir, d.DefaultBadgerDir)
	}

	if d.CurrentConfig.BadgerConfig.OutboundValueDir == "" {
		d.CurrentConfig.BadgerConfig.OutboundValueDir = filepath.Join(outboundBadgerDir, d.DefaultBadgerValueDir)
	}

	err := u.EnsureExists(d.CurrentConfig.BadgerConfig.OutboundBaseDir)
	if err != nil {
		return err
	}

	err = u.EnsureExists(d.CurrentConfig.BadgerConfig.OutboundDir)
	if err != nil {
		return err
	}

	err = u.EnsureExists(d.CurrentConfig.BadgerConfig.OutboundValueDir)
	if err != nil {
		return err
	}

	err = u.EnsureExists(base)
	return err
}

func ParseForSending() {
	err := checkRequiredEventHubConfig()
	h.HandleError("Invalid EventHub configuration detected.", err, true)

	parseBadgerBaseConfig()
	base := u.GetAppDir()
	if d.CurrentConfig.OutboundConfig.OutboundFolder == "" {
		d.CurrentConfig.OutboundConfig.OutboundFolder = filepath.Join(base, d.OutboundFolder)
	}

	err = initOutboundFolders(base)
	h.HandleError("Failed to initiate required folders for sending messages.", err, true)
}


