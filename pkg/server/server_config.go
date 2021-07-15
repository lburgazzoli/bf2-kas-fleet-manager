package server

import (
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/shared"
	"github.com/spf13/pflag"
)

type ServerConfig struct {
	BindAddress   string `json:"bind_address"`
	HTTPSCertFile string `json:"https_cert_file"`
	HTTPSKeyFile  string `json:"https_key_file"`
	EnableHTTPS   bool   `json:"enable_https"`
	EnableJWT     bool   `json:"enable_jwt"`
	JwksURL       string `json:"jwks_url"`
	JwksFile      string `json:"jwks_file"`
	// The public http host URL to access the service
	// For staging it is "https://api.stage.openshift.com"
	// For production it is "https://api.openshift.com"
	PublicHostURL         string `json:"public_url"`
	EnableTermsAcceptance bool   `json:"enable_terms_acceptance"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		BindAddress:   "localhost:8000",
		EnableHTTPS:   false,
		EnableJWT:     true,
		JwksURL:       "https://sso.redhat.com/auth/realms/redhat-external/protocol/openid-connect/certs",
		JwksFile:      "config/jwks-file.json",
		HTTPSCertFile: "",
		HTTPSKeyFile:  "",
		PublicHostURL: "http://localhost",
	}
}

func (s *ServerConfig) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.BindAddress, "api-server-bindaddress", s.BindAddress, "API server bind adddress")
	fs.StringVar(&s.HTTPSCertFile, "https-cert-file", s.HTTPSCertFile, "The path to the tls.crt file.")
	fs.StringVar(&s.HTTPSKeyFile, "https-key-file", s.HTTPSKeyFile, "The path to the tls.key file.")
	fs.BoolVar(&s.EnableHTTPS, "enable-https", s.EnableHTTPS, "Enable HTTPS rather than HTTP")
	fs.BoolVar(&s.EnableJWT, "enable-jwt", s.EnableJWT, "Enable JWT authentication validation")
	fs.BoolVar(&s.EnableTermsAcceptance, "enable-terms-acceptance", s.EnableTermsAcceptance, "Enable terms acceptance check")
	fs.StringVar(&s.JwksURL, "jwks-url", s.JwksURL, "The URL of the JSON web token signing certificates.")
	fs.StringVar(&s.JwksFile, "jwks-file", s.JwksFile, "File containing the the JSON web token signing certificates.")
	fs.StringVar(&s.PublicHostURL, "public-host-url", s.PublicHostURL, "Public http host URL of the service")
}

func (s *ServerConfig) ReadFiles() error {
	s.JwksFile = shared.BuildFullFilePath(s.JwksFile)

	return nil
}