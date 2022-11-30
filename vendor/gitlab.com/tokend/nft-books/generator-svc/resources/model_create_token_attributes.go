/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateTokenAttributes struct {
	// Address of a user who purchased this token
	Account string `json:"account"`
	ChainId int64  `json:"chain_id"`
	// Hash of a metadata file
	MetadataHash string      `json:"metadata_hash"`
	Status       TokenStatus `json:"status"`
	TokenId      int32       `json:"token_id"`
}
