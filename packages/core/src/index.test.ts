import { describe, expect, it } from "vitest";

import { describeCli, parseBuiltInCommand, renderHelpLines } from "./index.js";

const TestCli = {
  command: "jira",
  displayName: "Jira CLI",
  providerId: "jira",
  summary: "Scaffold only. Provider commands have not been implemented yet.",
  version: "0.1.0",
};

describe("describeCli", () => {
  it("formats provider metadata for CLI help output", () => {
    expect(describeCli(TestCli)).toBe("Jira CLI (jira) [jira]");
  });
});

describe("parseBuiltInCommand", () => {
  it("parses version aliases", () => {
    expect(parseBuiltInCommand(["node", "jira", "version"])).toBe("version");
    expect(parseBuiltInCommand(["node", "jira", "--version"])).toBe("version");
  });

  it("defaults to help", () => {
    expect(parseBuiltInCommand(["node", "jira"])).toBe("help");
    expect(parseBuiltInCommand(["node", "jira", "help"])).toBe("help");
  });
});

describe("renderHelpLines", () => {
  it("renders the built-in help output", () => {
    expect(renderHelpLines(TestCli)).toEqual([
      "jira",
      "Jira CLI (jira) [jira]",
      "",
      "Scaffold only. Provider commands have not been implemented yet.",
      "",
      "Available commands:",
      "  help",
      "  version",
    ]);
  });
});
