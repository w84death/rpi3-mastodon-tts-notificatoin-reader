# Mastodon TTS Notifier for Raspberry Pi 3

A headless notification system that reads Mastodon mentions through Text-to-Speech. Runs on Raspberry Pi 3 with local AI voice synthesis.

## Hardware Requirements
- Raspberry Pi 3 (1GB RAM recommended)
- 8GB+ microSD card
- USB sound card or HDMI audio output
- Speakers/headphones
- Stable internet connection

## Files Required
1. `mastodon_tts.py` - Main application script
2. `config.ini` - Configuration file
3. `requirements.txt` - Python dependencies
4. `mastodon-tts.service` - Systemd service file

## Setup Instructions

### 1. Operating System Setup
1. Flash [Raspberry Pi OS Lite](https://www.raspberrypi.com/software/) (32-bit)
2. Enable SSH: ```touch /boot/ssh```
3. Configure audio output: 
```
sudo raspi-config
Navigate to: *System Options > Audio > Select output device*
```

### 2. Clone Repository

```git clone https://your-repository-url.git /home/pi/mastodon-tts```

### 3. Configuration Files
Create these files in the project directory:

#### `config.ini`
```
[mastodon]
api_base_url = https://your.mastodon.instance
access_token = your_access_token_here

[tts]
model_path = /home/pi/models/en_us_hifi92_light_cpu.addon
speaker = 0
```

#### `mastodon-tts.service`
```
[Unit]
Description=Mastodon Notification TTS Reader
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=/usr/bin/python3 /home/pi/mastodon-tts.py
WorkingDirectory=/home/pi
Restart=always
RestartSec=10
User=pi

[Install]
WantedBy=multi-user.target
```
### 4. Mastodon API Setup
1. Create application at *Preferences > Development > New Application*
   - Scopes: `read:notifications`
   - Redirect URI: `urn:ietf:wg:oauth:2.0:oob`
2. Save generated access token in `config.ini`

### 5. Install Dependencies

```
sudo apt install ffmpeg python3-pip
pip3 install -r requirements.txt
```

### 6. Balacoon TTS Setup
1. Download voice model:
```
mkdir -p ~/models
wget https://huggingface.co/balacoon/tts/resolve/main/en_us_hifi92_light_cpu.addon -P ~/models
```
2. Verify model path in `config.ini`

### 7. Systemd Service Setup
```
sudo cp mastodon-tts.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable mastodon-tts
```

## First Run

Test audio output
```
aplay /usr/share/sounds/alsa/Front_Center.wav
```

Manual start (for testing)
```
python3 mastodon_tts.py
```

## Customization Options

### Change Check Interval
Modify in `mastodon_tts.py`:
```
schedule.every(15).minutes.do(check_and_speak)
```

### Voice Parameters
Add to `config.ini`:
```
[tts]
speed = 1.0 # 0.5-2.0
pitch = 1.0 # 0.5-1.5
```

### Notification Filtering
Edit exclusion list in `mastodon_tts.py`:
```exclude_types=["follow", "favourite", "reblog"]```

## Troubleshooting

**No Sound Output**

Test ALSA configuration
```speaker-test -t wav -c 2```
Check volume levels
```alsamixer```

text

**API Connection Issues**

```sudo journalctl -u mastodon-tts -f```


**TTS Failures**

Run in debug mode

```python3 mastodon_tts.py --debug```

## License
MIT License - Share modifications freely. Credit original authors if redistributed.
