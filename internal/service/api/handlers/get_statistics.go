package handlers

import (
	"github.com/dl-nft-books/tracker-svc/internal/service/api/requests"
	"github.com/dl-nft-books/tracker-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/kit/pgdb"
	"net/http"
	"strconv"
)

func GetStatistics(w http.ResponseWriter, r *http.Request) {
	var response resources.StatisticsResponse
	response.Data.Type = resources.STATISTICS
	request, err := requests.NewGetStatisticsRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// Amount pie chart
	bookStats, err := DB(r).Statistics().BookStatisticsQ.New().Select()
	for _, book := range bookStats {
		if book.BookId == 0 {
			response.Data.Attributes.AmountPieChart.Attributes.Total = book.Amount
			response.Data.Attributes.PricePieChart.Attributes.Total = book.UsdPrice
			continue
		}
		response.Data.Attributes.AmountPieChart.Attributes.Books =
			append(response.Data.Attributes.AmountPieChart.Attributes.Books, resources.AmountPieChartBooks{
				Key: resources.Key{},
				Attributes: resources.AmountPieChartBooksAttributes{
					BookId: book.BookId,
					Amount: 0,
				},
			})
	}

	// Price pie chart and
	// Token list
	tokenStats, err := DB(r).Statistics().TokenStatisticsQ.New().
		FilterByBookId(0).
		Page(pgdb.OffsetPageParams{
			Limit:      100,
			Order:      "desc",
			PageNumber: 0,
		}, "usd_price").Select()
	for i, token := range tokenStats {
		if i <= int(request.Limit) {
			response.Data.Attributes.TokenList =
				append(response.Data.Attributes.TokenList, resources.PricePieChartTokens{
					Key: resources.Key{},
					Attributes: resources.PricePieChartTokensAttributes{
						Name:           token.TokenSymbol,
						NativeCurrency: token.TokenPrice,
						Usd:            token.UsdPrice,
					},
				})
		}
		response.Data.Attributes.PricePieChart.Attributes.Tokens =
			append(response.Data.Attributes.PricePieChart.Attributes.Tokens, resources.PricePieChartTokens{
				Key: resources.Key{},
				Attributes: resources.PricePieChartTokensAttributes{
					Name:           token.TokenSymbol,
					NativeCurrency: token.TokenPrice,
					Usd:            token.UsdPrice,
				},
			})
	}

	// Chain pie chart
	chainStats, err := DB(r).Statistics().ChainStatisticsQ.New().FilterByBookId(0).Select()
	for _, chain := range chainStats {
		response.Data.Attributes.ChainPieChart.Attributes.Chains =
			append(response.Data.Attributes.ChainPieChart.Attributes.Chains, resources.ChainPieChartChains{
				Key: resources.Key{},
				Attributes: resources.ChainPieChartChainsAttributes{
					Amount:  chain.Amount,
					ChainId: chain.ChainId,
				},
			})
	}

	nftPayments, err := DB(r).Payments().Page(pgdb.OffsetPageParams{
		Limit:      request.Limit,
		Order:      "desc",
		PageNumber: 0,
	}, "price_token").FilterByType(int8(resources.NFT)).Select()
	for _, nftPayment := range nftPayments {
		response.Data.Attributes.NftList = append(response.Data.Attributes.NftList, resources.NftListItem{
			Key: resources.Key{
				ID:   strconv.Itoa(int(nftPayment.NftId)),
				Type: resources.NFT_LIST,
			},
			Attributes: resources.NftListItemAttributes{
				Address: nftPayment.TokenAddress,
			},
		})
	}
	ape.Render(w, response)
}
