package test

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/moby/buildkit-bench/util/testutil"
	"github.com/stretchr/testify/require"
)

func TestDaemon(t *testing.T) {
	testutil.Run(t, testutil.TestFuncs(
		testDaemonVersion,
	))
}

func BenchmarkDaemon(b *testing.B) {
	testutil.Run(b, testutil.BenchFuncs(
		benchmarkDaemonVersion,
		benchmarkDaemonSize,
	))
}

func testDaemonVersion(t *testing.T, sb testutil.Sandbox) {
	buildkitdPath := path.Join(sb.BinsDir(), sb.Name(), "buildkitd")

	output, err := exec.Command(buildkitdPath, "--version").Output()
	require.NoError(t, err)

	versionParts := strings.Fields(string(output))
	require.Len(t, versionParts, 4)
	require.Equal(t, "buildkitd", versionParts[0])
	t.Log("repo:", versionParts[1])
	t.Log("version:", versionParts[2])
	t.Log("commit:", versionParts[3])
}

func benchmarkDaemonVersion(b *testing.B, sb testutil.Sandbox) {
	for i := 0; i < b.N; i++ {
		buildkitdPath := path.Join(sb.BinsDir(), sb.Name(), "buildkitd")
		start := time.Now()
		require.NoError(b, exec.Command(buildkitdPath, "--version").Run())
		testutil.ReportMetricDuration(b, time.Since(start))
	}
}

func benchmarkDaemonSize(b *testing.B, sb testutil.Sandbox) {
	for i := 0; i < b.N; i++ {
		buildkitdPath := path.Join(sb.BinsDir(), sb.Name(), "buildkitd")
		fi, err := os.Stat(buildkitdPath)
		require.NoError(b, err)
		testutil.ReportMetric(b, float64(fi.Size()), testutil.MetricBytes)
	}
}