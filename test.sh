#!/usr/bin/env bash
# 运行全部单元测试脚本
# 用法: ./test.sh
set -u
cd "$(dirname "$0")"

# 部分环境下默认 go build 缓存目录不可写，自动回退到临时目录
if ! go env GOCACHE >/dev/null 2>&1 || [ ! -w "$(go env GOCACHE)" ]; then
	export GOCACHE="$(mktemp -d)/gocache"
	echo "GOCACHE 不可写，回退到: $GOCACHE"
fi

failed=0

echo "==> [1/4] go vet 静态检查"
go vet ./... || failed=1

echo "==> [2/4] 并发修复相关包竞态检测: utils downloader"
go test -race -count=1 -v ./utils/... ./downloader/... || failed=1

echo "==> [3/4] 其他本地包单元测试: config parser request"
go test -count=1 -v ./config/... ./parser/... ./request/... || failed=1

echo "==> [4/4] 全部单元测试（含 extractors，部分用例需要网络）"
go test -count=1 ./... || failed=1

if [ "$failed" -ne 0 ]; then
	echo "存在失败的测试，请检查上方输出"
	exit 1
fi
echo "全部测试通过"
