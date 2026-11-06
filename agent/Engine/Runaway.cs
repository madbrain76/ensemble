using System.Text;

namespace Ensemble;

public sealed partial class Engine
{
    // One instance belongs to one transport attempt: a retry starts fresh, and a
    // discarded attempt's bytes are never carried into the next decision.
    internal sealed class RunawayGuard
    {
        sealed class State
        {
            public char Run;
            public int RunLength;
            public readonly StringBuilder Tail = new();
        }

        private readonly Dictionary<int, State> parts = [];

        public bool Observe(int partId, ChunkKind kind, string? text)
        {
            if (kind == ChunkKind.ToolCall || string.IsNullOrEmpty(text))
            {
                return false;
            }
            if (!parts.TryGetValue(partId, out State? state))
            {
                state = parts[partId] = new State();
            }
            foreach (char character in text)
            {
                if (char.IsWhiteSpace(character))
                {
                    // A wide table pads with hundreds of contiguous spaces. Whitespace
                    // is never evidence of collapse, and it breaks an earlier run.
                    state.Run = character;
                    state.RunLength = 0;
                }
                else if (character == state.Run)
                {
                    state.RunLength++;
                }
                else
                {
                    state.Run = character;
                    state.RunLength = 1;
                }
                state.Tail.Append(character);
                if (state.Tail.Length > RunawayPhrase * 2)
                {
                    state.Tail.Remove(0, state.Tail.Length - RunawayPhrase * 2);
                }
                if (state.RunLength >= RunawayRunes)
                {
                    return true;
                }
            }
            return LooksCollapsed(state.Tail.ToString());
        }
    }

    // Wrap the caller's sink so the guard sees every chunk before it is delivered.
    internal static Action<Chunk>? WatchRunaway(Action<Chunk>? onChunk, RunawayGuard guard,
        Action trip)
    {
        if (onChunk is null)
        {
            return null;
        }
        return chunk =>
        {
            if (guard.Observe(chunk.PartID, chunk.Kind, chunk.Text))
            {
                trip();
            }
            onChunk(chunk);
        };
    }

    internal static partial bool LooksCollapsed(string tail)
    {
        if (tail.Length < RunawayPhrase)
        {
            return false;
        }
        for (int period = 2; period <= RunawayPeriod; period++)
        {
            if (tail.Length < period * 2)
            {
                break;
            }
            bool uniform = true;
            for (int index = tail.Length - period; index < tail.Length; index++)
            {
                if (tail[index] != tail[tail.Length - period])
                {
                    uniform = false;
                    break;
                }
            }
            if (uniform)
            {
                // A cycle of one repeated character belongs to the run rule instead.
                continue;
            }
            bool repeats = true;
            for (int index = tail.Length - RunawayPhrase; index < tail.Length; index++)
            {
                if (tail[index] != tail[tail.Length - period + (index - (tail.Length - period)) % period])
                {
                    repeats = false;
                    break;
                }
            }
            if (repeats)
            {
                return true;
            }
        }
        return false;
    }
}
