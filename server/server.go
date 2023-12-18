package server

import (
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
func Notify(s *lsp.Server, msg jsonrpc.ResponseMessage) {
	session := s.RpcServer.Session[0]
	err := session.Write(msg)
	if err!=nil {
		logs.Printf("Error occured onNotify: %s\n",err)
	}else{
		logs.Printf("Notified %s", msg)
	}
}
