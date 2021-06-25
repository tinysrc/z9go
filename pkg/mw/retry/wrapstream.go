package retry

import (
	"context"
	"io"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type wrappedClientStream struct {
	grpc.ClientStream
	buffers      []interface{}
	isClosedSend bool
	parentCtx    context.Context
	callOpts     *options
	streamer     func(ctx context.Context) (grpc.ClientStream, error)
	mu           sync.RWMutex
}

func (s *wrappedClientStream) setStream(clientStream grpc.ClientStream) {
	s.mu.Lock()
	s.ClientStream = clientStream
	s.mu.Unlock()
}

func (s *wrappedClientStream) getStream() grpc.ClientStream {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ClientStream
}

func (s *wrappedClientStream) SendMsg(m interface{}) error {
	s.mu.Lock()
	s.buffers = append(s.buffers, m)
	s.mu.Unlock()
	return s.getStream().SendMsg(m)
}

func (s *wrappedClientStream) CloseSend() error {
	s.mu.Lock()
	s.isClosedSend = true
	s.mu.Unlock()
	return s.getStream().CloseSend()
}

func (s *wrappedClientStream) Header() (metadata.MD, error) {
	return s.getStream().Header()
}

func (s *wrappedClientStream) Trailer() metadata.MD {
	return s.getStream().Trailer()
}

func (s *wrappedClientStream) RecvMsg(m interface{}) error {
	retry, err := s.receiveMsgAndIndicateRetry(m)
	if !retry {
		return err
	}
	for i := uint(1); i < s.callOpts.max; i++ {
		if err := waitRetryBackoff(i, s.parentCtx, s.callOpts); err != nil {
			return err
		}
		ctx := callContext(s.parentCtx, s.callOpts, i)
		stream, err := s.reestablishStreamAndResendBuffer(ctx)
		if err != nil {
			if isRetriable(err, s.callOpts) {
				continue
			}
			return err
		}
		s.setStream(stream)
		retry, err = s.receiveMsgAndIndicateRetry(m)
		if !retry {
			return err
		}
	}
	return err
}

func (s *wrappedClientStream) receiveMsgAndIndicateRetry(m interface{}) (bool, error) {
	err := s.getStream().RecvMsg(m)
	if err == nil || err == io.EOF {
		return false, err
	}
	if isCtxErr(err) {
		if s.parentCtx.Err() != nil {
			logTrace(s.parentCtx, "grpc retry context error=%v", s.parentCtx.Err())
			return false, err
		} else if s.callOpts.timeout != 0 {
			logTrace(s.parentCtx, "grpc retry context error from retry call")
			return true, err
		}
	}
	return isRetriable(err, s.callOpts), err
}

func (s *wrappedClientStream) reestablishStreamAndResendBuffer(ctx context.Context) (grpc.ClientStream, error) {
	s.mu.RLock()
	buffers := s.buffers
	s.mu.RUnlock()
	stream, err := s.streamer(ctx)
	if err != nil {
		logTrace(ctx, "grpc retry new stream failed error=%v", err)
		return nil, err
	}
	for _, m := range buffers {
		if err := stream.SendMsg(m); err != nil {
			logTrace(ctx, "grpc retry resend failed error=%v", err)
			return nil, err
		}
	}
	if err := stream.CloseSend(); err != nil {
		logTrace(ctx, "grpc retry close send failed error=%v", err)
		return nil, err
	}
	return stream, nil
}
