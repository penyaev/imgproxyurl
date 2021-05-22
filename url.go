package imgproxyurl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"sort"
	"strings"
)

var std = &Url{
	options: make(map[string]string),
}

type Url struct {
	key            []byte
	salt           []byte
	options        map[string]string
	sourceUrl      string
	plainSourceUrl bool
	format         string
	endpoint       string
	signatureSize  int
}

func New(sourceUrl string, options ...Option) (*Url, error) {
	result, err := std.WithOptions(options...)
	if err != nil {
		return nil, err
	}
	result.sourceUrl = sourceUrl
	return result, nil
}

func (u *Url) WithOptions(options ...Option) (*Url, error) {
	return u.clone(options)
}

func (u *Url) String() string {
	p := u.getPath()

	var signature string
	if u.key == nil || u.salt == nil {
		signature = "insecure"
	} else {
		signature = u.sign(p)
	}

	var result string
	signedPath := fmt.Sprintf("/%s%s", signature, p)

	if u.endpoint != "" {
		end := len(u.endpoint)
		if u.endpoint[end-1] == '/' {
			end--
		}
		result = u.endpoint[:end] + signedPath
	} else {
		result = signedPath
	}

	return result
}

func (u *Url) sign(str string) string {
	mac := hmac.New(sha256.New, u.key)
	mac.Write(u.salt)
	mac.Write([]byte(str))

	size := u.signatureSize
	if size == 0 {
		size = 32
	}
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)[:size])
}

func (u *Url) getPath() string {
	var urlParts []string
	for name, option := range u.options {
		s := name
		if option != "" {
			s += ":" + option
		}
		urlParts = append(urlParts, s)
	}
	// sort url parts to make sure the resulting url is stable (to make unit test work)
	sort.Slice(urlParts, func(i, j int) bool {
		return urlParts[i] < urlParts[j]
	})
	urlParts = append(urlParts, u.encodeSourceUrl())
	return "/" + strings.Join(urlParts, "/")
}

func (u *Url) encodeSourceUrl() string {
	var encodedUrl string
	if u.plainSourceUrl {
		encodedUrl = "plain/" + url.QueryEscape(u.sourceUrl)
		if u.format != "" {
			encodedUrl += "@" + u.format
		}
	} else {
		encodedUrl = base64.RawURLEncoding.EncodeToString([]byte(u.sourceUrl))
		if u.format != "" {
			encodedUrl += "." + u.format
		}
	}
	return encodedUrl
}

func (u *Url) applyOptions(options ...Option) error {
	for _, option := range options {
		switch option.(type) {
		case ProcessingOption:
			u.options[option.(ProcessingOption).Key()] = option.(ProcessingOption).String()
		case Format:
			u.format = option.(Format).Format
		case SourceUrl:
			u.sourceUrl = option.(SourceUrl).Url
		case PlainSourceUrl:
			u.plainSourceUrl = option.(PlainSourceUrl).Plain
		case Key:
			key := option.(Key).Key
			bytes, err := hex.DecodeString(key)
			if err != nil {
				return errors.WithMessage(err, "hexdecode")
			}

			u.key = bytes
		case Salt:
			salt := option.(Salt).Salt
			bytes, err := hex.DecodeString(salt)
			if err != nil {
				return errors.WithMessage(err, "hexdecode")
			}

			u.salt = bytes
		case KeyRaw:
			u.key = option.(KeyRaw).KeyRaw
		case SaltRaw:
			u.salt = option.(SaltRaw).SaltRaw
		case Endpoint:
			u.endpoint = option.(Endpoint).Endpoint
		case SignatureSize:
			u.signatureSize = option.(SignatureSize).SignatureSize
		}
	}

	return nil
}

func (u *Url) clone(addOptions []Option) (*Url, error) {
	clone := &Url{
		key:            u.key,
		salt:           u.salt,
		options:        make(map[string]string, len(u.options)),
		sourceUrl:      u.sourceUrl,
		plainSourceUrl: u.plainSourceUrl,
		format:         u.format,
		endpoint:       u.endpoint,
		signatureSize:  u.signatureSize,
	}
	for key, value := range u.options {
		clone.options[key] = value
	}
	err := clone.applyOptions(addOptions...)
	if err != nil {
		return nil, err
	}

	return clone, nil
}

func SetKeySalt(key string, salt string) error {
	return std.applyOptions(Key{key}, Salt{salt})
}

func SetKeySaltRaw(key []byte, salt []byte) error {
	return std.applyOptions(KeyRaw{key}, SaltRaw{salt})
}

func SetEndpoint(endpoint string) {
	_ = std.applyOptions(Endpoint{endpoint})
}
