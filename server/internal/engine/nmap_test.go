package engine

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/voidsgnl/bastion/server/internal/store"
)

const nmapOpenPortXML = `<?xml version="1.0"?>
<nmaprun>
  <host>
    <ports>
      <port protocol="tcp" portid="18080"><state state="open" reason="syn-ack"/></port>
      <port protocol="tcp" portid="443"><state state="open" reason="syn-ack"/></port>
      <port protocol="tcp" portid="18080"><state state="closed" reason="reset"/></port>
      <port protocol="udp" portid="18080"><state state="open" reason="udp-response"/></port>
    </ports>
  </host>
</nmaprun>`

type fakeNmapExecutor struct {
	binary string
	args   []string
	output []byte
	err    error
	called bool
}

func (f *fakeNmapExecutor) Run(_ context.Context, binary string, args ...string) ([]byte, error) {
	f.called = true
	f.binary = binary
	f.args = append([]string(nil), args...)
	return f.output, f.err
}

func validLabInput() (store.Scan, store.Target) {
	return store.Scan{ID: "scan_test", Profile: "safe"}, store.Target{
		ID: "target_test", Type: "cidr", Value: loopbackCIDR, Authorized: true,
	}
}

func TestNmapRunnerUsesFixedLoopbackArguments(t *testing.T) {
	scan, target := validLabInput()
	executor := &fakeNmapExecutor{output: []byte(nmapOpenPortXML)}
	runner := NmapRunner{
		Binary:   "nmap-test",
		Executor: executor,
	}

	findings, err := runner.Run(context.Background(), scan, target)
	if err != nil {
		t.Fatal(err)
	}
	wantArgs := []string{
		"-sT", "-Pn", "-n", "--open",
		"-p", "18080",
		"--max-retries", "1",
		"--host-timeout", "8s",
		"-oX", "-", "127.0.0.1",
	}
	if !executor.called || executor.binary != "nmap-test" || !reflect.DeepEqual(executor.args, wantArgs) {
		t.Fatalf("command = %q %q, want %q %q", executor.binary, executor.args, "nmap-test", wantArgs)
	}
	if len(findings) != 1 {
		t.Fatalf("findings = %d, want 1", len(findings))
	}
	if findings[0].Severity != "info" || findings[0].Confidence != "unverified" {
		t.Fatalf("finding classification = %s/%s, want info/unverified", findings[0].Severity, findings[0].Confidence)
	}
	if findings[0].SourceEngine != "nmap" {
		t.Fatalf("source engine = %q, want nmap", findings[0].SourceEngine)
	}
}

func TestNmapRunnerRejectsOutOfScopeInputsBeforeExecution(t *testing.T) {
	tests := []struct {
		name   string
		scan   store.Scan
		target store.Target
	}{
		{
			name: "public network",
			scan: store.Scan{Profile: "safe"},
			target: store.Target{
				Type: "cidr", Value: "203.0.113.0/24", Authorized: true,
			},
		},
		{
			name: "non-cidr",
			scan: store.Scan{Profile: "safe"},
			target: store.Target{
				Type: "domain", Value: "localhost", Authorized: true,
			},
		},
		{
			name: "non-safe profile",
			scan: store.Scan{Profile: "full"},
			target: store.Target{
				Type: "cidr", Value: loopbackCIDR, Authorized: true,
			},
		},
		{
			name: "unattested",
			scan: store.Scan{Profile: "safe"},
			target: store.Target{
				Type: "cidr", Value: loopbackCIDR,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			executor := &fakeNmapExecutor{}
			runner := NmapRunner{
				Executor: executor,
			}
			if _, err := runner.Run(context.Background(), test.scan, test.target); err == nil {
				t.Fatal("expected input to be rejected")
			}
			if executor.called {
				t.Fatal("nmap command executed for rejected input")
			}
		})
	}
}

func TestParseNmapXMLRejectsMalformedOutput(t *testing.T) {
	scan, target := validLabInput()
	if _, err := parseNmapXML([]byte("<nmaprun>"), scan, target); err == nil {
		t.Fatal("expected malformed XML to fail")
	}
}

func TestNmapRunnerHonorsCanceledContextWithoutExecution(t *testing.T) {
	scan, target := validLabInput()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	executor := &fakeNmapExecutor{}
	runner := NmapRunner{
		Executor: executor,
	}
	if _, err := runner.Run(ctx, scan, target); !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	if executor.called {
		t.Fatal("nmap command executed after cancellation")
	}
}

func TestNmapRunnerSurfacesCommandFailure(t *testing.T) {
	scan, target := validLabInput()
	runner := NmapRunner{
		Executor: &fakeNmapExecutor{err: errors.New("nmap unavailable")},
	}
	if _, err := runner.Run(context.Background(), scan, target); err == nil {
		t.Fatal("expected command error")
	}
}

func TestNmapRunnerAppliesTimeout(t *testing.T) {
	scan, target := validLabInput()
	runner := NmapRunner{
		Timeout:  time.Nanosecond,
		Executor: timeoutExecutor{},
	}
	if _, err := runner.Run(context.Background(), scan, target); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context deadline exceeded", err)
	}
}

type timeoutExecutor struct{}

func (timeoutExecutor) Run(ctx context.Context, _ string, _ ...string) ([]byte, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}
