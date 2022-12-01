package contract_reader

import "gitlab.com/distributed_lab/logan/v3/errors"

var (
	FromNotSpecifiedErr    = errors.New("starting block is not specified")
	AddressNotSpecifiedErr = errors.New("address is not specified")
)
