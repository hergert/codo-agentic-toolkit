package doctor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Item struct {
	Label  string
	Detail string
}

type Summary struct {
	Items []Item
}

func Collect(root string) Summary {
	items := []Item{
		checkGo(),
		checkGoImports(),
		checkNode(),
		checkTSDeps(root),
		checkPython3(),
		checkPythonTools(),
		checkDart(),
	}
	return Summary{Items: items}
}

func checkGo() Item {
	const label = "Go"
	path, err := exec.LookPath("go")
	if err != nil {
		return Item{Label: label, Detail: "missing `go` (install Go)"}
	}
	out, err := exec.Command(path, "version").Output()
	if err != nil {
		return Item{Label: label, Detail: "ok (`go` available)"}
	}
	return Item{Label: label, Detail: fmt.Sprintf("ok (%s)", strings.TrimSpace(string(out)))}
}

func checkGoImports() Item {
	const label = "Go"
	if _, err := exec.LookPath("goimports"); err != nil {
		return Item{
			Label:  label,
			Detail: "optional `goimports` for format (`go install golang.org/x/tools/cmd/goimports@latest`)",
		}
	}
	return Item{Label: label, Detail: "`goimports` available"}
}

func checkNode() Item {
	const label = "Node"
	nodePath, nodeErr := exec.LookPath("node")
	_, npxErr := exec.LookPath("npx")
	if nodeErr != nil && npxErr != nil {
		return Item{Label: label, Detail: "missing `node`/`npx`"}
	}
	if nodeErr != nil {
		return Item{Label: label, Detail: "missing `node` (install Node.js)"}
	}
	if npxErr != nil {
		return Item{Label: label, Detail: "missing `npx` (install Node.js >= 8)"}
	}
	out, err := exec.Command(nodePath, "--version").Output()
	if err != nil {
		return Item{Label: label, Detail: "ok (`node` available)"}
	}
	return Item{Label: label, Detail: fmt.Sprintf("ok (%s)", strings.TrimSpace(string(out)))}
}

func checkTSDeps(root string) Item {
	const label = "TS/JS devDeps (project)"
	pkgPath := filepath.Join(root, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Item{Label: label, Detail: "package.json not found (skipped)"}
		}
		return Item{Label: label, Detail: fmt.Sprintf("error reading package.json: %v", err)}
	}
	var pkg struct {
		Dependencies    map[string]interface{} `json:"dependencies"`
		DevDependencies map[string]interface{} `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return Item{Label: label, Detail: fmt.Sprintf("error parsing package.json: %v", err)}
	}
	hasDep := func(name string) bool {
		if pkg.Dependencies != nil {
			if _, ok := pkg.Dependencies[name]; ok {
				return true
			}
		}
		if pkg.DevDependencies != nil {
			if _, ok := pkg.DevDependencies[name]; ok {
				return true
			}
		}
		return false
	}
	missing := []string{}
	if !hasDep("prettier") {
		missing = append(missing, "prettier")
	}
	if !hasDep("eslint") {
		missing = append(missing, "eslint")
	}
	tsconfigPath := filepath.Join(root, "tsconfig.json")
	if fileExists(tsconfigPath) && !hasDep("typescript") {
		missing = append(missing, "typescript")
	}
	if len(missing) > 0 {
		return Item{Label: label, Detail: "missing devDeps: " + strings.Join(missing, ", ")}
	}
	return Item{Label: label, Detail: "ok"}
}

func checkPython3() Item {
	const label = "Python3"
	path, err := exec.LookPath("python3")
	if err != nil {
		return Item{Label: label, Detail: "missing `python3`"}
	}
	out, err := exec.Command(path, "--version").Output()
	if err != nil {
		return Item{Label: label, Detail: "ok (`python3` available)"}
	}
	return Item{Label: label, Detail: fmt.Sprintf("ok (%s)", strings.TrimSpace(string(out)))}
}

func checkPythonTools() Item {
	const label = "Python"
	missing := []string{}
	if _, err := exec.LookPath("ruff"); err != nil {
		missing = append(missing, "`ruff`")
	}
	if _, err := exec.LookPath("black"); err != nil {
		missing = append(missing, "`black`")
	}
	if len(missing) == 0 {
		return Item{Label: label, Detail: "optional formatters available (`ruff`, `black`)"}
	}
	return Item{Label: label, Detail: "optional " + strings.Join(missing, ", ")}
}

func checkDart() Item {
	const label = "Dart"
	if _, err := exec.LookPath("dart"); err != nil {
		return Item{Label: label, Detail: "missing `dart`"}
	}
	out, err := exec.Command("dart", "--version").CombinedOutput()
	if err != nil {
		return Item{Label: label, Detail: "ok (`dart` available)"}
	}
	return Item{Label: label, Detail: fmt.Sprintf("ok (%s)", strings.TrimSpace(string(out)))}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
