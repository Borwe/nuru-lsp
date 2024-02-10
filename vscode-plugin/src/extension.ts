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
  const command = (()=>{
    if(os.platform() === "win32"){
      return "nuru-lsp.exe"
    }
    return "nuru-lsp"
  })()

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
    window.showWarningMessage("Make sure nuru-lsp executable in your path, some things might not work if not")
  }

  const serverOptions: ServerOptions = {
    run: { command: command, transport: TransportKind.stdio },
    debug: {
      command: command,
      transport: TransportKind.stdio,
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
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
