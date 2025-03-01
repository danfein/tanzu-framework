// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"strconv"

	configv1alpha1 "github.com/vmware-tanzu/tanzu-framework/apis/config/v1alpha1"
)

func populateDefaultCliFeatureValues(c *configv1alpha1.ClientConfig, defaultCliFeatureFlags map[string]bool) error {
	for featureName, flagValue := range defaultCliFeatureFlags {
		plugin, flag, err := c.SplitFeaturePath(featureName)
		if err != nil {
			return err
		}
		addFeatureFlag(c, plugin, flag, flagValue)
	}
	return nil
}

func addFeatureFlag(c *configv1alpha1.ClientConfig, plugin, flag string, flagValue bool) {
	if c.ClientOptions == nil {
		c.ClientOptions = &configv1alpha1.ClientOptions{}
	}
	if c.ClientOptions.Features == nil {
		c.ClientOptions.Features = make(map[string]configv1alpha1.FeatureMap)
	}
	if c.ClientOptions.Features[plugin] == nil {
		c.ClientOptions.Features[plugin] = make(map[string]string)
	}
	c.ClientOptions.Features[plugin][flag] = strconv.FormatBool(flagValue)
}

// addDefaultFeatureFlagsIfMissing augments the given configuration object with any default feature flags that do not already have a value
// and returns TRUE if any were added (so the config can be written out to disk, if the caller wants to)
func addDefaultFeatureFlagsIfMissing(config *configv1alpha1.ClientConfig, defaultFeatureFlags map[string]bool) bool {
	added := false

	for featurePath, activated := range defaultFeatureFlags {
		plugin, feature, err := config.SplitFeaturePath(featurePath)
		if err == nil && !containsFeatureFlag(config, plugin, feature) {
			addFeatureFlag(config, plugin, feature, activated)
			added = true
		}
	}

	return added
}

// containsFeatureFlag returns true if the features section in the configuration object contains any value for the plugin.feature combination
func containsFeatureFlag(config *configv1alpha1.ClientConfig, plugin, feature string) bool {
	return config.ClientOptions != nil && config.ClientOptions.Features != nil && config.ClientOptions.Features[plugin] != nil &&
		config.ClientOptions.Features[plugin][feature] != ""
}
