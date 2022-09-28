/*
 * Connector Service Fleet Manager Private APIs
 *
 * Connector Service Fleet Manager apis that are used by internal services.
 *
 * API version: 0.0.3
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package private

// KnativeConnectionSettings Holds the configuration to connect to a Knative service.
type KnativeConnectionSettings struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}
