/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateNftRequest struct {
	Key
	Attributes CreateNftRequestAttributes `json:"attributes"`
}
type CreateNftRequestRequest struct {
	Data     CreateNftRequest `json:"data"`
	Included Included         `json:"included"`
}

type CreateNftRequestListRequest struct {
	Data     []CreateNftRequest `json:"data"`
	Included Included           `json:"included"`
	Links    *Links             `json:"links"`
}

// MustCreateNftRequest - returns CreateNftRequest from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateNftRequest(key Key) *CreateNftRequest {
	var createNftRequest CreateNftRequest
	if c.tryFindEntry(key, &createNftRequest) {
		return &createNftRequest
	}
	return nil
}
