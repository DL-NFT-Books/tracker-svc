package data

type BookQ interface {
	New() BookQ
	Get() (*Book, error)
	Select() ([]Book, error)
	FilterActual() BookQ
	FilterByID(id ...int64) BookQ
	FilterByContractAddress(address ...string) BookQ
}

type Book struct {
	ID              int64  `db:"id" structs:"-"`
	Title           string `db:"title" structs:"title"`
	Description     string `db:"description" structs:"description"`
	Price           string `db:"price" structs:"price"`
	ContractAddress string `db:"contract_address" structs:"contract_address"`
	ContractName    string `db:"contract_name" structs:"contract_name"`
	ContractVersion string `db:"contract_version" structs:"contract_version"`
	Banner          string `db:"banner" structs:"banner"`
	File            string `db:"file" structs:"file"`
	Deleted         bool   `db:"deleted" structs:"-"`
}
