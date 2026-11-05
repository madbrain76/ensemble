namespace Ensemble;

public sealed partial class Agent
{
    internal partial bool HasDataDirectory()
    {
        lock (ingress)
        {
            ObjectDisposedException.ThrowIf(destroyed, this);
            return dataDirectory is not null;
        }
    }

    private partial void AcquireWriterLock(string directory)
    {
        if (writerLock is not null)
        {
            return;
        }

        // Each launcher numbers its own events from the sequence it restored. Two
        // processes appending one recovery journal therefore write duplicate and
        // out-of-order records, and every later load of that actor then fails with
        // no way to start. The second writer refuses while the first is alive.
        // The operating system drops this handle when its owner exits, so a killed
        // launcher cannot strand the actor home behind a stale lock file.
        var files = Sandbox();
        files.CreateDirectory(directory);
        string path = files.Resolve(Path.Combine(directory, ".writer.lock"));
        try
        {
            writerLock = new FileStream(path, new FileStreamOptions
            {
                Mode = FileMode.OpenOrCreate,
                Access = FileAccess.ReadWrite,
                Share = FileShare.None,
            });
        }
        catch (IOException error)
        {
            throw new WriterConflictException(
                "Another Ensemble is already using actor home " + directory
                + "; close it before starting this one.", error);
        }
    }

    private partial void ReleaseWriterLock()
    {
        writerLock?.Dispose();
        writerLock = null;
    }

    public partial string DataDirectory()
    {
        lock (ingress)
        {
            ObjectDisposedException.ThrowIf(destroyed, this);
            return dataDirectory ?? throw new InvalidOperationException(
                "Construct Agent with a data directory before using child storage.");
        }
    }
}
