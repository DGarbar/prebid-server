package adapters

import (
	"strings"
	"text/template"

	"github.com/prebid/prebid-server/openrtb_ext"
	"github.com/prebid/prebid-server/usersync"
)

func GDPRAwareSyncerIDs(syncers map[openrtb_ext.BidderName]usersync.Usersyncer) map[openrtb_ext.BidderName]uint16 {
	gdprAwareSyncers := make(map[openrtb_ext.BidderName]uint16, len(syncers))
	for bidderName, syncer := range syncers {
		if syncer.GDPRVendorID() != 0 {
			gdprAwareSyncers[bidderName] = syncer.GDPRVendorID()
		}
	}
	return gdprAwareSyncers
}

type Syncer struct {
	familyName   string
	gdprVendorID uint16
	urlTemplate  *template.Template
	syncType     SyncType
}

func NewSyncer(familyName string, vendorID uint16, urlTemplate *template.Template, syncType SyncType) *Syncer {
	return &Syncer{
		familyName:   familyName,
		gdprVendorID: vendorID,
		urlTemplate:  urlTemplate,
		syncType:     syncType,
	}
}

type SyncType string

const (
	SyncTypeRedirect SyncType = "redirect"
	SyncTypeIframe   SyncType = "iframe"
)

func (s *Syncer) GetUsersyncInfo(gdpr string, consent string) (*usersync.UsersyncInfo, error) {
	sb := strings.Builder{}
	err := s.urlTemplate.Execute(&sb, TemplateValues{
		GDPR:        gdpr,
		GDPRConsent: consent,
	})
	if err != nil {
		return nil, err
	}

	return &usersync.UsersyncInfo{
		URL:         sb.String(),
		Type:        string(s.syncType),
		SupportCORS: false,
	}, err
}

type TemplateValues struct {
	GDPR        string
	GDPRConsent string
}

func (s *Syncer) FamilyName() string {
	return s.familyName
}

func (s *Syncer) GDPRVendorID() uint16 {
	return s.gdprVendorID
}