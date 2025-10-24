package components

import (
	"github.com/hypercommithq/libhtml"
	"github.com/hypercommithq/libhtml/attr"
)

func Celebration() html.Node {
	return html.Element("script",
		attr.Type("module"),
		html.Text(`
import confetti from 'https://cdn.jsdelivr.net/npm/canvas-confetti@1.9.3/+esm';

(function() {
	const duration = 3000;
	const animationEnd = Date.now() + duration;
	const defaults = { startVelocity: 30, spread: 360, ticks: 60, zIndex: 9999 };

	function randomInRange(min, max) {
		return Math.random() * (max - min) + min;
	}

	const interval = setInterval(function() {
		const timeLeft = animationEnd - Date.now();

		if (timeLeft <= 0) {
			return clearInterval(interval);
		}

		const particleCount = 50 * (timeLeft / duration);

		confetti({
			...defaults,
			particleCount,
			origin: { x: randomInRange(0.1, 0.3), y: Math.random() - 0.2 }
		});
		confetti({
			...defaults,
			particleCount,
			origin: { x: randomInRange(0.7, 0.9), y: Math.random() - 0.2 }
		});
	}, 250);
})();
		`),
	)
}
