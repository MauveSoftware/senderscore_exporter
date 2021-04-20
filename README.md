[![CircleCI](https://circleci.com/gh/MauveSoftware/senderscore_exporter.svg?style=shield)](https://circleci.com/gh/MauveSoftware/senderscore_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/mauvesoftware/senderscore_exporter)](https://goreportcard.com/report/github.com/mauvesoftware/senderscore_exporter)
[![Docker Build Status](https://img.shields.io/cloud/docker/build/mauvesoftware/senderscore_exporter.svg)](https://hub.docker.com/r/mauvesoftware/senderscore_exporter/builds)


# senderscore_exporter
Metrics exporter for senderscore.org scores to prometheus

## Install
```
go get -u github.com/MauveSoftware/senderscore_exporter
```

## Configuration
This is an example for a senderscore_exporter configuration file

```yaml
addresses:
- 185.138.52.255
- 2a07:a40:c0de::1
```

## Usage

### Binary
```bash
./senderscore_exporter
```

### Docker
```bash
docker run -d --restart always --name senderscore_exporter -v /etc/senderscore_exporter:/config -p 9665:9665 mauvesoftware/senderscore_exporter
```

## License
(c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under [MIT](LICENSE) license.

## Prometheus
see https://prometheus.io/

## Senderscore
see https://senderscore.org/
