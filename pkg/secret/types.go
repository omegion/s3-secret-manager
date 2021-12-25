package secret

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"text/tabwriter"
	"time"
)

const (
	writerPadding = 2
)

// Secret is a struct for secret management.
type Secret struct {
	Bucket       string
	Description  string
	Path         string
	LastModified *time.Time
	Value        map[string]string
	Tags         map[string]string
}

// Secrets is collection of Secret.
type Secrets struct {
	Items []*Secret
}

// EncodeTags encodes tags for S3 API.
func (s Secret) EncodeTags() string {
	encoded := url.Values{}
	for k, v := range s.GetTags() {
		encoded.Set(k, v)
	}

	return encoded.Encode()
}

// GetTags gets all tags both user defined and default.
func (s Secret) GetTags() map[string]string {
	defaultTags := s.defaultTags()
	for k, v := range s.Tags {
		defaultTags[k] = v
	}

	return defaultTags
}

func (s Secret) defaultTags() map[string]string {
	return map[string]string{"SecretPath": s.Path}
}

// EncodedValue encodes value.
func (s Secret) EncodedValue() ([]byte, error) {
	return json.Marshal(s.Value)
}

// GetValue gets value.
func (s Secret) GetValue(key string) (string, error) {
	if v, ok := s.Value[key]; ok {
		return v, nil
	}

	return "", NotFoundError{
		Key:    key,
		Secret: &s,
	}
}

// Print prints Secret details.
func (s Secret) Print() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, writerPadding, ' ', 0)
	fmt.Fprintf(w, "Key\tValue\n")
	fmt.Fprintf(w, "----\t----\n")

	for key, value := range s.Value {
		fmt.Fprintf(w, fmt.Sprintf("%s\t%s\n", key, value))
	}

	err := w.Flush()
	if err != nil {
		return err
	}

	return nil
}

// Print prints Secrets details.
func (s Secrets) Print() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, writerPadding, ' ', 0)
	fmt.Fprintf(w, "Secret\tLast Modified\n")
	fmt.Fprintf(w, "----\t----\n")

	for _, scrt := range s.Items {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\n", scrt.Path, scrt.LastModified))
	}

	err := w.Flush()
	if err != nil {
		return err
	}

	return nil
}
