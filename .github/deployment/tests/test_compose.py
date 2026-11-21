"""Production Compose invariants; run resolved Compose checks when available."""
import json
import os
import re
import shutil
import subprocess
import tempfile
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
COMPOSE = ROOT / 'docker-compose.prod.yml'
SHA = 'a' * 40  # synthetic SHA, not a live release
IMAGES = {
    'frontend': 'reshetovitsme/secret-frontend',
    'gateway': 'reshetovitsme/secret-gateway',
    'secret-assistant': 'reshetovitsme/secret-assistant',
}


def service(name, text):
    match = re.search(r'(?ms)^  ' + re.escape(name) + r':\n(.*?)(?=^  [a-z][\w-]*:|^\S|\Z)', text)
    if not match:
        raise AssertionError(f'missing service {name}')
    return match.group(1)


class ComposeStaticTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.config = COMPOSE.read_text()

    def test_only_frontend_is_published_on_loopback(self):
        front = service('frontend', self.config)
        self.assertIn('127.0.0.1:18080:80', front)
        for name in ('gateway', 'secret-assistant', 'redis'):
            self.assertNotIn('ports:', service(name, self.config))
        self.assertNotIn('"8080:80"', front)

    def test_sha_tags_only_not_mutable_latest(self):
        for name, repository in IMAGES.items():
            with self.subTest(name=name):
                self.assertIn('image: ' + repository + ':${IMAGE_TAG:?', service(name, self.config))
        self.assertNotIn(':latest', self.config)

    def test_redis_persistence_and_immutable_key_mount(self):
        redis = service('redis', self.config)
        self.assertIn('appendonly', redis)
        self.assertIn('appendfsync', redis)
        self.assertIn('always', redis)
        self.assertIn('redis_data:/data', redis)
        self.assertIn('redis_data:', self.config)
        self.assertIn('RSA_PRIVATE_KEY_FILE:?', self.config)
        self.assertIn('/run/secrets/secret_assistant_private_key', self.config)
        self.assertNotIn('./secrets/', self.config)
        self.assertIn('secret_assistant_private_key:', self.config)


class ComposeResolvedTests(unittest.TestCase):
    @unittest.skipUnless(shutil.which('docker'), 'Docker Compose unavailable here')
    def test_resolved_config_uses_only_sha_images_private_port_and_volume(self):
        if subprocess.run(['docker', 'compose', 'version'], capture_output=True).returncode:
            self.skipTest('Docker Compose plugin unavailable here')
        with tempfile.TemporaryDirectory() as temp:
            key = Path(temp) / 'synthetic-key.pem'
            key.write_text('synthetic placeholder, never for encryption')
            env = dict(os.environ, IMAGE_TAG=SHA, RSA_PRIVATE_KEY_FILE=str(key))
            result = subprocess.run(['docker', 'compose', '-f', str(COMPOSE),
                                     'config', '--format', 'json'], capture_output=True,
                                    text=True, env=env, cwd=ROOT)
            self.assertEqual(result.returncode, 0, result.stderr)
            config = json.loads(result.stdout)
            services = config['services']
            for name, repository in IMAGES.items():
                self.assertEqual(services[name]['image'], f'{repository}:{SHA}')
            self.assertEqual(len(services['frontend']['ports']), 1)
            port = services['frontend']['ports'][0]
            self.assertEqual((port['host_ip'], port['published'], port['target']),
                             ('127.0.0.1', '18080', 80))
            for name in ('gateway', 'secret-assistant', 'redis'):
                self.assertFalse(services[name].get('ports'), name)
            redis = services['redis']
            command = ' '.join(redis['command']) if isinstance(redis['command'], list) else redis['command']
            self.assertIn('appendfsync always', command)
            self.assertIn('appendonly yes', command)
            self.assertTrue(any(m['target'] == '/data' and m['type'] == 'volume'
                                for m in redis['volumes']))
            self.assertEqual(config['secrets']['secret_assistant_private_key']['file'], str(key))
            self.assertEqual(services['secret-assistant']['environment']['RSA_PRIVATE_KEY_PATH'],
                             '/run/secrets/secret_assistant_private_key')


if __name__ == '__main__':
    unittest.main()
