package internal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadUsedIPs(ipmapPath string) (map[string]bool, error) {
	used := make(map[string]bool)

	file, err := os.Open(ipmapPath)
	if err != nil {
		if os.IsNotExist(err) {
			return used, nil // empty map if file doesn't exist
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			used[parts[1]] = true
		}
	}

	return used, scanner.Err()
}

func FindNextFreeIP(base string, used map[string]bool) (string, error) {
	prefix := base[:strings.LastIndex(base, ".")]
	start, _ := strconv.Atoi(base[strings.LastIndex(base, ".")+1:])

	for i := start; i < 255; i++ {
		candidate := fmt.Sprintf("%s.%d", prefix, i)
		if !used[candidate] {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("no available IPs in subnet")
}

func AppendIPMapping(ipmapPath, domain, ip string) error {
	f, err := os.OpenFile(ipmapPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%s=%s\n", domain, ip))
	return err
}
