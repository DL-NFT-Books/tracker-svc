/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UpdateBookAttributes struct {
	Banner *Media `json:"banner,omitempty"`
	// Address of a contract corresponding to this book
	ContractAddress *string `json:"contract_address,omitempty"`
	// Contract name (in most cases coincides with a title field)
	ContractName *string `json:"contract_name,omitempty"`
	// status of a book deployment
	DeployStatus *DeployStatus `json:"deploy_status,omitempty"`
	// Book description
	Description *string `json:"description,omitempty"`
	File        *Media  `json:"file,omitempty"`
	// Price per one token ($)
	Price *string `json:"price,omitempty"`
	// Book title
	Title *string `json:"title,omitempty"`
	// Token symbol
	TokenSymbol *string `json:"token_symbol,omitempty"`
}
