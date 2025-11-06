# Enhanced Dubbing Features

## Overview

The Akashic Scribe dubbing feature has been significantly improved to provide professional-quality AI voice synthesis with extensive customization options. This document details the new capabilities and how to use them.

## New Features

### 1. OpenAI TTS Integration

The dubbing engine now integrates with OpenAI's Text-to-Speech API using the high-quality `tts-1-hd` model for superior audio generation.

**Supported Voice Models:**
- `alloy` - Balanced and neutral voice
- `echo` - Clear and resonant voice
- `fable` - Expressive storytelling voice
- `onyx` - Deep and authoritative voice
- `nova` - Warm and engaging voice
- `shimmer` - Bright and energetic voice

### 2. Advanced Voice Parameters

#### Voice Speed Control
- **Range**: 0.25x to 4.0x
- **Default**: 1.0x (normal speed)
- **Use Cases**:
  - Slow down for educational content (0.75x)
  - Speed up for time-compressed content (1.25x - 1.5x)
  - Dramatic emphasis with slower speech (0.5x)

#### Voice Pitch Adjustment
- **Range**: -20 to +20 semitones
- **Default**: 0 (no adjustment)
- **Use Cases**:
  - Lower pitch for masculine voices (-5 to -10)
  - Higher pitch for feminine voices (+5 to +10)
  - Character differentiation in narration

#### Voice Stability
- **Range**: 0.0 to 1.0
- **Default**: 0.5 (balanced)
- **Higher values** (0.7-1.0): More consistent, predictable voice
- **Lower values** (0.0-0.3): More expressive, varied intonation

### 3. Audio Format Support

Choose from multiple output formats based on your needs:

| Format | Type | Best For | Quality |
|--------|------|----------|---------|
| **MP3** | Lossy | Web streaming, general use | Good compression |
| **WAV** | Lossless | Professional editing | Uncompressed |
| **FLAC** | Lossless | Archival, high-quality | Compressed lossless |
| **AAC** | Lossy | Mobile devices, podcasts | Efficient compression |
| **OGG** | Lossy | Open-source applications | Good quality/size ratio |

### 4. Audio Quality Presets

| Quality | Bit Rate | Use Case |
|---------|----------|----------|
| **Low** | 64-128 kbps | Voice-only, minimal bandwidth |
| **Medium** | 128-192 kbps | Balanced quality/size |
| **High** | 192-320 kbps | Professional podcasts, presentations |
| **Lossless** | Variable | Studio-quality, archival |

### 5. Audio Processing Filters

#### Normalize Audio Levels
- **Default**: Enabled
- **Purpose**: Ensures consistent volume across the entire audio
- **Algorithm**: Uses FFmpeg's `loudnorm` filter for broadcast-standard normalization
- **Benefit**: Prevents volume fluctuations and clipping

#### Remove Long Silences
- **Default**: Disabled
- **Purpose**: Automatically detects and removes extended silent periods
- **Threshold**: -50dB detection threshold
- **Use Cases**:
  - Tightening pacing in presentations
  - Removing gaps in transcribed content
  - Creating more compact audio files

### 6. Custom Voice Support

Upload your own voice samples for voice cloning:

- **Supported Formats**: MP3, WAV, M4A, FLAC
- **Requirements**: Clear audio sample, 10+ seconds recommended
- **Use Cases**:
  - Brand consistency with specific voice talent
  - Celebrity voice impersonation (with proper rights)
  - Matching original speaker's voice characteristics

*Note: Custom voice synthesis requires voice cloning service integration (planned feature).*

## How to Use

### Basic Usage

1. **Enable Dubbing**
   - Check "Create Dubbing" in Step 2: The Incantation

2. **Select Voice Model**
   - Choose from 6 predefined OpenAI voices
   - Each voice has distinct characteristics

3. **Process Video**
   - Click "Begin Scribing" to generate dubbed audio
   - Audio will be saved in the output directory

### Advanced Configuration

1. **Access Advanced Options**
   - Enable dubbing feature
   - Click "Show Advanced Options"

2. **Adjust Voice Parameters**
   - Use sliders to fine-tune voice speed and pitch
   - Preview values are shown in real-time

3. **Configure Audio Settings**
   - Select desired output format
   - Choose quality preset
   - Enable/disable audio filters

4. **Custom Voice (Optional)**
   - Click "Use Custom Voice..."
   - Select a voice sample file
   - Clear selection with "Clear Custom Voice" if needed

## Technical Implementation

### Architecture

```
Translation Text
      ↓
OpenAI TTS API (tts-1-hd model)
      ↓
Raw MP3 Audio
      ↓
FFmpeg Processing Pipeline
  - Pitch adjustment (asetrate filter)
  - Audio normalization (loudnorm filter)
  - Silence removal (silenceremove filter)
  - Format conversion
  - Sample rate conversion
  - Channel configuration
      ↓
Final Dubbed Audio File
```

### FFmpeg Audio Filters

The system uses FFmpeg's advanced audio filter chain:

