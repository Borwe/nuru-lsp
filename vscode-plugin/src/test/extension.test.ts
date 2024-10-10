import * as assert from 'assert';
import * as vscode from 'vscode';
import { CMD, CommandType, getExtentionPath, getLatestReleaseVersion } from '../utils';
import * as path from 'path';

async function resetConfigs(){
	await vscode.workspace.getConfiguration("nuru-lsp").update("dbg", false)
	await vscode.workspace.getConfiguration("nuru-lsp").update("execPath","")
}

suite('Getting nuru-lsp file online', () => {

	////get an vscode.ExtensionContext instance
	const context = vscode.extensions.getExtension("BrianOrwe.nuru-lsp");
	context.activate()
	assert.strictEqual(context.isActive,true)
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

	//test('Checking if nuru-lsp not in dir', async () => {
	//	const exists = await vscode.commands.executeCommand("nuru.languageserver.is-installed")
	//	assert.strictEqual(exists, false);
	//});

	test("Test getting latest release num", async () => {
		const num = await getLatestReleaseVersion()
		assert.strict.notEqual(num, 0)
	})

	test("Test seeing if enabled will have lsp.log", async () => {
		let path: CommandType =await vscode.commands.executeCommand("nuru.languageserver.command")
		console.log("CMD:",path)
		assert.strictEqual(path.args.length, 0)
		await vscode.workspace.getConfiguration("nuru-lsp").update("dbg", true)
		path =await vscode.commands.executeCommand("nuru.languageserver.command")
		assert.strictEqual(path.args.length, 1)
		resetConfigs()
	})

	test("Test getting location of executable", async () => {
		let cmd: CommandType =await vscode.commands.executeCommand("nuru.languageserver.command")
		const defaultPath = path.join(context.extensionPath, CMD).replace(/\\/g, "/").toUpperCase()
		assert.strictEqual(cmd.cmd.toUpperCase(), defaultPath)
		await vscode.workspace.getConfiguration("nuru-lsp").update("execPath","/something/something/nuru-lsp.exe")
		cmd =await vscode.commands.executeCommand("nuru.languageserver.command")
		assert.notStrictEqual(cmd.cmd.toUpperCase(), defaultPath)
		resetConfigs()
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
