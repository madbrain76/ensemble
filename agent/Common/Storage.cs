namespace Ensemble;

public sealed partial class Ensemble
{
    // The application container, not any one actor's personal home. It holds
    // ensemble.json and agents/<name>/ for sibling top-level actors. A library
    // instance may omit storage when it only hosts transient Agents.
    private readonly string? dataDirectory;
    public partial string DataDirectory();
    // No filesystem effects or second Agent registry: derive a named home from
    // this parent's configured container. Unnamed transient Agents get no home.
    internal partial string? AgentDirectory(string? name);
    /// <summary>One-shot operator migration for state written by the Go
    /// implementation: reads a Go data directory, adapts its history documents to
    /// this codec and installs them here. Returns the process exit code.</summary>
    public static partial int RunGoImportCommand(string[] args, TextWriter output, TextWriter error);
}

public sealed partial class Agent
{
    // Set only during construction through Sandbox. History's save path may be
    // anywhere permitted by that Sandbox, including a session subdirectory; it
    // never supplies or changes this distinct data-home fact.
    private string? dataDirectory;
    /// <summary>The actor's configured home, not its save path or Jobs directory.
    /// Throws when a transient library Agent was constructed without storage.</summary>
    public partial string DataDirectory();
}
