#!/usr/bin/env python3
"""Exercise the installer in a disposable filesystem with fake systemd.

Only absolute installation paths and the root-only entry guard are replaced;
the tests neither need root nor touch the host's services or configuration.
"""
import hashlib
import io
import os
from pathlib import Path
import subprocess
import tarfile
import tempfile
import unittest

SOURCE = Path(__file__).resolve().parents[2]
OLD = b'#!/bin/sh\necho old-version\n'
NEW = b'#!/bin/sh\necho "ani tunnel v1.9.1 https://github.com/pesarakloo/ani-tunnel"\n'

class InstallTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        self.source = self.root / 'source'
        self.source.mkdir()
        self.config = self.root / 'etc/backpack'
        self.config.mkdir(parents=True)
        self.token = 'test-token-preserve-me'
        (self.config/'sample.toml').write_text('[server]\ntoken = "'+self.token+'"\n')
        self.bin = self.root / 'bin'
        self.bin.mkdir()
        self.live = self.bin/'backpack'
        self.live.write_bytes(OLD)
        self.live.chmod(0o755)
        self.install = self.root/'install'
        script = (SOURCE/'install.sh').read_text()
        script = script.replace('if [[ $EUID -ne 0 ]]; then', 'if false; then', 1)
        for old, new in [('/usr/local/bin/', str(self.bin)+'/'), ('/root/BackPack', str(self.install)), ('/etc/backpack', str(self.config))]:
            script = script.replace(old, new)
        (self.source/'install.sh').write_text(script)
        self.fake = self.root/'fake'
        self.fake.mkdir()
        (self.fake/'curl').write_text('#!/bin/sh\nexit 22\n')
        (self.fake/'systemctl').write_text('''#!/bin/sh
if [ "$1" = list-units ]; then
  echo 'backpack-sample.service loaded active running test'
  echo 'unrelated.service loaded active running ignore'
elif [ "$1" = restart ]; then
  echo "$2" >> "$ANI_TEST_ROOT/restarted"
  if [ -f "$ANI_TEST_ROOT/fail-restart" ] && ! grep -q old-version "$ANI_TEST_ROOT/bin/backpack"; then exit 1; fi
else exit 1; fi
''')
        for p in self.fake.iterdir():p.chmod(0o755)
        asset = self.source/'ani-tunnel_linux_amd64.tar.gz'
        with tarfile.open(asset,'w:gz') as archive:
            entry=tarfile.TarInfo('backpack'); entry.size=len(NEW); entry.mode=0o755
            archive.addfile(entry, io.BytesIO(NEW))
        self.asset = asset
        self.sums = self.source/'SHA256SUMS'
        self.sums.write_text(hashlib.sha256(asset.read_bytes()).hexdigest()+'  '+asset.name+'\n')

    def run_install(self):
        env=os.environ.copy()
        env.update(PATH=str(self.fake)+os.pathsep+env['PATH'], ANI_TEST_ROOT=str(self.root))
        env.pop('BP_BUILD_FROM_SOURCE',None); env.pop('BP_HELPERS',None)
        return subprocess.run(['bash',str(self.source/'install.sh')],input='',text=True,capture_output=True,env=env,timeout=20)

    def test_upgrade_preserves_config_and_keeps_previous_binary(self):
        result=self.run_install()
        self.assertEqual(result.returncode,0,result.stdout+result.stderr)
        self.assertEqual(self.live.read_bytes(),NEW)
        self.assertIn(self.token,(self.config/'sample.toml').read_text())
        backups=list((self.install/'backups').glob('ani-*'))
        self.assertEqual(len(backups),1)
        self.assertEqual((backups[0]/'backpack').read_bytes(),OLD)
        self.assertIn(self.token,(backups[0]/'config/sample.toml').read_text())
        self.assertEqual((self.root/'restarted').read_text(),'backpack-sample.service\n')
        self.assertEqual((self.bin/'ani-tunnel').resolve(),self.live)

    def test_failed_restart_restores_old_binary(self):
        (self.root/'fail-restart').touch()
        result=self.run_install()
        self.assertNotEqual(result.returncode,0)
        self.assertEqual(self.live.read_bytes(),OLD,result.stdout+result.stderr)
        self.assertIn(self.token,(self.config/'sample.toml').read_text())
        self.assertEqual((self.root/'restarted').read_text().splitlines(),['backpack-sample.service']*2)

    def test_bad_checksum_never_replaces_running_binary(self):
        self.sums.write_text('0'*64+'  '+self.asset.name+'\n')
        result=self.run_install()
        self.assertNotEqual(result.returncode,0)
        self.assertEqual(self.live.read_bytes(),OLD)
        self.assertFalse((self.root/'restarted').exists())

    def test_first_install_creates_command_alias(self):
        self.live.unlink()
        result=self.run_install()
        self.assertEqual(result.returncode,0,result.stdout+result.stderr)
        self.assertEqual((self.bin/'ani-tunnel').read_bytes(),NEW)
        self.assertEqual(list(self.install.glob('.ani-install.*')),[])

    def test_empty_archive_member_cannot_replace_binary(self):
        with tarfile.open(self.asset,'w:gz') as archive:
            entry=tarfile.TarInfo('backpack'); entry.size=0; entry.mode=0o755
            archive.addfile(entry, io.BytesIO(b''))
        self.sums.write_text(hashlib.sha256(self.asset.read_bytes()).hexdigest()+'  '+self.asset.name+'\n')
        result=self.run_install()
        self.assertNotEqual(result.returncode,0)
        self.assertEqual(self.live.read_bytes(),OLD)
        self.assertFalse((self.root/'restarted').exists())

if __name__ == '__main__': unittest.main(verbosity=2)
