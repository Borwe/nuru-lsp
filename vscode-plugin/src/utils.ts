import { readdirSync } from "fs";
import * as vscode from "vscode";
import * as os from "os"
import * as fs from "fs"
import { client, Context } from "./extension";
import { exec } from "child_process";
import path = require("path");
import { pipeline } from "stream";
import { promisify } from "util";

export const CMD: string = os.platform() == "win32" ? "nuru-lsp.exe" : "nuru-lsp"
export const VERSION = "0.0.07"
const OSTYPE = os.platform() === "win32" ? "windows" : os.platform() === "linux" ? "ubuntu" : os.platform() === "darwin" ? "macos" : "noooo"
const LINK_BASE = `https://github.com/Borwe/nuru-lsp/releases/download/v${VERSION}/nuru-lsp-${OSTYPE}-latest.zip`

type ReleaseInfo = {
    tag_name: string
}

export type CommandType = { cmd: string, args: string[] }

export function getExtentionPath(): CommandType {
    //check if debug enabled
    let args: string | undefined = undefined;
    const isDebug = vscode.workspace.getConfiguration("nuru-lsp").get<boolean>("dbg")
    if (isDebug) {
        args = path.join(Context.extensionPath, "lsp.log")
    }
    const cmd = getPathOfCMD()
    vscode.window.showInformationMessage("NURU-LSP server path: " +
        cmd + " " + args)
    return {
        cmd: cmd,
        args: args != undefined ? [args] : []
    }
}

export function isInstalled(): boolean {
    const extPath = Context.extensionPath
    if (readdirSync(extPath).find(f => f === CMD) == undefined) {
        vscode.window.showInformationMessage("NURU-LSP server not found, downloading...")
        return false
    }
    vscode.window.showInformationMessage("NURU-LSP server found")
    return true
}

async function showStatusBarMessage(msg: string): Promise<vscode.StatusBarItem> {
    const item = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 100)
    item.text = msg
    item.show()
    Context.subscriptions.push(item)
    return item
}

function parseVersionToNumber(version: string): number {
    const numvers = version.substring(1).split(".")
    let num = 0
    for (let i = numvers.length - 1, j = 1; i >= 0; i--, j *= 10) {
        num = num + parseInt(numvers[i]) * j
    }
    return num
}

export async function getLatestReleaseVersion(): Promise<number> {
    const res = await fetch("https://api.github.com/repos/borwe/nuru-lsp/releases/latest")
    const releaseObj: ReleaseInfo = JSON.parse(await res.text())
    return parseVersionToNumber(releaseObj.tag_name)
}

export async function downloadOrUpdate(): Promise<boolean> {
    if (isInstalled()) {
        const currentVersion: number = await new Promise((resolve, reject) => {
            exec(getPathOfCMD() + " --version", (err, stdout, stderr) => {
                if (err) {
                    vscode.window.showInformationMessage("FAILED Getting version info:" + err)
                    return 0
                }
                resolve(parseVersionToNumber(stdout))
            })
        })
        const releaseVer = await getLatestReleaseVersion()
        if (currentVersion >= releaseVer) {
            vscode.window.showInformationMessage("You are already using latest executable of nuru-lsp")
            return true
        }
    }

    return await getAndInstallLatest()
}

function getPathOfCMD(): string {
    const pathLoc = vscode.workspace.getConfiguration("nuru-lsp").get<string>("execPath").replace(/\\/g, "/")
    if (pathLoc == undefined || pathLoc.length == 0) {
        return path.join(Context.extensionPath, CMD).replace(/\\/g, "/")
    }
    return pathLoc
}

async function getAndInstallLatest(): Promise<boolean> {
    const downloadStatus = await showStatusBarMessage("Downloading nuru-lsp")
    const resp = await fetch(LINK_BASE)
    if (!resp.ok) {
        vscode.window.showErrorMessage("Failed to download zip file from: " + LINK_BASE)
        return false
    }

    const extPath = Context.extensionPath
    const zipPath = path.join(extPath, "nuru-lsp.zip").replace(/\\/g, "/")
    const fstream = fs.createWriteStream(zipPath)
    await promisify(pipeline)(resp.body, fstream)
    await new Promise(resolve => { fstream.close(resolve) })
    downloadStatus.hide()

    let cmd = `tar -xf ${zipPath} -C ${extPath}`
    if (OSTYPE !== "windows") {
        //cmd to extract zip on linux & mac
        cmd = `unzip ${zipPath} -d ${extPath}`
    }

    const extractStatus = await showStatusBarMessage("Extracting nuru-lsp.zip")
    return new Promise(resolve => exec(cmd, (err, stdout, stderr) => {
        if (err) {
            vscode.window.showErrorMessage(`Failed to extract zip file with commdand:\n ${cmd}`)
            resolve(false)
            extractStatus.hide()
            return
        }
        extractStatus.hide()
        resolve(true)
        vscode.window.showInformationMessage("NURU LSP server setup complete")
    }))
}

export async function handleLaunchingServer() {
    if (!await downloadOrUpdate()) {
        vscode.window.showErrorMessage("Initial Setup failed, couldn't start LSP, please check you have internet\n" +
            "then run-> Nuru LSP: Start again"
        )
        return;
    }
    if (!client.isRunning()) {
        client.start().catch((err) => {
            vscode.window.showErrorMessage(`Nuru Lsp failed to start error: ${err}`)
        });
    }
}