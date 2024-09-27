import { readdirSync } from "fs";
import * as vscode from "vscode";
import * as os from "os"
import * as fs from "fs"
import { Context } from "./extension";
import { exec } from "child_process";
import path = require("path");
import { pipeline } from "stream";
import { promisify } from "util";

const CMD: string = os.platform() == "win32" ? "nuru-lsp.exe" : "nuru-lsp"
export const VERSION = "0.0.06"
const OSTYPE = os.platform() === "win32" ? "windows" : os.platform() === "linux" ? "ubuntu" : os.platform() === "darwin" ? "macos" : "noooo"
const LINK_BASE = `https://github.com/Borwe/nuru-lsp/releases/download/v${VERSION}/nuru-lsp-${OSTYPE}-latest.zip`

type ReleaseInfo = {
    tag_name: string
}

export function isInstalled(): boolean {
    const extPath = Context.extensionPath
    if (readdirSync(extPath).find(f => f === CMD) == undefined) {
        return false
    }
    return true
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
        try {
            const currentVersion: number = await new Promise((resolve, reject) => {
                exec(getPathOfCMD() + " --version", (err, stdout, stderr) => {
                    if (err) {
                        throw err
                    }
                    resolve(parseVersionToNumber(stdout))
                })
            })
            const releaseVer = await getLatestReleaseVersion()
            if (currentVersion >= releaseVer) {
                vscode.window.showInformationMessage("You are already using latest executable of nuru-lsp")
                return true
            }
        } catch (e) {
            return false
        }
    }

    return await getAndInstallLatest()
}

const getPathOfCMD = () => path.join(Context.extensionPath, CMD)

async function getAndInstallLatest(): Promise<boolean> {
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
    let cmd = `tar -xf ${zipPath} -C ${extPath}`
    if(OSTYPE !== "windows"){
        //cmd to extract zip on linux & mac
        cmd = `unzip -q ${zipPath} -d ${extPath}`
    }
    return new Promise(resolve=>exec(`tar -xf ${zipPath} -C ${extPath}`, (err, stdout, stderr) => {
        if (err) {
            vscode.window.showErrorMessage(`Failed to extract zip file with commdand:\n ${cmd}`)
            resolve(false)
        }
        resolve(true)
    }))
}
