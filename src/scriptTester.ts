import * as vscode from "vscode";
import * as net from "net";
import * as os from "os";
import { URL } from "url";
import * as path from "path";
import {
  pluginJsonExist,
  setURL,
  runLocalServer,
  log,
  setParams,
} from "./helperModules";

async function testScript() {
  log("vbook-ext: testScript");

  const fileData = getOpeningFileContent();
  if (!fileData) {
    vscode.window.showInformationMessage("No file opened!");
    return;
  }

  if (!pluginJsonExist(fileData.path)) {
    vscode.window.showWarningMessage("Invalid workspace.");
    return;
  }

  const appIP = String(await setURL(fileData.path));
  if (!appIP) {
    vscode.window.showErrorMessage("IP not set");
    return;
  }

  const _url = new URL(appIP);
  const hostParts = _url.hostname.split(".");
  const itf = `${hostParts[0]}.${hostParts[1]}.`;

  const serverPort = Number(_url.port) - 10;

  const params = await setParams(fileData.name, fileData.path);
  log(`vbook-ext: Params: ${params}`);

  const extName = vscode.workspace.workspaceFolders?.[0].name;

  const data = {
    ip: getLocalIP(itf, serverPort),
    root: `${extName}/src`,
    language: "javascript",
    script: fileData.content,
    input: params?.trim().includes(",")
      ? [params.split(",").map((p) => p.trim())]
      : [params?.trim()],
  };

  const request = [
    "GET /test HTTP/1.1",
    `Host: ${_url.hostname}`,
    "Connection: close",
    `data: ${Buffer.from(JSON.stringify(data)).toString("base64")}`,
    "",
    "",
  ].join("\r\n");

  // return;

  let server = runLocalServer(serverPort, fileData.path);
  const client = net.createConnection(
    { host: _url.hostname, port: Number(_url.port) },
    () => {
      log(`vbook-ext: Connected to vbook: ${_url.hostname}:${_url.port}`);
      client.write(request);
    }
  );

  let chunks: Buffer[] = [];

  client.on("data", (chunk) => {
    if (!chunk) {
      client.destroy();
    }
    chunks.push(chunk);
  });

  client.on("end", () => {
    log("vbook-ext: Disconnected from server");
    server.close();

    const response = Buffer.concat(chunks).toString("utf-8");
    // log("vbook-ext: Response:", response);
    const bodyStart = response.indexOf("{");
    const body = response.slice(bodyStart);

    try {
      const json = JSON.parse(body);
      const result = json.result;

      if (result) {
        const parsed = JSON.stringify(JSON.parse(result), null, 2);
        log(`vbook-ext: Result:\n${parsed}`);
      } else {
        log("vbook-ext: Result not found");
      }
    } catch (err) {
      log(`vbook-ext: Failed to parse response: ${err}`);
    }
  });

  client.on("error", (err) => {
    log(`vbook-ext: Connection error: ${err.message} `);
    server.close();
  });
}

function getLocalIP(itf: string, port: number): string | null {
  const interfaces = os.networkInterfaces();

  for (const name of Object.keys(interfaces)) {
    for (const iface of interfaces[name]!) {
      const ip = iface.address;

      if (iface.family === "IPv4" && !iface.internal && ip.startsWith(itf)) {
        const localIp = `http://${ip}:${port}`;
        log(`vbook-ext: Local IP: ${localIp}`);
        return localIp;
      }
    }
  }

  return null;
}

function getOpeningFileContent(): {
  name: string;
  path: string;
  content: string;
} | null {
  const editor = vscode.window.activeTextEditor;
  if (!editor) {
    return null;
  }
  return {
    name: path.basename(editor.document.fileName),
    path: editor.document.fileName,
    content: editor.document.getText(),
  };
}

export { testScript };
