/*
 * Connector Catalog Service
 *
 * The endpoints of a Connector Catalog service are used to discover connector type APIs that can be used by the Connector Catalog Manager.
 *
 * API version: 0.0.2
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// ConnectorCatalogEntry struct for ConnectorCatalogEntry
type ConnectorCatalogEntry struct {
	Id            string                 `json:"id,omitempty"`
	Channel       string                 `json:"channel,omitempty"`
	ShardMetadata map[string]interface{} `json:"shard_metadata,omitempty"`
}