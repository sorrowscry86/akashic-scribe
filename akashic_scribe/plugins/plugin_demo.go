package main

import (
	"fmt"
	"os"
	"path/filepath"

	"akashic_scribe/core"
	"akashic_scribe/plugins/audio_effects"
	"akashic_scribe/plugins/format_converter"
	"akashic_scribe/plugins/subtitle_styler"
)

// This demo shows how to use the plugin system and all available plugins

func main() {
	fmt.Println("=== Akashic Scribe Plugin System Demo ===\n")

	// Create plugin manager
	dataDir := filepath.Join(os.TempDir(), "akashic_scribe_plugin_demo_data")
	cacheDir := filepath.Join(os.TempDir(), "akashic_scribe_plugin_demo_cache")
	manager := core.NewPluginManager(dataDir, cacheDir)

	fmt.Println("1. Plugin Manager created")
	fmt.Printf("   Data directory: %s\n", dataDir)
	fmt.Printf("   Cache directory: %s\n", cacheDir)
	fmt.Println()

	// Create plugins
	audioPlugin := audio_effects.NewAudioEffectsPlugin()
	formatPlugin := format_converter.NewFormatConverterPlugin()
	stylerPlugin := subtitle_styler.NewSubtitleStylerPlugin()

	// Register plugins
	fmt.Println("2. Registering plugins...")
	if err := manager.RegisterPlugin(audioPlugin); err != nil {
		fmt.Printf("   Error registering audio plugin: %v\n", err)
		return
	}
	fmt.Printf("   ✓ Registered: %s v%s\n", audioPlugin.Name(), audioPlugin.Version())

	if err := manager.RegisterPlugin(formatPlugin); err != nil {
		fmt.Printf("   Error registering format plugin: %v\n", err)
		return
	}
	fmt.Printf("   ✓ Registered: %s v%s\n", formatPlugin.Name(), formatPlugin.Version())

	if err := manager.RegisterPlugin(stylerPlugin); err != nil {
		fmt.Printf("   Error registering styler plugin: %v\n", err)
		return
	}
	fmt.Printf("   ✓ Registered: %s v%s\n", stylerPlugin.Name(), stylerPlugin.Version())
	fmt.Println()

	// Load plugins
	fmt.Println("3. Loading plugins...")
	if err := manager.LoadPlugin(audioPlugin); err != nil {
		fmt.Printf("   Error loading audio plugin: %v\n", err)
		return
	}
	fmt.Printf("   ✓ Loaded: %s\n", audioPlugin.Name())

	if err := manager.LoadPlugin(formatPlugin); err != nil {
		fmt.Printf("   Error loading format plugin: %v\n", err)
		return
	}
	fmt.Printf("   ✓ Loaded: %s\n", formatPlugin.Name())

	if err := manager.LoadPlugin(stylerPlugin); err != nil {
		fmt.Printf("   Error loading styler plugin: %v\n", err)
		return
	}
	fmt.Printf("   ✓ Loaded: %s\n", stylerPlugin.Name())
	fmt.Println()

	// Enable plugins
	fmt.Println("4. Enabling plugins...")
	if err := manager.EnablePlugin(audioPlugin.ID()); err != nil {
		fmt.Printf("   Error enabling audio plugin: %v\n", err)
		return
	}
	fmt.Printf("   ✓ Enabled: %s\n", audioPlugin.Name())

	if err := manager.EnablePlugin(formatPlugin.ID()); err != nil {
		fmt.Printf("   Error enabling format plugin: %v\n", err)
		return
	}
	fmt.Printf("   ✓ Enabled: %s\n", formatPlugin.Name())

	if err := manager.EnablePlugin(stylerPlugin.ID()); err != nil {
		fmt.Printf("   Error enabling styler plugin: %v\n", err)
		return
	}
	fmt.Printf("   ✓ Enabled: %s\n", stylerPlugin.Name())
	fmt.Println()

	// Query plugins
	fmt.Println("5. Query plugin system...")
	loadedPlugins := manager.GetLoadedPlugins()
	fmt.Printf("   Total loaded plugins: %d\n", len(loadedPlugins))

	availablePlugins := manager.GetAvailablePlugins()
	fmt.Printf("   Available plugins:\n")
	for _, info := range availablePlugins {
		fmt.Printf("     - %s (%s) - %s\n", info.Name, info.ID, info.Version)
		fmt.Printf("       Loaded: %v, Enabled: %v\n", info.Loaded, info.Enabled)
		fmt.Printf("       Capabilities: %v\n", info.Capabilities)
	}
	fmt.Println()

	// Query by capability
	fmt.Println("6. Query plugins by capability...")
	audioProcessors := manager.GetPluginsByCapability(core.CapabilityAudioProcessing)
	fmt.Printf("   Audio Processing plugins: %d\n", len(audioProcessors))
	for _, p := range audioProcessors {
		fmt.Printf("     - %s\n", p.Name())
	}

	formatConverters := manager.GetPluginsByCapability(core.CapabilityFormatConversion)
	fmt.Printf("   Format Conversion plugins: %d\n", len(formatConverters))
	for _, p := range formatConverters {
		fmt.Printf("     - %s\n", p.Name())
	}

	subtitleStylers := manager.GetPluginsByCapability(core.CapabilitySubtitleStyling)
	fmt.Printf("   Subtitle Styling plugins: %d\n", len(subtitleStylers))
	for _, p := range subtitleStylers {
		fmt.Printf("     - %s\n", p.Name())
	}
	fmt.Println()

	// Health checks
	fmt.Println("7. Performing health checks...")
	healthResults := manager.HealthCheckAll()
	for pluginID, err := range healthResults {
		if err == nil {
			fmt.Printf("   ✓ %s: healthy\n", pluginID)
		} else {
			fmt.Printf("   ✗ %s: %v\n", pluginID, err)
		}
	}
	fmt.Println()

	// Demonstrate plugin usage
	fmt.Println("8. Plugin Usage Examples...")
	demonstrateAudioPlugin(audioPlugin)
	demonstrateFormatPlugin(formatPlugin)
	demonstrateStylerPlugin(stylerPlugin)

	// Cleanup
	fmt.Println("\n9. Cleanup...")
	fmt.Println("   Disabling plugins...")
	manager.DisablePlugin(audioPlugin.ID())
	manager.DisablePlugin(formatPlugin.ID())
	manager.DisablePlugin(stylerPlugin.ID())

	fmt.Println("   Unloading plugins...")
	manager.UnloadPlugin(audioPlugin.ID())
	manager.UnloadPlugin(formatPlugin.ID())
	manager.UnloadPlugin(stylerPlugin.ID())

	fmt.Println("\n=== Demo Complete ===")
}

