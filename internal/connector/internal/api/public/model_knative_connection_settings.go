/*
 * Connector Management API
 *
 * Connector Management API is a REST API to manage connectors.
 *
 * API version: 0.1.0
 * Contact: rhosak-support@redhat.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package public

// KnativeConnectionSettings Holds the configuration to connect to a Knative service.
type KnativeConnectionSettings struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}
