package imgproxyurl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strings"
)

// Url represents an imgproxy url.
//
// Use url.Set* functions to set various options and finnaly use Url.Get() to get the url.
// If any error occurs, it is returned from Url.Get() method. This is so to make method chaining possible.
// Only the first occurred error is returned.
//
// A typical usage would look like:
// 		url, err := imgproxyurl.NewFromEnvironment("local:///0/D/0DJ2jdB5DDJa.jpg").
//						SetHeight(400).
//						SetWidth(300).
//						SetResizingType(imgproxyurl.ResizingTypeFit)
//						Get()
type Url struct {
	key               []byte
	salt              []byte
	processingOptions map[string]string
	imageUrl          string
	err               error
	plainImageUrl     bool
	extension         string
}

// New creates Url
func New(imageUrl string) *Url {
	return &Url{
		imageUrl:          imageUrl,
		processingOptions: make(map[string]string),
	}
}

// NewFromEnvironment creates Url while trying to get key and salt values
// from environment variables IMGPROXY_KEY and IMGPROXY_SALT
func NewFromEnvironment(imageUrl string) *Url {
	url := New(imageUrl)
	key := os.Getenv("IMGPROXY_KEY")
	salt := os.Getenv("IMGPROXY_SALT")

	if key != "" && salt != "" {
		url.SetKey(key)
		url.SetSalt(salt)
	} else if key != "" && salt == "" {
		url.setError(errors.New("salt is missing"))
	} else if key == "" && salt != "" {
		url.setError(errors.New("key is missing"))
	}

	return url
}

// Get returns an imgproxy url for the given imageUrl
//
// This is a shortnand for New(imageUrl).Get()
func Get(imageUrl string) (string, error) {
	return New(imageUrl).Get()
}

// GetFromEnvironment returns an imgproxy url for the given imageUrl,
// having taken the key and salt from the environment variables IMGPROXY_KEY and IMGPROXY_SALT.
//
// This is a shortnand for NewFromEnvironment(imageUrl).Get()
func GetFromEnvironment(imageUrl string) (string, error) {
	return NewFromEnvironment(imageUrl).Get()
}

func (url *Url) setError(err error) {
	if url.err != nil {
		return
	}

	url.err = err
}

// SetKey set the imgproxy key.
//
// key is expected to be a hex-encoded string
func (url *Url) SetKey(key string) *Url {
	bytes, err := hex.DecodeString(key)
	if err != nil {
		url.setError(errors.WithMessage(err, "key"))
		return url
	}
	url.key = bytes
	return url
}

// SetSalt set the imgproxy salt.
//
// salt is expected to be a hex-encoded string
func (url *Url) SetSalt(salt string) *Url {
	bytes, err := hex.DecodeString(salt)
	if err != nil {
		url.setError(errors.WithMessage(err, "salt"))
		return url
	}
	url.salt = bytes
	return url
}

// SetExtension set the the resulting image format
func (url *Url) SetExtension(extension string) *Url {
	url.extension = extension
	return url
}

func (url *Url) setOption(name string, arguments ...string) {
	url.processingOptions[name] = strings.Join(arguments, ":")
}

func (url *Url) unsetOption(name string) {
	delete(url.processingOptions, name)
}

func (url *Url) encodeImageUrl() string {
	var encodedUrl string
	if url.plainImageUrl {
		encodedUrl = "plain/" + url.imageUrl
		if url.extension != "" {
			encodedUrl += "@" + url.extension
		}
	} else {
		encodedUrl = base64.RawURLEncoding.EncodeToString([]byte(url.imageUrl))
		if url.extension != "" {
			encodedUrl += "." + url.extension
		}
	}
	return encodedUrl
}

func (url *Url) getPath() string {
	var urlParts []string
	for name, arguments := range url.processingOptions {
		urlParts = append(urlParts, name+":"+arguments)
	}
	urlParts = append(urlParts, url.encodeImageUrl())
	return "/" + strings.Join(urlParts, "/")
}

func (url *Url) sign(str string) string {
	mac := hmac.New(sha256.New, url.key)
	mac.Write(url.salt)
	mac.Write([]byte(str))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// Get returns the resulting imgproxy url.
//
// If a key and a salt were provided, the url will be signed.
// If any errors occured during the previous Set* functions calls, the first of them will be returned.
func (url *Url) Get() (string, error) {
	if url.imageUrl == "" {
		url.setError(errors.New("image url is missing"))
	}

	if url.err != nil {
		return "", url.err
	}

	path := url.getPath()

	var signature string
	if url.key != nil && url.salt != nil {
		signature = url.sign(path)
	} else {
		signature = "insecure"
	}

	return fmt.Sprintf("/%s%s", signature, path), nil
}

//GetAbsolute returns the resulting imgproxy absolute url by prepending suffix in front.
func (url *Url) GetAbsolute(prefix string) (string, error) {
	result, err := url.Get()
	if err != nil {
		return result, err
	}

	return prefix + result, nil
}
