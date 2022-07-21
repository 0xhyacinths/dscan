# DScan
A decentralized blockchain explorer!

***This is very early, so please bear with me as things are in flux.***

With the adoption of blockchain technology, it has become important to be able
to read the current state of the blockchain. DScan is a project to become a
decentralized explorer for Ethereum and Ethereum-Compatible chains with strong
delineation between what can be done entirely in-browser and what requires a
separate indexer server.

Allowing users to build their own indexer or even UIs that utilize these indexer
services around a specified API will provide many options for a block explorer
instead of one or two well-known ones.

## UI
The UI is just a bunch of static HTML, JS, and CSS. You can throw it behind any
web server that can handle this. However, because of how the UI works, data is
passed as query parameters. This means that it may not behave as expected when
using it with Metamask or other extensions that provide block explorer links.

This is still in the works, but there will be a server that will handle serving
the UI and managing the endpoints correctly so it can be integrated with these
extensions.


## Index Server
***This is still a to-do item***

Many explorers rely on continuous monitoring of the chain and as a result end up
requiring some form of central service. In existing implementations, the line
between the things that can be retrieved with an RPC and things that require an
indexer is a bit blurred. This project hopes to also delineate what can and
cannot be retrieved with a simple in-browser provider and provide a consistent
specification on retrieving that data from any indexer that they want to use.

This specification, as it materializes, will be available in a separate repo as
a gRPC proto file. Using gRPC gateways to facilitate browser use is a
[well-defined model][1] and also allows for developers to generate code to
interact with indexers or write their own indexers as they want.

### Caveats
Indexer operators may tamper with the data provided by their indexer, but they
cannot tamper with the blockchain itself. You should only use indexers you trust
to not secretly modify data as it may result in misleading information that does
not accurately reflect on-chain state.


## Notes
The CBOR library used does not work correctly normally. It needs to be patched
as per [this issue][2]. This is done automatically when you run `npm install`.

[1]: https://github.com/grpc-ecosystem/grpc-gateway
[2]: https://github.com/dignifiedquire/borc/issues/49
