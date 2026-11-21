"""Exercise synthetic one-time smoke requests without external services."""
import json
import subprocess
import sys
import threading
import unittest
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from urllib.parse import urlsplit

ROOT = Path(__file__).resolve().parents[3]
SCRIPT = ROOT / '.github/deployment/smoke.py'
SYNTHETIC_ID = '11111111-1111-4111-8111-111111111111'


class Handler(BaseHTTPRequestHandler):
    message = None
    gets = 0
    fail_first = False
    health_calls = 0
    unready_count = 0

    def log_message(self, *args):
        pass

    def respond(self, status, body=b'', content_type='text/plain'):
        self.send_response(status)
        self.send_header('Content-Type', content_type)
        self.end_headers()
        self.wfile.write(body)

    def do_POST(self):
        if self.path != '/api/message':
            return self.respond(404)
        body = self.rfile.read(int(self.headers['Content-Length']))
        Handler.message = json.loads(body)['message']
        self.respond(201, json.dumps({'messageId': SYNTHETIC_ID}).encode(), 'application/json')

    def do_GET(self):
        path = urlsplit(self.path).path
        if path == '/':
            return self.respond(200, b'<html>synthetic page</html>')
        if path == '/api/healthz':
            Handler.health_calls += 1
            if Handler.health_calls <= Handler.unready_count:
                return self.respond(503)
            return self.respond(200, b'ok')
        if path == '/api/message':
            Handler.gets += 1
            if Handler.fail_first and Handler.gets == 1:
                return self.respond(500)
            if Handler.message is None:
                return self.respond(404)
            message, Handler.message = Handler.message, None
            return self.respond(200, json.dumps({'message': message}).encode(), 'application/json')
        self.respond(404)


class SmokeTests(unittest.TestCase):
    def setUp(self):
        Handler.message = None
        Handler.gets = 0
        Handler.fail_first = False
        Handler.health_calls = 0
        Handler.unready_count = 0
        self.server = ThreadingHTTPServer(('127.0.0.1', 0), Handler)
        self.thread = threading.Thread(target=self.server.serve_forever, daemon=True)
        self.thread.start()
        self.addCleanup(self.server.server_close)
        self.addCleanup(self.server.shutdown)
        self.base_url = 'http://127.0.0.1:' + str(self.server.server_port)

    def run_smoke(self):
        self.assertTrue(SCRIPT.exists(), 'smoke script not implemented yet')
        return subprocess.run([sys.executable, str(SCRIPT), '--base-url', self.base_url],
                              capture_output=True, text=True)

    def test_create_consume_and_second_read_404_without_identifiers_in_output(self):
        result = self.run_smoke()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(Handler.gets, 2)
        self.assertIsNone(Handler.message)
        self.assertNotIn(SYNTHETIC_ID, result.stdout + result.stderr)
        self.assertNotIn(self.base_url, result.stdout + result.stderr)

    def test_waits_for_gateway_readiness_without_consuming_any_message(self):
        Handler.unready_count = 1
        result = self.run_smoke()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertGreaterEqual(Handler.health_calls, 2)
        self.assertEqual(Handler.gets, 2)

    def test_failure_does_not_retry_consuming_get_or_print_url(self):
        Handler.fail_first = True
        result = self.run_smoke()
        self.assertNotEqual(result.returncode, 0)
        self.assertEqual(Handler.gets, 1)
        self.assertNotIn(SYNTHETIC_ID, result.stdout + result.stderr)
        self.assertNotIn(self.base_url, result.stdout + result.stderr)


if __name__ == '__main__':
    unittest.main()
