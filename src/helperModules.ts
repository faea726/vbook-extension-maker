import * as http from "http";
import * as fs from "fs";
import * as path from "path";
import * as vscode from "vscode";
import * as querystring from "querystring";

// Check if plugin.json exists
function pluginJsonExist(inputScriptPath: string): boolean {
  const workspaceFolders = vscode.workspace.workspaceFolders;
  if (!workspaceFolders || workspaceFolders.length === 0) {
    return false;
  }

  const rootPath = path.resolve(inputScriptPath, "../../");
  const pluginJsonPath = path.join(rootPath, "plugin.json");

  return fs.existsSync(pluginJsonPath);
}

// Manage temporary data
let tempData: Record<string, any> = {};

function setValue(key: string, value: any, scriptPath: string) {
  if (!pluginJsonExist(scriptPath)) {
    vscode.window.showErrorMessage("Invalid workspace.");
    return null;
  }

  const rootPath = path.resolve(scriptPath, "../../");
  const tempPath = path.join(rootPath, "test.json");

  tempData[key] = value;
  fs.writeFileSync(tempPath, JSON.stringify(tempData, null, 2), "utf-8");
}

function getValue(key: string, scriptPath: string): any {
  if (!pluginJsonExist(scriptPath)) {
    vscode.window.showErrorMessage("Invalid workspace.");
    return null;
  }

  const rootPath = path.resolve(scriptPath, "../../");
  const tempPath = path.join(rootPath, "test.json");

  if (!fs.existsSync(tempPath)) {
    return null;
  }

  const data = fs.readFileSync(tempPath, "utf-8");
  tempData = JSON.parse(data);

  if (!tempData[key]) {
    return null;
  }

  return tempData[key];
}

async function setURL(scriptPath: string) {
  if (!pluginJsonExist(scriptPath)) {
    vscode.window.showWarningMessage("Inavlid workspace.");
    return;
  }

  const regexInput =
    /^(https?:\/\/[a-zA-Z0-9.-]+(:\d+)?|(?:\d{1,3}\.){3}\d{1,3})$/;
  const url = await vscode.window.showInputBox({
    prompt: "Vbook app IP",
    ignoreFocusOut: true,
    value: getValue("appIP", scriptPath),
    placeHolder: "http://192.168.1.7:8080",
    validateInput: (value) => {
      value = value.trim();
      return !regexInput.test(value) ? "URL invalid" : null;
    },
  });

  if (!url) {
    return null;
  }
  const normalizedUrl = normalizeHost(url);

  setValue("appIP", normalizedUrl, scriptPath);
  log(`vbook-ext: Set Vbook App IP to: ${normalizedUrl}`);
  return normalizedUrl;
}

async function setParams(scriptName: string, scriptPath: string) {
  if (!pluginJsonExist(scriptPath)) {
    vscode.window.showWarningMessage("Inavlid workspace.");
    return;
  }

  const params = await vscode.window.showInputBox({
    prompt: "Params for script",
    ignoreFocusOut: true,
    value: getValue(scriptName, scriptPath),
    placeHolder: "param_1, param_2,...",
  });

  setValue(scriptName, params, scriptPath);
  return params;
}

function normalizeHost(input: string): string {
  const patterns = {
    http_host_port: /\bhttp:\/\/(?:\d{1,3}\.){3}\d{1,3}:\d+\b/,
    https_host_port: /\bhttps:\/\/(?:\d{1,3}\.){3}\d{1,3}:\d+\b/,
    https_host: /\bhttps:\/\/(?:\d{1,3}\.){3}\d{1,3}\b/,
    http_host: /\bhttp:\/\/(?:\d{1,3}\.){3}\d{1,3}\b/,
    host_only: /\b(?:\d{1,3}\.){3}\d{1,3}\b/,
  };

  let matched: RegExpMatchArray | null = null;
  let caseType: string | null = null;

  for (const key in patterns) {
    matched = input.match(patterns[key as keyof typeof patterns]);
    if (matched) {
      caseType = key;
      break;
    }
  }

  if (!matched || !caseType) {
    return "";
  }

  const host = matched[0];

  switch (caseType) {
    case "http_host_port":
    case "https_host_port":
      return host;

    case "http_host":
    case "https_host":
      return host + ":8080";

    case "host_only":
      return "http://" + host + ":8080";

    default:
      return "";
  }
}

function runLocalServer(port: number, scriptPath: string): http.Server {
  const SRC_PATH = path.resolve(scriptPath, "../../");
  // log(`vbook-ext: SRC_PATH: ${SRC_PATH}`);
  const server = http.createServer((req, res) => {
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

    fs.readFile(filePath, "utf8", (err, data) => {
      if (err) {
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
    log(`vbook-ext: Server listening on port ${port}`);
  });
  return server;
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
  pluginJsonExist,
  setURL,
  setParams,
  runLocalServer,
  log,
  getOutputChannel,
};
