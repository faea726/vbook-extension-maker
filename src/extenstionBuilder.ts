import * as fs from "fs";
import * as path from "path";
import archiver = require("archiver");
import * as vscode from "vscode";
import { log } from "./helperModules";

function buildExtension() {
  log("vbook-ext: buildExtension");

  const editor = vscode.window.activeTextEditor;
  if (!editor) {
    vscode.window.showWarningMessage("Please leave a script open!");
    return null;
  }

  const scriptPath = editor.document.fileName;
  const rootPath = path.resolve(scriptPath, "../../");

  // Validate files exist
  const plugin_json_path = path.join(rootPath, "plugin.json");
  const icon_path = path.join(rootPath, "icon.png");
  const src_path = path.join(rootPath, "src");
  const output_path = path.join(rootPath, "plugin.zip");
  if (
    !fs.existsSync(plugin_json_path) ||
    !fs.existsSync(icon_path) ||
    !fs.existsSync(src_path) ||
    !fs.lstatSync(src_path).isDirectory()
  ) {
    vscode.window.showWarningMessage("Files not found.");
    return;
  }

  // Zip
  const output = fs.createWriteStream(output_path);
  const archive = archiver("zip", { zlib: { level: 9 } });

  output.on("close", () => {
    log(`vbook-ext: plugin.zip created: ${archive.pointer()} bytes`);
  });
  archive.on("error", (err) => {
    throw err;
  });
  archive.pipe(output);

  archive.file(plugin_json_path, { name: "plugin.json" });
  archive.file(icon_path, { name: "icon.png" });
  archive.directory(src_path, "src");

  archive.finalize();
}

export { buildExtension };
