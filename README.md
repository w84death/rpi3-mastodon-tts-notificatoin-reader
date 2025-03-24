# Mastodon TTS Notifier for Raspberry Pi 3

A headless notification system that reads Mastodon mentions through -to-Speech. Runs on Raspberry Pi 3 with local AI voice synthesis.

## Hardware Requirements
- Raspberry Pi 3 (1GB RAM recommended)
- 8GB+ microSD card
- USB sound card or HDMI audio output
- Speakers/headphones
- Stable internet connection

## Software Requirements

$ echo 'Welcome to the world of speech synthesis!' | piper --model en_US-danny-low.onnx --output-raw | aplay -r 16000 -f S16_LE -t raw -