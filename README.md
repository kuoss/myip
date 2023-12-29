# myip
"What's my IP" application

[![release](https://github.com/kuoss/myip/actions/workflows/release.yml/badge.svg)](https://github.com/kuoss/myip/actions)
[![pull-request](https://github.com/kuoss/myip/actions/workflows/pull-request.yml/badge.svg)](https://github.com/kuoss/myip/actions)
[![Coverage Status](https://coveralls.io/repos/github/kuoss/myip/badge.svg?branch=main)](https://coveralls.io/github/kuoss/myip?branch=main)
[![codecov](https://codecov.io/gh/kuoss/myip/branch/main/graph/badge.svg?token=EXPE6OS8HJ)](https://codecov.io/gh/kuoss/myip)
[![GitHub license](https://img.shields.io/github/license/kuoss/myip.svg)](https://github.com/kuoss/myip/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/kuoss/myip.svg)](https://github.com/kuoss/myip/stargazers)
[![contribuiton welcome](https://img.shields.io/badge/contributions-welcome-orange.svg)](https://github.com/kuoss/myip/blob/main/CONTRIBUTING.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/kuoss/myip)](https://goreportcard.com/report/github.com/kuoss/myip)

env         | description                                 | default | example
----------- | ------------------------------------------- | ------- | -------
APP_ADDR    | the TCP address for the server to listen on | `:80`   | `:8080`
APP_PROXIES | comma-separated trusted proxies             |         | `10.0.0.0/8,192.168.1.33`
