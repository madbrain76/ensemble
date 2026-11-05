using System.Text.RegularExpressions;

namespace Ensemble;

public sealed partial class Agent
{
    private static readonly string[] FirstPerson =
        ["i'll", "i will", "let me", "let's", "let us", "i'm going to", "i am going to",
         "i should", "i need to", "we'll", "i'd like to"];

    private static readonly string[] Connectors =
        ["now", "next", "then", "so", "first", "second", "finally", "also", "and", "but"];

    // The verbs that mean "I am about to do a thing", as opposed to describing one.
    private static readonly string[] ActionVerbs =
        ["writing", "running", "re-running", "rerunning", "launching", "creating", "adding",
         "applying", "updating", "fixing", "checking", "testing", "implementing", "proceeding",
         "executing", "issuing", "patching", "editing", "building", "wiring", "spawning",
         "dispatching", "committing", "querying", "fetching", "inspecting", "examining",
         "reviewing", "verifying", "confirming", "scanning", "probing", "gathering",
         "exploring", "summarizing", "drafting", "installing", "measuring", "benchmarking"];

    internal static partial bool PromisesAction(string? answer)
    {
        string text = (answer ?? "").Replace('\u2019', '\'').Trim();
        if (text.Length == 0)
        {
            return false;
        }
        string lower = text.ToLowerInvariant();
        foreach (string phrase in FirstPerson)
        {
            if (lower.Contains(phrase))
            {
                return true;
            }
        }
        bool labelled = text.EndsWith(':');
        if (labelled)
        {
            // A label colon is not a promise. "Fix order (cheapest risk-reduction
            // first):" followed by the list is a complete answer, so a colon only
            // counts when the same reply also names work in progress.
            return ActionVerbs.Any(verb => lower.Contains(verb));
        }
        foreach (string sentence in Regex.Split(text, @"[.!?]\s"))
        {
            string[] words = sentence.Split((char[]?)null, StringSplitOptions.RemoveEmptyEntries);
            int index = 0;
            while (index < words.Length
                   && Array.IndexOf(Connectors, words[index].ToLowerInvariant().Trim(',', ':')) >= 0)
            {
                index++;
            }
            // A parenthetical question must not hide the opener: the sentence
            // "Now inspecting the panics (are they checks?), the errors..." still
            // promises work.
            while (index < words.Length && words[index].StartsWith('('))
            {
                index++;
            }
            if (index >= words.Length)
            {
                continue;
            }
            if (Array.IndexOf(ActionVerbs, words[index].ToLowerInvariant().Trim(',', '.', ':', '(')) >= 0)
            {
                return true;
            }
        }
        return false;
    }

    internal partial string LastAgentAnswer()
    {
        var dialogue = History().View().Dialogue;
        for (int index = dialogue.Length - 1; index >= 0; index--)
        {
            if (dialogue[index].Actor != Actor.Agent)
            {
                continue;
            }
            return string.Join(' ', dialogue[index].Parts.OfType<TextPart>()
                .Select(part => part.Text));
        }
        return "";
    }
}
