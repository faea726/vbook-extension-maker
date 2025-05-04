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

  const regexInput =
    /^(https?:\/\/[a-zA-Z0-9.-]+(:\d+)?|(?:\d{1,3}\.){3}\d{1,3})$/;
  const url = await vscode.window.showInputBox({
    prompt: "Vbook app IP",
    ignoreFocusOut: true,
    validateInput: (value) => {
      value = value.trim();
      return !regexInput.test(value) ? "URL invalid" : null;
    },
  });

  if (!url) {
    setValue("url", undefined);
    return undefined;
  }
  const normalizedUrl = normalizeHost(url);

  setValue("url", normalizedUrl);
  log(`vbook-ext: Set Vbook App IP to: ${normalizedUrl}`);
  return normalizedUrl;
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

function runLocalServer(port: number): http.Server {
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
  checkPluginJson,
  setURL,
  getValue,
  runLocalServer,
  log,
  getOutputChannel,
};
