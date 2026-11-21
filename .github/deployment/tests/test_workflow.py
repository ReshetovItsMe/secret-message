"""Guard rails for publishing only tested revisions and pinned SSH deployment."""
import re
import unittest
from pathlib import Path

WORKFLOW = Path(__file__).resolve().parents[3] / '.github/workflows/build.yml'


def job(name, workflow):
    match = re.search(r'(?ms)^  ' + re.escape(name) + r':\n(.*?)(?=^  [a-z][\w-]*:|\Z)', workflow)
    if not match:
        raise AssertionError(f'missing job {name}')
    return match.group(1)


class WorkflowTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.workflow = WORKFLOW.read_text()

    def test_prs_and_main_keep_all_checks(self):
        self.assertIn('pull_request:', self.workflow)
        self.assertRegex(self.workflow, r'push:\s*branches:\s*- main')
        for name in ('frontend', 'gateway', 'secret_assistant', 'compose'):
            self.assertIn('runs-on: ubuntu-latest', job(name, self.workflow))
        self.assertIn('pnpm type-check', job('frontend', self.workflow))
        self.assertIn('make lint', job('gateway', self.workflow))
        self.assertIn('python3 -m unittest discover -s .github/deployment/tests', job('compose', self.workflow))
        for module in ('web', 'gateway', 'secret_assistant'):
            self.assertIn('docker build', job('compose', self.workflow))
            self.assertIn(' ' + module, job('compose', self.workflow))
        self.assertIn('--entrypoint nginx', job('compose', self.workflow))
        self.assertIn('--add-host gateway:127.0.0.1', job('compose', self.workflow))

    def test_publish_requires_checks_and_explicit_activation(self):
        published = job('publish', self.workflow)
        self.assertIn('needs: [frontend, gateway, secret_assistant, compose]', published)
        self.assertIn("github.event_name == 'push'", published)
        self.assertIn("vars.PRODUCTION_PUBLISH_ENABLED == 'true'", published)
        self.assertNotIn('vars.PRODUCTION_DEPLOY_ENABLED', published)
        self.assertIn('docker/login-action@', published)
        self.assertIn('secrets.DOCKERHUB_TOKEN', published)
        self.assertIn('docker/build-push-action@', published)
        for module, image in (('web', 'secret-frontend'), ('gateway', 'secret-gateway'),
                              ('secret_assistant', 'secret-assistant')):
            self.assertIn('context: ./' + module, published)
            self.assertIn('reshetovitsme/' + image + ':${{ github.sha }}', published)
        self.assertNotIn(':latest', published)
        self.assertNotIn('build-args:', published)

    def test_deploy_uses_fingerprint_and_serialization_without_password(self):
        deployment = job('deploy', self.workflow)
        self.assertIn('needs: [publish]', deployment)
        self.assertIn("github.event_name == 'push'", deployment)
        self.assertIn('queue: max', deployment)
        self.assertIn('fetch-depth: 0', deployment)
        self.assertIn('git merge-base --is-ancestor', deployment)
        self.assertIn('.deployed-sha', deployment)
        self.assertIn('cancel-in-progress: false', deployment)
        self.assertIn('StrictHostKeyChecking=yes', deployment)
        self.assertIn('secrets.DEPLOY_SSH_HOST_KEY', deployment)
        self.assertIn('secrets.DEPLOY_SSH_PRIVATE_KEY', deployment)
        self.assertIn('.github/deployment/rollout.sh', deployment)
        self.assertIn('.github/deployment/smoke.py', deployment)
        self.assertIn('deploy/rollout.sh', deployment)
        self.assertIn('CANDIDATE_COMPOSE_FILE=$stage/docker-compose.prod.yml', deployment)
        self.assertNotIn('install -m 0644 $stage/docker-compose.prod.yml /opt/secret-message/docker-compose.prod.yml', deployment)
        self.assertIn('status=$?', deployment)
        self.assertIn('Rollback or initial recovery failed', deployment)
        self.assertIn('${{ github.sha }}', deployment)
        self.assertNotIn('sshpass', deployment)
        self.assertNotIn('StrictHostKeyChecking=no', deployment)
        self.assertNotRegex(self.workflow, r'(?m)^\s*(?:host|password|privateKey):\s*(?:[0-9]{1,3}\.){3}')


if __name__ == '__main__':
    unittest.main()
