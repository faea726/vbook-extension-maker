import * as vscode from "vscode";
import { generateProject } from "./projectGenerator";
import { testScript } from "./scriptTester";
import { log } from "./helperModules";
import { buildExtension } from "./extenstionBuilder";
import { installExtension } from "./extensionInstaller";

export function activate(context: vscode.ExtensionContext) {
  log("vbook-ext: actived!");

  // Create project with template
  const createProjectCmd = vscode.commands.registerCommand(
    "vbook-extension-maker.createProject",
    generateProject,
  );
  context.subscriptions.push(createProjectCmd);

  // Test script
  const testScriptCmd = vscode.commands.registerCommand(
    "vbook-extension-maker.testScript",
    testScript,
  );
  context.subscriptions.push(testScriptCmd);

  // Build extension
  const buildExtensionCmd = vscode.commands.registerCommand(
    "vbook-extension-maker.buildExtension",
    buildExtension,
  );
  context.subscriptions.push(buildExtensionCmd);

  // Install extension
  const installExtensionCmd = vscode.commands.registerCommand(
    "vbook-extension-maker.installExtension",
    installExtension,
  );
  context.subscriptions.push(installExtensionCmd);
}

export function deactivate() {}
