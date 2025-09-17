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
  parseHttpResponse,
} from "./helperModules";

async function testScript() {
  log("\nvbook-ext: testScript");

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

  var _url: URL;
  try {
    _url = new URL(appIP);
  } catch (e) {
    log(`vbook-ext: Invalid App IP: ${appIP}`);
    return;
  }

  const host = _url.hostname;
  let interfaceIP = host;
  if (/^\d+\.\d+\.\d+\.\d+$/.test(host)) {
    const parts = host.split(".");
    interfaceIP = `${parts[0]}.${parts[1]}.${parts[2]}`; // e.g. "192.168"
  }
  const serverPort = Number(_url.port) - 10;

  const params = await setParams(fileData.name, fileData.path);
  log(`vbook-ext: Params: ${params}`);

  const extName = path.basename(path.resolve(fileData.path, "../../"));
  // log(`vbook-ext: extName: ${extName}`);

  const data = {
    ip: getLocalIP(serverPort, interfaceIP),
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
    },
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

    const rspStr = Buffer.concat(chunks).toString("utf-8");

    try {
      const body = parseHttpResponse(rspStr).body;

      log("\nResponse:");

      for (const [key, value] of Object.entries(body)) {
        if (typeof value === "object") {
          log(`\n${key}:`, JSON.stringify(value, null, 2));
        } else {
          if (value) {
            log(`\n${key}:`, value);
          }
        }
      }

      log("\nvbook-ext: Done");
    } catch (err) {
      log(`vbook-ext: ${err}`);
      log(`vbook-ext: Response:\n\n${rspStr}`);
    }
  });

  client.on("error", (err) => {
    log(`vbook-ext: Connection error: ${err.message} `);
    server.close();
  });
}

function getLocalIP(port: number, interfaceIP?: string): string | null {
  const interfaces = os.networkInterfaces();
  const prefix = interfaceIP?.trim();

  const isPrivate = (ip: string): boolean =>
    ip.startsWith("10.") ||
    ip.startsWith("192.168.") ||
    (ip.startsWith("172.") &&
      (() => {
        const second = parseInt(ip.split(".")[1], 10);
        return second >= 16 && second <= 31;
      })());

  const matchPrefix = (ip: string): boolean => {
    if (!prefix) {return false;}
    if (!/^\d+\.\d+\.\d+\.\d+$/.test(ip)) {return false;}
    const ipParts = ip.split(".");
    const prefParts = prefix.split(".");
    return prefParts.every((p, i) => ipParts[i] === p);
  };

  // Collect all IPv4 addresses
  const candidates: { ip: string; score: number }[] = [];

  for (const addrs of Object.values(interfaces)) {
    if (!addrs) {continue;}
    for (const addr of addrs) {
      if (addr.family !== "IPv4" || addr.internal) {continue;}

      let score = 0;
      if (isPrivate(addr.address)) {score += 10;} // prefer private LAN
      if (matchPrefix(addr.address)) {score += 5;} // prefer prefix
      // The higher the score, the more preferred

      candidates.push({ ip: addr.address, score });
    }
  }

  if (candidates.length === 0) {return null;}

  // Sort candidates by score descending
  candidates.sort((a, b) => b.score - a.score);

  const chosen = candidates[0].ip;
  const localIp = `http://${chosen}:${port}`;
  log(`vbook-ext: Local IP (chosen): ${localIp}`);
  return localIp;
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
