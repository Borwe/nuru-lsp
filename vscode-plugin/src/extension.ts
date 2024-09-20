import { exec, ExecException } from "child_process";
import * as os from "os";
import { workspace, ExtensionContext, window, commands } from "vscode";

import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";
import { downloadOrUpdate, isInstalled } from "./utils";

export let Context: ExtensionContext

let client: LanguageClient;
/** Hold information on location of lsp file to execute */
let command = "nuru-lsp"

const VERSION = "0.0.6"
const LINK_BASE = `https://github.com/Borwe/nuru-lsp/releases/download/v${VERSION}/`

async function downloadLSPExecutable(link: string) {
}

export function activate(context: ExtensionContext) {
  Context = context
  //register commands
  commands.registerCommand("nuru.languageserver.is-installed", isInstalled);
  commands.registerCommand("nuru.languageserver.download", downloadOrUpdate);
  commands.registerCommand("nuru.languageserver.restart", async ()=>{
    if(client.isRunning()){
      await client.stop()
    }
    client.start()
    window.showInformationMessage("Nuru LSP restarted")
  });
  commands.registerCommand("nuru.languageserver.start", async ()=>{
    if(!client.isRunning()){
      await client.start()
    }
    window.showInformationMessage("Nuru LSP started")
  });
  commands.registerCommand("nuru.languageserver.stop", async ()=>{
    if(client.isRunning()){
      await client.stop()
    }
    window.showInformationMessage("Nuru LSP stopped")
  });

  (async ()=>{
    let cmd: string = "nuru-lsp"
    let url = LINK_BASE+"nuru-lsp-ubuntu-latest.zip"
    let result: Promise<string> 
    if(os.platform() === "win32"){
      url = LINK_BASE + "nuru-lsp-windows-latest.zip"
      cmd = "nuru-lsp.exe"
    }
    if(os.platform() === "darwin"){
      url = LINK_BASE + "nuru-lsp-macos-latest.zip"
    }
    exec(cmd, async (err: ExecException|null, stdout: string, stderr: string)=>{
      if(err!=null){
        await downloadLSPExecutable(url)
      }
    })
    return result
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
