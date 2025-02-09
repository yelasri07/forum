package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadEnv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening .env file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove optional surrounding quotes
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			value = value[1 : len(value)-1]
		} else if strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`) {
			value = value[1 : len(value)-1]
		}

		os.Setenv(key, value)
	}
}
