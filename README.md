# Mastodon TTS Notifier for Raspberry Pi 3

A headless notification system that reads Mastodon notifications through Text-to-Speech. Runs on Raspberry Pi 3 with local AI voice synthesis using Piper.

## Overview

This application connects to your Mastodon account, monitors for new notifications (mentions, favorites, boosts, and follows), and reads them aloud using Piper TTS.

## Hardware Requirements
- Raspberry Pi 3 (1GB RAM recommended)
- 8GB+ microSD card
- USB sound card or HDMI audio output
- Speakers/headphones
- Stable internet connection

## Software Requirements
- Raspberry Pi OS (Lite or Desktop)
- Go programming language (1.16+)
- Piper TTS engine
- ALSA audio utilities
- Mastodon account with API access token

## Installation Guide

### 1. Set Up Your Raspberry Pi

1. Install Raspberry Pi OS following the [official guide](https://www.raspberrypi.org/documentation/installation/)
2. Ensure your Pi has internet connectivity
3. Update your system:
   ```
   sudo apt update && sudo apt upgrade -y
   ```

### 2. Install Dependencies

```bash
# Install a compatible version of Go (e.g., 1.20.x)
wget https://go.dev/dl/go1.20.8.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.20.8.linux-armv6l.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version

# Set up Go environment
mkdir -p ~/go/{bin,src,pkg}
echo "export GOPATH=$HOME/go" >> ~/.bashrc
echo "export PATH=$PATH:$GOPATH/bin" >> ~/.bashrc
source ~/.bashrc

# Install audio dependencies
sudo apt install alsa-utils

# Install git
sudo apt install git

# Install Piper TTS
sudo apt install python3-pip
pip3 install piper-tts
```

### 3. Download Voice Model

```bash
# Create directory for voice models
mkdir -p ~/piper-voices
cd ~/piper-voices

# Download English voice model
wget https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/danny/low/en_US-danny-low.onnx
wget https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/danny/low/en_US-danny-low.onnx.json
```

### 4. Get Mastodon API Access Token

1. Log into your Mastodon instance in a web browser
2. Go to Preferences > Development > New Application
3. Give it a name (e.g., "TTS Notifier")
4. Make sure the "read:notifications" scope is checked
5. Click "Submit"
6. Copy the "Access token" for later use

### 5. Build the Application

```bash
# Clone the repository
git clone https://github.com/yourusername/rpi3-mastodon-tts-notificatoin-reader.git
cd rpi3-mastodon-tts-notificatoin-reader

# Build the application
go build -o mastodon-tts-notifier main.go
```

### 6. Configure the Application

Create environment variables with your Mastodon credentials:

```bash
# Add these to your ~/.bashrc or ~/.profile
export MASTODON_INSTANCE_URL="https://your-instance.social"
export MASTODON_ACCESS_TOKEN="your-access-token"
```

Remember to reload your profile:
```bash
source ~/.bashrc  # Or source ~/.profile
```

### 7. Running the Application

Run the program manually:
```bash
./mastodon-tts-notifier
```

### 8. Set Up Autostart (Optional)

Create a systemd service to run on startup:

```bash
# Create service file
sudo nano /etc/systemd/system/mastodon-tts.service
```

Add the following content:
```
[Unit]
Description=Mastodon TTS Notification Reader
After=network.target

[Service]
ExecStart=/home/pi/rpi3-mastodon-tts-notificatoin-reader/mastodon-tts-notifier
WorkingDirectory=/home/pi/rpi3-mastodon-tts-notificatoin-reader
Environment=MASTODON_INSTANCE_URL=https://your-instance.social
Environment=MASTODON_ACCESS_TOKEN=your-access-token
User=pi
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start the service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable mastodon-tts
sudo systemctl start mastodon-tts
```

## Testing

Test your TTS setup with:
```bash
echo 'Welcome to the world of speech synthesis!' | piper --model ~/piper-voices/en_US-danny-low.onnx --output-raw | aplay -r 16000 -f S16_LE -t raw -
```

## Troubleshooting

- **No Audio**: Check if sound works with `aplay /usr/share/sounds/alsa/Front_Center.wav`
- **Program Crashes**: Make sure environment variables are set correctly
- **Build Errors**: Ensure Go is installed correctly with `go version`

## License

This project is licensed under the MIT License - see the LICENSE file for details.