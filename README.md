# Mastodon TTS Notification Reader

Reads Mastodon notifications via the API and uses Piper TTS + aplay to speak them aloud. Designed for Raspberry Pi or similar Linux systems.

## Prerequisites

*   Go compiler (for building)
*   Piper TTS (installed and model downloaded, e.g., `~/piper-voices/en_US-danny-low.onnx`)
*   `aplay` utility (usually part of `alsa-utils`)
*   `sh` (shell)
*   `dd` (usually part of `coreutils`)

## Setup

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd rpi3-mastodon-tts-notificatoin-reader
    ```
2.  **Build the application:**
    ```bash
    go build -o mastodon-tts-reader main.go
    ```
3.  **Set Environment Variables:**
    You need to provide your Mastodon instance URL and an access token with `read:notifications` scope. Create an access token in your Mastodon settings (Preferences -> Development -> New Application).

    Export the variables in your shell session or add them to your shell profile (e.g., `~/.bashrc` or `~/.profile`):
    ```bash
    export MASTODON_INSTANCE_URL="https://your.mastodon.instance"
    export MASTODON_ACCESS_TOKEN="YourAccessTokenHere"
    # Optional: Adjust the path to your Piper model if it's not the default
    # export PIPER_MODEL_PATH="/path/to/your/model.onnx"
    ```
    *Note: The current code hardcodes the model path `~/piper-voices/en_US-danny-low.onnx`. You might need to adjust the code or use the environment variable approach if you modify the Go program to read it.*

## Running the Application

You can run the application directly or set it up to run periodically using cron.

### Direct Execution

Ensure the environment variables are set, then run:

```bash
./mastodon-tts-reader
```

The application will fetch new notifications since the last run (tracked in `last_notification_id.txt`), speak them, and update the file.

### Running with Cron (Alternative)

To run the script periodically (e.g., checking for notifications every 5 minutes), you can use `cron`.

1.  **Create a wrapper script (Optional but recommended):**
    Create a script, for example, `/home/pi/run_mastodon_tts.sh`, to set environment variables and run the reader. Make it executable (`chmod +x`).

    ```bash
    #!/bin/bash
    # /home/pi/run_mastodon_tts.sh

    # Set the directory where the executable and last ID file are located
    cd /path/to/rpi3-mastodon-tts-notificatoin-reader || exit 1

    # Export necessary variables
    export MASTODON_INSTANCE_URL="https://your.mastodon.instance"
    export MASTODON_ACCESS_TOKEN="YourAccessTokenHere"
    # export PIPER_MODEL_PATH="/path/to/your/model.onnx" # If needed

    # Run the reader
    ./mastodon-tts-reader
    ```
    *Remember to replace placeholder paths and credentials.*

2.  **Edit your crontab:**
    Open the crontab editor:
    ```bash
    crontab -e
    ```

3.  **Add a cron job:**
    Add a line similar to the following, adjusting the schedule and paths as needed. This example runs the wrapper script every 5 minutes and logs output.

    ```crontab
    */5 * * * * /home/pi/run_mastodon_tts.sh >> /home/pi/mastodon_tts.log 2>&1
    ```

    Or, if you prefer the user's original example (runs once daily at 00:05), ensure the environment variables are somehow available to the `cron` job (e.g., defined within the crontab itself or sourced by the script):

    ```crontab
    # Example: Runs daily at 12:05 AM
    5 0 * * * /path/to/rpi3-mastodon-tts-notificatoin-reader/mastodon-tts-reader >> /home/pi/notifications.log 2>&1
    ```
    *Note: Running directly from cron without a wrapper script requires ensuring the `MASTODON_INSTANCE_URL` and `MASTODON_ACCESS_TOKEN` environment variables are defined within the crontab file before the command, or that the Go program is modified to read them from a config file.*
    *Also ensure the working directory is correct so `last_notification_id.txt` is found/saved in the right place, or use absolute paths in the Go code.*

## How it Works

*   Fetches notifications from the Mastodon API using the provided credentials.
*   Uses `since_id` parameter to only fetch notifications newer than the last processed one (ID stored in `last_notification_id.txt`).
*   Extracts text content from HTML in notifications.
*   Constructs a sentence based on the notification type (mention, favourite, reblog, follow).
*   Pipes the text to `piper` for Text-to-Speech synthesis.
*   Prepends a short silence using `dd` to avoid audio cutoff.
*   Pipes the raw audio data to `aplay` for playback.
*   Updates `last_notification_id.txt` with the ID of the newest notification processed.

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
Download the latest piper release for ARM64. Says Pi4 but works fine on Pi3.

```
wget https://github.com/rhasspy/piper/releases/download/v1.2.0/piper_arm64.tar.gz
tar -zxv piper_arm64.tar.gz
```

Add ~/piper/ to the PATH.

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
RestartSec=300

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