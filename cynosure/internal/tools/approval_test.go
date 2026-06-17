package tools

import "testing"

func TestRequiresApprovalReadOnlyTools(t *testing.T) {
	readOnly := []string{"read_file", "grep", "glob", "ls", "web_search", "web_fetch", "load_skill", "todo_write", "todo_list", "read_persisted_output", "spawn_subagent"}
	for _, name := range readOnly {
		need, _ := RequiresApproval(name, map[string]any{})
		if need {
			t.Errorf("tool %s should not require approval", name)
		}
	}
}

func TestRequiresApprovalWriteTools(t *testing.T) {
	cases := map[string]string{
		"write_file": "write_file *",
		"edit_file":  "edit_file *",
		"multi_edit": "multi_edit *",
	}
	for name, wantRule := range cases {
		need, rule := RequiresApproval(name, map[string]any{})
		if !need {
			t.Errorf("tool %s should require approval", name)
		}
		if rule != wantRule {
			t.Errorf("tool %s rule = %q, want %q", name, rule, wantRule)
		}
	}
}

func TestRequiresApprovalBashReadOnly(t *testing.T) {
	readOnly := []string{
		"ls -la",
		"cat a.txt",
		"grep foo .",
		"ls && pwd",
		"head -n 5 file; tail -n 5 file",
		"FOO=bar env",
		"/bin/ls /tmp",
	}
	for _, command := range readOnly {
		need, _ := RequiresApproval("bash", map[string]any{"command": command})
		if need {
			t.Errorf("bash %q should NOT require approval", command)
		}
	}
}

func TestRequiresApprovalBashMutating(t *testing.T) {
	cases := []struct {
		command  string
		wantRule string
	}{
		{"curl -s https://x.com", "curl *"},
		{"/usr/bin/python script.py", "python *"},
		{"FOO=bar curl https://x.com", "curl *"},
		{"rm -rf build", "rm *"},
		{"ls && curl x", "ls *"},
		{"cat a > b", "cat *"},
		{"echo hi; rm x", "echo *"},
		{"git commit -m x", "git *"},
		{"sed -i s/a/b/ f", "sed *"},
		{"foobar baz", "foobar *"},
		{"", "bash *"},
	}
	for _, c := range cases {
		need, rule := RequiresApproval("bash", map[string]any{"command": c.command})
		if !need {
			t.Errorf("bash %q should require approval", c.command)
		}
		if rule != c.wantRule {
			t.Errorf("bash %q rule = %q, want %q", c.command, rule, c.wantRule)
		}
	}
}