func demonstrateAudioPlugin(plugin core.Plugin) {
	fmt.Println("\n   Audio Effects Plugin:")
	fmt.Printf("     - Description: %s\n", plugin.Description())

	if audioPlugin, ok := plugin.(core.AudioProcessor); ok {
		formats := audioPlugin.GetSupportedFormats()
		fmt.Printf("     - Supported formats: %v\n", formats)

		// Show available presets if it's the actual implementation
		if ap, ok := plugin.(*audio_effects.AudioEffectsPlugin); ok {
			presets := ap.GetEffectPresets()
			fmt.Printf("     - Available presets: %v\n", getKeys(presets))
		}
	}

	fmt.Println("     - Example usage:")
	fmt.Println("       audioPlugin.ProcessAudio(\"input.mp3\", \"output.mp3\", map[string]interface{}{")
	fmt.Println("         \"noise_reduction\": true,")
	fmt.Println("         \"normalization\": true,")
	fmt.Println("         \"compression\": true,")
	fmt.Println("       })")
}

func demonstrateFormatPlugin(plugin core.Plugin) {
	fmt.Println("\n   Format Converter Plugin:")
	fmt.Printf("     - Description: %s\n", plugin.Description())

	if formatPlugin, ok := plugin.(core.FormatConverter); ok {
		inputFormats := formatPlugin.GetSupportedInputFormats()
		outputFormats := formatPlugin.GetSupportedOutputFormats()
		fmt.Printf("     - Input formats: %v\n", inputFormats)
		fmt.Printf("     - Output formats: %v\n", outputFormats)
	}

	fmt.Println("     - Example usage:")
	fmt.Println("       formatPlugin.ConvertFormat(\"input.srt\", \"output.vtt\", \"vtt\", nil)")
}

func demonstrateStylerPlugin(plugin core.Plugin) {
	fmt.Println("\n   Subtitle Styler Plugin:")
	fmt.Printf("     - Description: %s\n", plugin.Description())

	if stylerPlugin, ok := plugin.(*subtitle_styler.SubtitleStylerPlugin); ok {
		themes := stylerPlugin.GetAvailableThemes()
		fmt.Printf("     - Available themes: %v\n", themes)
		formats := stylerPlugin.GetSupportedFormats()
		fmt.Printf("     - Supported formats: %v\n", formats)
	}

	fmt.Println("     - Example usage:")
	fmt.Println("       stylerPlugin.ProcessSubtitles(\"input.srt\", \"output.ass\", map[string]interface{}{")
	fmt.Println("         \"theme\": \"cinema\",")
	fmt.Println("         \"position\": \"bottom\",")
	fmt.Println("       })")
}

func getKeys(m map[string]map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
