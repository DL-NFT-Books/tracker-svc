package postgres

import (
	"database/sql"
	"fmt"
	"github.com/dl-nft-books/tracker-svc/internal/data"
	"github.com/fatih/structs"

	"github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	blocksTable = "blocks"

	blocksIdColumn                = "id"
	blocksChainIdColumn           = "chain_id"
	blocksContractAddressColumn   = "contract_address"
	tokenPurchasedBlockColumn     = "token_purchased_block"
	tokenExchangedBlockColumn     = "token_exchanged_block"
	nftRequestCreatedBlockColumn  = "nft_request_created_block"
	nftRequestCanceledBlockColumn = "nft_request_canceled_block"
)

type blocksQ struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
}

func NewBlocksQ(db *pgdb.DB) data.BlocksQ {
	return &blocksQ{
		db:       db,
		selector: squirrel.Select(fmt.Sprintf("%s.*", blocksTable)).From(blocksTable),
	}
}

func (q *blocksQ) New() data.BlocksQ {
	return NewBlocksQ(q.db.Clone())
}

func (q *blocksQ) FilterById(id ...int64) data.BlocksQ {
	q.selector = q.selector.Where(squirrel.Eq{blocksIdColumn: id})
	return q
}

func (q *blocksQ) FilterByChainId(id ...int64) data.BlocksQ {
	q.selector = q.selector.Where(squirrel.Eq{blocksChainIdColumn: id})
	return q
}

func (q *blocksQ) FilterByContractAddress(contractAddress ...string) data.BlocksQ {
	q.selector = q.selector.Where(squirrel.Eq{blocksContractAddressColumn: contractAddress})
	return q
}

func (q *blocksQ) Insert(blocks data.Blocks) (id int64, err error) {
	statement := squirrel.Insert(blocksTable).
		Suffix("returning id").
		SetMap(structs.Map(&blocks))

	err = q.db.Get(&id, statement)
	return id, err
}

func (q *blocksQ) UpdateTokenPurchasedBlockColumn(newMintBlock uint64, chainId int64) error {
	statement := squirrel.Update(blocksTable).
		Set(tokenPurchasedBlockColumn, newMintBlock).
		Where(squirrel.Eq{blocksChainIdColumn: chainId})
	return q.db.Exec(statement)
}

func (q *blocksQ) UpdateTokenExchangedBlockColumn(newExchangedBlock uint64, id int64) error {
	statement := squirrel.Update(blocksTable).
		Set(tokenExchangedBlockColumn, newExchangedBlock).
		Where(squirrel.Eq{blocksIdColumn: id})
	return q.db.Exec(statement)
}

func (q *blocksQ) UpdateNFTRequestCreatedBlockColumn(newNftRequestCreatedBlock uint64, id int64) error {
	statement := squirrel.Update(blocksTable).
		Set(nftRequestCreatedBlockColumn, newNftRequestCreatedBlock).
		Where(squirrel.Eq{blocksIdColumn: id})
	return q.db.Exec(statement)
}

func (q *blocksQ) UpdateNFTRequestCanceledBlockColumn(newNftRequestCanceledBlock uint64, id int64) error {
	statement := squirrel.Update(blocksTable).
		Set(nftRequestCanceledBlockColumn, newNftRequestCanceledBlock).
		Where(squirrel.Eq{blocksIdColumn: id})
	return q.db.Exec(statement)
}

func (q *blocksQ) Get() (*data.Blocks, error) {
	var blocks data.Blocks

	err := q.db.Get(&blocks, q.selector)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &blocks, err
}

func (q *blocksQ) Select() (blocks []data.Blocks, err error) {
	err = q.db.Select(&blocks, q.selector)
	return
}

func (q *blocksQ) Delete(id int64) error {
	statement := squirrel.Delete(blocksTable).Where(squirrel.Eq{blocksIdColumn: id})
	return q.db.Exec(statement)
}
