/* 
 * Vendor API V1
 *
 * Create, list, promote, update and archive releases.
 *
 * OpenAPI spec version: 1.0.0
 * 
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 */

package swagger

// LicenseCounts is a struct to hold license count information
type LicenseCounts struct {

	Active map[string]int64 `json:"active,omitempty"`

	Airgap map[string]int64 `json:"airgap,omitempty"`

	Inactive map[string]int64 `json:"inactive,omitempty"`

	Total map[string]int64 `json:"total,omitempty"`
}
