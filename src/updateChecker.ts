import * as vscode from "vscode";
import { log } from "./helperModules";

export async function checkForUpdate() {
  const ext = vscode.extensions.getExtension("noir.vbook-extension-maker");
  const currentVersion = ext?.packageJSON.version;

  if (!currentVersion) {
    log("vbook-ext: Current version unknown");
    return;
  }

  try {
    const res = await fetch(
      "https://api.github.com/repos/faea726/vbook-extension-maker/releases/latest",
      {
        headers: { "User-Agent": "vscode-extension" },
      },
    );

    const data = (await res.json()) as { tag_name?: string; html_url?: string };
    const latest = data.tag_name?.replace(/^v/, "");

    if (latest && latest !== currentVersion) {
      vscode.window
        .showInformationMessage(
          `Latest release: ${latest}. Local version: ${currentVersion}\n\n`,
          "View Release",
        )
        .then((choice) => {
          if (choice === "View Release" && data.html_url) {
            vscode.env.openExternal(vscode.Uri.parse(data.html_url));
          }
        });
    }
  } catch (err) {
    log("vbook-ext: Failed to check for updates\n", err);
  }
}
