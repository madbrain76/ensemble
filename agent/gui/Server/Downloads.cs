using Ensemble;
using Microsoft.AspNetCore.Routing;

namespace EnsembleGui;

public static partial class Gui
{
    private static partial void MapDownload(IEndpointRouteBuilder routes, Agent source,
        int port, string? publicOrigin, CancellationToken stopping)
    {
        routes.MapGet("/download", async (HttpContext context) =>
        {
            context.Response.Headers.CacheControl = "no-store";
            context.Response.Headers["X-Content-Type-Options"] = "nosniff";
            context.Response.Headers.ContentSecurityPolicy = "default-src 'none'; sandbox";
            var request = context.Request;
            string? expectedOrigin = RequestOrigin(context, port, publicOrigin);
            bool origin = expectedOrigin is not null && MatchesOrigin(request, expectedOrigin);
            bool sameSite = request.Headers["Sec-Fetch-Site"].ToString()
                is "" or "none" or "same-origin";
            if (!origin || !sameSite)
            {
                return Results.StatusCode(StatusCodes.Status403Forbidden);
            }
            var query = request.Query;
            if (query.Keys.Any(key => key is not ("agent" or "path" or "download"))
                || query["agent"].Count != 1 || query["path"].Count != 1
                || (query.ContainsKey("download") && query["download"] != "1"))
            {
                return Results.BadRequest(new { error = "Expected Agent identity and one file path." });
            }
            string path = query["path"].ToString();
            string name = System.IO.Path.GetFileName(path.TrimEnd('/', '\\'));
            if (name.Length == 0 || path.Contains("..", StringComparison.Ordinal))
            {
                return Results.BadRequest(new { error = "Work product must be one workspace file." });
            }
            string mime;
            try
            {
                // Refuse anything the browser might execute or sniff before the file
                // is even opened. This is the same whitelist the tool publishes.
                mime = Tool.DownloadMime(name);
            }
            catch (ArgumentException)
            {
                return Results.StatusCode(StatusCodes.Status415UnsupportedMediaType);
            }
            try
            {
                string address = query["agent"].ToString();
                Agent? owner = address.Length == 0 ? source : source.Parent().FindAgentAddress(address);
                if (owner is null)
                {
                    return Results.NotFound();
                }
                using var linked = CancellationTokenSource.CreateLinkedTokenSource(
                    stopping, context.RequestAborted);
                linked.Token.ThrowIfCancellationRequested();
                // The sandbox, not the caller, decides what "inside the workspace"
                // means, and the regular-file check rejects directories and FIFOs.
                await using var stream = owner.Sandbox().OpenRegularRead(path);
                long remaining = stream.Length;
                const long browserLimit = 64 * 1024 * 1024;
                if (remaining > browserLimit)
                {
                    return Results.Json(new
                    {
                        error = "Work product exceeds the 64 MiB browser download limit."
                    }, statusCode: 413);
                }
                context.Response.ContentType = mime;
                context.Response.ContentLength = remaining;
                context.Response.Headers.ContentDisposition =
                    $"attachment; filename=\"{name.Replace("\"", "")}\"";
                byte[] buffer = new byte[64 * 1024];
                while (remaining > 0)
                {
                    int count = (int)Math.Min(buffer.Length, remaining);
                    int read = await stream.ReadAsync(buffer.AsMemory(0, count), linked.Token);
                    if (read == 0)
                    {
                        throw new EndOfStreamException("Work product changed while being read.");
                    }
                    await context.Response.Body.WriteAsync(buffer.AsMemory(0, read), linked.Token);
                    remaining -= count;
                }
                return Results.Empty;
            }
            catch (Exception error) when (!context.Response.HasStarted && error is
                (ArgumentException or IOException or UnauthorizedAccessException
                    or ObjectDisposedException or InvalidOperationException))
            {
                // Never disclose whether a path exists outside the workspace.
                return Results.NotFound();
            }
        });
    }
}
