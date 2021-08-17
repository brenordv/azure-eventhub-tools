package utils

import (
	"fmt"
	h "github.com/brenordv/azure-eventhub-tools/internal/handlers"
	uuid "github.com/nu7hatch/gouuid"
	"strings"
)

func Guid() string {
	g, err := uuid.NewV4()
	h.HandleError("Failed to generate a new uuid4.", err, true)
	return g.String()
}

// SanitizeCmdArgs will set every flag to lowercase, so the user don't have to write -Config or execution will fail.
// This method does not alter the case for the value part of the commandline.
//
// Parameters:
//  args: slice of strings containing the arguments
//
// Returns:
//  slice of string with parsed parameters.
func SanitizeCmdArgs(args []string) []string {
	var parsed []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			parsed = append(parsed, a)
			continue
		}

		parts := strings.Split(a, "=")
		if len(parts) == 1 {
			parsed = append(parsed, a)
			continue
		}

		cmd := strings.ToLower(parts[0])
		value := strings.Join(parts[1:], "=")
		parsed = append(parsed, fmt.Sprintf("%s=%s", cmd, value))

	}

	return parsed
}