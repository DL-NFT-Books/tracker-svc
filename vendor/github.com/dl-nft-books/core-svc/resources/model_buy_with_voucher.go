/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type BuyWithVoucher struct {
	Key
	Attributes BuyWithVoucherAttributes `json:"attributes"`
}
type BuyWithVoucherRequest struct {
	Data     BuyWithVoucher `json:"data"`
	Included Included       `json:"included"`
}

type BuyWithVoucherListRequest struct {
	Data     []BuyWithVoucher `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
}

// MustBuyWithVoucher - returns BuyWithVoucher from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustBuyWithVoucher(key Key) *BuyWithVoucher {
	var buyWithVoucher BuyWithVoucher
	if c.tryFindEntry(key, &buyWithVoucher) {
		return &buyWithVoucher
	}
	return nil
}
