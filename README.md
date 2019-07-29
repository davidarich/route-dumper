#route-dumper

Statically analyzes [Symfony](https://github.com/symfony/symfony) Controller classes to extract route annotations then dumps to JSON. Intended for use with [FriendsOfSymfony/FOSJsRoutingBundle](https://github.com/FriendsOfSymfony/FOSJsRoutingBundle) as a container friendly `fos:js-routing:dump --format=json` replacement. 

## Usage
```
Usage:
  route-dumper [PATH] [flags]

Flags:
  -h, --help   help for route-dumper

requires path argument

```