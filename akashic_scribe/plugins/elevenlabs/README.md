# ElevenLabs Voice Synthesis Plugin

**VoidCat RDC - Akashic Scribe**

## Overview

The ElevenLabs plugin integrates ElevenLabs' state-of-the-art AI voice synthesis technology into Akashic Scribe, providing natural-sounding, multilingual text-to-speech capabilities with advanced customization options.

## Features

- **High-Quality Voice Synthesis**: Generate natural-sounding speech using ElevenLabs' advanced AI models
- **Multiple Voice Models**: Access to a wide range of professional voices across different styles and languages
- **Advanced Customization**: Fine-tune voice characteristics including stability, similarity, and style
- **Multilingual Support**: Generate speech in multiple languages using multilingual models
- **Fast Generation**: Turbo mode for low-latency applications
- **Voice Cloning**: Create custom voices from audio samples (API support)
- **Usage Tracking**: Monitor API usage and subscription limits

## Capabilities

- `voice_synthesis` - Text-to-speech generation
- `api_integration` - ElevenLabs API integration

## Requirements

- ElevenLabs API key (get one at [elevenlabs.io](https://elevenlabs.io))
- Internet connectivity
- Valid ElevenLabs subscription

## Installation

### Set API Key

```bash
# Set as environment variable (recommended)
export ELEVENLABS_API_KEY="your_api_key_here"

# Or set programmatically
plugin.SetAPIKey("your_api_key_here")
```

### Register Plugin

```go
import "akashic_scribe/plugins/elevenlabs"

plugin := elevenlabs.NewElevenLabsPlugin()
manager.RegisterPlugin(plugin)
manager.LoadPlugin(plugin)
manager.EnablePlugin(plugin.ID())
```

## Usage

### Basic Text-to-Speech

```go
plugin := elevenlabs.NewElevenLabsPlugin()

// Simple generation
err := plugin.GenerateSpeech(
    "Hello, this is a test of ElevenLabs voice synthesis.",
    "output.mp3",
    map[string]interface{}{
        "voice_id": "21m00Tcm4TlvDq8ikWAM", // Rachel voice
    },
)
```

### With Custom Settings

```go
err := plugin.GenerateSpeech(
    "This speech has custom voice settings.",
    "output.mp3",
    map[string]interface{}{
        "voice_id":         "21m00Tcm4TlvDq8ikWAM",
        "model_id":         "eleven_multilingual_v2",
        "stability":        0.7,   // Higher = more consistent
        "similarity_boost": 0.8,   // Higher = closer to original voice
        "style":            0.5,   // Style intensity (v2 models only)
        "speaker_boost":    true,  // Enhance speaker characteristics
    },
)
```

### Using Presets

```go
// Natural preset (balanced)
plugin.GenerateWithPreset(text, voiceID, "natural", "output.mp3")

// Expressive preset (more emotion)
plugin.GenerateWithPreset(text, voiceID, "expressive", "output.mp3")

// Stable preset (consistent voice)
plugin.GenerateWithPreset(text, voiceID, "stable", "output.mp3")

// Fast preset (low latency)
plugin.GenerateWithPreset(text, voiceID, "fast", "output.mp3")
```

### Get Available Voices

```go
// Get all voices
voices, err := plugin.GetVoices()
for _, voice := range voices {
    fmt.Printf("Voice: %s (%s)\n", voice.Name, voice.VoiceID)
    fmt.Printf("  Category: %s\n", voice.Category)
    fmt.Printf("  Description: %s\n", voice.Description)
}

// Get specific voice by name
voice, err := plugin.GetVoiceByName("Rachel")
if err == nil {
    fmt.Printf("Found voice ID: %s\n", voice.VoiceID)
}
```

### Check Usage

```go
usage, err := plugin.GetUsage()
if err == nil {
    fmt.Printf("Character count: %v\n", usage["character_count"])
    fmt.Printf("Character limit: %v\n", usage["character_limit"])
}
```

## Configuration

### Voice Settings Parameters

| Parameter | Type | Range | Default | Description |
|-----------|------|-------|---------|-------------|
| `voice_id` | string | - | Required | Voice identifier |
| `model_id` | string | - | `eleven_monolingual_v1` | TTS model to use |
| `stability` | float64 | 0.0-1.0 | 0.5 | Voice consistency (higher = more stable) |
| `similarity_boost` | float64 | 0.0-1.0 | 0.75 | Similarity to original (higher = closer match) |
| `style` | float64 | 0.0-1.0 | 0.0 | Style intensity (v2 models only) |
| `speaker_boost` | bool | - | true | Enhance speaker characteristics |

### Available Models

| Model | Description | Languages | Speed | Best For |
|-------|-------------|-----------|-------|----------|
| `eleven_monolingual_v1` | English only | English | Medium | English content, good quality |
| `eleven_multilingual_v1` | Multiple languages | 20+ | Medium | International content |
| `eleven_multilingual_v2` | Improved multilingual | 20+ | Medium | Best quality, style control |
| `eleven_turbo_v2` | Fastest generation | Multiple | Fast | Low-latency applications |

### Presets

| Preset | Stability | Similarity | Style | Model | Use Case |
|--------|-----------|------------|-------|-------|----------|
| `natural` | 0.5 | 0.75 | 0.0 | v2 | Balanced, natural speech |
| `expressive` | 0.3 | 0.8 | 0.6 | v2 | Emotional, varied intonation |
| `stable` | 0.8 | 0.5 | - | v1 | Consistent, predictable |
| `fast` | 0.5 | 0.75 | - | turbo | Fastest generation |

## Popular Voice IDs

Some commonly used voice IDs (check API for current list):

- `21m00Tcm4TlvDq8ikWAM` - Rachel (Female, American)
- `AZnzlk1XvdvUeBnXmlld` - Domi (Female, American)
- `EXAVITQu4vr4xnSDxMaL` - Bella (Female, American)
- `ErXwobaYiN019PkySvjV` - Antoni (Male, American)
- `MF3mGyEYCl7XYWbV9V6O` - Elli (Female, American)
- `TxGEqnHWrfWFTfGW9XjX` - Josh (Male, American)
- `VR6AewLTigWG4xSOukaG` - Arnold (Male, American)
- `pNInz6obpgDQGcFmaJgB` - Adam (Male, American)
- `yoZ06aMxZJJ28mfd3POQ` - Sam (Male, American)

## Examples

### Generate Podcast Narration

```go
text := "Welcome to the Akashic Scribe podcast. In today's episode, we'll explore the fascinating world of AI-powered voice synthesis."

err := plugin.GenerateSpeech(text, "podcast_intro.mp3", map[string]interface{}{
    "voice_id":         "pNInz6obpgDQGcFmaJgB", // Adam
    "model_id":         "eleven_multilingual_v2",
    "stability":        0.6,
    "similarity_boost": 0.8,
    "style":            0.3,
})
```

### Generate Multilingual Content

```go
// French
frenchText := "Bonjour, bienvenue dans Akashic Scribe."
plugin.GenerateSpeech(frenchText, "french.mp3", map[string]interface{}{
    "voice_id": "21m00Tcm4TlvDq8ikWAM",
    "model_id": "eleven_multilingual_v2",
})

// Spanish
spanishText := "Hola, bienvenido a Akashic Scribe."
plugin.GenerateSpeech(spanishText, "spanish.mp3", map[string]interface{}{
    "voice_id": "21m00Tcm4TlvDq8ikWAM",
    "model_id": "eleven_multilingual_v2",
})
```

### Character Voice Acting

```go
// Narrator voice
plugin.GenerateWithPreset(narratorText, "TxGEqnHWrfWFTfGW9XjX", "natural", "narrator.mp3")

// Character 1 (excited)
plugin.GenerateSpeech(character1Text, "char1.mp3", map[string]interface{}{
    "voice_id":         "AZnzlk1XvdvUeBnXmlld",
    "stability":        0.3,
    "similarity_boost": 0.9,
    "style":            0.7,
})

// Character 2 (serious)
plugin.GenerateSpeech(character2Text, "char2.mp3", map[string]interface{}{
    "voice_id":         "ErXwobaYiN019PkySvjV",
    "stability":        0.8,
    "similarity_boost": 0.6,
    "style":            0.2,
})
```

## Advanced Features

### Voice Caching

```go
// Cache frequently used voice info
voice, _ := plugin.GetVoiceByName("Rachel")
plugin.SaveVoiceToCache(*voice)

// Load from cache (faster)
cachedVoice, _ := plugin.LoadVoiceFromCache(voice.VoiceID)
```

### Error Handling

```go
err := plugin.GenerateSpeech(text, output, options)
if err != nil {
    if strings.Contains(err.Error(), "API key not configured") {
        // Handle missing API key
        plugin.SetAPIKey(newKey)
    } else if strings.Contains(err.Error(), "status 401") {
        // Handle authentication error
    } else if strings.Contains(err.Error(), "status 429") {
        // Handle rate limiting
        time.Sleep(time.Second * 5)
    }
}
```

## Performance

### Generation Times (Approximate)

| Text Length | Model | Time |
|-------------|-------|------|
| 100 chars | Turbo | 0.5-1s |
| 100 chars | Standard | 1-2s |
| 500 chars | Turbo | 2-3s |
| 500 chars | Standard | 3-5s |
| 1000 chars | Turbo | 4-6s |
| 1000 chars | Standard | 6-10s |

### Optimization Tips

1. **Use Turbo Model**: For real-time applications
2. **Batch Processing**: Process multiple texts in parallel
3. **Cache Voices**: Save voice info to reduce API calls
4. **Appropriate Settings**: Lower stability = faster generation

## API Costs

ElevenLabs pricing (as of 2024):

- **Free Tier**: 10,000 characters/month
- **Starter**: $5/month for 30,000 characters
- **Creator**: $22/month for 100,000 characters
- **Pro**: $99/month for 500,000 characters
- **Scale**: $330/month for 2M characters

## Limitations

- Requires active internet connection
- API rate limits apply
- Character limits based on subscription tier
- Voice cloning requires voice samples and higher tier
- Some voices may not be available on all tiers

## Troubleshooting

### API Key Issues

```
Error: API key not configured
```

**Solution**: Set `ELEVENLABS_API_KEY` environment variable or call `SetAPIKey()`

### Authentication Errors

```
Error: API error (status 401)
```

**Solution**: Verify API key is correct and active

### Rate Limiting

```
Error: API error (status 429)
```

**Solution**: Implement retry logic with exponential backoff

### Character Limit Exceeded

```
Error: quota exceeded
```

**Solution**: Check usage with `GetUsage()` and upgrade plan if needed

## Integration with Other Plugins

### With Format Converter

```go
// Convert subtitles and generate audio
formatPlugin.ConvertFormat("input.srt", "output.txt", "txt", nil)

// Read text and generate speech
text, _ := os.ReadFile("output.txt")
elevenLabsPlugin.GenerateSpeech(string(text), "audio.mp3", options)
```

### With Audio Effects

```go
// Generate speech
elevenLabsPlugin.GenerateSpeech(text, "raw.mp3", options)

// Apply effects
audioPlugin.ProcessAudio("raw.mp3", "final.mp3", map[string]interface{}{
    "normalization":   true,
    "noise_reduction": true,
})
```

## Best Practices

1. **API Key Security**: Never hardcode API keys, use environment variables
2. **Error Handling**: Always handle API errors gracefully
3. **Usage Monitoring**: Regularly check usage to avoid unexpected charges
4. **Voice Selection**: Test different voices for your content type
5. **Settings Tuning**: Adjust stability/similarity based on content needs
6. **Caching**: Cache voice information to reduce API calls
7. **Batch Processing**: Process multiple texts efficiently
8. **Rate Limiting**: Implement retry logic for rate limit errors

## Resources

- **ElevenLabs Website**: https://elevenlabs.io
- **API Documentation**: https://docs.elevenlabs.io
- **Voice Lab**: https://elevenlabs.io/voice-lab
- **Pricing**: https://elevenlabs.io/pricing
- **Support**: support@elevenlabs.io

## License

This plugin is part of Akashic Scribe and follows the same license terms.

---

**© 2025 VoidCat RDC. All rights reserved.**

*Powered by ElevenLabs AI*
