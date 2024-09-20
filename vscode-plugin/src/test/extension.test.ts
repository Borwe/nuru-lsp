import * as assert from 'assert';
import * as vscode from 'vscode';
import * as extension from '../extension'
import { getLatestReleaseVersion } from '../utils';

suite('Getting nuru-lsp file', () => {

	test('Checking if nuru-lsp not in dir', async () => {
		const exists = await vscode.commands.executeCommand("nuru.languageserver.is-installed")
		assert.strictEqual(exists, false);
	});

	test("Test getting latest release", async()=>{
		const num = await getLatestReleaseVersion()
		assert.strict.notEqual(num, 0)
	})

	test('Checking if nuru-lsp in dir after downloading', async () => {
		const downloaded = await vscode.commands.executeCommand("nuru.languageserver.download")
		assert.strictEqual(downloaded, true);
		const exists = await vscode.commands.executeCommand("nuru.languageserver.is-installed")
		assert.strictEqual(exists, true);
	});
});
