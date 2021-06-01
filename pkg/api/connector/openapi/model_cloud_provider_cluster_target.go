/*
 * Connector Service Fleet Manager
 *
 * Connector Service Fleet Manager is a Rest API to manage connectors.
 *
 * API version: 0.0.2
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// CloudProviderClusterTarget Targets workloads to a cloud provider
type CloudProviderClusterTarget struct {
	Kind          string `json:"kind"`
	CloudProvider string `json:"cloud_provider,omitempty"`
	Region        string `json:"region,omitempty"`
	MultiAz       bool   `json:"multi_az,omitempty"`
}