package server

import (
	"encoding/json"

	"github.com/Borwe/go-lsp/jsonrpc"
	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp"
	"github.com/Borwe/go-lsp/lsp/defines"
)

func createNewServerOpts() *lsp.Options {
	providesValue := true
	return &lsp.Options{
		CompletionProvider: &defines.CompletionOptions{
			TriggerCharacters: &[]string{"."},
			WorkDoneProgressOptions: defines.WorkDoneProgressOptions{
				WorkDoneProgress: &providesValue,
			},
		},
		HoverProvider: &defines.HoverOptions{
			WorkDoneProgressOptions: defines.WorkDoneProgressOptions{
				WorkDoneProgress: &providesValue,
			},
		},
	}
}

var Server = lsp.NewServer(createNewServerOpts())

// Send notifications to the client
func Notify(s *lsp.Server, method string, result interface{}) {
	params, err := json.Marshal(result)
	msg := jsonrpc.NotificationMessage{
		Method:      method,
		BaseMessage: jsonrpc.BaseMessage{Jsonrpc: "2.0"},
		Params:      params,
	}
	session := s.RpcServer.Session[0]
	err = session.Notify(msg)
	if err != nil {
		logs.Printf("Error occured onNotify: %s\n", err)
	}
}
