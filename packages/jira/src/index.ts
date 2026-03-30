import { runCli, type CliDescriptor } from "@mistle-tools/core";

export const JiraCli: CliDescriptor = {
  command: "jira",
  displayName: "Jira CLI",
  providerId: "jira",
  summary: "Scaffold only. Provider commands have not been implemented yet.",
  version: "0.1.0",
};

function main(argv: readonly string[]): void {
  runCli(argv, JiraCli);
}

if (import.meta.main) {
  try {
    main(Bun.argv);
  } catch (error: unknown) {
    const message = error instanceof Error ? error.message : "Unknown Jira CLI error";
    console.error(message);
    process.exit(1);
  }
}
