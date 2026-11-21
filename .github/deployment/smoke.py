#!/usr/bin/env python3
"""One synthetic end-to-end check. Never print a message, link or ID."""
import argparse
import json
import secrets
import sys
import time
import uuid
from urllib import error, parse, request


class NoRedirects(request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None


def check(base_url):
    parsed = parse.urlsplit(base_url)
    if parsed.scheme != 'https' and not (parsed.scheme == 'http' and
                                          parsed.hostname in ('localhost', '127.0.0.1')):
        raise ValueError('invalid base URL')
    if parsed.path not in ('', '/') or parsed.query or parsed.fragment or not parsed.hostname:
        raise ValueError('invalid base URL')

    opener = request.build_opener(NoRedirects())
    origin = base_url.rstrip('/')

    def fetch(path, method='GET', body=None):
        req = request.Request(origin + path, method=method, data=body,
                              headers={'Accept': 'application/json',
                                       'Cache-Control': 'no-store',
                                       'Content-Type': 'application/json'})
        try:
            with opener.open(req, timeout=10) as response:
                return response.status, response.read(1 << 20)
        except error.HTTPError as exc:
            return exc.code, b''

    # Retry only readiness: never retry a GET that consumes a message.
    for attempt in range(30):
        try:
            if fetch('/')[0] == 200 and fetch('/api/healthz') == (200, b'ok'):
                break
        except error.URLError:
            pass
        if attempt < 29:
            time.sleep(1)
    else:
        raise ValueError('frontend or gateway unavailable')
    message = 'synthetic-smoke-' + secrets.token_hex(16)
    status, payload = fetch('/api/message', 'POST',
                            json.dumps({'message': message}).encode())
    if status != 201:
        raise ValueError('create failed')
    identifier = json.loads(payload)['messageId']
    uuid.UUID(identifier)
    path = '/api/message?' + parse.urlencode({'messageId': identifier})
    status, payload = fetch(path)
    if status != 200 or json.loads(payload).get('message') != message:
        raise ValueError('first consume failed')
    if fetch(path)[0] != 404:
        raise ValueError('second consume was not rejected')


def main():
    parser = argparse.ArgumentParser(description='Synthetic one-time message check')
    parser.add_argument('--base-url', required=True)
    args = parser.parse_args()
    try:
        check(args.base_url)
    except Exception as exc:
        # Errors (especially urllib errors) can contain URLs; never print them.
        print('Synthetic smoke check failed (' + type(exc).__name__ + ')', file=sys.stderr)
        return 1
    print('Synthetic smoke check passed')
    return 0


if __name__ == '__main__':
    sys.exit(main())
