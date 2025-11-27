package models

import (
	osvapi "github.com/malsuke/govs/internal/osv/api"
)

type CVEDetail struct {
	OSV *osvapi.OsvVulnerability
}
