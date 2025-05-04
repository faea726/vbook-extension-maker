import * as http from "http";
import * as fs from "fs";
import * as path from "path";
import * as vscode from "vscode";
import * as querystring from "querystring";

// Check if plugin.json exists
function checkPluginJson(): boolean {
  const workspaceFolders = vscode.workspace.workspaceFolders;
  if (!workspaceFolders || workspaceFolders.length === 0) {
    return false;
  }

  const rootPath = workspaceFolders[0].uri.fsPath;
  const pluginJsonPath = path.join(rootPath, "plugin.json");

  return fs.existsSync(pluginJsonPath);
}

// Manage temporary data
let tempData: Record<string, any> = {};
function setValue(key: string, value: any) {
  tempData[key] = value;
}
function getValue(key: string): any {
  return tempData[key];
}
function resetStore() {
  tempData = {};
}

async function setURL() {
  if (!checkPluginJson()) {
    vscode.window.showWarningMessage("Inavlid workspace.");
    return;
  }

  if (getValue("url")) {
    return getValue("url");
  }

  const regex = /^http:\/\/[^\/:]+:\d+$/;

  const url = await vscode.window.showInputBox({
    prompt: "Vbook app URL",
    ignoreFocusOut: true,
    validateInput: (value) => {
      value = value.trim();
      return !regex.test(value) ? "URL invalid" : null;
    },
  });

  if (!url) {
    setValue("url", undefined);
    return undefined;
  }

  setValue("url", url.trim());
  log(`Set URL to: ${url}`);
  return url.trim();
}

function runLocalServer(port: number) {
  const rootPath = vscode.workspace.workspaceFolders?.[0].uri.fsPath + "";
  const server = http.createServer((req, res) => {
    // log("vbook-ext: Request received:", req);
    const SRC_PATH = path.dirname(rootPath);

    const url = new URL(req.url!, `http://${req.headers.host}`);
    const queryString = url.searchParams.toString();

    const params = querystring.parse(queryString);

    const file = params["file"] as string;
    const root = params["root"] as string;

    if (!file || !root) {
      res.writeHead(400, { "Content-Type": "text/plain" });
      res.end("Missing required query parameters: file and root");
      return;
    }

    const filePath = path.join(SRC_PATH, root, file);
    // log("vbook-ext: filePath:", filePath);

    fs.readFile(filePath, "utf8", (_, data) => {
      if (_) {
        res.writeHead(500, { "Content-Type": "text/plain" });
        res.end("Error reading the file.");
        return;
      }

      const base64Data = Buffer.from(data).toString("base64");

      // Send the headers first
      res.writeHead(200, {
        "Content-Length": base64Data.length,
        "Content-Type": "text/plain", // or the appropriate type depending on your content
      });

      // Now, send the body data
      res.end(base64Data);
    });
  });

  server.listen(port, () => {
    log(`Server listening on port ${port}`);
  });
}

// Logger
let outputChannel: vscode.OutputChannel;
function getOutputChannel(): vscode.OutputChannel {
  if (!outputChannel) {
    outputChannel = vscode.window.createOutputChannel("Vbook Extension Maker");
  }
  outputChannel.show(true);
  return outputChannel;
}

function log(message: string) {
  getOutputChannel().appendLine(message);
}

export {
  checkPluginJson,
  setURL,
  getValue,
  runLocalServer,
  log,
  getOutputChannel,
};
