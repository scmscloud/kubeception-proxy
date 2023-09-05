package encryption

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"io/ioutil"
)

func Compress(raw *[]byte) (*string, error) {

	b := &bytes.Buffer{}

	w, err := flate.NewWriter(b, flate.BestCompression)
	if err != nil {
		return nil, err
	}

	if _, err := w.Write(*raw); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	s := base64.RawURLEncoding.EncodeToString(b.Bytes())
	return &s, nil
}

func Decompress(raw *string) (*[]byte, error) {

	z, err := base64.RawURLEncoding.DecodeString(*raw)
	if err != nil {
		return nil, err
	}

	flr := flate.NewReader(bytes.NewReader(z))

	r, err := ioutil.ReadAll(flr)
	if err != nil {
		return nil, err
	}

	if err := flr.Close(); err != nil {
		return nil, err
	}

	return &r, nil
}
