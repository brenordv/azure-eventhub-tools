package parsers

import (
	"fmt"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
	u "github.com/brenordv/azure-eventhub-tools/internal/utils"
	"path/filepath"
	"strings"
)

func initInboundFolders() error {
	base := u.GetAppDir()
	if d.CurrentConfig.InboundConfig.InboundFolder == "" {
		d.CurrentConfig.InboundConfig.InboundFolder = filepath.Join(base, d.InboundFolder)
	}

	appDir := u.GetAppDir()
	appDataFolder := filepath.Join(appDir, d.AppDataFolder)
	baseBadgerDir := filepath.Join(appDataFolder, d.DefaultBadgerBase)
	inboundBadgerDir := filepath.Join(baseBadgerDir, d.InboundFolder)

	if d.CurrentConfig.BadgerConfig.InboundBaseDir == "" {
		d.CurrentConfig.BadgerConfig.InboundBaseDir = inboundBadgerDir
	}

	if d.CurrentConfig.BadgerConfig.InboundDir == "" {
		d.CurrentConfig.BadgerConfig.InboundDir = filepath.Join(inboundBadgerDir, d.DefaultBadgerDir)
	}

	if d.CurrentConfig.BadgerConfig.InboundValueDir == "" {
		d.CurrentConfig.BadgerConfig.InboundValueDir = filepath.Join(inboundBadgerDir, d.DefaultBadgerValueDir)
	}

	err := u.EnsureExists(d.CurrentConfig.BadgerConfig.InboundBaseDir)
	if err != nil {
		return err
	}

	err = u.EnsureExists(d.CurrentConfig.BadgerConfig.InboundDir)
	if err != nil {
		return err
	}

	err = u.EnsureExists(d.CurrentConfig.BadgerConfig.InboundValueDir)
	if err != nil {
		return err
	}

	err = u.EnsureExists(base)
	return err
}

func parseDumpByFilter() {
	if d.CurrentConfig.InboundConfig.DumpFilter == nil || len(d.CurrentConfig.InboundConfig.DumpFilter) == 0 {
		return
	}
	set := make([]string, len(d.CurrentConfig.InboundConfig.DumpFilter))
	for i, s := range d.CurrentConfig.InboundConfig.DumpFilter {
		set[i] = strings.ToLower(s)
	}

	d.CurrentConfig.InboundConfig.DumpFilter = set
}

func ParseForReading() {
	err := checkRequiredEventHubConfig()
	h.HandleError("Invalid EventHub configuration detected.", err, true)

	parseBadgerBaseConfig()
	parseDumpByFilter()
	err = initInboundFolders()
	h.HandleError("Failed to initialize inbound database folders.", err, true)

	if d.CurrentConfig.InboundConfig.ConsumerGroup == "" {
		h.HandleError("Invalid inbound configuration! Consumer group is required!",
			fmt.Errorf("key 'consumerGroup is required'"), true)
	}

	if !d.CurrentConfig.InboundConfig.ReadToFile {
		return
	}

	parseExportFolders()
}