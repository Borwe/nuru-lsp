import * as path from "path";
import * as fs from "fs";
import { workspace, ExtensionContext, window } from "vscode";

import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

let client: LanguageClient;

export function activate(context: ExtensionContext) {
  const command = "nuru-lsp"

  let foundNuruExecutable = false
  const paths = process.env.PATH?.split(path.delimiter).map(p=>path.join(p,command))
   || []
  for(const p of paths){
    if(fs.existsSync(p)){
      foundNuruExecutable = true
      break;
    }
  }

  if(!foundNuruExecutable){
    window.showWarningMessage("Missing nuru-lsp executable","Couldn't find nuru-lsp executable in your path, some things might not work")
  }

  const serverOptions: ServerOptions = {
    run: { command: command, transport: TransportKind.ipc },
    debug: {
      command: command,
      transport: TransportKind.ipc,
    },
  };

  // Options to control the language client
  const clientOptions: LanguageClientOptions = {
    // Register the server for all documents by default
    documentSelector: [{ scheme: "file", language: "nuru", pattern: "*.{nr,sr}" }],
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
  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
