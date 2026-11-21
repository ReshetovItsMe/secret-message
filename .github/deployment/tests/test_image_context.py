"""Security boundaries for public Docker build contexts."""
import re
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
MODULES = ('web', 'gateway', 'secret_assistant')
REQUIRED_EXCLUSIONS = {
    '.env*', '**/.env*', '*.pem', '**/*.pem', '*.key', '**/*.key',
    'secrets/', '**/secrets/', '*.log', '**/*.log',
    'node_modules/', '**/node_modules/', 'dist/', '**/dist/',
    'bin/', '**/bin/', '.git/', '**/.git/',
}


class ImageContextTests(unittest.TestCase):
    def test_each_build_context_excludes_sensitive_files_and_artifacts(self):
        for module in MODULES:
            with self.subTest(module=module):
                ignore = ROOT / module / '.dockerignore'
                self.assertTrue(ignore.is_file(), str(ignore))
                lines = {line.strip() for line in ignore.read_text().splitlines()
                         if line.strip() and not line.startswith('#')}
                self.assertTrue(REQUIRED_EXCLUSIONS <= lines,
                                f'{module} missing {REQUIRED_EXCLUSIONS - lines}')
                self.assertFalse(any(line.startswith('!') for line in lines))

    def test_crypto_final_image_contains_only_binary_and_runtime(self):
        dockerfile = (ROOT / 'secret_assistant/Dockerfile').read_text()
        stages = re.split(r'(?im)^FROM\s+', dockerfile)
        self.assertGreaterEqual(len(stages), 3)
        build, runtime = stages[-2:]
        self.assertRegex(build, r'(?i)^golang:1\.27')
        self.assertNotRegex(runtime, r'(?i)^golang:')
        self.assertRegex(runtime, r'(?m)^COPY\s+--from=build\s+/secret-assistant\s+')
        self.assertNotRegex(runtime, r'(?m)^COPY\s+\.\s+\.')
        self.assertNotIn('go build', runtime)
        self.assertIn('ENV RSA_PRIVATE_KEY_PATH=/run/secrets/secret_assistant_private_key', runtime)


if __name__ == '__main__':
    unittest.main()
