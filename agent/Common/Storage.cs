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
    /// <summary>Find history the Go implementation left in a directory it was once
    /// started from and install it here, so resuming those conversations needs no
    /// operator action. Silent and idempotent: returns a sentence describing what it
    /// restored, or null when there was nothing to do. It never throws and never
    /// writes to the directory it reads.</summary>
    public static partial string? RestoreGoHistory(string destination, TextWriter diagnostics);
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
