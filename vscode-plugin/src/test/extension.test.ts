import * as assert from 'assert';
import * as vscode from 'vscode';
import { getLatestReleaseVersion } from '../utils';
import * as path from 'path';
import * as os from 'os';
import * as fs from 'fs';

suite('Getting nuru-lsp file online', () => {

	////get an vscode.ExtensionContext instance
	//const context = vscode.extensions.getExtension("BrianOrwe.nuru-lsp");
	////clear the nuru-lsp.zip and nuru-lsp executable
	//const extPath = context.extensionPath.replace(/\\/g, "/")
	//const zip = path.join(extPath, "nuru-lsp.zip").replace(/\\/g, "/")
	//const exe = path.join(extPath,
	//	`nuru-lsp${os.platform() === "win32" ? ".exe" : ""}`).replace(/\\/g, "/")
	//console.log("EXE:",exe)
	//if (fs.existsSync(exe)) {
	//	fs.unlinkSync(exe)
	//}
	//if (fs.existsSync(zip)) {
	//	console.log("ZIP:",zip)
	//	fs.unlinkSync(zip)
	//}
	//assert.strictEqual(fs.existsSync(exe), false)
	//assert.strictEqual(fs.existsSync(zip), false)

	console.log("YOLO!!")
	

	//test('Checking if nuru-lsp not in dir', async () => {
	//	const exists = await vscode.commands.executeCommand("nuru.languageserver.is-installed")
	//	assert.strictEqual(exists, false);
	//});

	test("Test getting latest release", async () => {
		const num = await getLatestReleaseVersion()
		assert.strict.notEqual(num, 0)
	})

	//test('Checking if nuru-lsp in dir after downloading', async () => {
	//	const downloaded = await vscode.commands.executeCommand("nuru.languageserver.download")
	//	assert.strictEqual(downloaded, true);
	//	const exists = await vscode.commands.executeCommand("nuru.languageserver.is-installed")
	//	assert.strictEqual(exists, true);
	//});
});

suite('Testing LSP workings', () => {
	//test("nuru-lsp is not running", async()=>{
	//	const isRunning = await vscode.commands.executeCommand("nuru.languageserver.is-running")
	//	assert.strictEqual(isRunning, false)
	//})

	//test("nuru-lSP attatches to nr file once opened", () => {
	//	assert.strictEqual(false, true, "NOT IMPLEMENTED")
	//})
})
