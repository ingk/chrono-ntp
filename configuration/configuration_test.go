package configuration

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func createTempDir(t *testing.T, name string) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", name)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })
	return tempDir
}

func mustWriteFile(t *testing.T, path string, data []byte, perm os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to mkdirall %s: %v", path, err)
	}
	if err := os.WriteFile(path, data, perm); err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
}

func withUserConfigDir(t *testing.T, dir string) {
	t.Helper()
	old := userConfigDir
	userConfigDir = func() (string, error) { return dir, nil }
	t.Cleanup(func() { userConfigDir = old })
}

// Tests
func TestParseConfiguration(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg, err := parseConfiguration(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Server != defaultConfiguration().Server {
			t.Fatalf("expected default server %q, got %q", defaultConfiguration().Server, cfg.Server)
		}
	})

	t.Run("content", func(t *testing.T) {
		raw := []byte(`
server = "pool.example-time-server.org"
time-zone = "Europe/Berlin"
hide-status-bar = true
hide-date = true
show-time-zone = true
time-format = "12h_AM_PM"
beep-pattern = "greenwich"
offline = true
`)
		cfg, err := parseConfiguration(raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Server != "pool.example-time-server.org" {
			t.Fatalf("bad server: %q", cfg.Server)
		}
		if cfg.TimeZone != "Europe/Berlin" {
			t.Fatalf("bad timezone: %q", cfg.TimeZone)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		if _, err := parseConfiguration([]byte("not valid toml")); err == nil {
			t.Fatal("expected parse error")
		}
	})
}

func TestCanonicalizeDir(t *testing.T) {
	if _, err := canonicalizeDir(""); err == nil {
		t.Fatal("expected error for empty path")
	}
	if _, err := canonicalizeDir("relative/path"); err == nil {
		t.Fatal("expected error for relative path")
	}
	abs := string(filepath.Separator) + "tmp"
	got, err := canonicalizeDir(abs)
	if err != nil {
		t.Fatalf("expected no error for absolute path: %v", err)
	}
	if got != abs {
		t.Fatalf("expected %q, got %q", abs, got)
	}
}

func TestGetConfigurationContents(t *testing.T) {
	tmp := createTempDir(t, "cfg-contents")
	defer os.RemoveAll(tmp)

	path := filepath.Join(tmp, "config.toml")
	mustWriteFile(t, path, []byte("server = \"ok\""), 0o644)

	data, err := getConfigurationContents(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "server = \"ok\"" {
		t.Fatalf("unexpected contents: %s", string(data))
	}

	missing, err := getConfigurationContents(filepath.Join(tmp, "missing.toml"))
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}
	if missing != nil {
		t.Fatal("expected nil for missing file")
	}

	blocked := filepath.Join(tmp, "blocked")
	if err := os.Mkdir(blocked, 0); err != nil {
		t.Fatalf("failed to create blocked dir: %v", err)
	}
	defer func() { os.Chmod(blocked, 0o755) }()

	if _, err := getConfigurationContents(filepath.Join(blocked, "f")); err == nil {
		t.Fatal("expected error for blocked dir")
	}
}

func TestConfigurationPathsAndWriteLoad(t *testing.T) {
	t.Run("search paths include user and home", func(t *testing.T) {
		tmp := createTempDir(t, "cfg-search")
		defer os.RemoveAll(tmp)
		withUserConfigDir(t, tmp)
		t.Setenv("HOME", tmp)

		ps := configurationSearchPaths()
		if len(ps) < 2 {
			t.Fatalf("expected >=2 paths, got %d", len(ps))
		}
	})

	t.Run("config write path falls back to HOME", func(t *testing.T) {
		tmp := createTempDir(t, "cfg-write")
		defer os.RemoveAll(tmp)
		// force userConfigDir to fail
		old := userConfigDir
		userConfigDir = func() (string, error) { return "", errors.New("fail") }
		t.Cleanup(func() { userConfigDir = old })
		t.Setenv("HOME", tmp)

		got, err := configWritePath()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := filepath.Join(tmp, defaultConfigFileName)
		if got != want {
			t.Fatalf("expected %q got %q", want, got)
		}
	})

	t.Run("config write errors without HOME", func(t *testing.T) {
		old := userConfigDir
		userConfigDir = func() (string, error) { return "", errors.New("fail") }
		t.Cleanup(func() { userConfigDir = old })
		os.Unsetenv("HOME")

		if _, err := configWritePath(); err == nil {
			t.Fatal("expected error with no HOME and no user config dir")
		}
	})

	t.Run("write permissions and load fallback", func(t *testing.T) {
		tmp := createTempDir(t, "cfg-writeperm")
		defer os.RemoveAll(tmp)
		withUserConfigDir(t, tmp)
		t.Setenv("HOME", tmp)

		cfg := Configuration{Server: "perm.server", TimeZone: "UTC"}
		path, err := WriteConfiguration(cfg)
		if err != nil {
			t.Fatalf("write error: %v", err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("stat error: %v", err)
		}
		if info.Mode().Perm() != 0o600 {
			t.Fatalf("expected 0600 got %v", info.Mode().Perm())
		}
	})
}

func TestLoadAndWriteEdgeCases(t *testing.T) {
	t.Run("load from configdir precedence", func(t *testing.T) {
		tmp := createTempDir(t, "cfg-load-cfgdir")
		defer os.RemoveAll(tmp)
		t.Setenv("HOME", tmp)
		t.Setenv("XDG_CONFIG_HOME", tmp)

		userCfgDir, err := os.UserConfigDir()
		if err != nil {
			t.Fatalf("userConfigDir failed: %v", err)
		}
		cfgDir := filepath.Join(userCfgDir, appConfigDirName)
		cfgPath := filepath.Join(cfgDir, appConfigFileName)
		mustWriteFile(t, cfgPath, []byte("server = \"configdir.server\"\ntime-zone = \"UTC\"\n"), 0o644)
		mustWriteFile(t, filepath.Join(tmp, defaultConfigFileName), []byte("server = \"home.server\"\n"), 0o644)

		cfg, err := LoadConfiguration()
		if err != nil {
			t.Fatalf("load error: %v", err)
		}
		if cfg.Server != "configdir.server" {
			t.Fatalf("expected configdir.server got %q", cfg.Server)
		}
	})

	t.Run("load fallback from HOME", func(t *testing.T) {
		tmp := createTempDir(t, "cfg-load-home")
		defer os.RemoveAll(tmp)
		t.Setenv("HOME", tmp)
		mustWriteFile(t, filepath.Join(tmp, defaultConfigFileName), []byte("server = \"mocked.server\"\ntime-zone = \"UTC\"\n"), 0o644)

		cfg, err := LoadConfiguration()
		if err != nil {
			t.Fatalf("load error: %v", err)
		}
		if cfg.Server != "mocked.server" {
			t.Fatalf("expected mocked.server got %q", cfg.Server)
		}
	})

	t.Run("load returns error for invalid toml", func(t *testing.T) {
		tmp := createTempDir(t, "cfg-load-invalid")
		defer os.RemoveAll(tmp)
		t.Setenv("HOME", tmp)
		mustWriteFile(t, filepath.Join(tmp, defaultConfigFileName), []byte("invalid-toml"), 0o644)

		if _, err := LoadConfiguration(); err == nil {
			t.Fatal("expected error for invalid toml")
		}
	})

	t.Run("write fallback and load", func(t *testing.T) {
		tmp := createTempDir(t, "cfg-write-fallback")
		defer os.RemoveAll(tmp)
		t.Setenv("HOME", tmp)
		old := userConfigDir
		userConfigDir = func() (string, error) { return "", errors.New("fail") }
		t.Cleanup(func() { userConfigDir = old })

		cfg := Configuration{Server: "fallback.server", TimeZone: "UTC"}
		got, err := WriteConfiguration(cfg)
		if err != nil {
			t.Fatalf("write failed: %v", err)
		}
		if got != filepath.Join(tmp, defaultConfigFileName) {
			t.Fatalf("unexpected write path: %q", got)
		}
		loaded, err := LoadConfiguration()
		if err != nil {
			t.Fatalf("load failed: %v", err)
		}
		if loaded.Server != "fallback.server" {
			t.Fatalf("expected fallback.server got %q", loaded.Server)
		}
	})

	t.Run("write errors when configdir is blocked", func(t *testing.T) {
		tmp := createTempDir(t, "cfg-write-blocked")
		defer os.RemoveAll(tmp)
		t.Setenv("HOME", tmp)
		t.Setenv("XDG_CONFIG_HOME", tmp)

		userCfgDir, err := os.UserConfigDir()
		if err != nil {
			t.Fatalf("user configdir failed: %v", err)
		}
		blocked := filepath.Join(userCfgDir, appConfigDirName)
		// create a file where a directory is expected
		if err := os.MkdirAll(filepath.Dir(blocked), 0o755); err != nil {
			t.Fatalf("mkdirall failed: %v", err)
		}
		if err := os.WriteFile(blocked, []byte("not-a-dir"), 0o644); err != nil {
			t.Fatalf("write failed: %v", err)
		}

		if _, err := WriteConfiguration(Configuration{}); err == nil {
			t.Fatal("expected write error when dir blocked")
		}
	})

	t.Run("toml marshal error bubbles up", func(t *testing.T) {
		old := tomlMarshal
		tomlMarshal = func(v interface{}) ([]byte, error) { return nil, errors.New("marshal fail") }
		t.Cleanup(func() { tomlMarshal = old })

		tmp := createTempDir(t, "cfg-marshal-fail")
		defer os.RemoveAll(tmp)
		withUserConfigDir(t, tmp)
		t.Setenv("HOME", tmp)

		if _, err := WriteConfiguration(Configuration{Server: "s"}); err == nil {
			t.Fatal("expected marshal error")
		}
	})

	t.Run("home mkdirall error", func(t *testing.T) {
		tmp := createTempDir(t, "cfg-home-mkdirfail")
		defer os.RemoveAll(tmp)
		// create a file at HOME so MkdirAll on files fails
		homeFile := filepath.Join(tmp, "homefile")
		mustWriteFile(t, homeFile, []byte("block"), 0o644)
		old := userConfigDir
		userConfigDir = func() (string, error) { return "", errors.New("fail") }
		t.Cleanup(func() { userConfigDir = old })
		t.Setenv("HOME", homeFile)

		if _, err := configWritePath(); err == nil {
			t.Fatal("expected error when HOME path not writable")
		}
	})

	t.Run("load returns error when config path unreadable", func(t *testing.T) {
		tmp := createTempDir(t, "cfg-blocked-read")
		defer os.RemoveAll(tmp)
		blocked := filepath.Join(tmp, "chrono-ntp")
		if err := os.WriteFile(blocked, []byte("not-a-dir"), 0o644); err != nil {
			t.Fatalf("write failed: %v", err)
		}
		withUserConfigDir(t, tmp)
		t.Setenv("HOME", tmp)

		if _, err := LoadConfiguration(); err == nil {
			t.Fatal("expected error when config path unreadable")
		}
	})
}

func TestLoadConfiguration_NoFiles(t *testing.T) {
	tmp := createTempDir(t, "cfg-nofiles")
	defer os.RemoveAll(tmp)

	// Ensure userConfigDir returns a valid dir but no config files exist.
	withUserConfigDir(t, tmp)
	t.Setenv("HOME", tmp)

	cfg, err := LoadConfiguration()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Server != defaultConfiguration().Server {
		t.Fatalf("expected default server %q, got %q", defaultConfiguration().Server, cfg.Server)
	}
}
