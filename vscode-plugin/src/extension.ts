import { exec, ExecException } from "child_process";
import * as os from "os";
import { workspace, ExtensionContext, window, commands } from "vscode";

import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";
import { downloadOrUpdate, getExtentionPath, isInstalled, handleLaunchingServer, VERSION } from "./utils";

export let Context: ExtensionContext

export let client: LanguageClient;

export function activate(context: ExtensionContext) {
  /** Hold information on location of lsp file to execute */
  Context = context
  const command = getExtentionPath()
  //register commands
  commands.registerCommand("nuru.languageserver.is-installed", isInstalled);
  commands.registerCommand("nuru.languageserver.download", downloadOrUpdate);
  commands.registerCommand("nuru.languageserver.command", getExtentionPath)
  commands.registerCommand("nuru.languageserver.is-running", () => {
    if (client && client.isRunning()) {
      return true
    }
    return false
  })
  commands.registerCommand("nuru.languageserver.restart", async () => {
    if (client.isRunning()) {
      await client.stop()
    }
    handleLaunchingServer()
    window.showInformationMessage("Nuru LSP restarted")
  });
  commands.registerCommand("nuru.languageserver.start", async () => {
    if (!client.isRunning()) {
      handleLaunchingServer()
    }
  });
  commands.registerCommand("nuru.languageserver.stop", async () => {
    if (client.isRunning()) {
      await client.stop()
    }
    window.showInformationMessage("Nuru LSP stopped")
  });


  const serverOptions: ServerOptions = {
    run: { command: command, transport: TransportKind.stdio },
    debug: { command: command, transport: TransportKind.stdio },
  };

  // Options to control the language client
  const clientOptions: LanguageClientOptions = {
    // Register the server for all documents by default
    documentSelector: [{ language: "nr", scheme:"file" }],
    synchronize: {
      // Notify the server about file changes to '.clientrc files contained in the workspace
      fileEvents: workspace.createFileSystemWatcher("**/*.{nr,sr}"),
    },
  };

  // Create the language client and start the client.
  client = new LanguageClient(
    "nuru-lsp",
    "nuru-lsp",
    serverOptions,
    clientOptions
  );

  handleLaunchingServer()
}


export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
