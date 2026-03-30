import { describe, expect, it } from "vitest";

import { renderHelpLines } from "@mistle-tools/core";

import { JiraCli } from "./index.js";

describe("jira cli scaffold", () => {
  it("declares the Jira CLI metadata", () => {
    expect(JiraCli.command).toBe("jira");
    expect(JiraCli.providerId).toBe("jira");
  });

  it("uses the shared core help output", () => {
    expect(renderHelpLines(JiraCli)).toContain(
      "Scaffold only. Provider commands have not been implemented yet.",
    );
  });
});
