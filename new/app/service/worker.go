package service

import (
	"bytes"
	"sort"
	"strconv"
)

//NewWorker returns a new service worker.
func NewWorker() *Worker {
	return &Worker{
		Assets: make(map[string]bool),
	}
}

//Worker is a service worker.
type Worker struct {
	Version string
	Assets  map[string]bool
}

func (worker Worker) renderMap(b *bytes.Buffer, mapping map[string]bool) {
	//Deterministic render
	keys := make([]string, 0, len(mapping))
	for i := range mapping {
		keys = append(keys, i)
	}
	sort.Strings(keys)

	var i = 0
	for key := range keys {
		asset := keys[key]

		if asset == "" {
			continue
		}

		b.WriteString(strconv.Quote(asset))
		if i < len(keys)-1 {
			b.WriteString(", ")
		}
		i++
	}
}

//Render the service worker to JS.
func (worker Worker) Render() []byte {
	var b bytes.Buffer

	b.WriteString(`const version = "`)
	b.WriteString(worker.Version)
	b.WriteString(`";`)

	b.WriteString(`self.addEventListener('install', function(event) {
		self.skipWaiting();
		caches.delete("dynamic");
  event.waitUntil(
    caches.open("assets").then(function(cache) {
      return cache.addAll(
        ["/", `)

	worker.renderMap(&b, worker.Assets)

	b.WriteString(`]
      );
    }).catch(function(e) {
		console.log("Couldn't install because: ", e);
	})
  );
});

self.addEventListener('fetch', event => event.respondWith(cacheThenNetwork(event)));

async function cacheThenNetwork(event) {
	let request = event.request;

	const assets = await caches.open("assets");

	//Try load a cached asset first.
	const CachedAsset = await assets.match(request);
	if (CachedAsset) return CachedAsset;

	//Get the request from the network.
	try {
		let clone = request.clone();
		clone.url = request.url+"?="+Math.random();
		const NetworkReponse = await fetch(clone, {cache: "no-store"});
		if (request.method == "GET" && NetworkReponse.status == 200) {
			const dynamic = await caches.open("dynamic");
			dynamic.put(request, NetworkReponse.clone());
		}
		return NetworkReponse;
	} catch (e) {
		//Try the dynamic cache.
		if (request.method == "GET") {
			const dynamic = await caches.open("dynamic");
			const CachedDynamic = await dynamic.match(request);
			if (CachedDynamic) return CachedDynamic;
		}

		return new Response("404 not found", {
			status: 404,
		})
	}
}
`)

	return b.Bytes()
}
