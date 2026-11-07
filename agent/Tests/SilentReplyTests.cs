using System.Net;
using System.Net.Sockets;
using System.Text;
using Ensemble;

// A model that answers with neither prose nor a tool call shows the user nothing,
// yet the provider reports success. The classifier cannot catch this: "" promises
// nothing, and the output-limit ladder only fires when the provider says the
// budget ran out. These runs observe what the actor does about the silence.
internal static class SilentReplyTests
{
    private static void Check(bool condition, string message)
    {
        if (!condition)
        {
            throw new InvalidOperationException(message);
        }
    }

    // Answers with the given contents in order, repeating the last, and tallies
    // every request so a runaway is a failure rather than a long wait.
    private static HttpListener StartProvider(string[] contents, int[] tally)
    {
        var probe = new TcpListener(IPAddress.Loopback, 0);
        probe.Start();
        int port = ((IPEndPoint)probe.LocalEndpoint).Port;
        probe.Stop();
        var listener = new HttpListener();
        listener.Prefixes.Add($"http://127.0.0.1:{port}/");
        listener.Start();
        _ = Task.Run(async () =>
        {
            while (listener.IsListening)
            {
                HttpListenerContext context;
                try { context = await listener.GetContextAsync(); }
                catch (Exception) { return; }
                using var reader = new StreamReader(context.Request.InputStream, Encoding.UTF8);
                _ = await reader.ReadToEndAsync();
                int index = Interlocked.Increment(ref tally[0]) - 1;
                string content = index < contents.Length ? contents[index] : contents[^1];
                byte[] json = Encoding.UTF8.GetBytes(
                    "{\"id\":\"c\",\"object\":\"chat.completion\",\"created\":1,"
                    + "\"model\":\"silent-model\",\"choices\":[{\"index\":0,"
                    + "\"finish_reason\":\"stop\",\"message\":{\"role\":\"assistant\","
                    + "\"content\":" + content + "}}]}");
                // Chunked framing makes the client give up, which would look like a
                // transport fault rather than a policy decision.
                context.Response.ContentType = "application/json";
                context.Response.ContentLength64 = json.Length;
                context.Response.OutputStream.Write(json, 0, json.Length);
                context.Response.Close();
            }
        });
        return listener;
    }

    private static string Url(HttpListener listener)
        => listener.Prefixes.First().TrimEnd('/') + "/v1";

    public static async Task Run(string directory)
    {
        // Every reply is empty. The turn must not simply vanish.
        var tally = new int[1];
        var listener = StartProvider(["\"\""], tally);
        var app = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: directory);
        try
        {
            var agent = new Agent(app, new Configuration(Vendor.OpenAI, "silent-model",
                Url(listener), "test-key"), name: "silent-empty");
            agent.InstallTools();
            var turn = agent.SubmitAsync([new TextPart("run the build")]);
            Task done = await Task.WhenAny(turn, Task.Delay(TimeSpan.FromSeconds(20)));
            Check(done == turn,
                $"the empty exchange never finished after {tally[0]} requests");
            Check(tally[0] == 1 + Agent.PromissoryLimit,
                $"{tally[0]} requests for an always-empty model; the silence was not "
                + "corrected within the budget");
            // System corrections are recorded facts but not dialogue, so inspect
            // the event stream the way the browser's projection does.
            Check(agent.History().Events().Count(fact => fact.Kind == EventKind.MessageReceived
                && fact.Actor == Actor.System
                && fact.Parts.OfType<TextPart>().Any(part => part.Text == Agent.SilentCorrection))
                == Agent.PromissoryLimit,
                "the empty reply was not corrected once per budgeted attempt");
            Check(agent.History().Events().Any(fact => fact.Kind == EventKind.ErrorOccurred
                && fact.Error == Agent.SilentStopNotice),
                "the turn closed silently; the user was never told the model said nothing");
            agent.Destroy();
        }
        finally
        {
            listener.Stop();
            listener.Close();
            app.Destroy();
        }

        // One empty reply followed by a real answer: the user gets the answer, and
        // exactly one correction is spent.
        var secondTally = new int[1];
        var second = StartProvider(["\"\"", "\"Built and tested; both hosts are updated.\""],
            secondTally);
        var home = new Ensemble.Ensemble(TextWriter.Null, dataDirectory: directory + "-2");
        try
        {
            var agent = new Agent(home, new Configuration(Vendor.OpenAI, "silent-model",
                Url(second), "test-key"), name: "silent-recovers");
            agent.InstallTools();
            var turn = agent.SubmitAsync([new TextPart("build both hosts")]);
            Task done = await Task.WhenAny(turn, Task.Delay(TimeSpan.FromSeconds(20)));
            Check(done == turn, "the recovering exchange never finished");
            string answer = await turn;
            Check(answer.Contains("both hosts are updated"),
                "the corrected answer never reached the user: " + answer);
            Check(secondTally[0] == 2,
                $"expected one correction, saw {secondTally[0]} provider requests");
            Check(!agent.History().Events().Any(fact => fact.Kind == EventKind.ErrorOccurred),
                "a turn that recovered was still reported as an empty stop");
            agent.Destroy();
        }
        finally
        {
            second.Stop();
            second.Close();
            home.Destroy();
        }
    }
}
