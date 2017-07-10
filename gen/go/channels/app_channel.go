/* 
 * Vendor API V1
 *
 * List, create, update, delete and archive channels.
 *
 * OpenAPI spec version: 1.0.0
 * 
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 */

package swagger

// An app channel belongs to an app. It contains references to the top (current) release in the channel.
type AppChannel struct {

	Adoption ChannelAdoption `json:"Adoption,omitempty"`

	// Description that will be shown during license installation
	Description string `json:"Description"`

	// The ID of the channel
	Id string `json:"Id"`

	LicenseCounts LicenseCounts `json:"LicenseCounts,omitempty"`

	// The name of channel
	Name string `json:"Name"`

	// The position for which the channel occurs in a list
	Position int64 `json:"Position,omitempty"`

	// The label of the current release sequence
	ReleaseLabel string `json:"ReleaseLabel,omitempty"`

	// Release notes for the current release sequence
	ReleaseNotes string `json:"ReleaseNotes,omitempty"`

	// A reference to the current release sequence
	ReleaseSequence int64 `json:"ReleaseSequence,omitempty"`
}
