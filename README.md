PlurGo
======

Plurk API in GO

config.json
----
You will need to save ConsumerToken/ConsumerSecret in config.json, the
format is JSON. The example is below:

> {
>   "ConsumerToken": "<FILL ME IN>",
>   "ConsumerSecret": "<FILL ME IN>",
>   "AccessToken": "",
>   "AccessSecret": ""
> }

Example
----
```
% go run plurgo.go -config=config.json
```

Build Status
----
[![Build Status](https://travis-ci.org/clsung/plurgo.svg?branch=master)](https://travis-ci.org/clsung/plurgo)

Meta
----

* Code: `git clone git://github.com/clsung/plurgo.git`
* Home: <http://github.com/clsung/plurgo>
* Bugs: <http://github.com/clsung/plurgo/issues>

Author
------

Cheng-Lung Sung :: clsung@gmail.com :: @clsung
