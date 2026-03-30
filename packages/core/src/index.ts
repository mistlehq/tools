export type CliDescriptor = {
  command: string;
  displayName: string;
  providerId: string;
  summary: string;
  version: string;
};

export type BuiltInCommand = "help" | "version";

export function describeCli(input: CliDescriptor): string {
  return `${input.displayName} (${input.command}) [${input.providerId}]`;
}

export function parseBuiltInCommand(argv: readonly string[]): BuiltInCommand {
  const firstArgument = argv[2];

  if (firstArgument === "--version" || firstArgument === "version") {
    return "version";
  }

  return "help";
}

export function renderHelpLines(input: CliDescriptor): string[] {
  return [
    input.command,
    describeCli(input),
    "",
    input.summary,
    "",
    "Available commands:",
    "  help",
    "  version",
  ];
}

export function runCli(
  argv: readonly string[],
  input: CliDescriptor,
  writeLine: (line: string) => void = console.log,
): void {
  const command = parseBuiltInCommand(argv);

  if (command === "version") {
    writeLine(input.version);
    return;
  }

  for (const line of renderHelpLines(input)) {
    writeLine(line);
  }
}
