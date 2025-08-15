#!/usr/bin/env pwsh

# $NATS_CONTAINER_NAME = $args[0]
# $SLI_NATS_CONTAINER_NAME = $args[1]

$NATS_CONTAINER_NAME = "sl-nats"
$SLI_NATS_CONTAINER_NAME = "sl-nats-box"

docker cp "config/nats-server.conf" ${NATS_CONTAINER_NAME}:"nats-server.conf"
docker restart $NATS_CONTAINER_NAME

docker exec -it $SLI_NATS_CONTAINER_NAME mkdir /config

$jetstream = Get-Content -Path "config\jetstream_config.json" | ConvertFrom-Json
foreach ($item in $jetstream) {
    $configValue = $item.config
    $server = $item.server
    
    $newJsonContent = $configValue | ConvertTo-Json

    $configFile = "jetstream_config.json"
    $newJsonContent | Out-File -FilePath $configFile -Encoding utf8

    docker cp "jetstream_config.json" ${SLI_NATS_CONTAINER_NAME}:"/config/jetstream_config.json"
    docker exec -it $SLI_NATS_CONTAINER_NAME sed -i '1s/^.//' /config/jetstream_config.json

    Remove-Item $configFile

    docker exec -it $SLI_NATS_CONTAINER_NAME nats --server=${server}:4222 stream add --config=/config/jetstream_config.json
    # docker exec -it $SLI_NATS_CONTAINER_NAME rm /config/jetstream_config.json
}


# $consumer = Get-Content -Path "jetstream\jetstream_consumer_config.json" | ConvertFrom-Json
# foreach ($item in $consumer) {
#     $configValue = $item.details.config
#     $server = $item.server
#     $streamName = $item.details.stream_name
#     $name = $item.details.name

#     $newJsonContent = $configValue | ConvertTo-Json

#     $configFile = "jetstream_consumer_config.json"
#     $newJsonContent | Out-File -FilePath $configFile -Encoding utf8

#     docker cp "jetstream_consumer_config.json" ${SLI_NATS_CONTAINER_NAME}:"/config/jetstream_consumer_${name}.json"
#     docker exec -it $SLI_NATS_CONTAINER_NAME sed -i '1s/^.//' /config/jetstream_consumer_${name}.json

#     Remove-Item $configFile

#     docker exec -it $SLI_NATS_CONTAINER_NAME nats --server=${server}:4222 consumer add $streamName --config=/config/jetstream_consumer_${name}.json
#     # docker exec -it $SLI_NATS_CONTAINER_NAME rm /config/jetstream_consumer_${name}.json
# }

# docker exec -it $SLI_NATS_CONTAINER_NAME rm -r /config

Write-Host "completed"