/*
 * Kafka Service Fleet Manager
 *
 * Kafka Service Fleet Manager APIs that are used by internal services e.g kas-fleetshard operators.
 *
 * API version: 1.5.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package private

// DatePlaneClusterUpdateStatusRequestResizeInfo struct for DatePlaneClusterUpdateStatusRequestResizeInfo
type DatePlaneClusterUpdateStatusRequestResizeInfo struct {
	NodeDelta *int32                                              `json:"nodeDelta,omitempty"`
	Delta     *DatePlaneClusterUpdateStatusRequestResizeInfoDelta `json:"delta,omitempty"`
}
