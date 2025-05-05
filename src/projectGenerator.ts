import * as vscode from "vscode";
import * as fs from "fs";
import * as path from "path";
import { log } from "./helperModules";

async function generateProject() {
  log("\nvbook-ext: generateProject");
  const projectName = await vscode.window.showInputBox({
    prompt: "Enter a name for your new project",
    ignoreFocusOut: true,
    validateInput: (value) => {
      return value.trim().length === 0 ? "Project name cannot be empty" : null;
    },
  });

  if (!projectName) {
    vscode.window.showErrorMessage(
      "Project creation cancelled. No name was given."
    );
    return;
  }

  const targetUri = await vscode.window.showOpenDialog({
    canSelectFolders: true,
    openLabel: "Select folder to create the project in",
  });

  if (!targetUri || targetUri.length === 0) {
    vscode.window.showErrorMessage("No folder selected.");
    return;
  }

  const destinationPath = path.join(targetUri[0].fsPath, projectName);
  const templatePath = path.join(__dirname, "..", "template");

  if (fs.existsSync(destinationPath)) {
    vscode.window.showErrorMessage(
      `A folder named "${projectName}" already exists.`
    );
    return;
  }

  try {
    copyDirectory(templatePath, destinationPath);
    vscode.window.showInformationMessage(
      `Project "${projectName}" created successfully!`
    );

    const uri = vscode.Uri.file(destinationPath);
    await vscode.commands.executeCommand("vscode.openFolder", uri, false);
  } catch (error) {
    vscode.window.showErrorMessage(`Failed to create project: ${error}`);
  }
}

function copyDirectory(src: string, dest: string) {
  if (!fs.existsSync(dest)) {
    fs.mkdirSync(dest, { recursive: true });
  }

  const entries = fs.readdirSync(src, { withFileTypes: true });

  for (const entry of entries) {
    const srcPath = path.join(src, entry.name);
    const destPath = path.join(dest, entry.name);

    if (entry.isDirectory()) {
      copyDirectory(srcPath, destPath);
    } else {
      fs.copyFileSync(srcPath, destPath);
    }
  }
}

export { generateProject };
