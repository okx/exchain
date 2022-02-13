package client

//import (
//	govclient "github.com/okex/exchain/libs/cosmos-sdk/x/gov/client"
//	"net/http"
//)
//
//var (
//	UpdateClientProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpdateClientProposal, emptyRestHandler)
//	UpgradeProposalHandler      = govclient.NewProposalHandler(cli.NewCmdSubmitUpgradeProposal, emptyRestHandler)
//)
//
//func emptyRestHandler(client.Context) govrest.ProposalRESTHandler {
//	return govrest.ProposalRESTHandler{
//		SubRoute: "unsupported-ibc-client",
//		Handler: func(w http.ResponseWriter, r *http.Request) {
//			rest.WriteErrorResponse(w, http.StatusBadRequest, "Legacy REST Routes are not supported for IBC proposals")
//		},
//	}
//}
