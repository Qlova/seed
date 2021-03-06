package popup

import (
	"qlova.org/seed"
	"qlova.org/seed/client"
)

func init() {
	client.RegisterRenderer(func(c seed.Seed) []byte {
		return []byte(`

seed.CurrentPopup = null;

seed.show = async function(id, args) {
	let popup = q.get(id);
	if (!popup) {
		console.error("seed.show: invalid popup ", id);
		return;
	}
	popup.template = popup.parent;

	popup.parent.parentElement.appendChild(popup);

	seed.CurrentPopup = popup;
	popup.args = args;

	if (popup.onshow) await popup.onshow();

	let promises = [];

	if (seed.goto.in) {
		promises.push(seed.goto.in);
		seed.goto.in = null;
	}

	if (seed.goto.out) {
		promises.push(seed.goto.out);
		seed.goto.out = null;
	}

	for (let promise of promises) {
		await promise;
	}
};

seed.hide = async function(id) {
	let popup = q.get(id);
	if (!popup) {
		console.error("seed.show: invalid popup ", id);
		return;
	}

	if (!popup.template) return;

	if (popup.onhide) await popup.onhide();

	let promises = [];

	if (seed.goto.in) {
		promises.push(seed.goto.in);
		seed.goto.in = null;
	}

	if (seed.goto.out) {
		promises.push(seed.goto.out);
		seed.goto.out = null;
	}

	for (let promise of promises) {
		await promise;
	}

	seed.CurrentPopup = null;

	popup.template.content.appendChild(popup);
};`)
	})
}
