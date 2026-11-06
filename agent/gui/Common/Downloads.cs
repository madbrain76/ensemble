using Ensemble;
using Microsoft.AspNetCore.Routing;

namespace EnsembleGui;

public static partial class Gui
{
    // Serving a work product the Agent explicitly produced. The path arrives as text,
    // so it is re-opened through the owning Agent's sandbox rather than trusted.
    private static partial void MapDownload(IEndpointRouteBuilder routes, Agent source,
        int port, string? publicOrigin, CancellationToken stopping);
}
