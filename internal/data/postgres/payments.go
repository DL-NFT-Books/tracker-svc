package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/nft-books/contract-tracker/internal/data"
)

const (
	paymentsTable = "payments"

	paymentsId           = "id"
	paymentsContractId   = "contract_id"
	paymentsPayerAddress = "payer_address"
	paymentsTokenAddress = "token_address"
	paymentsAmount       = "amount"
	paymentsPrice        = "price"
)

type paymentsQ struct {
	database *pgdb.DB
	selector squirrel.SelectBuilder
}

func NewPaymentsQ(database *pgdb.DB) data.PaymentsQ {
	return &paymentsQ{
		database: database,
		selector: squirrel.Select(fmt.Sprintf("%s.*", paymentsTable)).From(paymentsTable),
	}
}

func (q *paymentsQ) New() data.PaymentsQ {
	return NewPaymentsQ(q.database.Clone())
}

func (q *paymentsQ) Page(page pgdb.OffsetPageParams) data.PaymentsQ {
	q.selector = page.ApplyTo(q.selector, "id")
	return q
}

func (q *paymentsQ) Sort(sort pgdb.Sorts) data.PaymentsQ {
	q.selector = sort.ApplyTo(q.selector, map[string]string{
		"id": "id",
	})

	return q
}

func (q *paymentsQ) FilterById(id ...int64) data.PaymentsQ {
	q.selector = q.selector.Where(squirrel.Eq{paymentsId: id})
	return q
}

func (q *paymentsQ) FilterByPayer(payer ...string) data.PaymentsQ {
	q.selector = q.selector.Where(squirrel.Eq{paymentsPayerAddress: payer})
	return q
}

func (q *paymentsQ) FilterByTokenAddress(tokenAddress ...string) data.PaymentsQ {
	q.selector = q.selector.Where(squirrel.Eq{paymentsTokenAddress: tokenAddress})
	return q
}

func (q *paymentsQ) FilterByContractId(contractId ...int64) data.PaymentsQ {
	q.selector = q.selector.Where(squirrel.Eq{paymentsContractId: contractId})
	return q
}

func (q *paymentsQ) Select() (payments []data.Payment, err error) {
	err = q.database.Select(&payments, q.selector)
	return
}

func (q *paymentsQ) Get() (*data.Payment, error) {
	var payment data.Payment
	err := q.database.Get(&payment, q.selector)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &payment, err
}

func (q *paymentsQ) Insert(payment data.Payment) (id int64, err error) {
	statement := squirrel.Insert(paymentsTable).
		Suffix("returning id").
		SetMap(structs.Map(&payment))

	err = q.database.Get(&id, statement)
	return id, err
}

func (q *paymentsQ) Delete(id int64) error {
	statement := squirrel.Delete(paymentsTable).Where(squirrel.Eq{paymentsId: id})
	return q.database.Exec(statement)
}
