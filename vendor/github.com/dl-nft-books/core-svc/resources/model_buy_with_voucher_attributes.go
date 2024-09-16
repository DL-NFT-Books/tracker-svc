/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "math/big"

type BuyWithVoucherAttributes struct {
	// end timestamp of signature
	EndSigTimestamp *big.Int  `json:"end_sig_timestamp,omitempty"`
	PermitSig       Signature `json:"permit_sig"`
	// id of task
	TaskId int64 `json:"task_id"`
	// Voucher address
	VoucherAddress string `json:"voucher_address"`
}
