package ensemble

import (
	"context"
	"ensemble/internal/common"
	"ensemble/internal/mcp"
	"io"
	"net"
	"sync"
)

type MCPServer = common.MCPServer
type MCPTransport = common.MCPTransport
type MCPCallResult = common.MCPCallResult
type MCPContent = common.MCPContent
type SplitReadWriteCloser = common.SplitReadWriteCloser

// NewMCPTransport adapts split streams with the shared newline codec.
func NewMCPTransport(reader io.Reader, writer io.Writer, closer io.Closer) MCPTransport {
	var closeFunc func() error
	if closer != nil {
		closeFunc = closer.Close
	}
	return mcp.NewStream(SplitReadWriteCloser{Reader: reader, Writer: writer, CloseFunc: closeFunc})
}

// MCPPipe joins two MCP peers without a listener.
func MCPPipe() (MCPTransport, MCPTransport) {
	left, right := net.Pipe()
	return mcp.NewStream(left), mcp.NewStream(right)
}

// NewMCPStream adapts one full-duplex stream or a split stdin/stdout pair. The
// returned transport owns closure; callers must not mix protocol and UI bytes.
func NewMCPStream(rw io.ReadWriteCloser) MCPTransport { return mcp.NewStream(rw) }
func ServeMCP(ctx context.Context, a Agent, t MCPTransport) error {
	return a.ServeMCPConnection(ctx, t)
}

// ListenMCP serves dedicated TCP sessions until cancellation. It does not start
// the GUI HTTP server or consume conversational stdin.
func ListenMCP(ctx context.Context, a Agent, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	stop := context.AfterFunc(ctx, func() { _ = listener.Close() })
	defer stop()
	var workers sync.WaitGroup
	defer workers.Wait()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		workers.Add(1)
		go func() {
			defer workers.Done()
			if err := ServeMCP(ctx, a, NewMCPStream(conn)); err != nil && ctx.Err() == nil {
				a.Parent().Logger().Printf("MCP connection: %v", err)
			}
		}()
	}
}
