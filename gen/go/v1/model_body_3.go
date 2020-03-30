/*
 * Vendor API V1
 *
 * Apps documentation
 *
 * API version: 1.0.0
 * Contact: info@replicated.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type BodyCreateChannel struct {
	AirgapEnabled bool   `json:"airgap_enabled,omitempty"`
	Description   string `json:"description,omitempty"`
	Name          string `json:"name"`
}
