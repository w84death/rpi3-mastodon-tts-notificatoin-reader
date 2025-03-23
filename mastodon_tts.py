from mastodon import Mastodon, MastodonNetworkError, MastodonAPIError
import time
import schedule
import os
import logging
from datetime import datetime
from balacoon_tts import TTS
import wave
import subprocess

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Configuration
LAST_NOTIFICATION_FILE = "last_notification_id.txt"
MODEL_PATH = "/path/to/balacoon-model.addon"
API_BASE_URL = 'https://your.mastodon.instance'
ACCESS_TOKEN = 'your_access_token'

class NotificationSystem:
    def __init__(self):
        self.mastodon = Mastodon(
            access_token=ACCESS_TOKEN,
            api_base_url=API_BASE_URL
        )
        self.tts = TTS(MODEL_PATH)
        self.speaker = self.tts.get_speakers()[0]

    def resilient_api_call(self, func, *args, **kwargs):
        """Handles API call retries with exponential backoff"""
        max_retries = 3
        base_delay = 2
        
        for attempt in range(max_retries):
            try:
                return func(*args, **kwargs)
            except (MastodonNetworkError, MastodonAPIError) as e:
                if attempt < max_retries - 1:
                    delay = base_delay * (2 ** attempt)
                    logger.warning(f"API error: {e}. Retry {attempt+1}/{max_retries} in {delay}s")
                    time.sleep(delay)
                else:
                    logger.error(f"API call failed after {max_retries} attempts: {e}")
                    return None

    def get_notifications(self):
        """Retrieve notifications with resilient error handling"""
        last_id = self._read_last_id()
        
        try:
            response = self.resilient_api_call(
                self.mastodon.notifications,
                since_id=last_id,
                exclude_types=["follow", "favourite", "reblog"]
            )
            
            if response is None:
                return []
                
            if response:
                self._save_last_id(response[-1]['id'])
                
            return response[::-1]  # Return oldest first
            
        except Exception as e:
            logger.error(f"Notification processing failed: {e}")
            return []

    def process_notification(self, notification):
        """Convert notification to speech"""
        try:
            if notification['type'] == 'mention':
                author = notification['account']['display_name'] or notification['account']['username']
                content = self._clean_content(notification['status']['content'])
                return f"New mention from {author}: {content}"
        except KeyError as e:
            logger.error(f"Malformed notification: {e}")
        return None

    def speak(self, text):
        """Convert text to speech using Balacoon TTS"""
        try:
            samples = self.tts.synthesize(text, self.speaker)
            with wave.open("/tmp/tts.wav", "w") as f:
                f.setparams((1, 2, self.tts.get_sampling_rate(), len(samples), "NONE", "NONE"))
                f.writeframes(samples)
            subprocess.run(["aplay", "/tmp/tts.wav"], check=True)
            return True
        except Exception as e:
            logger.error(f"TTS failed: {e}")
            return False

    def _clean_content(self, html):
        """Basic HTML to plaintext conversion"""
        return html.replace('<p>', '').replace('</p>', '\n').replace('<br />', '\n')

    def _read_last_id(self):
        try:
            with open(LAST_NOTIFICATION_FILE, 'r') as f:
                return f.read().strip()
        except FileNotFoundError:
            return None

    def _save_last_id(self, notification_id):
        with open(LAST_NOTIFICATION_FILE, 'w') as f:
            f.write(str(notification_id))

def main():
    system = NotificationSystem()
    
    def check_and_speak():
        logger.info(f"Checking notifications at {datetime.now()}")
        notifications = system.get_notifications()
        
        if notifications:
            logger.info(f"Processing {len(notifications)} new notifications")
            for notification in notifications:
                text = system.process_notification(notification)
                if text:
                    logger.info(f"Speaking notification: {text[:80]}...")
                    system.speak(text)
                    time.sleep(1)
        else:
            logger.info("No new notifications")

    # Initial check
    check_and_speak()
    
    # Scheduled checks
    schedule.every(15).minutes.do(check_and_speak)

    while True:
        schedule.run_pending()
        time.sleep(1)

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        logger.info("Shutting down notification system")
