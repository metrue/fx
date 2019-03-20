fx
------
            ┌────────┐
            │fx init │       fx━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
            └────────┘       ┃          ┌───────────────────────┐                                   ┃
     ────────────────────────╋─────────▶│Environment initialize │                                   ┃
            ┌──────┐         ┃          │* proxy docker sock    │                                   ┃
            │fx up │         ┃          │* pull fx base docker  │                                   ┃
    ┌ ─ ─ ─ ┴──────┘─ ─ ┐    ┃          └───────────────────────┘                                   ┃
       Function Source       ┃          ┌──────────────┐       ┌─────────────────────────────┐      ┃
    └ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┘────╋──┬──────▶│     Pack     │       │                             │      ┃
                             ┃  │       └──────┬───────┘       │                             │      ┃
            ┌────────┐       ┃  │       ┌──────▼───────┐       │                             │      ┃
            │fx call │       ┃  │       │Build Service │◀─────▶│                             │      ┃
            └────────┘       ┃  │       └──────┬───────┘       │                             │      ┃
    ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┐    ┃  │       ┌──────▼───────┐       │                             │      ┃
       Function Source       ┃  │       │ Run Service  │◀─────▶│                             │      ┃
    │   (with params)   │────╋──┤       └──────────────┘       │                             │      ┃
     ─ ─ ─ ─ ─ ─ ─ ─ ─ ─     ┃  │                              │                             │      ┃
                             ┃  │                              │                             │      ┃
                             ┃  │       ┌──────────────┐       │         Docker API          │      ┃
           ┌────────┐        ┃  └──────▶│ Call Service │       │                             │      ┃
           │fx down │        ┃          │    (http)    │       │                             │      ┃
           └────────┘        ┃          └──────────────┘       │                             │      ┃
     ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─     ┃          ┌──────────────┐       │                             │      ┃
         Service Name   │────╋─────────▶│ Stop Service │◀─────▶│                             │      ┃
     └ ─ ─ ─ ─ ─ ─ ─ ─ ─     ┃          └──────────────┘       │                             │      ┃
          ┌────────┐         ┃                                 │                             │      ┃
          │fx list │         ┃                                 │                             │      ┃
          └────────┘         ┃                                 │                             │      ┃
     ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─     ┃          ┌──────────────┐       │                             │      ┃
         Service Name   │────╋─────────▶│List Services │◀─────▶│                             │      ┃
     └ ─ ─ ─ ─ ─ ─ ─ ─ ─     ┃          └──────────────┘       └─────────────────────────────┘      ┃
                             ┃                                                                      ┃
                             ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛


Poor man's function as a service.
<br/>
![build](https://circleci.com/gh/metrue/fx.svg?style=svg&circle-token=bd62abac47802f8504faa4cf8db43e4f117e7cd7)
[![codecov](https://codecov.io/gh/metrue/fx/branch/master/graph/badge.svg)](https://codecov.io/gh/metrue/fx)
[![Go Report Card](https://goreportcard.com/badge/github.com/metrue/fx?style=flat-square)](https://goreportcard.com/report/github.com/metrue/fx)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/metrue/fx)
![](https://img.shields.io/github/license/metrue/fx.svg)
[![Release](https://img.shields.io/github/release/metrue/fx.svg?style=flat-square)](https://github.com/metrue/fx/releases/latest)

fx is a tool to help you do Function as a Service on your own server. fx can make your stateless function a service in seconds. The most exciting thing is that you can write your functions with most programming languages,  Installation, Usage and Contribution, see [wiki](https://github.com/metrue/fx/wiki)

| Language      | Status        | Contributor   |
| ------------- |:-------------:|:-------------:|
| Go            | Supported     | fx            |
| Node          | Supported     | fx            |
| Python        | Supported     | fx            |
| Ruby          | Supported     | fx            |
| Java          | Supported     | fx            |
| PHP           | Supported     | [@chlins](https://github.com/chlins)|
| Julia         | Supported     | [@mbesancon](https://github.com/mbesancon)|
| D             | Supported     | [@andre2007](https://github.com/andre2007)|
| R             | Working on [need your help](https://github.com/metrue/fx/issues/31)   | |

You can refer to the [doc](https://github.com/metrue/fx/blob/master/NEW_LANGUAGE_SUPPORT.md) to make fx support the language not listed above. Welcome to tweet [me](https://twitter.com/_metrue) or [Buy me a coffee](https://www.paypal.me/minghe).

Thank you to all the people who already contributed to fx!

<table>
  <tbody>
    <tr>
        <a href="https://github.com/metrue" target="_blank">
            <img alt="metrue" src="https://avatars2.githubusercontent.com/u/1001246?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/pplam" target="_blank">
            <img alt="pplam" src="https://avatars2.githubusercontent.com/u/12783579?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/muka" target="_blank">
            <img alt="muka" src="https://avatars2.githubusercontent.com/u/1021269?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/xwjdsh" target="_blank">
            <img alt="xwjdsh" src="https://avatars2.githubusercontent.com/u/11025519?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/mbesancon" target="_blank">
            <img alt="mbesancon" src="https://avatars2.githubusercontent.com/u/7623090?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/avelino" target="_blank">
            <img alt="avelino" src="https://avatars2.githubusercontent.com/u/31996?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/DaidoujiChen" target="_blank">
            <img alt="DaidoujiChen" src="https://avatars0.githubusercontent.com/u/670441?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/chlins" target="_blank">
            <img alt="chlins" src="https://avatars2.githubusercontent.com/u/31262637?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/andre2007" target="_blank">
            <img alt="andre2007" src="https://avatars1.githubusercontent.com/u/1451047?s=50&v=4" width="50">
        </a>
        <a href="https://github.com/steventhanna" target="_blank">
            <img alt="andre2007" src="https://avatars1.githubusercontent.com/u/2541678?s=50&v=4" width="50">
        </a>
    </tr>
  </tbody>
</table>
