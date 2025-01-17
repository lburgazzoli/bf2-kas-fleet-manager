/*
 * Connector Service Fleet Manager Private APIs
 *
 * Connector Service Fleet Manager apis that are used by internal services.
 *
 * API version: 0.0.3
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package private

// ConnectorDeploymentWatchEvent struct for ConnectorDeploymentWatchEvent
type ConnectorDeploymentWatchEvent struct {
	Type   string              `json:"type"`
	Error  Error               `json:"error,omitempty"`
	Object ConnectorDeployment `json:"object,omitempty"`
}
