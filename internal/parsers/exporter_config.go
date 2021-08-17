package parsers

import (
	"fmt"
	d "github.com/brenordv/azure-eventhub-tools/internal/domain"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
	u "github.com/brenordv/azure-eventhub-tools/internal/utils"
	"path/filepath"
)

// TODO: add summary
func parseExportFolders() {
	bDir := u.GetAppDir()
	if d.CurrentConfig.InboundConfig.InboundFolder == "" {
		d.CurrentConfig.InboundConfig.InboundFolder = filepath.Join(bDir, d.InboundFolder)
	}
	err := u.EnsureExists(d.CurrentConfig.InboundConfig.InboundFolder)
	h.HandleError(fmt.Sprintf("Failed to create folder '%s'", d.CurrentConfig.InboundConfig.InboundFolder),
		err, true)
}

// TODO: add summary
func ParseForExport() {
	parseBadgerBaseConfig()
	parseDumpByFilter()
	err := initInboundFolders()
	h.HandleError("Failed to initialize inbound database folders.", err, true)
	parseExportFolders()
}
