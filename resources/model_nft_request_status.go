package resources

type NftRequestStatus int64

const (
	NftRequestNone NftRequestStatus = iota + 1
	NftRequestPending
	NftRequestCanceled
	NftRequestAccepted
)

var nftRequestStates = map[NftRequestStatus]string{
	NftRequestNone:     "none",
	NftRequestPending:  "pending",
	NftRequestCanceled: "canceled",
	NftRequestAccepted: "accepted",
}

func (s NftRequestStatus) String() string {
	return nftRequestStates[s]
}
