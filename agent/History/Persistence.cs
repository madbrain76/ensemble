using System.Collections.Immutable;
using System.Text;
namespace Ensemble;

public sealed partial class History
{
    private partial SaveDocument ReadSaveFile(string path)
    {
        using var input = Parent().Sandbox().OpenFile(path, FileMode.Open, FileAccess.Read);
        using var reader = new StreamReader(input);
        var document = JsonCodec.ReadSave(reader.ReadToEnd());
        return document;
    }

    public partial bool LoadSave(string path)
    {
        path = Parent().Sandbox().Resolve(path);
        lock (saveGate)
        {
            lock (gate)
            {
                ObjectDisposedException.ThrowIf(destroyed, this);
                if (journalPath is not null)
                {
                    throw new InvalidOperationException("Load before opening History writers");
                }
            }
            // A journal the app itself cannot trust used to abort the launch from
            // outside the candidate loop below: nothing was logged, and no backup
            // was tried, so the actor could never start again. The records before the
            // damage are contiguous and durable, so recovery keeps the snapshot and
            // that prefix, preserves the bytes for inspection, and says plainly which
            // file was rejected and why. Reading stays lazy: only the prefix that the
            // selected snapshot can continue from is ever held in hand.
            string journalFile = path + ".journal.jsonl";
            InvalidDataException? journalFault = null;
            int salvaged = 0;
            IEnumerable<Event> JournalPrefix()
            {
                foreach (var fact in ReadJournalPrefix(journalFile, error => journalFault = error))
                {
                    salvaged++;
                    yield return fact;
                }
            }
            HistorySnapshot? recovered = null;
            bool found = false;
            foreach (string candidate in new[] { path, path + ".bak" })
            {
                try
                {
                    var document = ReadSaveFile(candidate);
                    recovered = ReplayJournal(document.Context, document.Sequence,
                        JournalPrefix(), document.LatestHandoff, document.PreviousHandoff);
                    if (candidate != path)
                    {
                        Parent().Parent().Log($"persistence: recovered backup and journal at sequence {recovered.Sequence}");
                    }
                    break;
                }
                catch (Exception error) when (error is FileNotFoundException or DirectoryNotFoundException)
                {
                    continue;
                }
                catch (Exception error) when (error is IOException or InvalidDataException or UnauthorizedAccessException or FormatException)
                {
                    found = true;
                    Parent().Parent().Log($"persistence: {candidate} rejected "
                        + $"({error.GetType().Name}: {error.Message}); trying backup");
                }
            }
            if (recovered is null)
            {
                if (found)
                {
                    Parent().Parent().Log($"persistence: no valid save candidate under {path}; "
                        + "files left untouched");
                    throw new InvalidDataException("No valid save candidate; files left untouched");
                }
                // With no usable snapshot the journal is the only surviving record, so a
                // damaged tail cannot be moved aside: refuse and leave every file alone.
                InvalidDataException? onlyRecord = null;
                recovered = ReplayJournal(new Context(), 0,
                    ReadJournalPrefix(journalFile, error => onlyRecord = error));
                if (onlyRecord is not null)
                {
                    Parent().Parent().Log($"persistence: {journalFile} rejected "
                        + $"({onlyRecord.GetType().Name}: {onlyRecord.Message}); no saved "
                        + "snapshot remains to continue from; files left untouched");
                    throw onlyRecord;
                }
                if (recovered.Sequence == 0)
                {
                    return false;
                }
            }
            if (journalFault is not null)
            {
                // Continuing requires that the unusable journal stop being the
                // active one: StartJournal rewrites the active file, so the rejected
                // records are moved aside whole under a distinct name first. A failed
                // move must not silently fall through to that rewrite.
                string preserved = $"{journalFile}.rejected-"
                    + DateTimeOffset.UtcNow.ToString("yyyyMMdd-HHmmss");
                try
                {
                    Parent().Sandbox().MoveFile(journalFile, preserved);
                }
                catch (Exception error) when (error is IOException or UnauthorizedAccessException)
                {
                    Parent().Parent().Log($"persistence: {journalFile} rejected "
                        + $"({journalFault.GetType().Name}: {journalFault.Message}) and could not be "
                        + $"preserved as {preserved}; files left untouched");
                    throw journalFault;
                }

                Parent().Parent().Log($"persistence: {journalFile} rejected "
                    + $"({journalFault.GetType().Name}: {journalFault.Message}); preserved as "
                    + $"{preserved}; continued from the saved snapshot at sequence "
                    + $"{recovered.Sequence} with {salvaged} journal record(s) before the damage");
            }
            lock (gate)
            {
                events.Clear();
                events.AddRange(recovered.Events);
                context = recovered.Context;
                foldedSequence = recovered.Sequence;
                latestHandoff = recovered.LatestHandoff;
                previousHandoff = recovered.PreviousHandoff;
            }
            return true;
        }
    }

