using System.Net;
using System.Net.Sockets;
using System.Text;
using Ensemble;

// The correction for a reply that promises work is itself a request, so a bound
// that resets at the wrong moment becomes a self-amplifying queue. This file only
// observes how many times the provider was asked.
internal static class PromissoryLoopTests
{
    private static void Check(bool condition, string message)
    {
        if (!condition)
        {
            throw new InvalidOperationException(message);
        }
    }

    public static async Task Run(string directory)
    {
        var probe = new TcpListener(IPAddress.Loopback, 0);
        probe.Start();
        int port = ((IPEndPoint)probe.LocalEndpoint).Port;
        probe.Stop();
        var listener = new HttpListener();
        listener.Prefixes.Add($"http://127.0.0.1:{port}/");
        listener.Start();
        int requests = 0;
        var serving = Task.Run(async () =>
        {
            while (listener.IsListening)
            {
                HttpListenerContext context;
                try { context = await listener.GetContextAsync(); }
                catch (Exception) { return; }
                using var reader = new StreamReader(context.Request.InputStream, Encoding.UTF8);
                _ = await reader.ReadToEndAsync();
                Interlocked.Increment(ref requests);
                // The same stalling reply every time: it announces work, delivers no
                // tool call, and so satisfies the correction test forever.
                byte[] json = Encoding.UTF8.GetBytes("""
                    {"id":"c","object":"chat.completion","created":1,"model":"loop-model","choices":[
                    {"index":0,"finish_reason":"stop","message":{"role":"assistant",
                    "content":"Now I'll run the full suite and report back."}}]}
                    """);
                // Content length must be declared: chunked framing makes the client
                // give up, and a broken pipe here would look like an engine fault.
                context.Response.ContentType = "application/json";
                context.Response.ContentLength64 = json.Length;
                context.Response.OutputStream.Write(json, 0, json.Length);
                context.Response.Close();
            }
        });
        var app = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: directory);
        try
        {
            var agent = new Agent(app, new Configuration(Vendor.OpenAI, "loop-model",
                $"http://127.0.0.1:{port}/v1", "test-key"), name: "loop");
            agent.InstallTools();
            var turn = agent.SubmitAsync([new TextPart("run the build")]);
            // The GUI does not wait: typing while the agent is busy delivers a hint
            // into a live exchange, which is the shape that ran away.
            var hint = agent.PostAsync(new AgentInput(InputKind.Hint, "still waiting"));
            Task done = await Task.WhenAny(turn, Task.Delay(TimeSpan.FromSeconds(20)));
            Check(done == turn, $"the exchange never finished after {requests} provider requests");
            // One request, then at most PromissoryLimit corrections. Anything larger
            // means the queue was refilling itself: each correction enqueues work.
            _ = hint;
            Check(requests <= 2 + Agent.PromissoryLimit,
                $"{requests} requests for one prompt; the correction budget re-armed itself");
            agent.Destroy();

            // A human speaking again must earn a fresh budget, or the guard would
            // silently stop correcting the agent for the rest of the session.
            var again = new Agent(app, new Configuration(Vendor.OpenAI, "loop-model",
                $"http://127.0.0.1:{port}/v1", "test-key"), name: "loop2");
            again.InstallTools();
            int before = requests;
            var second = again.SubmitAsync([new TextPart("and now the tests")]);
            done = await Task.WhenAny(second, Task.Delay(TimeSpan.FromSeconds(20)));
            Check(done == second, "the second exchange never finished");
            int spent = requests - before;
            Check(spent >= 2, $"a new human turn was never corrected: only {spent} requests");
            Check(spent <= 1 + Agent.PromissoryLimit,
                $"the second exchange also ran away with {spent} requests");
            again.Destroy();
        }
        finally
        {
            app.Destroy();
            try { listener.Stop(); listener.Close(); } catch (Exception) { }
            await serving;
        }
    }
}
