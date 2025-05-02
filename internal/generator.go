package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/joho/godotenv"
)

type TemplateData struct {
	Domain      string
	Prefix      string
	ProjectIP   string
	NetworkName string
}

func GenerateProject(domain string) error {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	network := os.Getenv("NETWORK_NAME")
	baseIP := os.Getenv("PROJECT_START_IP")
	ipmapPath := ".ipmap.env"
	templateDir := "templates"
	projectDir := filepath.Join("domains", domain)
	prefix := strings.Split(domain, ".")[0]
	hostsPath := "/mnt/c/Windows/System32/drivers/etc/hosts"
	reverseProxyDir := "shared-services"
	reverseProxyName := "nginx-reverse-proxy"

	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		return fmt.Errorf("project already exists: %s", projectDir)
	}
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return err
	}

	usedIPs, err := LoadUsedIPs(ipmapPath)
	if err != nil {
		return err
	}
	projectIP, err := FindNextFreeIP(baseIP, usedIPs)
	if err != nil {
		return err
	}
	if err := AppendIPMapping(ipmapPath, domain, projectIP); err != nil {
		return err
	}

	data := TemplateData{
		Domain: domain, Prefix: prefix, ProjectIP: projectIP, NetworkName: network,
	}

	if err := RenderTemplate(
		filepath.Join(templateDir, "docker-compose.yml.tmpl"),
		filepath.Join(projectDir, "docker-compose.yml"),
		data,
	); err != nil {
		return err
	}

	confDir := filepath.Join(projectDir, "conf", "nginx")
	if err := os.MkdirAll(confDir, 0755); err != nil {
		return err
	}
	if err := RenderTemplate(
		filepath.Join(templateDir, "nginx.conf.tmpl"),
		filepath.Join(confDir, "default.conf"),
		data,
	); err != nil {
		return err
	}

	appDstDir := filepath.Join(projectDir, "app")
	if err := os.MkdirAll(appDstDir, 0755); err != nil {
		return err
	}
	if err := RenderTemplate(
		filepath.Join(templateDir, "app", "index.html"),
		filepath.Join(appDstDir, "index.html"),
		data,
	); err != nil {
		return err
	}

	sitesDir := filepath.Join(reverseProxyDir, "sites")
	if err := os.MkdirAll(sitesDir, 0755); err != nil {
		return err
	}
	siteConf := filepath.Join(sitesDir, domain+".conf")
	if _, err := os.Stat(siteConf); os.IsNotExist(err) {
		if err := RenderTemplate(
			filepath.Join(templateDir, "site.conf.tmpl"),
			siteConf,
			data,
		); err != nil {
			return err
		}
		fmt.Println("Created reverse proxy config:", siteConf)
	} else {
		fmt.Println("Reverse proxy config already exists:", siteConf)
	}

	fmt.Println("Starting shared-services...")
	if err := runDockerComposeUp(reverseProxyDir); err != nil {
		return err
	}

	fmt.Println("Starting project containers...")
	if err := runDockerComposeUp(projectDir); err != nil {
		return err
	}

	fmt.Println("Reloading reverse proxy config...")
	err = exec.Command("docker", "exec", reverseProxyName, "nginx", "-s", "reload").Run()
	if err != nil {
		fmt.Println("Reload failed, restarting container...")
		err = exec.Command("docker", "restart", reverseProxyName).Run()
		if err != nil {
			return fmt.Errorf("failed to reload or restart reverse proxy: %w", err)
		}
	}

	if err := updateWindowsHosts(domain, hostsPath); err != nil {
		return err
	}

	return nil
}

func DeleteProject(domain string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Are you sure you want to delete domain '%s'? [y/N]: ", domain)
	answer, _ := reader.ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer != "y" {
		fmt.Println("Aborted.")
		return
	}

	projectPath := filepath.Join("domains", domain)
	composeFile := filepath.Join(projectPath, "docker-compose.yml")

	if _, err := os.Stat(composeFile); err == nil {
		fmt.Println("Stopping containers for", domain, "...")
		stopCmd := exec.Command("docker", "compose", "down")
		stopCmd.Dir = projectPath
		stopCmd.Stdout = os.Stdout
		stopCmd.Stderr = os.Stderr
		if err := stopCmd.Run(); err != nil {
			fmt.Println("Warning: failed to stop containers:", err)
		}
	}

	if err := os.RemoveAll(projectPath); err != nil {
		fmt.Println("Failed to remove domain folder:", err)
	} else {
		fmt.Println("Deleted folder:", projectPath)
	}

	ipmap := ".ipmap.env"
	lines, err := os.ReadFile(ipmap)
	if err == nil {
		var kept []string
		for _, line := range strings.Split(string(lines), "\n") {
			if !strings.HasPrefix(line, domain+"=") {
				kept = append(kept, line)
			}
		}
		os.WriteFile(ipmap, []byte(strings.Join(kept, "\n")), 0644)
		fmt.Println("Updated:", ipmap)
	}

	sitePath := filepath.Join("shared-services", "sites", domain+".conf")
	if err := os.Remove(sitePath); err == nil {
		fmt.Println("Removed reverse proxy config:", sitePath)
	}

	hosts := "/mnt/c/Windows/System32/drivers/etc/hosts"
	hfile, err := os.ReadFile(hosts)
	if err == nil {
		var out []string
		for _, line := range strings.Split(string(hfile), "\n") {
			if !strings.Contains(line, domain) {
				out = append(out, line)
			}
		}
		os.WriteFile(hosts, []byte(strings.Join(out, "\n")), 0644)
		fmt.Println("Updated Windows hosts file.")
	}

	fmt.Println("Domain", domain, "was successfully deleted.")
}

func runDockerComposeUp(dir string) error {
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func updateWindowsHosts(domain, path string) error {
	hostsEntry := fmt.Sprintf("127.0.0.1 %s", domain)

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), domain) {
			fmt.Println("Hosts entry already exists.")
			return nil
		}
	}

	if _, err := file.WriteString("\n" + hostsEntry + "\n"); err != nil {
		return err
	}

	fmt.Println("Domain added to Windows hosts file.")
	return nil
}
