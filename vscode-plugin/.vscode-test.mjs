import { defineConfig } from '@vscode/test-cli';
import * as dotenv from "dotenv"

function getInstallationConfig(){
	dotenv.config()
	const path = process.env.VSCODE_PATH
	if(path){
		return {
			fromPath: path
		}
	}
	return undefined
}

export default defineConfig({
	files: 'out/test/**/*.test.js',
	useInstallation: getInstallationConfig(),
	mocha: {
		ui: "tdd",
		timeout: 5000000000
	}
});
