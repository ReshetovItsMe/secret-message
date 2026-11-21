"""Static privacy checks for the bundled nginx layer (no live secrets)."""
import re
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]


class WebSecurityTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.config = (ROOT / 'web/nginx.conf').read_text()

    def test_no_request_uri_or_referer_is_logged(self):
        config = self.config
        self.assertRegex(config, r'(?m)^\s*access_log\s+off;')
        self.assertNotRegex(config, r'(?m)^\s*access_log\s+/(?!dev/null)')
        self.assertRegex(config, r'(?m)^\s*error_log\s+/dev/null\s+(?:crit|alert|emerg);')
        self.assertNotRegex(config, r'(?m)^\s*error_log\s+/var/log/')

    def test_all_locations_receive_no_store_and_no_referrer_even_on_errors(self):
        # In nginx, a location add_header overrides http/server add_header; don't
        # silently lose security headers on /secret/ and /api/ or their 404s.
        config = self.config
        self.assertRegex(config, r'add_header\s+Cache-Control\s+"no-store"\s+always;')
        self.assertRegex(config, r'add_header\s+Referrer-Policy\s+"no-referrer"\s+always;')
        self.assertEqual(len(re.findall(r'(?m)^\s*add_header\s+', config)), 2)
        self.assertLess(config.index('add_header Cache-Control'), config.index('server {'))

    def test_spa_and_one_time_api_routes_are_unchanged(self):
        self.assertIn('try_files $uri $uri/ /index.html;', self.config)
        self.assertIn('location /secret/ {', self.config)
        self.assertRegex(self.config, r'rewrite\s+\^/api/\(\.\*\)\$\s+/\$1\s+break;')
        self.assertIn('proxy_pass http://gateway:3000;', self.config)
        self.assertIn('proxy_request_buffering off;', self.config)
        self.assertIn('proxy_buffering off;', self.config)
        self.assertNotRegex(self.config, r'(?m)^\s*proxy_cache\b')


if __name__ == '__main__':
    unittest.main()
