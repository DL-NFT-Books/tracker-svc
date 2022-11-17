package external

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/nft-books/contract-tracker/resources"
)

type Task struct {
	Id               int64                `db:"id" structs:"-" json:"-"`
	BookId           int64                `db:"book_id" structs:"book_id" json:"book_id"`
	TokenId          int64                `db:"token_id" structs:"token_id"`
	Account          string               `db:"account" structs:"account"`
	Signature        string               `db:"signature" structs:"signature"`
	FileIpfsHash     string               `db:"file_ipfs_hash" structs:"file_ipfs_hash"`
	MetadataIpfsHash string               `db:"metadata_ipfs_hash" structs:"metadata_ipfs_hash"`
	Uri              string               `db:"uri" structs:"uri"`
	Status           resources.TaskStatus `db:"status" structs:"status"`
}

// TaskSelector is a structure for all applicable filters and params on tasksQ `Select`
type TaskSelector struct {
	PageParams pgdb.CursorPageParams
	Account    *string
	IpfsHash   *string
	Status     *resources.TaskStatus
}

type TasksQ interface {
	New() TasksQ

	GetById(id int64) (*Task, error)
	GetByHash(ipfsHash string) (*Task, error)
	Select(selector TaskSelector) ([]Task, error)

	Sort(sort pgdb.Sorts) TasksQ
	Page(page pgdb.OffsetPageParams) TasksQ

	Update(task Task, id int64) error
	Insert(task Task) (id int64, err error)
	Delete(id int64) error
	Transaction(fn func(q TasksQ) error) error

	UpdateFileIpfsHash(newIpfsHash string, id int64) error
	UpdateMetadataIpfsHash(newIpfsHash string, id int64) error
	UpdateTokenId(newTokenId, id int64) error
	UpdateStatus(newStatus resources.TaskStatus, id int64) error
}
