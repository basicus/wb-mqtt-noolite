# Кросс-сборка wb-mqtt-noolite и mtrf_tool под контроллеры Wiren Board,
# и сборка .deb пакетов для Debian 12 (bookworm) / 13 (trixie).
#
# Соответствие ревизий контроллера и цели сборки (см. README.md):
#   armv5 - Wiren Board 3.5 и старше (ARM926EJ-S)
#   armv7 - Wiren Board rev. 6.3-6.6, 6.7, 7.4 (NXP i.MX 6ULL / Cortex-A7, 32 бита)
#   arm64 - Wiren Board 8.x (Cortex-A53, 64 бита)
#
# Использование:
#   make armv7             # собрать wb-mqtt-noolite и mtrf_tool под armv7 в build/armv7/
#   make all               # собрать под все три цели
#   make deb                # .deb для armhf и arm64, Debian 12 и 13, в build/deb/
#   make clean              # удалить build/

BUILD_DIR := build
DEB_DIR := $(BUILD_DIR)/deb
VERSION := $(shell cat VERSION 2>/dev/null || echo 0.0.0)
# Suite'ы Debian, для которых собирается пакет - см. packaging/build-deb.sh
DEB_SUITES := bookworm trixie

# CGO_ENABLED=0 - все зависимости чистые Go, статическая линковка без зависимости от libc
# целевой системы (не нужен ARM-тулчейн с C-кросс-компилятором).
# -trimpath - убрать из бинарника локальные пути сборки.
# -s -w - убрать таблицу символов и отладочную информацию DWARF (см. README, "Оптимизация").
export CGO_ENABLED := 0
GOFLAGS := -trimpath
LDFLAGS := -s -w

.PHONY: all armv5 armv7 arm64 clean build-target deb deb-armhf deb-arm64

all: armv5 armv7 arm64

armv5:
	@$(MAKE) --no-print-directory build-target TARGET=armv5 GOARCH=arm GOARM=5

armv7:
	@$(MAKE) --no-print-directory build-target TARGET=armv7 GOARCH=arm GOARM=7

arm64:
	@$(MAKE) --no-print-directory build-target TARGET=arm64 GOARCH=arm64 GOARM=

# build-target - внутренняя цель, собирает оба бинарника под TARGET/GOARCH/GOARM.
# Не вызывать напрямую - используйте armv5/armv7/arm64.
build-target:
	@mkdir -p $(BUILD_DIR)/$(TARGET)
	GOOS=linux GOARCH=$(GOARCH) GOARM=$(GOARM) go build $(GOFLAGS) -ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(TARGET)/wb-mqtt-noolite ./cmd/wb-mqtt-noolite
	GOOS=linux GOARCH=$(GOARCH) GOARM=$(GOARM) go build $(GOFLAGS) -ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(TARGET)/mtrf_tool ./cmd/mtrf_tool
	@echo "--> $(BUILD_DIR)/$(TARGET)/"
	@ls -la $(BUILD_DIR)/$(TARGET)/

deb: deb-armhf deb-arm64

# Debian arch "armhf" собирается из armv7-сборки (тот же 32-битный Cortex-A7,
# просто другое имя таргета в Debian, см. README, раздел про .deb пакеты).
deb-armhf: armv7
	@for suite in $(DEB_SUITES); do \
		packaging/build-deb.sh armhf $$suite $(BUILD_DIR)/armv7 $(VERSION) $(DEB_DIR) || exit 1; \
	done

deb-arm64: arm64
	@for suite in $(DEB_SUITES); do \
		packaging/build-deb.sh arm64 $$suite $(BUILD_DIR)/arm64 $(VERSION) $(DEB_DIR) || exit 1; \
	done

clean:
	rm -rf $(BUILD_DIR)
