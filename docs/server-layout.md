# Server layout (file locations)

Everything lives in a tidy, predictable layout. You can also see this any time
from **Manage → File Locations** in the CLI.

| Path | What |
|------|------|
| `/root/ani-tunnel` | The release bundle and downloaded archives. |
| `/root/ani-tunnel/backups` | [Backup](backup-restore.md) `.tar.gz` files. |
| `/etc/ani-tunnel` | Tunnel configs (one `.toml` per tunnel) and runtime state. |
| `/usr/local/bin/ani-tunnel` | The binary itself. |
| `ani-tunnel-<name>.service` | A systemd unit per tunnel. |
| `ani-tunnel-monitor.service` | The [monitor service](monitor-service.md). |

The install directory is recorded in `/etc/ani-tunnel/install_path`, which is what
the uninstaller reads to know what to remove.

---

<div dir="rtl">

## خلاصهٔ فارسی

همه‌چیز در یک ساختار مرتب و قابل‌پیش‌بینی است. همین را هر وقت خواستی از
`Manage → File Locations` در CLI هم می‌بینی.

`‎/root/ani-tunnel` بستهٔ ریلیز و آرشیوهای دانلودشده ·
`‎/root/ani-tunnel/backups` فایل‌های [پشتیبان](backup-restore.md) ·
`‎/etc/ani-tunnel` کانفیگ تونل‌ها (برای هر تونل یک فایل `.toml`) و وضعیت اجرا ·
`‎/usr/local/bin/ani-tunnel` خودِ باینری ·
`ani-tunnel-<name>.service` یک یونیت systemd به‌ازای هر تونل ·
`ani-tunnel-monitor.service` [سرویس مانیتور](monitor-service.md).

مسیر نصب در `/etc/ani-tunnel/install_path` ثبت می‌شود و حذف‌کنندهٔ داخلی از روی
همان می‌فهمد چه چیزی را پاک کند.

</div>

---
[← Back to the docs index](README.md)

---

*Last verified against upstream v1.9.0.*
