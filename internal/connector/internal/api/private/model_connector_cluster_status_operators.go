/*
 * Connector Service Fleet Manager Private APIs
 *
 * Connector Service Fleet Manager apis that are used by internal services.
 *
 * API version: 0.0.3
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package private

// ConnectorClusterStatusOperators struct for ConnectorClusterStatusOperators
type ConnectorClusterStatusOperators struct {
	Operator ConnectorOperator `json:"operator,omitempty"`
	// the namespace to which the operator has been installed
	Namespace string `json:"namespace,omitempty"`
	// the status of the operator
	Status string `json:"status,omitempty"`
}