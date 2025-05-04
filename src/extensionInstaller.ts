import * as vscode from "vscode";
import * as fs from "fs";
import * as path from "path";
import { setURL, getValue, checkPluginJson } from "./helperModules";

async function installExtension() {
  console.log("vbook-ext: installExtension");
  const rootPath = vscode.workspace.workspaceFolders?.[0].uri.fsPath;
  if (!rootPath || !checkPluginJson()) {
    vscode.window.showWarningMessage("Inavlid workspace.");
    return;
  }

  const data = preparePluginData(rootPath);
  console.log("vbook-ext: data:", data);

  var url: string;
  url = getValue("url");
  if (!url) {
    url = await setURL();
  }

  try {
    console.log(`vbook-ext: Connect to: ${url}/install`);
    await fetch(`${url}/install`, {
      method: "GET",
      headers: {
        data: Buffer.from(JSON.stringify(data)).toString("base64"),
      },
    });
  } catch (error) {
    console.log("vbook-ext: done installation process.");
  }
}

function preparePluginData(pluginDir: string): any {
  const pluginDetailPath = path.join(pluginDir, "plugin.json");
  const iconPath = path.join(pluginDir, "icon.png");

  if (!fs.existsSync(pluginDetailPath) || !fs.existsSync(iconPath)) {
    throw new Error("invalid plugin");
  }

  const pluginDetail = JSON.parse(fs.readFileSync(pluginDetailPath, "utf-8"));
  let data: any = { ...pluginDetail.metadata, ...pluginDetail.script };

  data.id = "debug-" + pluginDetail.metadata.source;

  const iconBuffer = fs.readFileSync(iconPath);
  data.icon = `data:image/*;base64,${iconBuffer.toString("base64")}`;

  data.enabled = true;
  data.debug = true;
  data.data = {};

  // Read the plugin scripts from the src folder and add them to the data
  const pluginScripts = fs
    .readdirSync(path.join(pluginDir, "src"))
    .filter((file) => file.endsWith(".js"));

  for (const script of pluginScripts) {
    const scriptPath = path.join(pluginDir, "src", script);
    if (fs.existsSync(scriptPath)) {
      console.log(script);
      data.data[script] = fs.readFileSync(scriptPath, "utf-8");
    }
  }
  data.data = JSON.stringify(data.data);

  return data;
}

export { installExtension };
