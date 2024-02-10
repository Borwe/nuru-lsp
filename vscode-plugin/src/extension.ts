import * as path from "path";
import * as fs from "fs";
import * as os from "os";
import { workspace, ExtensionContext, window, commands } from "vscode";

import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

let client: LanguageClient;

export function activate(context: ExtensionContext) {
  //register commands
  commands.registerCommand("nuru.languageserver.restart", async ()=>{
    if(client.isRunning()){
      await client.stop()
    }
    client.start()
    window.showInformationMessage("Nuru LSP restarted")
  })
  commands.registerCommand("nuru.languageserver.stop", async ()=>{
    if(client.isRunning()){
      await client.stop()
    }
    window.showInformationMessage("Nuru LSP stopped")
  })
  commands.registerCommand("nuru.languageserver.start", async ()=>{
    if(!client.isRunning()){
      await client.start()
    }
    window.showInformationMessage("Nuru LSP started")
  })


  const command = (()=>{
    if(os.platform() === "win32"){
      return "nuru-lsp.exe"
    }
    return "nuru-lsp"
  })()

  const serverOptions: ServerOptions = {
    run: { command: command, transport: TransportKind.stdio },
    debug: { command: command, transport: TransportKind.stdio },
  };

  // Options to control the language client
  const clientOptions: LanguageClientOptions = {
    // Register the server for all documents by default
    documentSelector: [{ language: "nuru" }],
    synchronize: {
      // Notify the server about file changes to '.clientrc files contained in the workspace
      fileEvents: workspace.createFileSystemWatcher("**/.clientrc"),
    },
  };

  // Create the language client and start the client.
  client = new LanguageClient(
    "nuru-lsp",
    "nuru-lsp",
    serverOptions,
    clientOptions
  );

  // Start the client. This will also launch the server
  client.start().catch((err)=>{
    window.showErrorMessage(`Nuru Lsp failed to start error: ${err}`)
  });
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
