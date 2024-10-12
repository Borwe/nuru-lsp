import * as vscode from "vscode";
import * as os from "os"
import * as fs from "fs"
import { client, Context } from "./extension";
import { exec } from "child_process";
import path = require("path");
import { pipeline } from "stream";
import { promisify } from "util";

const GIT_API_URL = "https://api.github.com/repos/borwe/nuru-lsp/releases/latest"
export const CMD: string = os.platform() == "win32" ? "nuru-lsp.exe" : "nuru-lsp"
const OSTYPE = os.platform() === "win32" ? "windows" : os.platform() === "linux" ? "ubuntu" : os.platform() === "darwin" ? "macos" : "noooo"

async function getLatestReleaserVersionString(): Promise<ReleaseInfo> {
    try {
        const res = await fetch(GIT_API_URL)
        if (!res.ok) {
            throw "Nuru LSP: Internet not reachable"
        }
        const info: ReleaseInfo = JSON.parse(await res.text())
        return info
    } catch (err) {
        throw "Nuru LSP: Internet not reachable"
    }
}

async function getLatestLinkBase(): Promise<string> {
    const info = await getLatestReleaserVersionString()
    return `https://github.com/Borwe/nuru-lsp/releases/download/${info.tag_name}/nuru-lsp-${OSTYPE}-latest.zip`
}

type ReleaseInfo = {
    tag_name: string
}

export type CommandType = { cmd: string, args: string[] }

export function getExtentionPath(): CommandType {
    //check if debug enabled
    let args: string | undefined = undefined;
    const isDebug = vscode.workspace.getConfiguration("nuru-lsp").get<boolean>("dbg")
    if (isDebug) {
        args = path.join(Context.extensionPath, "lsp.log").replace(/\\/g, "/")
    }
    const cmd = getPathOfCMD()
    return {
        cmd: cmd,
        args: args != undefined ? [args] : []
    }
}

export function isInstalled(): boolean {
    const extPath = getPathOfCMD()
    if (fs.existsSync(extPath) == false) {
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
    const split = version.split(" ")
    const versionD = split.length > 1 ? split[1] : split[0]
    const numvers = versionD.substring(1).split(".")
    let num = 0
    for (let i = numvers.length - 1, j = 1; i >= 0; i--, j *= 10) {
        num = num + parseInt(numvers[i]) * j
    }
    return num
}

export async function getLatestReleaseVersion(): Promise<number> {
    const releaseObj: ReleaseInfo = await getLatestReleaserVersionString()
    return parseVersionToNumber(releaseObj.tag_name)
}

export async function downloadOrUpdate(): Promise<boolean> {
    try {
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
            vscode.window.showInformationMessage(`VERSIONS CURRENT: ${currentVersion} ONLINE: ${releaseVer}`)
            if (currentVersion >= releaseVer) {
                vscode.window.showInformationMessage("You are already using latest executable of nuru-lsp")
                return true
            }
            return true
        }

        return await getAndInstallLatest()
    } catch (err) {
        vscode.window.showErrorMessage("Error: " + err)
    }
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
    const LINK_BASE = await getLatestLinkBase()
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
    try {
        if (!await downloadOrUpdate()) {
            vscode.window.showErrorMessage("Initial Setup failed, couldn't start LSP, please check you have internet\n" +
                "then run-> Nuru LSP: Start again"
            )
            return;
        }
    } catch (err) {
        vscode.window.showErrorMessage(`Nuru Lsp error: ${err}`)
    }
    if (!client.isRunning()) {
        client.start().catch((err) => {
            vscode.window.showErrorMessage(`Nuru Lsp failed to start error: ${err}`)
        });
    }
}

export async function openLogFileIfDebug() {
    if (vscode.workspace.getConfiguration("nuru-lsp").get<boolean>("dbg")) {
        const logFile = vscode.Uri.file(getExtentionPath().args[0])
        const doc = await vscode.workspace.openTextDocument(logFile)
        vscode.window.showTextDocument(doc, { preview: false })
        return
    }
    vscode.window.showInformationMessage("Didn't enable debug mode in Nuru-lsp, not opening log file")
}
