#!/bin/sh

rm -r internal/mocks/
rm -r pkg/mocks/

mockgen -source=internal/core/iface.go -destination=internal/mocks/core/iface.go
mockgen -source=internal/repository/iface.go -destination=internal/mocks/repository/iface.go
mockgen -source=internal/repository/compression/iface.go -destination=internal/mocks/repository/compression/iface.go
mockgen -source=internal/repository/api/http/iface.go -destination=internal/mocks/repository/api/http/iface.go
mockgen -source=pkg/redis/pool_iface.go -destination=pkg/mocks/redis/mock_redis_pool.go

AWS_SDK_GO_FILES=()
for file in "$(go env GOMODCACHE)"/github.com/aws/aws-sdk-go@*
do
  AWS_SDK_GO_FILES+=("${file}")
done

AWS_SDK_GO_CURRENT=${AWS_SDK_GO_FILES[${#AWS_SDK_GO_FILES[@]}-1]}
mockgen -source="${AWS_SDK_GO_CURRENT}/service/sqs/sqsiface/interface.go" -destination=pkg/mocks/sqs/mock_sqs.go

REDIGO_FILES=()
for file in "$(go env GOMODCACHE)"/github.com/gomodule/redigo@*
do
  REDIGO_FILES+=("${file}")
done

REDIGO_CURRENT=${REDIGO_FILES[${#REDIGO_FILES[@]}-1]}
mockgen -source="${REDIGO_CURRENT}/redis/redis.go" -destination=pkg/mocks/redis/mock_redigo.go
