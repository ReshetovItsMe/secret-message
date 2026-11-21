"""App-only rollout invariants with fake compose and no live credentials."""
import os
import shutil
import subprocess
import tempfile
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
SHA = 'a' * 40
OLD = 'b' * 40


class RolloutTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        (self.root / 'deploy').mkdir()
        source = ROOT / '.github/deployment/rollout.sh'
        if source.exists():
            shutil.copy2(source, self.root / 'deploy/rollout.sh')
        shutil.copy2(ROOT / 'docker-compose.prod.yml', self.root)
        self.active_config = self.root / 'docker-compose.prod.yml'
        self.stage = self.root / 'staging'
        (self.stage / SHA).mkdir(parents=True)
        (self.stage / OLD).mkdir(parents=True)
        self.candidate_config = self.stage / SHA / 'docker-compose.prod.yml'
        self.candidate_config.write_text(self.active_config.read_text() + '\n# synthetic candidate app config\n')
        self.previous_config = self.stage / OLD / 'docker-compose.prod.yml'
        self.previous_config.write_text(self.active_config.read_text() + '\n# synthetic previous app config\n')
        (self.root / 'key.pem').write_text('synthetic placeholder')
        (self.root / 'deploy.conf').write_text(
            'APP_HOST=example.invalid\nRSA_PRIVATE_KEY_FILE=' + str(self.root / 'key.pem') + '\n')
        (self.root / 'bin').mkdir()
        docker = self.root / 'bin/docker'
        docker.write_text('#!/bin/sh\n'
            'printf "tag=%s %s\\n" "$IMAGE_TAG" "$*" >> "$FAKE_CALLS"\n'
            'case "$*" in\n'
            '  *"ps --status running -q redis"*) echo synthetic-container ;;\n'
            '  *" pull "*) [ "$FAKE_FAIL_PULL" != 1 ] || exit 42 ;;\n'
            '  *" up -d --no-deps "*) if [ "$FAKE_FAIL_UP" = 1 ] && [ "$IMAGE_TAG" = "$FAKE_CANDIDATE_SHA" ]; then exit 44; fi ;;\n'
            'esac\n')
        docker.chmod(0o755)
        df = self.root / 'bin/df'
        df.write_text('#!/bin/sh\n'
            'printf "Filesystem 1024-blocks Used Available Capacity Mounted\\n"\n'
            'printf "synthetic 3000000 1 %s 1%% /\\n" "$FAKE_DISK_KB"\n')
        df.chmod(0o755)
        python = self.root / 'bin/python3'
        python.write_text('#!/bin/sh\n'
            'printf "%s\\n" "$*" >> "$FAKE_CALLS"\n'
            'if [ "$FAKE_FAIL_SMOKE" = 1 ]; then exit 43; fi\n'
            'if [ "$FAKE_FAIL_SMOKE" = first ] && [ ! -e "$FAKE_SMOKE_COUNT" ]; then\n'
            '  touch "$FAKE_SMOKE_COUNT"; exit 43; fi\n')
        python.chmod(0o755)
        self.calls = self.root / 'calls'
        self.env = dict(os.environ, PATH=str(self.root / 'bin') + ':' + os.environ['PATH'],
                        DEPLOY_SETTINGS_FILE=str(self.root / 'deploy.conf'),
                        FAKE_CALLS=str(self.calls), FAKE_SMOKE_COUNT=str(self.root / 'smoke-count'),
                        FAKE_DISK_KB='3000000',
                        FAKE_FAIL_PULL='0', FAKE_FAIL_UP='0', FAKE_FAIL_SMOKE='0',
                        FAKE_CANDIDATE_SHA=SHA, CANDIDATE_COMPOSE_FILE=str(self.candidate_config))

    def deploy(self, sha=SHA):
        script = self.root / 'deploy/rollout.sh'
        self.assertTrue(script.exists(), 'rollout script not implemented yet')
        return subprocess.run(['bash', str(script), sha], cwd=self.root, env=self.env,
                              capture_output=True, text=True)

    def calls_text(self):
        return self.calls.read_text() if self.calls.exists() else ''

    def test_rejects_non_sha_without_touching_containers(self):
        result = self.deploy('latest')
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse(self.calls.exists())

    def test_low_disk_stops_before_any_pull(self):
        self.env['FAKE_DISK_KB'] = '1000'
        result = self.deploy()
        self.assertNotEqual(result.returncode, 0)
        self.assertNotIn(' pull ', self.calls_text())
        self.assertFalse((self.root / '.deployed-sha').exists())

    def test_pull_failure_does_not_replace_previous(self):
        (self.root / '.deployed-sha').write_text(OLD + '\n')
        self.env['FAKE_FAIL_PULL'] = '1'
        result = self.deploy()
        self.assertNotEqual(result.returncode, 0)
        self.assertEqual((self.root / '.deployed-sha').read_text().strip(), OLD)
        self.assertNotIn(' up ', self.calls_text())

    def test_failed_smoke_rolls_back_only_app_services(self):
        (self.root / '.deployed-sha').write_text(OLD + '\n')
        self.env['FAKE_FAIL_SMOKE'] = '1'
        result = self.deploy()
        self.assertNotEqual(result.returncode, 0)
        calls = self.calls_text()
        self.assertGreaterEqual(calls.count(' up -d --no-deps frontend gateway secret-assistant'), 2)
        self.assertIn('tag=' + SHA, calls)
        self.assertIn('tag=' + OLD, calls)
        self.assertIn(str(self.previous_config), calls)
        self.assertNotIn('synthetic candidate app config', self.active_config.read_text())
        self.assertNotIn(' up -d --no-deps redis', calls)
        self.assertNotIn(' down', calls)
        self.assertNotIn(' --volumes', calls)
        self.assertEqual((self.root / '.deployed-sha').read_text().strip(), OLD)

    def test_healthy_rollback_restores_previous_revision(self):
        (self.root / '.deployed-sha').write_text(OLD + '\n')
        self.env['FAKE_FAIL_SMOKE'] = 'first'
        result = self.deploy()
        self.assertNotEqual(result.returncode, 0)
        self.assertIn('previous app images restored', result.stderr)
        self.assertEqual((self.root / '.deployed-sha').read_text().strip(), OLD)
        self.assertIn('tag=' + OLD, self.calls_text())

    def test_partial_compose_up_failure_rolls_back_apps(self):
        (self.root / '.deployed-sha').write_text(OLD + '\n')
        self.env['FAKE_FAIL_UP'] = '1'
        result = self.deploy()
        self.assertEqual(result.returncode, 1, result.stderr)
        calls = self.calls_text()
        self.assertIn('tag=' + SHA, calls)
        self.assertIn('tag=' + OLD, calls)
        self.assertIn(str(self.previous_config), calls)
        self.assertNotIn(' up -d --no-deps redis', calls)
        self.assertEqual((self.root / '.deployed-sha').read_text().strip(), OLD)

    def test_rollback_failure_has_distinct_status(self):
        (self.root / '.deployed-sha').write_text(OLD + '\n')
        self.env['FAKE_FAIL_SMOKE'] = '1'
        result = self.deploy()
        self.assertEqual(result.returncode, 2, result.stderr)
        self.assertIn('manual recovery required', result.stderr)

    def test_first_failed_smoke_requires_manual_recovery(self):
        self.env['FAKE_FAIL_SMOKE'] = '1'
        result = self.deploy()
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse((self.root / '.deployed-sha').exists())
        self.assertEqual(self.calls_text().count(' up -d --no-deps frontend gateway secret-assistant'), 1)

    def test_success_records_sha_atomically(self):
        result = self.deploy()
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual((self.root / '.deployed-sha').read_text().strip(), SHA)
        self.assertIn('synthetic candidate app config', self.active_config.read_text())
        self.assertFalse(list(self.root.glob('.deployed-sha.*')))


if __name__ == '__main__':
    unittest.main()
