# Flatrtree - Go

Flatrtree is a serialization format and set of libraries for reading and writing R-trees. It's directly inspired by [Flatbush](https://github.com/mourner/flatbush) and [FlatGeobuf](https://github.com/flatgeobuf/flatgeobuf), and aims to make tiny, portable R-trees accessible in new contexts.

- Store R-trees to disk or transport them over a network.
- Build R-trees in one language and query them in another.
- Query R-trees before reading the data they index.

## Installation

Requires Go 1.18 or later.

```console
$ go get github.com/flatrtree/flatrtree-go
```

## Usage

Flatrtree separates building and querying behavior. The builder doesn’t know how to query an index and the index doesn’t know how it was built. This is inspired by [FlatBuffers](https://google.github.io/flatbuffers/).

```golang

```

## Serialization

Flatrtree uses [protocol buffers](https://protobuf.dev/) for serialization, taking advantage of varint encoding to reduce the output size in bytes. There are many tradeoffs to explore for serialization and this seems like a good place to start. It wouldn’t be hard to roll your own format with something like FlatBuffers if that better fit your needs.

```golang

```
