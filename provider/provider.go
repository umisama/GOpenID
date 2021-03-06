package provider

import (
	"github.com/GehirnInc/GOpenID"
	"time"
)

type Provider struct {
	store         gopenid.Store
	signer        *Signer
	endpoint      string
	assocLifetime time.Duration
}

func NewProvider(endpoint string, store gopenid.Store, lifetime time.Duration) *Provider {
	signer := NewSigner(store, lifetime)

	return &Provider{
		store:         store,
		signer:        signer,
		endpoint:      endpoint,
		assocLifetime: lifetime,
	}
}

func (p *Provider) EstablishSession(msg gopenid.Message) (Session, error) {
	return SessionFromMessage(p, msg)
}

func (p *Provider) getAssocExpires() time.Time {
	return time.Now().Add(p.assocLifetime)
}

func (p *Provider) GetYadisProviderIdentifier() Response {
	et := &gopenid.XRDSDocument{
		XRD: gopenid.XRDSXRDElement{
			Services: []gopenid.XRDSServiceElement{
				gopenid.XRDSServiceElement{
					Priority: 1,
					Type: []string{
						gopenid.NsOpenID20Server.String(),
					},
					URI: p.endpoint,
				},
			},
		},
	}

	return NewYadisResponse(et)
}

func (p *Provider) GetYadisClaimedIdentifier(localid string) Response {
	et := &gopenid.XRDSDocument{
		XRD: gopenid.XRDSXRDElement{
			Services: []gopenid.XRDSServiceElement{
				gopenid.XRDSServiceElement{
					Priority: 1,
					Type: []string{
						gopenid.NsOpenID20Signon.String(),
					},
					URI:     p.endpoint,
					LocalID: localid,
				},
			},
		},
	}

	return NewYadisResponse(et)
}