    internal partial void ConfigureRuntime(ImmutableArray<SkillPart> skills)
    {
        Append(new Event(EventKind.RuntimeConfigured, parts: skills));
    }

    private partial void ReplaceFile(string path, string text)
    {
        Parent().Sandbox().ReplaceFile(path, new UTF8Encoding(false).GetBytes(text));
    }

    internal partial PersistenceOptions NormalizePersistence(PersistenceOptions options)
    {
        ArgumentNullException.ThrowIfNull(options);
        ArgumentException.ThrowIfNullOrWhiteSpace(options.Path);
        var files = Parent().Sandbox();
        string path = files.Resolve(options.Path);
        return new PersistenceOptions(path, options.Audit);
    }

    public partial void Save(PersistenceOptions options)
    {
        options = NormalizePersistence(options);
        lock (saveGate)
        {
            if (appendFaulted)
            {
                throw new IOException("Cannot save after an uncertain History append");
            }
            if (journalPath is not null && journalPath != options.Path + ".journal.jsonl")
            {
                throw new ArgumentException("Save path must match the active recovery journal");
            }
            SaveDocument candidate;
            lock (gate)
            {
                ObjectDisposedException.ThrowIf(destroyed, this);
                // A state-only save does not allocate a copy of working events.
                var configuration = Parent().Engine().SavedConfiguration(context);
                candidate = new SaveDocument(context, foldedSequence, configuration,
                    latestHandoff, previousHandoff);
            }
            string text = JsonCodec.WriteSave(candidate);
            _ = JsonCodec.ReadSave(text); // Validate before touching disk.

            string? previous = null;
            long backupSequence = 0;
            try
            {
                using var input = Parent().Sandbox().OpenFile(options.Path, FileMode.Open, FileAccess.Read);
                using var reader = new StreamReader(input);
                var validated = JsonCodec.ReadSave(reader.ReadToEnd());
                backupSequence = validated.Sequence;
                // Rotate recovered state, not historical encoding. A legacy
                // primary may contain a large folded log or an unfolded tail;
                // decoding preserves the tail and writing removes the duplicate.
                previous = JsonCodec.WriteSave(validated);
            }
            catch (Exception error) when (error is FileNotFoundException or DirectoryNotFoundException or FormatException)
            {
                // Missing or invalid primaries must not replace a good backup.
                previous = null;
            }
            if (previous is not null)
            {
                ReplaceFile(options.Path + ".bak", previous);
            }
            else
            {
                try
                {
                    backupSequence = ReadSaveFile(options.Path + ".bak").Sequence;
                }
                catch (Exception error) when (error is IOException or InvalidDataException or UnauthorizedAccessException or FormatException)
                {
                    // No validated backup boundary: retain the whole journal.
                    backupSequence = 0;
                }
            }
            ReplaceFile(options.Path, text);
            if (!options.Audit)
            {
                TrimJournal(backupSequence);
            }
        }
    }
}
