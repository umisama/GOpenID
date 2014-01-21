package main

import (
	"github.com/GehirnInc/GOpenID"
	"github.com/GehirnInc/GOpenID/provider"
	"log"
	"net/http"
	"net/url"
)

const (
	ASSOCIATION_LIFETIME = 60 * 60 * 12 // 12 hours
)

type FileStore struct {
}

func (s *FileStore) StoreAssociation(assoc *gopenid.Association) error {
	return nil
}

func (s *FileStore) GetAssociation(assocHandle string, isStateless bool) (*gopenid.Association, error) {
	return nil, nil
}

func (s *FileStore) DeleteAssociation(assoc *gopenid.Association) error {
	return nil
}

func (s *FileStore) IsKnownNonce(nonce string) (bool, error) {
	return false, nil
}

func (s *FileStore) StoreNonce(nonce string) error {
	return nil
}

func main() {
	store := &FileStore{}
	p := provider.NewProvider("http://localhost:6543/", store, ASSOCIATION_LIFETIME)

	http.HandleFunc("/openid", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		msg, err := gopenid.MessageFromQuery(r.Form)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		session, err := p.EstablishSession(msg)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		switch ret := session.(type) {
		case *provider.CheckIDSession:
			ret.Accept("", "")
			res, err := ret.GetResponse()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			resMsg := res.GetMessage()

			returnTo, ok := msg.GetArg(gopenid.NewMessageKey(msg.GetOpenIDNamespace(), "return_to"))
			if ok {
				u, _ := url.Parse(returnTo.String())
				u.RawQuery = resMsg.ToQuery().Encode()
				http.Redirect(w, r, u.String(), 302)
				return
			}
		case *provider.AssociateSession:
			res, err := ret.GetResponse()
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			resMsg := res.GetMessage()
			kv, err := resMsg.ToKeyValue([]string{"openid.ns", "openid.assoc_handle", "openid.session_type", "openid.assoc_type", "openid.expires_in", "openid.mac_key"})
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Write(kv)
			return
		case *provider.CheckAuthenticationSession:
		}
	})

	err := http.ListenAndServe(":6543", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
