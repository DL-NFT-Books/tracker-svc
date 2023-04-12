/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateTaskAttributes struct {
	// payer account
	Account string `json:"account"`
	Banner  Media  `json:"banner"`
	// id of book
	BookId int64 `json:"book_id"`
	// id of network chain
	ChainId int64 `json:"chain_id"`
	// name of future token
	TokenName string `json:"token_name"`
}