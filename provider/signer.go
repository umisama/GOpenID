package provider

import (
	"crypto/rand"
	"errors"
	"github.com/GehirnInc/GOpenID"
	"strings"
	"time"
)

var (
	ErrAlreadySigned         = errors.New("response has been signed")
	ErrNotNeedsSigning       = errors.New("response does not need signing")
	ErrIdentityNotSet        = errors.New("identity not set")
	ErrIdentitySet           = errors.New("identity set")
	ErrIdentityNotMatched    = errors.New("identity not matched")
	ErrMessageNotSigned      = errors.New("message is not signed")
	ErrVerifyingNotSupported = errors.New("verifying not supported")
)

type Signer struct {
	store    gopenid.Store
	lifetime time.Duration
}

func NewSigner(store gopenid.Store, lifetime time.Duration) *Signer {
	return &Signer{
		store:    store,
		lifetime: lifetime,
	}
}

func (s *Signer) Invalidate(handle string, isStateless bool) {
	assoc, ok := s.store.GetAssociation(handle, isStateless)
	if ok {
		s.store.DeleteAssociation(assoc)
	}

	return
}

func (s *Signer) Verify(req Request, isStateless bool) (ok bool, err error) {
	var (
		assocHandle gopenid.MessageValue
		signed      gopenid.MessageValue
		sig         gopenid.MessageValue
	)

	switch ret := req.(type) {
	case *CheckAuthenticationRequest:
		assocHandle = ret.assocHandle
		signed = ret.signed
		sig = ret.sig
	default:
		err = ErrVerifyingNotSupported
		return
	}

	assoc, ok := s.store.GetAssociation(assocHandle.String(), isStateless)
	if !ok {
		err = gopenid.ErrAssociationNotFound
		return
	}

	// signing
	msg := req.GetMessage()
	verify := msg.Copy()
	if err = assoc.Sign(verify, strings.Split(signed.String(), ",")); err != nil {
		return
	}

	expected, _ := verify.GetArg(
		gopenid.NewMessageKey(verify.GetOpenIDNamespace(), "sig"),
	)
	ok = sig == expected

	return
}

func (s *Signer) Sign(res *OpenIDResponse, assocHandle string, order []string) (err error) {
	var assoc *gopenid.Association

	if assocHandle == "" {
		assoc, err = gopenid.CreateAssociation(
			rand.Reader,
			gopenid.DefaultAssoc,
			s.getExpires(),
			true,
		)
	} else {
		var ok bool

		assoc, ok = s.store.GetAssociation(assocHandle, false)
		if !ok || !assoc.IsValid() {
			res.AddArg(
				gopenid.NewMessageKey(res.GetNamespace(), "invalidate_handle"),
				gopenid.MessageValue(assocHandle),
			)

			assoc, err = gopenid.CreateAssociation(
				rand.Reader,
				gopenid.DefaultAssoc,
				s.getExpires(),
				true,
			)
		}
	}

	if err != nil {
		return
	}

	if assoc.IsStateless() {
		s.store.StoreAssociation(assoc)
	} else {
		s.store.DeleteAssociation(assoc)
	}

	return assoc.Sign(res.message, order)
}

func (s *Signer) getExpires() time.Time {
	return time.Now().Add(s.lifetime)
}
