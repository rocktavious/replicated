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

type LicenseVersions struct {
	InstalledAppVersion *InstalledAppVersion `json:"InstalledAppVersion,omitempty"`
	ReplicatedVersions  map[string][]string  `json:"ReplicatedVersions,omitempty"`
}
