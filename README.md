# KRE Runners

This repo contains all language runners for [KRE](https://github.com/konstellation-io/kre).

## Runners coverage

|      Component       |                               Coverage                               |                           Bugs                           |               Maintainability Rating               |
| :------------------: | :------------------------------------------------------------------: | :------------------------------------------------------: | :------------------------------------------------: |
|    KRE Entrypoint    | [![coverage][kre-entrypoint-coverage]][kre-entrypoint-coverage-link] | [![bugs][kre-entrypoint-bugs]][kre-entrypoint-bugs-link] | [![mr][kre-entrypoint-mr]][kre-entrypoint-mr-link] |
|        KRE Go        |         [![coverage][kre-go-coverage]][kre-go-coverage-link]         |         [![bugs][kre-go-bugs]][kre-go-bugs-link]         |         [![mr][kre-go-mr]][kre-go-mr-link]         |
|        KRE Py        |         [![coverage][kre-py-coverage]][kre-py-coverage-link]         |         [![bugs][kre-py-bugs]][kre-py-bugs-link]         |         [![mr][kre-py-mr]][kre-py-mr-link]         |
| KRT Files Downloader |         [![coverage][krt-fd-coverage]][krt-fd-coverage-link]         |         [![bugs][krt-fd-bugs]][krt-fd-bugs-link]         |         [![mr][krt-fd-mr]][krt-fd-mr-link]         |

[kre-py-coverage]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_py&metric=coverage
[kre-py-coverage-link]: https://sonarcloud.io/dashboard?id=konstellation-io_kre_py
[kre-py-bugs]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_py&metric=bugs
[kre-py-bugs-link]: https://sonarcloud.io/dashboard?id=konstellation-io_kre_py
[kre-py-loc]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_py&metric=ncloc
[kre-py-loc-link]: https://sonarcloud.io/dashboard?id=konstellation-io_kre_py
[kre-py-mr]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_py&metric=sqale_rating
[kre-py-mr-link]: https://sonarcloud.io/dashboard?id=konstellation-io_kre_py
[kre-go-coverage]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_go&metric=coverage
[kre-go-coverage-link]: https://sonarcloud.io/dashboard?id=konstellation-io_kre_go
[kre-go-bugs]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_go&metric=bugs
[kre-go-bugs-link]: https://sonarcloud.io/dashboard?id=konstellation-io_kre_go
[kre-go-loc]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_go&metric=ncloc
[kre-go-loc-link]: https://sonarcloud.io/dashboard?id=konstellation-io_kre_go
[kre-go-mr]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_go&metric=sqale_rating
[kre-go-mr-link]: https://sonarcloud.io/dashboard?id=konstellation-io_kre_go
[kre-entrypoint-coverage]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_entrypoint&metric=coverage
[kre-entrypoint-coverage-link]: https://sonarcloud.io/dashboard?id=konstellation-io_kre_entrypoint
[kre-entrypoint-bugs]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_entrypoint&metric=bugs
[kre-entrypoint-bugs-link]: https://sonarcloud.io/dashboard?id=konstellation-io_kre_entrypoint
[kre-entrypoint-loc]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_entrypoint&metric=ncloc
[kre-entrypoint-loc-link]: https://sonarcloud.io/dashboard?id=konstellation-io_kre_entrypoint
[kre-entrypoint-mr]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_kre_entrypoint&metric=sqale_rating
[kre-entrypoint-mr-link]: https://sonarcloud.io/dashboard?id=konstellation-io_kre_entrypoint
[krt-fd-coverage]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_krt_files_downloader&metric=coverage
[krt-fd-coverage-link]: https://sonarcloud.io/dashboard?id=konstellation-io_krt_files_downloader
[krt-fd-bugs]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_krt_files_downloader&metric=bugs
[krt-fd-bugs-link]: https://sonarcloud.io/dashboard?id=konstellation-io_krt_files_downloader
[krt-fd-loc]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_krt_files_downloader&metric=ncloc
[krt-fd-loc-link]: https://sonarcloud.io/dashboard?id=konstellation-io_krt_files_downloader
[krt-fd-mr]: https://sonarcloud.io/api/project_badges/measure?project=konstellation-io_krt_files_downloader&metric=sqale_rating
[krt-fd-mr-link]: https://sonarcloud.io/dashboard?id=konstellation-io_krt_files_downloader

## Protobuf

All components receive and send a `KreNatsMessage` protobuf.
To generate the protobuf code, the `protoc` compiler must be installed.
Use the following command to generate the code:

```
./generate_protobuf_code.sh
```
