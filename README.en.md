# wb-mqtt-noolite

[Русский](README.md) | [English](README.en.md)

wb-mqtt-noolite integrates Noolite devices (remotes, switches and sensors) into a [Wiren Board](https://wirenboard.com/) automation controller.
This is achieved by following the [Wiren Board MQTT Conventions](https://github.com/wirenboard/conventions) and [WB-STD-001 "MQTT Topic Identifiers in the Wiren Board Ecosystem"](https://github.com/wirenboard/wb-standards), and by using the MTRF-64 adapter, which handles communication with Noolite(-F) devices.
_The service has been tested with the [MTRF-64-USB-A](https://noo.by/adapter-mtrf-64-usb-a.html) adapter._

**Features:**
* Receiving control commands from Wiren Board and forwarding them to Noolite devices
* Communicating with the NooLite MTRF-64 adapter, handling device responses and publishing status updates to MQTT
* Polling Noolite(-F) devices on a schedule and sending the results to Wiren Board

> **Upgrading from 1.x?** See ["Upgrading from 1.x to 2.0"](#upgrading-from-1x-to-20) below.

### Building for Wirenboard

Controller revision to Go build target (`GOARCH`/`GOARM`) mapping:

| Controller | CPU | Build target | `GOARCH` | `GOARM` |
|---|---|---|---|---|
| Wiren Board 3.5 and older | ARM926EJ-S | `armv5` | `arm` | `5` |
| Wiren Board rev. 6.3-6.6 | NXP i.MX 6ULL, Cortex-A7, 800 MHz | `armv7` | `arm` | `7` |
| Wiren Board 6.7 | NXP i.MX 6ULL, Cortex-A7, 800 MHz | `armv7` | `arm` | `7` |
| Wiren Board 7.4 | Cortex-A7, 4 cores, 1.2 GHz | `armv7` | `arm` | `7` |
| Wiren Board 8.5 | Cortex-A53, 4 cores, 1.5 GHz (64-bit) | `arm64` | `arm64` | — |

All the Cortex-A7 controllers (rev. 6.3-6.6, 6.7, 7.4) build with the exact same command — `GOARM` only depends on the CPU's instruction set (ARMv7), not on core count or clock speed. A separate build target is only needed for Wiren Board 8.x, which has a 64-bit Cortex-A53 CPU.

#### Using the Makefile (recommended)
```shell
make armv5     # Wiren Board 3.5 and older -> build/armv5/
make armv7     # Wiren Board rev. 6.3-6.6, 6.7, 7.4 -> build/armv7/
make arm64     # Wiren Board 8.5 -> build/arm64/
make all       # all three targets at once
make deb        # .deb packages for Debian 12/13 (armhf, arm64) -> build/deb/
make clean     # remove build/
```
Each target builds both binaries (`wb-mqtt-noolite` and `mtrf_tool`) with the same optimization flags (see ["Size and memory optimization"](#size-and-memory-optimization) below). See ["Building a .deb package (Debian 12/13)"](#building-a-deb-package-debian-1213) below for more on `make deb`.

#### Manually (without the Makefile)
```shell
# Wiren Board 3.5 and older
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -trimpath -ldflags="-s -w" -o wb-mqtt-noolite-armv5 ./cmd/wb-mqtt-noolite

# Wiren Board rev. 6.3-6.6, 6.7, 7.4
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -ldflags="-s -w" -o wb-mqtt-noolite-armv7 ./cmd/wb-mqtt-noolite

# Wiren Board 8.5
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o wb-mqtt-noolite-arm64 ./cmd/wb-mqtt-noolite
```
Replace `./cmd/wb-mqtt-noolite` with `./cmd/mtrf_tool` to build the pairing utility for the same controller.

## Size and memory optimization

Build flags used both in the Makefile and in the commands above:
- **`CGO_ENABLED=0`** — all project dependencies are pure Go, so cross-compiling needs no ARM toolchain with a C compiler, and the resulting binary is fully static (no dependency on the target system's libc).
- **`-trimpath`** — strips the build machine's local filesystem paths from the binary.
- **`-ldflags="-s -w"`** — strips the symbol table and DWARF debug information.

With all three flags combined, the `wb-mqtt-noolite` armv7 binary shrinks from about 12.4 MB to 8.6 MB (unoptimized build vs. optimized build); the arm64 binary comes in at about 8.3 MB.

To bound RAM usage on weaker controllers (mainly Wiren Board rev. 6.3-6.6/6.7 with the NXP i.MX 6ULL), set the Go runtime's environment variables when starting the service — no rebuild required:
- **`GOMEMLIMIT`** — a soft memory limit: the Go garbage collector tries on its own not to exceed the given value (e.g. `GOMEMLIMIT=48MiB`).
- **`GOGC`** — the heap growth ratio between GC cycles (default 100 = the heap can double before the next GC); a lower value (e.g. `GOGC=50`) reduces peak memory use at the cost of more frequent GC and slightly higher CPU load.

Neither requires code changes — both can be set via `Environment=` in the systemd unit file (see below, `packaging/wb-mqtt-noolite.service`) and tuned by watching actual usage (`systemctl status wb-mqtt-noolite`, `ps -o rss -p $(pidof wb-mqtt-noolite)`).

## Running as a systemd service

`packaging/wb-mqtt-noolite.service` contains a ready-made unit file. Installing it on the controller (after building the binary for the right target and copying it over, e.g. via `scp`):
```shell
install -m 0755 wb-mqtt-noolite /usr/bin/wb-mqtt-noolite
install -m 0644 packaging/wb-mqtt-noolite.service /etc/systemd/system/wb-mqtt-noolite.service
systemctl daemon-reload
systemctl enable --now wb-mqtt-noolite
```
Checking status and logs:
```shell
systemctl status wb-mqtt-noolite
journalctl -u wb-mqtt-noolite -f
```

By default the unit file:
- restarts the service on failure (`Restart=on-failure`);
- sets `GOMEMLIMIT=48MiB` and a hard cgroup limit `MemoryMax=64M` — lower both on weaker controllers (rev. 6.3-6.6/6.7) or raise them on more powerful ones (7.4/8.5), see the section above;
- runs as `root` to access the adapter's serial port (`/dev/ttyUSB0`) with no extra setup; an alternative (a dedicated user in the `dialout` group) is documented in a comment inside the unit file itself.

## Building a .deb package (Debian 12/13)

```shell
make deb
```
builds the binaries (`make armv7`/`make arm64`) and packages them into `build/deb/`:
```
wb-mqtt-noolite_2.0.0~bookworm_armhf.deb
wb-mqtt-noolite_2.0.0~trixie_armhf.deb
wb-mqtt-noolite_2.0.0~bookworm_arm64.deb
wb-mqtt-noolite_2.0.0~trixie_arm64.deb
```
The version comes from the `VERSION` file at the repo root. The binary is static (`CGO_ENABLED=0`), so the bookworm and trixie packages are identical in content — the `~bookworm`/`~trixie` version suffix only exists so both can be kept side by side in an apt repository with several suites, not because a different build is actually needed per Debian version. `armhf` corresponds to the `armv7` target (Wiren Board rev. 6.3-6.6, 6.7, 7.4), `arm64` to Wiren Board 8.5. There's no separate package for `armv5` (Wiren Board 3.5 and older) — such controllers aren't supported by Debian 12/13.

To build just one package: `make deb-armhf` or `make deb-arm64`; to build manually for a specific suite directly (e.g. a different one): `packaging/build-deb.sh armhf bookworm build/armv7 2.0.0 build/deb`. Building only needs `dpkg-deb` and `fakeroot` (the `fakeroot` package on Debian/Ubuntu) — a full `debhelper`/`dpkg-buildpackage` setup isn't required, since the package is built from already-compiled binaries rather than from source.

The package includes:
- `/usr/bin/wb-mqtt-noolite`, `/usr/bin/mtrf_tool`;
- `/etc/wb-mqtt-noolite.json` — default configuration (port `/dev/ttyUSB0`, broker `127.0.0.1:1883`);
- `/etc/wb-mqtt-noolite-templates.json` — the stock device templates (see `templates/templates.json`);
- `/etc/wb-mqtt-noolite-devices.json` — an **empty list** `[]`, to be filled in after installation;
- `/lib/systemd/system/wb-mqtt-noolite.service`.

All three files under `/etc` are marked as conffiles — on a package upgrade dpkg will not silently overwrite your edits, it will ask what to do (keep yours / take the new one / view a diff).

Installing on the controller:
```shell
apt install ./wb-mqtt-noolite_2.0.0~bookworm_armhf.deb
```
After installation the service is **not started automatically** — `postinst` only reloads the systemd configuration and prints a reminder. Before the first start, edit `/etc/wb-mqtt-noolite.json` and fill in `/etc/wb-mqtt-noolite-devices.json` (see ["Device list (devices.json)"](#device-list-devicesjson) and ["mtrf_tool utility"](#mtrf_tool-utility-pairing-and-unpairing-devices)), then:
```shell
systemctl enable --now wb-mqtt-noolite
```

## Configuration and running the service
wb-mqtt-noolite is configured via a JSON configuration file.
Main parameters: the Noolite adapter's serial port and MQTT broker credentials.
The `device_config` section must also list the paths to separate JSON files describing device model templates and the actual devices.
Templates are located in the `templates` directory and will grow as the project develops.
The `example` directory contains a sample configuration file and device list.

**Command-line parameters:** `--config` or `-c`, pointing to the service configuration file path

Example configuration file:
````json
{
  "serial_port": "/dev/ttyUSB0",
  "timezone": "Europe/Moscow",
  "loglevel": "info",
  "mqtt": {
    "host": "127.0.0.1",
    "port": 1883,
    "username": "",
    "password": ""
  },
  "device_config": {
    "templates": "./templates/templates.json",
    "devices": "./example/devices.json"
  }
}
````

### Device list (devices.json)
The service needs information about which devices are present in the system and which model template each of them uses.
**Example device list:**
```json
[
  {
    "name": "Underfloor heating",
    "noolite_type": "txf",
    "ch": 1,
    "template": "srf-1-3000-t"
  }
]
```
_**Fields:**_
- **name**: Device name. If set, it is published as the device's `title` in MQTT (see below); if not set, the model name from the template is used instead
- **noolite_type**: Noolite type (TX, TX-F, RX, RX-F)
- **ch**: Channel, 1-63
- **address**: An address may optionally be specified
- **template**: Device template

### Device templates (templates.json)
A device template describes a supported device model, its controls, and the Noolite commands sent to the device.
Scheduled execution of specific commands is supported, with the schedule given in crontab format.
A template consists of a model name (`name`), a localized title (`title`), and a set of controls (`controls`).

Each control has:
- **name**: the channel identifier (the `/controls/<name>` topic segment), formed per the [WB-STD-001](https://github.com/wirenboard/wb-standards) rules (lowercase Latin, `snake_case`)
- **title**: a localized channel name, an object of the form `{"en": "...", "ru": "..."}`
- **type**: the control type, see [Wiren Board MQTT Conventions](https://github.com/wirenboard/conventions). For physical quantities outside the standard type list (`switch`, `alarm`, `range`, `rgb`, `text`), use the generic `value` type together with the `units` key (e.g. `"type": "value", "units": "deg C"`) — the specific types `temperature`/`rel_humidity`/`power`/`power_consumption` are deprecated
- **order**: sort order
- **readonly**: whether the value can be changed from the UI
- **get_command**: the command executed on a schedule, e.g. _"ReadState 0"_
- **set_command**: the Noolite command sent when a command to set this control is received
- **polling**: whether to run the command on a schedule
- **polling_cron**: the schedule for the command given in `get_command`
- **min**, **max**: for the `range` type, the minimum and maximum accepted value
- **units**: the unit of measurement (machine form per the Wiren Board MQTT Conventions dictionary, e.g. `deg C`, `%`)
- **precision**: precision
- **dont_use_retain**: ignore the retained MQTT value when restoring state after a reconnect

Example:
```json
{
  "templates": [
    {
      "name": "srf-1-3000-t",
      "title": {
        "en": "SRF-1-3000-T thermostat",
        "ru": "Термостат SRF-1-3000-T"
      },
      "controls": [
        {
          "name": "on",
          "type": "switch",
          "title": {"en": "Heating", "ru": "Нагрев"},
          "order": 1,
          "initial_value": "0",
          "readonly": false,
          "set_command": "SetSwitch"
        },
        {
          "name": "setpoint_t_ambient",
          "type": "range",
          "title": {"en": "Setpoint Temperature", "ru": "Заданная температура"},
          "order": 2,
          "min": 5,
          "max": 30,
          "initial_value": "23",
          "set_command": "SetTemperature"
        },
        {
          "name": "temperature",
          "type": "value",
          "units": "deg C",
          "title": {"en": "Temperature", "ru": "Температура"},
          "order": 3,
          "readonly": true,
          "get_command": "ReadState 0",
          "polling": true,
          "polling_cron": "*/1 * * * *"
        },
        {
          "name": "model",
          "type": "text",
          "title": {"en": "Model", "ru": "Модель"},
          "order": 4,
          "readonly": true
        },
        {
          "name": "address",
          "type": "text",
          "title": {"en": "Address", "ru": "Адрес"},
          "order": 5,
          "readonly": true
        }
      ]
    }
  ]
}
```

#### Supported commands
Each control specifies the commands that convert an MQTT message into a Noolite command and back.\
**Supported commands are listed below:**

|Command|Description and parameters|Example|
|---|---|---|
|SetOn|Turn on, no parameters.|SetOn|
|SetOff|Turn off, no parameters.|SetOff|
|SetSwitch|Turn on or off; 1 or 0 is passed respectively.|SetSwitch|
|SetTemperature|Set the temperature (for srf-1-3000-t).|SetTemperature|
|ReadState|Request status, parameter `fmt` (see the MTRF documentation)|ReadState 0|

*The command list will be extended.

## MQTT topic structure

The service publishes topics per the [Wiren Board MQTT Conventions](https://github.com/wirenboard/conventions) and [WB-STD-001](https://github.com/wirenboard/wb-standards).

For a device of type `txf` on channel `1` (default `device_prefix` is `/devices/mtrf_`):

| Topic | Retained | Description |
|---|---|---|
| `/devices/mtrf_txf_1/meta` | yes | JSON `{"driver": "wb-mqtt-noolite", "title": {"en": "...", "ru": "..."}}` |
| `/devices/mtrf_txf_1/meta/error` | yes | Device error flag: `r`/`w`/`p`, or an empty string (no error) |
| `/devices/mtrf_txf_1/controls/<control>` | yes | Current channel value |
| `/devices/mtrf_txf_1/controls/<control>/meta` | yes | JSON `{"type", "title", "units", "order", "min", "max", "precision", "readonly"}` |
| `/devices/mtrf_txf_1/controls/<control>/meta/error` | yes | Channel error flag: `r`/`w`/`p`, or an empty string |
| `/devices/mtrf_txf_1/controls/<control>/on` | no | Command to change the value (published by the client/UI) |

Error flags (`meta/error`) use letter codes: `r` — read/communication error, `w` — write error, `p` — a missed poll period. An empty string means there is no error.

## Upgrading from 1.x to 2.0

Version 2.0 brings the service in line with the current revision of the [Wiren Board MQTT Conventions](https://github.com/wirenboard/conventions) and [WB-STD-001](https://github.com/wirenboard/wb-standards). The changes only affect the MQTT topic format and the `templates.json` file — **`devices.json` does not need any changes**.

**What changed:**
1. Device and control metadata is now published as a single JSON topic (`.../meta`, `.../controls/<c>/meta`) instead of a set of separate sub-topics (`meta/name`, `meta/type`, `meta/order`, `meta/min`, `meta/max`, `meta/units`, `meta/precision`, `meta/readonly`) — the old sub-topics are no longer published.
2. Devices and controls now have localized titles (`title`, en/ru) — previously only internal identifiers were published to MQTT.
3. `meta/error` topics are now retained and use letter error codes (`r`/`w`/`p`) instead of free-form text.
4. The specific control types `temperature`, `rel_humidity`, `power`, `power_consumption` were replaced with the generic `value` type plus a `units` key (these types are deprecated in the Wiren Board MQTT Conventions).
5. Channel identifiers were renamed to match the [WB-STD-001](https://github.com/wirenboard/wb-standards) canon:
   - `value` → `temperature` (current temperature of `srf-1-3000-t`)
   - `setting` → `setpoint_t_ambient` (thermostat setpoint of `srf-1-3000-t`)
   - `low_battery` → `alarm_low_battery`, type `switch` → `alarm` (low battery indicator)

**What you need to do:**
- If you use the stock `templates/templates.json` from this repository, just update it along with the binary — no further action is needed.
- If you copied and customized `templates.json` (added your own devices/templates), run it through the `configconvert` tool:
  ```shell
  # back up the file first
  cp /etc/wb-mqtt-noolite-templates.json /etc/wb-mqtt-noolite-templates.json.bak

  go run ./cmd/configconvert -in /etc/wb-mqtt-noolite-templates.json
  ```
  The converter replaces deprecated types and channel identifiers with the canonical ones and prints a list of the changes it made. Channels and templates without a `title` get a placeholder (equal to the identifier) — translate/fix it manually after conversion.
- If your automations (wb-rules, dashboards, etc.) reference the topics `.../controls/value`, `.../controls/setting` or `.../controls/low_battery` directly, update them to the new names (`temperature`, `setpoint_t_ambient`, `alarm_low_battery` respectively), and remove any direct dependency on `meta/name`, `meta/type` or other `meta/*` sub-topics if you were using them.

## mtrf_tool utility: pairing and unpairing devices

The package includes the `mtrf_tool` utility, which sends commands directly to the MTRF-64-USB-A adapter: pairing and unpairing devices, setting the thermostat temperature and sensor mode, turning the load on/off, and reading status.

### Building
```shell
go build -o mtrf_tool ./cmd/mtrf_tool
```

### Main parameters
| Flag | Description |
|---|---|
| `-d`, `--device` | Adapter serial port, default `/dev/ttyUSB0` |
| `-c`, `--channel` | Channel (pairing slot on the adapter), **1-63**. `0` cannot be used — in `mtrf_tool` it is reserved as "channel not set" (with `--channel 0` the tool just prints usage) |
| `-m`, `--mode` | Protocol and direction: `tx`/`txf` — power units controlled by the adapter; `rx`/`rxf` — switches, remotes and sensors the adapter receives commands from |
| `--command` | `bind`, `unbind`, `clear_all`, `on`, `off`, `status`, `status_output`, `thermostat_mode`, `poweron_state`, `temperature` |

**Important:** the channel passed to `mtrf_tool` (`-c`) must match the `ch` value of the corresponding device in `devices.json` — this is how the wb-mqtt-noolite service maps a physically paired device to its template and configuration.

The pairing scenarios below are based on the MTRF-64-USB-A operation manual (`docs/rukovodstvo-po-ekspluatatsii-mtrf-64-usb-a.pdf`, sections 5-6) and the `noolite/commands.go` code.

### Pairing power units (TX / TX-F): relays, dimmers, outlets, thermostat

Power units (relays, dimmers, the SRF-1-3000-T thermostat, etc.) are actuators controlled by the adapter. They use "manual pairing" (manual §5.1):

1. Put the power unit into pairing mode (usually by holding its service button; see that specific product's own manual for the exact procedure).
2. Pair it to a free channel, e.g. 1, using `txf` mode (nooLite-F) or `tx` (classic nooLite):
   ```shell
   ./mtrf_tool -d /dev/ttyUSB0 -c 1 -m txf --command bind
   ```
3. Successful pairing is indicated by the unit's LED turning off; `mtrf_tool`'s log will show the adapter's response with `CTR=3` ("pairing successful") and, for nooLite-F, the device's 32-bit address.
4. Once pairing is confirmed, `mtrf_tool` prints a ready-made `devices.json` fragment — just copy it and fill in `name` (the nooLite-F address is already filled in, and `template` is auto-filled when the device model is recognized):
   ```json
   {
     "noolite_type": "txf",
     "ch": 1,
     "name": "TODO",
     "address": "00017526",
     "template": "srf-1-3000-t"
   }
   ```
   The model is only recognized for a small set of stock templates — if it can't be recognized, `template` stays `"TODO"` and needs to be filled in manually.

Unpairing a power unit requires confirmation on the device itself (pressing its service button, manual §6.1/§6.2):
```shell
./mtrf_tool -d /dev/ttyUSB0 -c 1 -m txf --command unbind
```

### Pairing switches, remotes and sensors (RX / RX-F): handling button presses

Remotes (PU212-2), sensors (PT111, WS1) and similar devices transmit commands over the air on their own — the adapter must be switched to receive mode and "listen" on the channel for a limited time (manual §5.3):

1. Enable pairing on the channel the device will be paired to, e.g. 5, using `rx` mode (classic nooLite — the mode used by all the stock `pt111`, `ws1` and `pu212-2` templates):
   ```shell
   ./mtrf_tool -d /dev/ttyUSB0 -c 5 -m rx --command bind
   ```
2. The adapter enters a listening window of **40 seconds**. During that time, press the pairing service button/combination on the device being paired (the exact procedure depends on the model — see that product's own manual; typically it's holding a button until an indicator starts blinking in a characteristic pattern).
3. If pairing does not happen within 40 seconds, the adapter automatically exits the listening window — repeat the `bind` command if needed.
4. Once pairing is confirmed, `mtrf_tool` prints a `devices.json` fragment just like for power units above — for classic nooLite the address is not available (the protocol doesn't carry one), so the fragment has no `address` field and `template` must be filled in manually.
5. **Subsequent button presses are handled fully automatically**: the wb-mqtt-noolite service receives the adapter's commands (on/off, toggle, sensor state changes, low battery, etc.) with no further action required and immediately publishes them to the corresponding `/controls/<control>` topics per the device's template. Scheduled polling is not needed for such devices — that's why the stock `pu212-2`, `ws1` and `pt111` templates don't use it. See ["Incoming commands from switches and remotes (RX)"](#incoming-commands-from-switches-and-remotes-rx) below for exactly which remote commands are supported.

For nooLite-F devices (`rxf`), the manual does not describe a separate pairing procedure, but per the adapter's protocol it works the same way — `MODE=3` instead of `MODE=1`, except the adapter's response includes the device's 32-bit address; `mtrf_tool` supports this mode (`--mode rxf`).

Unpairing (clearing one channel) does not require confirmation on the device itself (manual §6.4):
```shell
./mtrf_tool -d /dev/ttyUSB0 -c 5 -m rx --command unbind
```

A full adapter memory wipe (all channels at once, manual §6.5) is done without `-c` — it only works for `rx`/`rxf` modes:
```shell
./mtrf_tool -d /dev/ttyUSB0 -m rx --command clear_all
```

## Incoming commands from switches and remotes (RX)

Besides explicit on/off commands, Noolite remotes and switches (e.g. PU212-2) can send other protocol commands depending on how a given button is programmed. The service supports:

- **Switch (toggle, CMD=4)** — the button doesn't report a target state, it just "toggles" the load. The service keeps the last known value of the `on` control and inverts it (`0↔1`) when a `Switch` command arrives, then publishes the new value as usual. This is how a single "on/off" button can be wired to one remote instead of separate "on" and "off" buttons.
- **Load_Preset (recall a scene, CMD=7)** — the button is programmed to recall a stored scene/preset. The value is published to a separate `preset` control (added to the `pu212-2` template) — this is exactly the case that's convenient to use as a scene-trigger button in wb-rules (e.g. "away"/"home", if different buttons on the remote are programmed with different presets).

  ⚠️ The MTRF-64-USB-A manual does not document the data format of the `Load_Preset` command for remotes — the `preset` control is published with the packet's `D0` byte as-is. What a given value actually means depends on how the remote was programmed, and should be verified empirically (watch the value in MQTT while pressing the button you care about) before wiring an automation to it.

Both commands are only supported in `rx` mode (classic nooLite, matching all other command reception in this version of the service) — the data format was not verified for `rxf`.
