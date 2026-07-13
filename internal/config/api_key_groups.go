package config

import (
	"fmt"
	"strings"
)

// NormalizeAndValidateAPIKeyGroups trims group values and ensures every mapped
// downstream key remains present in api-keys.
func (cfg *Config) NormalizeAndValidateAPIKeyGroups() error {
	if cfg == nil || len(cfg.APIKeyGroups) == 0 {
		return nil
	}
	keys := make(map[string]struct{}, len(cfg.APIKeys))
	for _, key := range cfg.APIKeys {
		if key != "" {
			keys[key] = struct{}{}
		}
	}
	normalized := make(map[string]string, len(cfg.APIKeyGroups))
	for key, group := range cfg.APIKeyGroups {
		if _, ok := keys[key]; !ok {
			return fmt.Errorf("api-key-groups key %q is not present in api-keys", key)
		}
		normalized[key] = strings.TrimSpace(group)
	}
	cfg.APIKeyGroups = normalized
	return nil
}
