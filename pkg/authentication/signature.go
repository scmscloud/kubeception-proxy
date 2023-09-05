package authentication

import (
	"encoding/hex"
	"io"
	"log"

	"github.com/things-go/go-socks5"
	"github.com/things-go/go-socks5/statute"
	"scaleship.io/kubernetes/kubeception-proxy/pkg/encryption"
)

type Signature struct{}

func (a Signature) Authenticate(r io.Reader, w io.Writer, usraddr string) (*socks5.AuthContext, error) {

	if _, err := w.Write([]byte{statute.VersionSocks5, statute.MethodUserPassAuth}); err != nil {
		return nil, err
	}

	usr, err := statute.ParseUserPassRequest(r)
	if err != nil {
		return nil, err
	}

	gzp := string(usr.User) + string(usr.Pass)
	sign, err := encryption.Decompress(&gzp)
	if err != nil {
		log.Printf("[WARNING] Unable to decompress signatures from %q in SOCKS5 authentication process.", usraddr)
		if _, err := w.Write([]byte{statute.UserPassAuthVersion, statute.AuthFailure}); err != nil {
			return nil, err
		}
		return nil, statute.ErrUserAuthFailed
	}

	signature := hex.EncodeToString(*sign)

	if _, err := w.Write([]byte{statute.UserPassAuthVersion, statute.AuthSuccess}); err != nil {
		return nil, err
	}

	return &socks5.AuthContext{
		Method: statute.MethodUserPassAuth,
		Payload: map[string]string{
			"signature": signature,
		},
	}, nil
}

func (a Signature) GetCode() uint8 {
	return statute.MethodUserPassAuth
}
