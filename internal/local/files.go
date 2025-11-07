package local

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// UpdateEnvFile updates the .env file with generated secrets
func UpdateEnvFile(envPath string, secrets *Secrets, ports *Ports) error {
	// Read the file
	content, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	text := string(content)

	// Update secrets
	text = replaceEnvVar(text, "POSTGRES_PASSWORD", secrets.PostgresPassword)
	text = replaceEnvVar(text, "JWT_SECRET", secrets.JWTSecret)
	text = replaceEnvVar(text, "ANON_KEY", secrets.AnonKey)
	text = replaceEnvVar(text, "SERVICE_ROLE_KEY", secrets.ServiceRoleKey)
	text = replaceEnvVar(text, "DASHBOARD_USERNAME", "supabase")
	text = replaceEnvVar(text, "DASHBOARD_PASSWORD", secrets.DashboardPassword)
	text = replaceEnvVar(text, "VAULT_ENC_KEY", secrets.VaultEncKey)

	// Update ports
	text = replaceEnvVar(text, "KONG_HTTP_PORT", fmt.Sprintf("%d", ports.API))
	text = replaceEnvVar(text, "KONG_HTTPS_PORT", fmt.Sprintf("%d", ports.KongHTTPS))
	text = replaceEnvVar(text, "POSTGRES_PORT", fmt.Sprintf("%d", ports.DB))

	// Write back
	if err := os.WriteFile(envPath, []byte(text), 0600); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	return nil
}

// replaceEnvVar replaces an environment variable value in the text
func replaceEnvVar(text, key, value string) string {
	pattern := fmt.Sprintf(`^%s=.*`, regexp.QuoteMeta(key))
	re := regexp.MustCompile("(?m)" + pattern)
	replacement := fmt.Sprintf("%s=%s", key, value)
	return re.ReplaceAllString(text, replacement)
}

// UpdateDockerComposeFile updates the docker-compose.yml file with project-specific configuration
func UpdateDockerComposeFile(composePath, projectID string, ports *Ports) error {
	// Read the file
	content, err := os.ReadFile(composePath)
	if err != nil {
		return fmt.Errorf("failed to read docker-compose.yml: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var updatedLines []string

	for _, line := range lines {
		// Update container names
		if strings.Contains(line, "container_name:") {
			// Extract indentation and existing name
			parts := strings.SplitN(line, "container_name:", 2)
			if len(parts) == 2 {
				indent := parts[0]
				existingName := strings.TrimSpace(parts[1])
				// Prepend project ID
				newName := fmt.Sprintf("%s-%s", projectID, existingName)
				line = fmt.Sprintf("%scontainer_name: %s", indent, newName)
			}
		}

		// Update port mappings
		line = updatePortMapping(line, 8000, ports.API)
		line = updatePortMapping(line, 5432, ports.DB)
		line = updatePortMapping(line, 3000, ports.Studio)
		line = updatePortMapping(line, 9000, ports.Inbucket)
		line = updatePortMapping(line, 4000, ports.Analytics)
		line = updatePortMapping(line, 8443, ports.KongHTTPS)

		updatedLines = append(updatedLines, line)
	}

	// Write back
	result := strings.Join(updatedLines, "\n")
	if err := os.WriteFile(composePath, []byte(result), 0644); err != nil {
		return fmt.Errorf("failed to write docker-compose.yml: %w", err)
	}

	return nil
}

// updatePortMapping updates a port mapping in a docker-compose line
func updatePortMapping(line string, containerPort, hostPort int) string {
	// Match patterns like "- 8000:8000" or "  - '8000:8000'" or "- \"8000:8000\""
	pattern := fmt.Sprintf(`^(\s*-\s*['"]?)\d+(:%d['"]?.*)$`, containerPort)
	re := regexp.MustCompile(pattern)
	if re.MatchString(line) {
		replacement := fmt.Sprintf("${1}%d${2}", hostPort)
		return re.ReplaceAllString(line, replacement)
	}
	return line
}

// UpdateConfigToml updates the config.toml file with project-specific configuration
func UpdateConfigToml(configPath, projectID string, ports *Ports) error {
	// Read the file
	file, err := os.Open(configPath)
	if err != nil {
		// config.toml might not exist in newer versions, which is fine
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read config.toml: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	currentSection := ""

	for scanner.Scan() {
		line := scanner.Text()

		// Track current section
		if strings.HasPrefix(strings.TrimSpace(line), "[") {
			sectionMatch := regexp.MustCompile(`\[(.*?)\]`).FindStringSubmatch(line)
			if len(sectionMatch) > 1 {
				currentSection = sectionMatch[1]
			}
		}

		// Update project_id
		if strings.HasPrefix(strings.TrimSpace(line), "project_id") {
			line = fmt.Sprintf("project_id = \"%s\"", projectID)
		}

		// Update ports based on section
		switch currentSection {
		case "db":
			if strings.HasPrefix(strings.TrimSpace(line), "port =") {
				line = fmt.Sprintf("port = %d", ports.DB)
			}
			if strings.HasPrefix(strings.TrimSpace(line), "shadow_port =") {
				line = fmt.Sprintf("shadow_port = %d", ports.Shadow)
			}
		case "studio":
			if strings.HasPrefix(strings.TrimSpace(line), "port =") {
				line = fmt.Sprintf("port = %d", ports.Studio)
			}
		case "inbucket":
			if strings.HasPrefix(strings.TrimSpace(line), "port =") {
				line = fmt.Sprintf("port = %d", ports.Inbucket)
			}
			if strings.HasPrefix(strings.TrimSpace(line), "smtp_port =") {
				line = fmt.Sprintf("smtp_port = %d", ports.SMTP)
			}
			if strings.HasPrefix(strings.TrimSpace(line), "pop3_port =") {
				line = fmt.Sprintf("pop3_port = %d", ports.POP3)
			}
		case "db.pooler":
			if strings.HasPrefix(strings.TrimSpace(line), "port =") {
				line = fmt.Sprintf("port = %d", ports.Pooler)
			}
		case "analytics":
			if strings.HasPrefix(strings.TrimSpace(line), "port =") {
				line = fmt.Sprintf("port = %d", ports.Analytics)
			}
		default:
			// For the root section or API section
			if strings.HasPrefix(strings.TrimSpace(line), "port =") && currentSection == "" {
				line = fmt.Sprintf("port = %d", ports.API)
			}
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading config.toml: %w", err)
	}

	// Write back
	result := strings.Join(lines, "\n")
	if err := os.WriteFile(configPath, []byte(result), 0644); err != nil {
		return fmt.Errorf("failed to write config.toml: %w", err)
	}

	return nil
}