```bash
# Example filter chain
-af "asetrate=44100*2^(pitch/12),aresample=44100,loudnorm,silenceremove=..."
```

**Filter Components:**
1. **asetrate**: Adjusts pitch by modifying sample rate
2. **aresample**: Resamples to target rate
3. **loudnorm**: Normalizes audio to broadcast standards
4. **silenceremove**: Removes silence at start/end

### API Requirements

**Environment Variable:**
```bash
export OPENAI_API_KEY="your-api-key-here"
```

**Dependencies:**
- `ffmpeg` - Audio processing (required)
- `yt-dlp` - Video download (required for URLs)
- OpenAI API access (required for TTS)

## Configuration Reference

### ScribeOptions Structure

```go
type ScribeOptions struct {
    // Dubbing configuration
    CreateDubbing     bool    // Enable dubbing
    VoiceModel        string  // Voice model name
    UseCustomVoice    bool    // Use custom voice file
    CustomVoicePath   string  // Path to custom voice

    // Voice parameters
    VoiceSpeed        float64 // 0.25 - 4.0
    VoicePitch        float64 // -20 to +20 semitones
    VoiceStability    float64 // 0.0 - 1.0

    // Audio settings
    AudioFormat       string  // mp3, wav, flac, aac, ogg
    AudioQuality      string  // low, medium, high, lossless
    AudioSampleRate   int     // 8000, 16000, 22050, 44100, 48000
    AudioBitRate      int     // 64-320 kbps
    AudioChannels     int     // 1 (mono), 2 (stereo)

    // Audio filters
    NormalizeAudio    bool    // Enable normalization
    RemoveSilence     bool    // Remove long silences
}
```

## Best Practices

### Voice Selection

1. **Educational Content**: Use `fable` (expressive) or `echo` (clear)
2. **Professional Presentations**: Use `alloy` (neutral) or `onyx` (authoritative)
3. **Engaging Narratives**: Use `nova` (warm) or `shimmer` (energetic)
4. **News/Documentary**: Use `onyx` (deep) or `echo` (clear)

### Quality Settings

- **Podcasts**: High quality, MP3 format, 192 kbps
- **Streaming**: Medium quality, AAC format, 128 kbps
- **Archival**: Lossless quality, FLAC format
- **Quick previews**: Low quality, MP3 format, 96 kbps

### Processing Tips

1. **Normalize audio** for consistent volume (recommended)
2. **Avoid extreme pitch adjustments** (±10 semitones max for natural sound)
3. **Use standard speed** (0.75x - 1.25x) for best quality
4. **Remove silence** for tighter pacing in presentations
5. **Higher sample rates** (48000 Hz) for music or high-fidelity content

## Error Handling

### Common Issues

**"OPENAI_API_KEY environment variable not set"**
- Set your OpenAI API key as an environment variable
- Required for TTS generation

**"Voice model must be specified when dubbing is enabled"**
- Select a voice model from the dropdown
- Or upload a custom voice file

**"Custom voice file not found"**
- Verify the file path is correct
- Ensure file format is supported (MP3, WAV, M4A, FLAC)

**"ffmpeg processing failed"**
- Ensure FFmpeg is installed and in PATH
- Check audio filter parameters are valid

### Validation

The system automatically validates:
- Voice speed range (0.25 - 4.0)
- Pitch range (-20 to +20 semitones)
- Audio format compatibility
- Sample rate values
- Bit rate range (64-320 kbps)

## Performance Considerations

### Processing Time

Dubbing adds processing time based on:
- Translation text length
- Selected audio quality
- Applied filters
- Output format

**Typical Times:**
- Short video (1-5 min): 20-60 seconds
- Medium video (5-15 min): 1-3 minutes
- Long video (15+ min): 3-10 minutes

### API Costs

OpenAI TTS API pricing (as of 2024):
- **tts-1-hd**: $15.00 per 1 million characters
- Average cost per video: $0.01 - $0.10

### Disk Space

Estimated file sizes:
- MP3 (high): ~1.5 MB per minute
- WAV: ~10 MB per minute
- FLAC: ~5 MB per minute
- AAC (medium): ~1 MB per minute

## Roadmap

### Planned Features

- [ ] Voice cloning service integration
- [ ] Real-time preview of voice parameters
- [ ] Batch processing for multiple videos
- [ ] Voice emotion control
- [ ] Multi-speaker support
- [ ] Background music mixing
- [ ] Video re-encoding with dubbed audio
- [ ] Cloud-based processing option

### Future Voice Engines

- Google Cloud Text-to-Speech
- Azure Speech Services
- ElevenLabs integration
- Local TTS models (Coqui TTS, Bark)

## Support

For issues or questions about the dubbing feature:

- **GitHub Issues**: [Report bugs](https://github.com/sorrowscry86/akashic-scribe/issues)
- **Discussions**: [Ask questions](https://github.com/sorrowscry86/akashic-scribe/discussions)
- **Email**: SorrowsCry86@voidcat.org

---

**Last Updated**: 2025-11-05
**Version**: 1.1.0 (Enhanced Dubbing Release)
**Author**: VoidCat RDC Development Team
