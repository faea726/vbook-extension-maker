import * as vscode from "vscode";
import * as net from "net";
import * as os from "os";
import { URL } from "url";
import {
  checkPluginJson,
  setURL,
  getValue,
  runLocalServer,
} from "./helperModules";

async function testScript() {
  console.log("vbook-ext: testScript");
  if (!checkPluginJson()) {
    vscode.window.showWarningMessage("Invalid workspace.");
    return;
  }

  const fileContent = getOpeningFileContent();
  if (!fileContent) {
    vscode.window.showInformationMessage("No file opened!");
    return;
  }

  var url: string;
  url = getValue("url");
  if (!url) {
    url = await setURL();
  }
  const _url = new URL(url);

  const params = await vscode.window.showInputBox({
    prompt: "Params for the script",
    ignoreFocusOut: true,
  });
  const ext_name = vscode.workspace.workspaceFolders?.[0].name;

  const data = {
    ip: getLocalIP(Number(_url.port)),
    root: `${ext_name}/src`,
    language: "javascript",
    script: fileContent,
    input: params?.trim().includes(",")
      ? [params.split(",").map((p) => p.trim())]
      : [params],
  };

  const request = [
    "GET /test HTTP/1.1",
    `Host: ${_url.hostname}`,
    "Connection: close",
    `data: ${Buffer.from(JSON.stringify(data)).toString("base64")}`,
    "",
    "",
  ].join("\r\n");

  runLocalServer(Number(_url.port));

  const client = net.createConnection(
    { host: _url.hostname, port: Number(_url.port) },
    () => {
      console.log(
        `vbook-ext: Connected to vbook: ${_url.hostname}:${_url.port}`
      );
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
    console.log("vbook-ext: Disconnected from server");

    const response = Buffer.concat(chunks).toString("utf-8");
    // console.log("vbook-ext: Response:", response);
    const bodyStart = response.indexOf("{");
    const body = response.slice(bodyStart);

    try {
      const json = JSON.parse(body);
      const result = json.result;

      if (result) {
        const parsed = JSON.parse(result);
        console.log("vbook-ext: Parsed Result:", parsed);
      } else {
        console.log("vbook-ext: Result not found");
      }
    } catch (err) {
      console.error("vbook-ext: Failed to parse response:", err);
    }
  });

  client.on("error", (err) => {
    console.error("vbook-ext: Connection error:", err.message);
  });
}

function getLocalIP(port: number): string | null {
  const interfaces = os.networkInterfaces();

  for (const name of Object.keys(interfaces)) {
    for (const iface of interfaces[name]!) {
      const ip = iface.address;

      if (
        iface.family === "IPv4" &&
        !iface.internal &&
        (ip.startsWith("192.168.") ||
          ip.startsWith("10.") ||
          ip.match(/^172\.(1[6-9]|2\d|3[0-1])\./))
      ) {
        const localIp = `http://${ip}:${port}`;
        console.log("vbook-ext: Local IP:", localIp);
        return localIp;
      }
    }
  }

  return null;
}

function getOpeningFileContent(): string {
  const editor = vscode.window.activeTextEditor;
  if (!editor) {
    return "";
  }
  return editor.document.getText();
}

export { testScript };
