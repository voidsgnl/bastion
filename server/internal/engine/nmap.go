package engine

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/voidsgnl/bastion/server/internal/store"
)

const (
	loopbackCIDR = "127.0.0.1/32"
	loopbackIP   = "127.0.0.1"
	labPort      = 18080
)

type NmapRunner struct {
	Binary   string
	Timeout  time.Duration
	Executor CommandExecutor
}

type CommandExecutor interface {
	Run(context.Context, string, ...string) ([]byte, error)
}

func NewNmapRunner() NmapRunner {
	return NmapRunner{Binary: "nmap", Timeout: 10 * time.Second}
}

func (r NmapRunner) Run(ctx context.Context, scan store.Scan, target store.Target) ([]store.Finding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if !target.Authorized {
		return nil, fmt.Errorf("target is not authorized")
	}
	if target.Type != "cidr" || target.Value != loopbackCIDR {
		return nil, fmt.Errorf("nmap lab runner only accepts target %s", loopbackCIDR)
	}
	if scan.Profile != "safe" {
		return nil, fmt.Errorf("nmap lab runner only accepts the safe profile")
	}

	timeout := r.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	scanContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	binary := r.Binary
	if binary == "" {
		binary = "nmap"
	}
	args := []string{
		"-sT", "-Pn", "-n", "--open",
		"-p", strconv.Itoa(labPort),
		"--max-retries", "1",
		"--host-timeout", "8s",
		"-oX", "-", loopbackIP,
	}

	executor := r.Executor
	if executor == nil {
		executor = execNmapCommand{}
	}
	output, err := executor.Run(scanContext, binary, args...)
	if err != nil {
		return nil, fmt.Errorf("run nmap lab probe: %w", err)
	}
	return parseNmapXML(output, scan, target)
}

type execNmapCommand struct{}

func (execNmapCommand) Run(ctx context.Context, binary string, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, binary, args...)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err == nil {
		return output, nil
	}
	if ctx.Err() != nil {
		return nil, fmt.Errorf("nmap command canceled or timed out: %w", ctx.Err())
	}
	if details := strings.TrimSpace(stderr.String()); details != "" {
		return nil, fmt.Errorf("nmap command failed: %w: %s", err, details)
	}
	return nil, fmt.Errorf("nmap command failed: %w", err)
}

type nmapDocument struct {
	XMLName xml.Name   `xml:"nmaprun"`
	Hosts   []nmapHost `xml:"host"`
}

type nmapHost struct {
	Ports []nmapPort `xml:"ports>port"`
}

type nmapPort struct {
	Protocol string `xml:"protocol,attr"`
	PortID   int    `xml:"portid,attr"`
	State    struct {
		Name string `xml:"state,attr"`
	} `xml:"state"`
}

func parseNmapXML(output []byte, scan store.Scan, target store.Target) ([]store.Finding, error) {
	var document nmapDocument
	if err := xml.Unmarshal(output, &document); err != nil {
		return nil, fmt.Errorf("parse nmap XML: %w", err)
	}
	if document.XMLName.Local != "nmaprun" {
		return nil, fmt.Errorf("parse nmap XML: unexpected root element %q", document.XMLName.Local)
	}

	findings := make([]store.Finding, 0)
	now := time.Now().UTC()
	for _, host := range document.Hosts {
		for _, port := range host.Ports {
			if port.Protocol != "tcp" || port.PortID != labPort || port.State.Name != "open" {
				continue
			}
			findings = append(findings, store.Finding{
				ID:           fmt.Sprintf("f_%s_tcp_%d", scan.ID, labPort),
				ScanID:       scan.ID,
				AssetID:      "asset_" + target.ID,
				SourceEngine: "nmap",
				Title:        fmt.Sprintf("TCP port %d is open on the local lab target", labPort),
				Category:     "network",
				Severity:     "info",
				Confidence:   "unverified",
				Status:       "open",
				Evidence: []store.Evidence{{
					Type:    "other",
					Summary: fmt.Sprintf("Nmap TCP connect probe found TCP/%d open on %s.", labPort, loopbackIP),
				}},
				Risk: store.Risk{
					Score:       0,
					Explanation: "An open port is an observation, not a vulnerability assessment.",
				},
				FirstSeen: now,
				LastSeen:  now,
			})
		}
	}
	return findings, nil
}
