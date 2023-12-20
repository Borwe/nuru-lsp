package server

import (
	"encoding/json"

	"github.com/Borwe/go-lsp/jsonrpc"
	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp"
	"github.com/Borwe/go-lsp/lsp/defines"
)

var Server = lsp.NewServer(&lsp.Options{
	CompletionProvider: &defines.CompletionOptions{
		TriggerCharacters: &[]string{"."},
	},
})

// Send notifications to the client
func Notify(s *lsp.Server, method string, result interface{}) {
	params, err := json.Marshal(result)
	logs.Printf("Notifying with: %s\n",params)
	msg := jsonrpc.NotificationMessage{
		Method: method,
		BaseMessage: jsonrpc.BaseMessage{ Jsonrpc: "2.0"},
		Params: params,
	}
	session := s.RpcServer.Session[0]
	err = session.Notify(msg)
	if err!=nil {
		logs.Printf("Error occured onNotify: %s\n",err)
	}else{
		logs.Printf("Notified %s", msg)
	}
}
